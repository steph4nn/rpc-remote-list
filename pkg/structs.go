package remotelist

import (
	"sync"
)

type AppendArgs struct {
	ListId int
	Value  int
}

type RemoveArgs struct {
	ListId int
}

type GetArgs struct {
	ListId int
	Index  int
}

type SizeArgs struct {
	ListId int
}

type RemoteList struct {
	mu         sync.RWMutex 
	list       map[int][]int
	logCounter int64
}

type Snapshot struct {
	Timestamp      int64
	LastLogApplied int64
	List           map[int][]int
}

type snapFile struct {
	name      string
	timestamp int64
}

type LogEntry struct {
	Id        int64
	Operation string
	ListId    int
	Value     int
}
