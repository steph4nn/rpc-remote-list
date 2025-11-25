package remotelist

import (
	"fmt"
)

func (l *RemoteList) Append(args AppendArgs, reply *bool) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	l.logCounter++
	AppendLog(l.logCounter, args)
	l.list[args.ListId] = append(l.list[args.ListId], args.Value)
	*reply = true
	SaveState(l)
	return nil
}

func (l *RemoteList) Remove(args RemoveArgs, reply *int) error {
	l.mu.Lock()
	defer l.mu.Unlock()

	lastIndex := len(l.list[args.ListId]) - 1
	*reply = l.list[args.ListId][lastIndex]
	removedValue := *reply
	l.logCounter++
	RemoveLog(l.logCounter, args, removedValue)
	l.list[args.ListId] = l.list[args.ListId][:lastIndex]
	if len(l.list[args.ListId]) == 0 {
		delete(l.list, args.ListId)
	}
	SaveState(l)
	return nil
}

func NewRemoteList() *RemoteList {
	list := getState()
	logCounter := getLastLogId()
	fmt.Printf("RemoteList inicializada. Último log ID: %d\n", logCounter)
	return &RemoteList{
		list:       list,
		logCounter: logCounter,
	}
}

func (l *RemoteList) Get(args GetArgs, reply *int) error {
	l.mu.RLock()
	defer l.mu.RUnlock()

	*reply = l.list[args.ListId][args.Index]
	return nil
}

func (l *RemoteList) Size(args SizeArgs, reply *int) error {
	l.mu.RLock()
	defer l.mu.RUnlock()
	list := l.list[args.ListId]
	*reply = len(list)
	return nil
}
