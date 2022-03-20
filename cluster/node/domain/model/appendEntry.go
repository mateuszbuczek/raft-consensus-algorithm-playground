package model

type AppendEntryArgs struct {
	Term             int
	LeaderId         int
	Entries          []LogEntry
	PreviousLogIndex int
	PreviousLogTerm  int
	LeaderCommit     int
}

func (args AppendEntryArgs) HasEntries() bool {
	if len(args.Entries) != 0 {
		return true
	} else {
		return false
	}
}

type AppendEntryReply struct {
	Term          int
	Success       bool
	ConflictIndex int
	ConflictTerm  int
}
