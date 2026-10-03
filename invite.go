// Package invite 是对局邀请：创建去重、有效期、人数上限与接受幂等。
// 缺陷：同一个邀请码可以被重复创建、邀请不过期、人数上限不生效、
// 接受不幂等、查询遍历全部邀请。
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
	members    map[string][]string
	memo       map[string]bool
	scanned    int
}

// New 建邀请台账。
func New(ttlMs int64, maxPlayers int) *Invites {
	return &Invites{ttlMs: ttlMs, maxPlayers: maxPlayers, codes: map[string]entry{},
		members: map[string][]string{}, memo: map[string]bool{}}
}

// Create 创建一个邀请。
func (i *Invites) Create(code, from, room string, nowMs int64) bool {
	i.scanned += len(i.codes)
	i.codes[code] = entry{from: from, room: room, at: nowMs}
	return true
}

// Accept 接受一个邀请。
func (i *Invites) Accept(code, player string, nowMs int64, reqID string) bool {
	i.scanned += len(i.codes)
	i.memo[reqID] = true
	record, ok := i.codes[code]
	if !ok {
		return false
	}
	i.members[record.room] = append(i.members[record.room], player)
	return true
}

// Pending 返回还在有效期内的邀请码。
func (i *Invites) Pending(nowMs int64) []string {
	i.scanned += len(i.codes)
	out := []string{}
	for code := range i.codes {
		out = append(out, code)
	}
	return out
}

// Members 返回房间当前人数。
func (i *Invites) Members(room string) int { return len(i.members[room]) }

// Scanned 是累计扫描量（规模观测）。
func (i *Invites) Scanned() int { return i.scanned }
