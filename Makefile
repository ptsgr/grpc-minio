.PHONY: build-client run-server run-minio

build-client:
	cd cmd/client/ && go build -o ../../
run-server:
	go run cmd/server/server.go
run-minio:
	docker-compose up -d