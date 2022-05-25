.PHONY: dev dev-down proto server-go client-go

dev:
	docker-compose up -d

dev-down:
	docker-compose down

server-go:
	air cmd/server/main.go

client-go:
	go run cmd/client/main.go

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
  proto/*.proto