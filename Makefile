run:
	@go run cmd/main.go
run-dev:
	@redis-cli flushall && go run cmd/main.go