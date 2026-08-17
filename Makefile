.PHONY: lint
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --enable gosec

.PHONY: mock
mock:
	@go run go.uber.org/mock/mockgen@latest -source port/service.go -destination mock/service_mock.go -package mock
	@go run go.uber.org/mock/mockgen@latest -source port/repository.go -destination mock/repository_mock.go -package mock

.PHONY: start
start:
	@docker compose up --build

.PHONY: stop
stop:
	@docker compose down -v

.PHONY: test
test:
	@go test -v -cover ./...