// Package invite 是对局邀请：创建去重、有效期、人数上限与接受幂等。
package invite

type entry struct {
	from string
	room string
	at   int64
}

// Invites 是邀请台账。
type Invites struct {
	ttlMs      int64
	maxPlayers int
	codes      map[string]entry
	order      []string // 按创建时间排列的待处理邀请码，头部惰性弹出过期项
	head       int
	members    map[string][]string
	accepted   map[string]map[string]bool // code -> 已接受的玩家
	memo       map[string]bool            // 已生效的 reqID
	scanned    int
}

// New 建邀请台账。
func New(ttlMs int64, maxPlayers int) *Invites {
	return &Invites{ttlMs: ttlMs, maxPlayers: maxPlayers, codes: map[string]entry{},
		members: map[string][]string{}, accepted: map[string]map[string]bool{},
		memo: map[string]bool{}}
}

// Create 创建一个邀请；同一个邀请码只创建一次，重复创建返回 false 且不改动原邀请。
func (i *Invites) Create(code, from, room string, nowMs int64) bool {
	if _, ok := i.codes[code]; ok {
		return false
	}
	i.codes[code] = entry{from: from, room: room, at: nowMs}
	i.order = append(i.order, code)
	return true
}

// Accept 接受一个邀请。过期、满员、同人重复、reqID 重复都返回 false。
func (i *Invites) Accept(code, player string, nowMs int64, reqID string) bool {
	if i.memo[reqID] {
		return false
	}
	i.memo[reqID] = true
	record, ok := i.codes[code]
	if !ok || nowMs-record.at >= i.ttlMs {
		return false
	}
	if i.accepted[code][player] {
		return false
	}
	if len(i.members[record.room]) >= i.maxPlayers {
		return false
	}
	if i.accepted[code] == nil {
		i.accepted[code] = map[string]bool{}
	}
	i.accepted[code][player] = true
	i.members[record.room] = append(i.members[record.room], player)
	return true
}

// Pending 返回还在有效期内的邀请码。只从头部弹出过期项，摊还 O(1)。
func (i *Invites) Pending(nowMs int64) []string {
	for i.head < len(i.order) {
		i.scanned++
		record := i.codes[i.order[i.head]]
		if nowMs-record.at < i.ttlMs {
			break
		}
		delete(i.codes, i.order[i.head])
		i.head++
	}
	if i.head == len(i.order) {
		i.order = nil
		i.head = 0
	}
	out := make([]string, 0, len(i.order)-i.head)
	return append(out, i.order[i.head:]...)
}

// Members 返回房间当前人数。
func (i *Invites) Members(room string) int { return len(i.members[room]) }

// Scanned 是累计扫描量（规模观测）。
func (i *Invites) Scanned() int { return i.scanned }
