package service

import (
	"errors"
	"raft-node/domain/model"
)

func (rs *RaftService) GetEntriesExternal() ([]model.LogEntry, error) {
	if rs.isLeader() {
		return rs.log[:rs.commitIndex], nil
	} else {
		return nil, errors.New("requested node is not a current leader")
	}
}

func (rs *RaftService) AppendEntryExternal(value interface{}) error {
	if rs.isLeader() {
		rs.log = append(rs.log, model.LogEntry{
			Term:  rs.getLastLogTerm(),
			Index: rs.getLastLogIndex() + 1,
			Value: value,
		})
		return nil
	} else {
		return errors.New("requested node is not a current leader")
	}
}
