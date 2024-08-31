srv:
	go run ./cmd/server -a localhost:8081
agent:
	go run ./cmd/agent -a localhost:8081
test:
	go test -v -coverpkg=./... -cover  ./... 
