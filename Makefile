.PHONY: test coverage lint clean setup help
SHELL := '/bin/bash'
.DEFAULT_GOAL := help

test: clean ## run all tests
	@GOOS=js GOARCH=wasm go test -v -covermode=atomic -coverprofile=coverage.out \
		-v -exec="$(shell go env GOROOT)/misc/wasm/go_js_wasm_exec" github.com/chumaumenze/gjs/...

coverage: ## run coverage tool
	@go tool cover -func=./coverage.out
	@go tool cover -html=./coverage.out -o coverage.html

lint: ## lint go files in current directory
	@GOOS=js GOARCH=wasm golangci-lint run ./...

clean: ## remove build artefacts
	@go clean
	@rm -rf ./coverage.out ./coverage.html

setup: ## install dev tools
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50

# got from :https://marmelab.com/blog/2016/02/29/auto-documented-makefile.html
# but disallow underscore in command names as we want some private to have format "_command-name"
help:  ## print command reference
	@printf "  Welcome to \033[36mGJS\033[0m command reference.\n"
	@printf "  If you wish to contribute, please follow guide at top section of \033[36mMakefile\033[0m.\n\n"
	@printf "  Usage:\n    \033[36mmake <target> [..arguments]\033[0m\n\n  Targets:\n"
	@grep -E '^[a-zA-Z-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "    \033[36m%-20s\033[0m %s\n", $$1, $$2}'
