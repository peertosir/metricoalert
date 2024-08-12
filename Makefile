srv:
	go run ./cmd/server -a localhost:45625
agent:
	go run ./cmd/agent -a localhost:45625
test:
	go test -v -coverpkg=./... -cover  ./... 
