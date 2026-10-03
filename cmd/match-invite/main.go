// Command match-invite 跑对局邀请样例。
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"example.com/invite"
)

// Run 执行一次命令行调用，返回退出码。
func Run(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("match-invite", flag.ContinueOnError)
	flags.SetOutput(stderr)
	sample := flags.String("sample", "dup", "dup / ttl / cap / idem / scan")
	if err := flags.Parse(args); err != nil {
		return 2
	}
	switch *sample {
	case "dup":
		world := invite.New(100000, 4)
		world.Create("c1", "p1", "room", 0)
		fmt.Fprintf(stdout, "second=%v\n", world.Create("c1", "p1", "room", 10))
	case "ttl":
		world := invite.New(1000, 4)
		world.Create("c1", "p1", "room", 0)
		fmt.Fprintf(stdout, "accepted=%v\n", world.Accept("c1", "p2", 1001, "r1"))
	case "cap":
		world := invite.New(100000, 2)
		world.Create("c1", "p1", "room", 0)
		world.Accept("c1", "p2", 0, "r1")
		world.Accept("c1", "p3", 0, "r2")
		fmt.Fprintf(stdout, "third=%v members=%d\n", world.Accept("c1", "p4", 0, "r3"), world.Members("room"))
	case "idem":
		world := invite.New(100000, 8)
		world.Create("c1", "p1", "room", 0)
		world.Accept("c1", "p2", 0, "same")
		fmt.Fprintf(stdout, "second=%v members=%d\n", world.Accept("c1", "p2", 0, "same"), world.Members("room"))
	case "scan":
		world := invite.New(100000, 8)
		for index := 0; index < 3000; index++ {
			world.Create(fmt.Sprintf("c-%d", index), "p1", "room", 0)
		}
		for index := 0; index < 3000; index++ {
			world.Pending(0)
		}
		fmt.Fprintf(stdout, "scanned=%d\n", world.Scanned())
	default:
		fmt.Fprintln(stderr, "需要 --sample dup|ttl|cap|idem|scan")
		return 2
	}
	return 0
}

func main() {
	os.Exit(Run(os.Args[1:], os.Stdout, os.Stderr))
}
