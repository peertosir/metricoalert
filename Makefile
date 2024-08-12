srv:
	go run ./cmd/server
agent:
	go run ./cmd/agent
test:
	go test -v -coverpkg=./... -cover  ./... 
