package service

import (
	"log"
	"raft-node/domain/model"
	"raft-node/infra/cluster"
)

func (rs *RaftService) RequestVote(args model.RequestVoteArgs, reply *model.RequestVoteReply) error {
	if rs.isLastLogTermAhead(args.LastLogTerm) {
		reply.VoteGranted = false
		reply.Term = rs.getLastLogTerm()
		return nil
	}

	if rs.isLastLogIndexAhead(args.LastLogIndex) {
		reply.VoteGranted = false
		reply.Term = rs.getLastLogTerm()
		return nil
	}

	if rs.canVote {
		log.Printf("voted for %d\n", args.CandidateId)
		rs.canVote = false
		reply.VoteGranted = true
		reply.Term = rs.getLastLogTerm()
	}
	return nil
}

func (rs *RaftService) broadcastRequestVote() {
	nodes := rs.fetchCurrentNodes()

	if len(nodes) == 1 {
		rs.becomeLeaderChan <- true
		return
	}

	for _, node := range nodes {
		if node.Id == rs.NodeId {
			continue
		}

		var args = model.RequestVoteArgs{
			Term:         rs.getLastLogTerm() + 1,
			CandidateId:  rs.NodeId,
			LastLogIndex: rs.getLastLogIndex(),
			LastLogTerm:  rs.getLastLogTerm(),
		}

		go rs.callSingleRequestVote(nodes, node, args)
	}
}

func (rs *RaftService) callSingleRequestVote(nodes []*cluster.Node, node *cluster.Node, args model.RequestVoteArgs) {
	reply, nodeError := cluster.CallRequestVote(node.Host, args)
	if nodeError != nil {
		log.Printf("was not able to call RequestVote in node %d, error %v", node.Id, nodeError)
		return
	}
	if reply.VoteGranted {
		rs.votesCount++
		if rs.votesCount >= len(nodes)/2+1 {
			rs.becomeLeaderChan <- true
		}
	} else {
		rs.becomeFollower()
	}
}
