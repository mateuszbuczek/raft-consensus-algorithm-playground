package rpc

import (
	"raft-node/domain/model"
	"raft-node/domain/service"
)

type Handler struct {
	RaftService *service.RaftService
}

func (h *Handler) AppendEntry(args model.AppendEntryArgs, reply *model.AppendEntryReply) error {
	err := h.RaftService.AppendEntry(args, reply)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (h *Handler) RequestVote(args model.RequestVoteArgs, reply *model.RequestVoteReply) error {
	err := h.RaftService.RequestVote(args, reply)
	if err != nil {
		return err
	} else {
		return nil
	}
}

func (h *Handler) AppendEntryExternal(args AppendEntryExternalArgs, reply *AppendEntryExternalReply) error {
	err := h.RaftService.AppendEntryExternal(args.Value)
	if err != nil {
		reply.Success = false
		reply.ErrorMessage = err.Error()
	} else {
		reply.Success = true
	}
	return nil
}

func (h *Handler) GetEntriesExternal(args GetEntriesExternalArgs, reply *GetEntriesExternalReply) error {
	response, err := h.RaftService.GetEntriesExternal()
	if err != nil {
		reply.Success = false
		reply.ErrorMessage = err.Error()
	} else {
		reply.Success = true
		reply.Entries = response
	}
	return nil
}

type AppendEntryExternalArgs struct {
	Value interface{}
}

type AppendEntryExternalReply struct {
	Success      bool
	ErrorMessage string
}

type GetEntriesExternalArgs struct {
}

type GetEntriesExternalReply struct {
	Success      bool
	ErrorMessage string
	Entries      []model.LogEntry
}
