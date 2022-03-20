package rpc

import (
	"log"
	"net"
	"net/rpc"
	"raft-node/domain/service"
)

type Server struct {
	RaftService *service.RaftService
}

func (server *Server) Start(port string) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatal("error while starting rpc server", err)
	}

	go func() {
		handler := &Handler{RaftService: server.RaftService}
		rpcRegisterError := rpc.Register(handler)
		if rpcRegisterError != nil {
			log.Fatal(rpcRegisterError)
		}

		for {
			rpc.Accept(listener)
		}
	}()
	log.Printf("RPC server started on port: %s", port)
}
