package main

import (
	"flag"
	"log"
	"net/rpc"
	"time"
)

func main() {
	leaderAddress := flag.String("leaderAddress", "127.0.0.1:8000", "leader url address")
	flag.Parse()

	go func() {
		for {
			var reply GetEntriesExternalReply
			client, err := rpc.Dial("tcp", *leaderAddress)
			if err != nil {
				log.Println(err)
			} else {
				err = client.Call("Handler.GetEntriesExternal", GetEntriesExternalArgs{}, &reply)
				if err != nil {
					log.Println(err)
				} else {
					entriesLength := len(reply.Entries)

					if entriesLength == 0 {
						log.Printf("GetEntriesExternalReply length: %d", entriesLength)
					} else {
						lastIndex := reply.Entries[entriesLength-1].Index
						lastTerm := reply.Entries[entriesLength-1].Term
						log.Printf("GetEntriesExternalReply length: %d, lastTerm: %d, lastIndex: %d \n",
							entriesLength, lastTerm, lastIndex)
					}
				}
				client.Close()
			}
			time.Sleep(5000 * time.Millisecond)
		}
	}()

	select {}
}

type GetEntriesExternalArgs struct {
}

type GetEntriesExternalReply struct {
	Success      bool
	ErrorMessage string
	Entries      []LogEntry
}

type LogEntry struct {
	Term  int
	Index int
	Value interface{}
}
