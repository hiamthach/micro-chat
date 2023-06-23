proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
	--go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	proto/*.proto

protojs:
	protoc --proto_path=proto --js_out=import_style=commonjs:client/ \
	--grpc-web_out=import_style=commonjs,mode=grpcwebtext:client/ \
	proto/*.proto

server:
	go run main.go

client:
	go run client/client.go

.PHONY: proto server client