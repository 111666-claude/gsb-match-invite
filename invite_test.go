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

func TestCreateDuplicateRejected(t *testing.T) {
	world := New(1000, 4)
	if !world.Create("c1", "p1", "room-a", 0) {
		t.Fatal("首次创建应该成功")
	}
	if world.Create("c1", "p2", "room-b", 10) {
		t.Fatal("重复创建同一个邀请码应该返回 false")
	}
	if len(world.Pending(10)) != 1 {
		t.Fatal("重复创建不该改动待处理邀请")
	}
	if !world.Accept("c1", "p9", 20, "r1") {
		t.Fatal("原邀请应该仍然有效")
	}
	if world.Members("room-a") != 1 || world.Members("room-b") != 0 {
		t.Fatal("重复创建不该改动原邀请的房间")
	}
}

func TestInviteExpires(t *testing.T) {
	world := New(1000, 4)
	world.Create("c1", "p1", "room", 0)
	if len(world.Pending(999)) != 1 {
		t.Fatal("有效期内应该能看到邀请")
	}
	if len(world.Pending(1000)) != 0 {
		t.Fatal("达到 ttlMs 后邀请应该失效")
	}
	if world.Accept("c1", "p2", 1000, "r1") {
		t.Fatal("过期邀请不该能接受")
	}
}

func TestAcceptBeforeExpiry(t *testing.T) {
	world := New(1000, 4)
	world.Create("c1", "p1", "room", 0)
	if !world.Accept("c1", "p2", 999, "r1") {
		t.Fatal("有效期内应该能接受邀请")
	}
}

func TestRoomCapacity(t *testing.T) {
	world := New(100000, 2)
	world.Create("c1", "p1", "room", 0)
	if !world.Accept("c1", "p2", 0, "r1") || !world.Accept("c1", "p3", 0, "r2") {
		t.Fatal("人数未满时应该能接受")
	}
	if world.Accept("c1", "p4", 0, "r3") {
		t.Fatal("房间满员后应该拒绝接受")
	}
	if world.Members("room") != 2 {
		t.Fatal("房间人数不该超过上限")
	}
}

func TestCapacitySharedAcrossInvites(t *testing.T) {
	world := New(100000, 1)
	world.Create("c1", "p1", "room", 0)
	world.Create("c2", "p1", "room", 0)
	if !world.Accept("c1", "p2", 0, "r1") {
		t.Fatal("第一个接受应该成功")
	}
	if world.Accept("c2", "p3", 0, "r2") {
		t.Fatal("房间满员后其他邀请也该被拒绝")
	}
}

func TestSamePlayerOnce(t *testing.T) {
	world := New(100000, 4)
	world.Create("c1", "p1", "room", 0)
	if !world.Accept("c1", "p2", 0, "r1") {
		t.Fatal("首次接受应该成功")
	}
	if world.Accept("c1", "p2", 0, "r2") {
		t.Fatal("同一个玩家不该重复接受同一个邀请")
	}
	if world.Members("room") != 1 {
		t.Fatal("重复接受不该增加人数")
	}
}

func TestSameReqIDOnce(t *testing.T) {
	world := New(100000, 4)
	world.Create("c1", "p1", "room", 0)
	if !world.Accept("c1", "p2", 0, "same") {
		t.Fatal("首次接受应该成功")
	}
	if world.Accept("c1", "p3", 0, "same") {
		t.Fatal("同一个 reqID 不该生效两次")
	}
	if world.Members("room") != 1 {
		t.Fatal("重放请求不该增加人数")
	}
}

func TestFailedAcceptDoesNotConsumeReqID(t *testing.T) {
	world := New(100000, 4)
	if world.Accept("ghost", "p2", 0, "r1") {
		t.Fatal("不存在的邀请码不该能接受")
	}
	world.Create("c1", "p1", "room", 0)
	if !world.Accept("c1", "p2", 0, "r1") {
		t.Fatal("失败的接受不该消耗 reqID")
	}
}

func TestRecreateAfterExpiry(t *testing.T) {
	world := New(100, 4)
	world.Create("c1", "p1", "room", 0)
	if !world.Create("c1", "p2", "room", 200) {
		t.Fatal("原邀请过期后应该能重新创建")
	}
	if !world.Accept("c1", "p3", 250, "r1") {
		t.Fatal("重新创建的邀请应该能接受")
	}
}

func TestScannedStaysConstant(t *testing.T) {
	world := New(100000, 8)
	for index := 0; index < 100; index++ {
		world.Create(string(rune(index)), "p1", "room", 0)
	}
	for index := 0; index < 100; index++ {
		world.Pending(0)
	}
	if world.Scanned() > 200 {
		t.Fatalf("Scanned 不该随邀请数线性增长: %d", world.Scanned())
	}
}
