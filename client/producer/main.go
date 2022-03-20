package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"net/rpc"
	"time"
)

func main() {
	clientId := rand.Intn(100)
	leaderAddress := flag.String("leaderAddress", "127.0.0.1:8000", "leader url address")
	flag.Parse()

	go func() {
		for {
			var reply AppendEntryExternalReply
			client, err := rpc.Dial("tcp", *leaderAddress)
			if err != nil {
				log.Println(err)
			} else {
				err = client.Call("Handler.AppendEntryExternal", AppendEntryExternalArgs{
					Value: fmt.Sprintf("%d %s", clientId, time.Now().String())}, &reply)
				if err != nil {
					log.Println(err)
				} else {
					log.Printf("AppendEntryExternalReply: %+v\n", reply)
				}
				client.Close()
			}

			time.Sleep(5 * time.Millisecond)
		}
	}()

	select {}
}

type AppendEntryExternalArgs struct {
	Value interface{}
}

type AppendEntryExternalReply struct {
	Success      bool
	ErrorMessage string
}
