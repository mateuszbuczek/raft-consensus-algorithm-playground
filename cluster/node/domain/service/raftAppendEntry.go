package service

import (
	"log"
	"raft-node/domain/model"
	"raft-node/infra/cluster"
	"time"
)

func (rs *RaftService) AppendEntry(args model.AppendEntryArgs, reply *model.AppendEntryReply) error {
	if rs.isLastLogTermAhead(args.Term) {
		reply.Success = false
		reply.Term = args.Term
		reply.ConflictTerm = rs.getLastLogTerm()
		return nil
	}

	rs.heartbeatChan <- true

	if !args.HasEntries() {
		reply.Success = true
		reply.Term = rs.getLastLogTerm()
		return nil
	}

	if rs.isLastLogIndexAhead(args.PreviousLogIndex) {
		reply.Success = false
		reply.Term = rs.getLastLogTerm()
		reply.ConflictIndex = rs.getLastLogIndex()
		return nil
	}

	rs.log = append(rs.log, args.Entries...)
	rs.commitIndex = rs.getLastLogIndex()
	reply.Success = true
	reply.Term = rs.getLastLogTerm()
	return nil
}

func (rs *RaftService) broadcastAndProcessAppendEntry() {
	nodes := rs.fetchCurrentNodes()
	appendEntryRepliesSuccessChan := make(chan *model.AppendEntryReply)
	appendEntryRepliesSuccess := 0

	for _, node := range nodes {
		if node.Id == rs.NodeId {
			continue
		}

		followerIndex := rs.getFollowerIndex(node)
		args := rs.createAppendEntryArgs(followerIndex)

		go rs.callAndProcessSingleAppendEntry(node, args, appendEntryRepliesSuccessChan)
	}

	for i := 0; i < len(nodes)-1; i++ {
		select {
		case <-appendEntryRepliesSuccessChan:
			appendEntryRepliesSuccess++
		case <-time.After(300 * time.Millisecond):
		}
	}

	if appendEntryRepliesSuccess > len(nodes)/2+1 {
		rs.commitIndex = rs.getLastLogIndex()
	}
}

func (rs *RaftService) callAndProcessSingleAppendEntry(node *cluster.Node, args model.AppendEntryArgs, successReplyChan chan *model.AppendEntryReply) {
	reply, nodeError := cluster.CallAppendEntry(node.Host, args)
	if nodeError != nil {
		log.Printf("was not able to call AppendEntry in node %d, error %v", node.Id, nodeError)
		return
	}
	if reply.Success {
		rs.followersIndex.Store(node.Id, args.LeaderCommit)
		successReplyChan <- reply
	} else if reply.ConflictTerm != 0 {
		rs.becomeFollower()
	} else if reply.ConflictIndex != 0 {
		rs.setFollowerIndex(node, reply.ConflictIndex)
	}
}

func (rs *RaftService) createAppendEntryArgs(followerIndex int) model.AppendEntryArgs {
	var args model.AppendEntryArgs
	args.Term = rs.getLastLogTerm()
	args.LeaderId = rs.NodeId
	args.LeaderCommit = rs.commitIndex

	if rs.isLastLogIndexAhead(followerIndex) {
		args.PreviousLogIndex = followerIndex
		args.PreviousLogTerm = rs.log[followerIndex].Term
		args.Entries = rs.log[followerIndex:]
	}
	return args
}
