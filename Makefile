proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

server:
	go run main.go

client:
	go run client/client.go

.PHONY: proto server client