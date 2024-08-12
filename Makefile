srv:
	go run ./cmd/server -a :8081
agent:
	go run ./cmd/agent
test:
	go test -v -coverpkg=./... -cover  ./... 
