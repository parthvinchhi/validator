.PHONY: run test build docker vet
run:
	go run ./cmd/api
test:
	go test ./... -race -cover
vet:
	go vet ./...
build:
	CGO_ENABLED=0 go build -o bin/validator-api ./cmd/api
docker:
	docker build -t validator-api .
