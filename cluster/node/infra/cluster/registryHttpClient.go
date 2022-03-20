package cluster

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
)

type NodeRegistryClient struct {
	Url string
}

func (client *NodeRegistryClient) Register(id int, host string) error {
	nodeJson, err := json.Marshal(Node{Id: id, Host: host})
	if err != nil {
		return err
	}

	response, err := http.Post("http://"+client.Url+"/nodes", "application/json", bytes.NewBuffer(nodeJson))
	if err != nil {
		return err
	}

	if response.StatusCode >= 200 && response.StatusCode < 300 {
		return nil
	} else {
		return errors.New("wasnt able to register server. status code: " + response.Status)
	}
}

func (client *NodeRegistryClient) GetNodes() ([]*Node, error) {
	var nodes []*Node
	response, err := http.Get("http://" + client.Url + "/nodes")
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	err = json.NewDecoder(response.Body).Decode(&nodes)
	if err != nil {
		return nil, err
	}

	return nodes, nil
}

type Node struct {
	Id   int    `json:"id"`
	Host string `json:"host"`
}
