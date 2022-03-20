package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

var NODE_LIST = make([]*Node, 0, 50)
var LOCK_PROVIDER = sync.RWMutex{}

func main() {
	http.HandleFunc("/nodes", nodesHandler)

	http.ListenAndServe(":8090", nil)
}

func nodesHandler(w http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		jsonNodeList, err := json.Marshal(NODE_LIST)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		_, err = w.Write(jsonNodeList)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else if req.Method == http.MethodPost {
		var node Node
		err := json.NewDecoder(req.Body).Decode(&node)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}

		LOCK_PROVIDER.Lock()
		defer LOCK_PROVIDER.Unlock()

		NODE_LIST = append(NODE_LIST, &node)
		fmt.Printf("NODE_REGISTRY: new node registered %v\n", node)
		w.WriteHeader(http.StatusAccepted)
	} else {
		w.WriteHeader(http.StatusNotFound)
	}
}

type Node struct {
	Id   int32  `json:"id"`
	Host string `json:"host"`
}
