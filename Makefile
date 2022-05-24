.PHONY: dev dev-down go proto

dev:
	docker-compose up -d

dev-down:
	docker-compose down

go:
	air

proto:
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
  --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
  proto/*.proto