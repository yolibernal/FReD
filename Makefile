PROJECT_NAME := "fred"
PKG := "gitlab.tu-berlin.de/mcc-fred/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v /ext/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v /ext/ | grep -v _test.go)

.PHONY: all dep build clean test coverage coverhtml lint staticcheck container docs

all: build

lint: ## Lint the files
	@GO111MODULE=off go get -u golang.org/x/lint/golint
	@golint -set_exit_status ${PKG_LIST}

test: ## Run unittests
	@rm -rf pkg/leveldb/test.db
	@rm -rf pkg/badgerdb/test.db
	@go test -short ${PKG_LIST}
	@rm -rf pkg/leveldb/test.db
	@rm -rf pkg/badgerdb/test.db

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

msan: dep ## Run memory sanitizer
	@go test -msan -short ${PKG_LIST}

coverage: ## Generate global code coverage report
	@sh ./ci/tools/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	@sh ./ci/tools/coverage.sh html;

dep: ## Get the dependencies
	@go get -d ./...

build: dep ## Build the binary file
	@go build -i -v $(PKG)/cmd/frednode

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

container: ## Create a Docker container
	@docker build . -t gitlab-registry.tubit.tu-berlin.de/mcc-fred/fred/fred

staticcheck: ## Do a static code check
	@GO111MODULE=off go get -u honnef.co/go/tools/cmd/staticcheck
	@staticcheck ${PKG_LIST}

docs: ## Build the FogStore documentation
	@mdpdf docs/doc.md 

help: ## Display this help screen
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
