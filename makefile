gen:
	protoc --go-grpc_out=./pb --go_out=./pb ./proto/*.proto 

clean:
	rm -f pb/pb/*.go

server:
	go run cmd/server/main.go -port 8080

client:
	go run cmd/client/main.go -address 0.0.0.0:8080

test:
	go test -cover -race 	./...

.PHONY: gen clean server client test