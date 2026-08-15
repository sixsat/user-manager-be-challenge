.PHONY: lint
lint:
	@go run github.com/golangci/golangci-lint/cmd/golangci-lint@latest run --enable gosec

.PHONY: start
start:
	@docker compose up --build

.PHONY: stop
stop:
	@docker compose down -v

.PHONY: test
test:
	@go test -v -cover ./...