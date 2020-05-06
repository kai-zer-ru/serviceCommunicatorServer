PROJECT_NAME := "serviceCommunicatorServer"
PKG := "github.com/kaizer666/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v color)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

.PHONY: all build clean test coverage coverhtml lint

all: build

lint: ## Lint the files
	go get -u github.com/mgechev/revive
	@revive -config ~/linterConf.toml

test: ## Run unittests
	@go test -short ${PKG_LIST}

race: ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	~/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	~/coverage.sh html;

build: ## Build the binary file
	@go install

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'