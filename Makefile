run-cluster:
	go run cluster/node-registry/main.go &

	cd cluster/node && go run main.go --nodeId 1 --nodeHost 127.0.0.1:8001 --nodeRegistryHost 127.0.0.1:8090 &
	cd cluster/node && go run main.go --nodeId 2 --nodeHost 127.0.0.1:8002 --nodeRegistryHost 127.0.0.1:8090 &
	cd cluster/node && go run main.go --nodeId 3 --nodeHost 127.0.0.1:8003 --nodeRegistryHost 127.0.0.1:8090

run-consumer:
	go run client/consumer/main.go --leaderAddress 127.0.0.1:8003

run-producer:
	go run client/producer/main.go --leaderAddress 127.0.0.1:8002

