.PHONY: vendor
.EXPORT_ALL_VARIABLES:
GO111MODULE=on

vendor:
	@echo "\033[92mINFO: Updating dependencies\033[0m"
	@go mod tidy
	@go mod vendor

test:
	@echo "\033[92mINFO: Running tests\033[0m"
	@go test -mod=vendor -race -covermode=atomic .
