runsv:
	go run ./cmd/server
runagent:
	go run ./cmd/agent
runtest:
	go test -v -coverpkg=./... -cover  ./... 
