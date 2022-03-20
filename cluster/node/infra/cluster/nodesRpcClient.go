package cluster

import (
	"net/rpc"
	"raft-node/domain/model"
)

func CallAppendEntry(url string, args model.AppendEntryArgs) (*model.AppendEntryReply, error) {
	var reply model.AppendEntryReply
	client, err := rpc.Dial("tcp", url)
	if err != nil {
		return nil, err
	}
	defer client.Close()
	err = client.Call("Handler.AppendEntry", args, &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}

func CallRequestVote(url string, args model.RequestVoteArgs) (*model.RequestVoteReply, error) {
	var reply model.RequestVoteReply
	client, err := rpc.Dial("tcp", url)
	if err != nil {
		return nil, err
	}
	defer client.Close()

	err = client.Call("Handler.RequestVote", args, &reply)
	if err != nil {
		return nil, err
	}

	return &reply, nil
}
