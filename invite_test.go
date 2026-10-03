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
