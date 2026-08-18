.PHONY: lint
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --enable gosec

.PHONY: mock
mock:
	@go run go.uber.org/mock/mockgen@latest -source port/service.go -destination mock/service_mock.go -package mock
	@go run go.uber.org/mock/mockgen@latest -source port/repository.go -destination mock/repository_mock.go -package mock

.PHONY: proto
proto:
	@protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/user.proto

.PHONY: start
start:
	@docker compose up --build

.PHONY: stop
stop:
	@docker compose down -v

.PHONY: test
test:
	@go test -v -cover ./...