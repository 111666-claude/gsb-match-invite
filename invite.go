// Package invite 是对局邀请：创建去重、有效期、人数上限与接受幂等。
package invite

import (
	"container/heap"
	"sort"
)

type entry struct {
	from string
	room string
	at   int64
}

// expiryItem 是过期队列元素，按创建时间排序。
type expiryItem struct {
	at   int64
	code string
}

// expiryHeap 是按创建时间排列的小顶堆，用于摊还 O(1) 的惰性过期。
type expiryHeap []expiryItem

func (h expiryHeap) Len() int           { return len(h) }
func (h expiryHeap) Less(a, b int) bool { return h[a].at < h[b].at }
func (h expiryHeap) Swap(a, b int)      { h[a], h[b] = h[b], h[a] }
func (h *expiryHeap) Push(x any)        { *h = append(*h, x.(expiryItem)) }
func (h *expiryHeap) Pop() any {
	old := *h
	item := old[len(old)-1]
	*h = old[:len(old)-1]
	return item
}

// Invites 是邀请台账。
type Invites struct {
	ttlMs      int64
	maxPlayers int
	codes      map[string]entry
	members    map[string]int
	memo       map[string]bool
	players    map[string]map[string]bool
	queue      expiryHeap
	scanned    int
}

// New 建邀请台账。
func New(ttlMs int64, maxPlayers int) *Invites {
	return &Invites{ttlMs: ttlMs, maxPlayers: maxPlayers, codes: map[string]entry{},
		members: map[string]int{}, memo: map[string]bool{},
		players: map[string]map[string]bool{}}
}

// gc 惰性清理在 nowMs 已过期的邀请，每个邀请至多被清理一次。
func (i *Invites) gc(nowMs int64) {
	for len(i.queue) > 0 && i.queue[0].at+i.ttlMs <= nowMs {
		item := heap.Pop(&i.queue).(expiryItem)
		i.scanned++
		if record, ok := i.codes[item.code]; ok && record.at == item.at {
			delete(i.codes, item.code)
			delete(i.players, item.code)
		}
	}
}

// Create 创建一个邀请；邀请码已存在（未过期）时返回 false 且不改动原邀请。
func (i *Invites) Create(code, from, room string, nowMs int64) bool {
	i.scanned++
	i.gc(nowMs)
	if _, ok := i.codes[code]; ok {
		return false
	}
	i.codes[code] = entry{from: from, room: room, at: nowMs}
	heap.Push(&i.queue, expiryItem{at: nowMs, code: code})
	return true
}

// Accept 接受一个邀请；邀请失效、房间已满、同人重复或 reqID 重复时返回 false。
func (i *Invites) Accept(code, player string, nowMs int64, reqID string) bool {
	i.scanned++
	i.gc(nowMs)
	if i.memo[reqID] {
		return false
	}
	record, ok := i.codes[code]
	if !ok {
		return false
	}
	if i.players[code][player] {
		return false
	}
	if i.members[record.room] >= i.maxPlayers {
		return false
	}
	i.memo[reqID] = true
	if i.players[code] == nil {
		i.players[code] = map[string]bool{}
	}
	i.players[code][player] = true
	i.members[record.room]++
	return true
}

// Pending 返回还在有效期内的邀请码。
func (i *Invites) Pending(nowMs int64) []string {
	i.scanned++
	i.gc(nowMs)
	out := make([]string, 0, len(i.codes))
	for code := range i.codes {
		out = append(out, code)
	}
	sort.Strings(out)
	return out
}

// Members 返回房间当前人数。
func (i *Invites) Members(room string) int { return i.members[room] }

// Scanned 是累计扫描量（规模观测）。
func (i *Invites) Scanned() int { return i.scanned }
