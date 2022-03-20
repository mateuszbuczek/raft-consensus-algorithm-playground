package service

import (
	"log"
	"math/rand"
	"raft-node/domain/model"
	"raft-node/infra/cluster"
	"sync"
	"time"
)

type state int

const (
	LEADER = iota
	CANDIDATE
	FOLLOWER
)

type RaftService struct {
	NodeId int
	log    []model.LogEntry
	state  state

	// LEADER
	followersIndex sync.Map
	commitIndex    int

	// CANDIDATE
	canVote          bool
	votesCount       int
	becomeLeaderChan chan bool

	// FOLLOWER
	heartbeatChan chan bool

	NodeRegistryClient *cluster.NodeRegistryClient
}

func (rs *RaftService) RunNodeAsync() {
	rs.state = FOLLOWER
	rs.heartbeatChan = make(chan bool)
	rs.becomeLeaderChan = make(chan bool)
	rs.followersIndex = sync.Map{}

	heartbeatInterval := 1000

	go rs.startNode(heartbeatInterval)
}

func (rs *RaftService) startNode(heartbeatInterval int) {
	for {
		switch rs.state {
		case FOLLOWER:
			select {
			case <-rs.heartbeatChan:
				log.Println("heartbeat received")
			case <-time.After(time.Duration(rand.Intn(heartbeatInterval*5)+heartbeatInterval*5) * time.Millisecond):
				log.Println("heartbeat have not been received")
				rs.state = CANDIDATE
			}
		case CANDIDATE:
			log.Println("becomes candidate")
			rs.votesCount = 1
			rs.canVote = false

			log.Printf("canVote %v", rs.canVote)
			go rs.broadcastRequestVote()

			select {
			case <-rs.becomeLeaderChan:
				log.Println("became leader")
				rs.state = LEADER
				rs.followersIndex = sync.Map{}
			case <-time.After(time.Duration(rand.Intn(heartbeatInterval*5)+heartbeatInterval*5) * time.Millisecond):
				log.Printf("did not become a leader, votes: %d", rs.votesCount)
				rs.becomeFollower()
			}
		case LEADER:
			rs.broadcastAndProcessAppendEntry()
			time.Sleep(time.Duration(heartbeatInterval) * time.Millisecond)
		}
	}
}

func (rs *RaftService) fetchCurrentNodes() []*cluster.Node {
	nodes, err := rs.NodeRegistryClient.GetNodes()
	if err != nil {
		log.Fatal("was not able to fetch nodes list:", err)
	}
	return nodes
}

func (rs *RaftService) becomeFollower() {
	rs.state = FOLLOWER
	rs.canVote = true
}

func (rs RaftService) isLastLogTermAhead(term int) bool {
	if rs.getLastLogTerm() > term {
		return true
	} else {
		return false
	}
}

func (rs RaftService) isLastLogIndexAhead(index int) bool {
	if rs.getLastLogIndex() > index {
		return true
	} else {
		return false
	}
}

func (rs *RaftService) isLeader() bool {
	return rs.state == LEADER
}

func (rs RaftService) getLastLogIndex() int {
	entriesLength := len(rs.log)

	if entriesLength == 0 {
		return 0
	} else {
		return rs.log[entriesLength-1].Index
	}
}

func (rs RaftService) getLastLogTerm() int {
	entriesLength := len(rs.log)

	if entriesLength == 0 {
		return 0
	} else {
		return rs.log[entriesLength-1].Term
	}
}

func (rs *RaftService) getFollowerIndex(node *cluster.Node) int {
	var followerIndex int
	foundValue, found := rs.followersIndex.Load(node.Id)
	if found {
		followerIndex = foundValue.(int)
	}
	return followerIndex
}

func (rs *RaftService) setFollowerIndex(node *cluster.Node, index int) {
	rs.followersIndex.Store(node.Id, index)
}
