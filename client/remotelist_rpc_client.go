package main

import (
	"fmt"
	remotelist "ifpb/remotelist/pkg"
	"net/rpc"
)

func main() {
	client, err := rpc.Dial("tcp", ":5000")
	if err != nil {
		fmt.Print("dialing:", err)
	}

	// Synchronous call
	var reply bool
	var reply_i int
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 0, Value: 10}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 0, Value: 20}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 0, Value: 30}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 0, Value: 40}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 0, Value: 50}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 0, Value: 60}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 0, Value: 70}, &reply)

	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 1, Value: 110}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 1, Value: 120}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 1, Value: 130}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 2, Value: 240}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 2, Value: 250}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 2, Value: 260}, &reply)
	err = client.Call("RemoteList.Append", remotelist.AppendArgs{ListId: 3, Value: 370}, &reply)

	err = client.Call("RemoteList.Remove", remotelist.RemoveArgs{ListId: 0}, &reply_i)
	if err != nil {
		fmt.Print("Error:", err)
	} else {
		fmt.Println("Elemento retirado:", reply_i)
	}
	err = client.Call("RemoteList.Remove", remotelist.RemoveArgs{ListId: 0}, &reply_i)
	if err != nil {
		fmt.Print("Error:", err)
	} else {
		fmt.Println("Elemento retirado:", reply_i)
	}
}
