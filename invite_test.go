package invite

import "testing"

func TestEmptyPending(t *testing.T) {
	if len(New(1000, 4).Pending(0)) != 0 {
		t.Fatal("没有邀请时不该有待处理码")
	}
}

func TestCreateThenPending(t *testing.T) {
	world := New(1000, 4)
	world.Create("c1", "p1", "room", 0)
	if len(world.Pending(0)) != 1 {
		t.Fatal("创建之后应该有一个待处理邀请")
	}
}

func TestAcceptUnknownCode(t *testing.T) {
	if New(1000, 4).Accept("ghost", "p1", 0, "r1") {
		t.Fatal("不存在的邀请码不该能接受")
	}
}

func TestMembersStartsAtZero(t *testing.T) {
	if New(1000, 4).Members("room") != 0 {
		t.Fatal("房间初始人数应该是 0")
	}
}

func TestScannedStartsAtZero(t *testing.T) {
	if New(1000, 4).Scanned() != 0 {
		t.Fatal("Scanned 初始应该是 0")
	}
}

func TestCreateDuplicateKeepsOriginal(t *testing.T) {
	world := New(1000, 4)
	if !world.Create("c1", "p1", "room-a", 0) {
		t.Fatal("首次创建应该成功")
	}
	if world.Create("c1", "p9", "room-b", 10) {
		t.Fatal("重复创建应该返回 false")
	}
	if !world.Accept("c1", "p2", 0, "r1") {
		t.Fatal("重复创建不该覆盖原邀请")
	}
	if world.Members("room-a") != 1 || world.Members("room-b") != 0 {
		t.Fatalf("重复创建不该改动原邀请：a=%d b=%d", world.Members("room-a"), world.Members("room-b"))
	}
}

func TestTTLExpiry(t *testing.T) {
	world := New(1000, 4)
	world.Create("c1", "p1", "room", 0)
	if len(world.Pending(999)) != 1 {
		t.Fatal("ttl 前一刻邀请应该有效")
	}
	if len(world.Pending(1000)) != 0 {
		t.Fatal("达到 ttl 后 Pending 不该返回该邀请")
	}
	if world.Accept("c1", "p2", 1000, "r1") {
		t.Fatal("达到 ttl 后邀请不可接受")
	}
}

func TestRoomCapacity(t *testing.T) {
	world := New(1000, 2)
	world.Create("c1", "p1", "room", 0)
	if !world.Accept("c1", "p2", 0, "r1") || !world.Accept("c1", "p3", 0, "r2") {
		t.Fatal("未满员时接受应该成功")
	}
	if world.Accept("c1", "p4", 0, "r3") {
		t.Fatal("达到人数上限后 Accept 应该返回 false")
	}
	if members := world.Members("room"); members != 2 {
		t.Fatalf("房间人数不该超过上限：%d", members)
	}
}

func TestSamePlayerAcceptsOnce(t *testing.T) {
	world := New(1000, 4)
	world.Create("c1", "p1", "room", 0)
	if !world.Accept("c1", "p2", 0, "r1") {
		t.Fatal("首次接受应该成功")
	}
	if world.Accept("c1", "p2", 0, "r2") {
		t.Fatal("同一玩家重复接受应该返回 false")
	}
	if world.Members("room") != 1 {
		t.Fatalf("重复接受不该增加人数：%d", world.Members("room"))
	}
}

func TestSameReqIDAppliesOnce(t *testing.T) {
	world := New(1000, 4)
	world.Create("c1", "p1", "room", 0)
	if !world.Accept("c1", "p2", 0, "same") {
		t.Fatal("首次请求应该生效")
	}
	if world.Accept("c1", "p3", 0, "same") {
		t.Fatal("同一 reqID 不该再次生效")
	}
	if world.Members("room") != 1 {
		t.Fatalf("重复 reqID 不该增加人数：%d", world.Members("room"))
	}
}

func TestScannedDoesNotGrowWithInviteCount(t *testing.T) {
	world := New(100000, 4)
	for index := 0; index < 100000; index++ {
		world.Create(codeFor(index), "p1", "room", 0)
	}
	for index := 0; index < 3000; index++ {
		world.Pending(0)
	}
	if scanned := world.Scanned(); scanned > 6000 {
		t.Fatalf("Scanned 不该随邀请数线性增长：%d", scanned)
	}
}

func codeFor(index int) string {
	if index == 0 {
		return "c0"
	}
	digits := []byte{}
	for index > 0 {
		digits = append([]byte{byte('0' + index%10)}, digits...)
		index /= 10
	}
	return "c" + string(digits)
}
