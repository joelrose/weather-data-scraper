generate: # Run generate command
	@go run main.go generate

migrate: # Run migrate command
	@go run main.go migrate

test: # Run tests
	@go test -coverpkg=./... -race -covermode=atomic ./...

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
