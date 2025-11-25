package main

import (
	"fmt"
	remotelist "ifpb/remotelist/pkg"
	"net"
	"net/rpc"
)

func main() {
	list := remotelist.NewRemoteList()
	remotelist.StartSnapshotRoutine(list)
	rpcs := rpc.NewServer()
	rpcs.Register(list)
	l, e := net.Listen("tcp", "[localhost]:5000")
	if e != nil {
		fmt.Println("listen error:", e)
		return
	}
	defer l.Close()
	for {
		conn, err := l.Accept()
		if err == nil {
			go rpcs.ServeConn(conn)
		} else {
			break
		}
	}
}
