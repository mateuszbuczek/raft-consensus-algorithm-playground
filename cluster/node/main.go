package main

import (
	"flag"
	"fmt"
	"log"
	"raft-node/api/rpc"
	"raft-node/domain/service"
	"raft-node/infra/cluster"
	"strings"
)

func main() {
	nodeId := flag.Int("nodeId", 1, "node Id")
	nodeHostString := flag.String("nodeHost", "127.0.0.1:9000", "rpc listen port")
	nodeRegistryHostString := flag.String("nodeRegistryHost", "127.0.0.1:8090", "node registry url")
	flag.Parse()

	nodeRegistryClient := cluster.NodeRegistryClient{Url: *nodeRegistryHostString}

	logger := log.Default()
	logger.SetPrefix(fmt.Sprintf("[NODE_ID %d] [NODE_HOST %s] ", *nodeId, *nodeHostString))
	logger.SetFlags(log.Ldate | log.Lmicroseconds)

	raftService := service.RaftService{
		NodeId:             *nodeId,
		NodeRegistryClient: &nodeRegistryClient,
	}

	rpcRaftController := rpc.Server{RaftService: &raftService}
	rpcRaftController.Start(strings.Split(*nodeHostString, ":")[1])

	raftService.RunNodeAsync()

	err := nodeRegistryClient.Register(*nodeId, *nodeHostString)
	if err != nil {
		errorMessage, _ := fmt.Printf("was not able to register node: %+v", err.Error())
		panic(errorMessage)
	}

	select {}
}
