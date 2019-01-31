.DEFAULT_GOAL := help

.PHONY: all
all:

.PHONY: build
build: ## Build for local environment
	@go build

.PHONY: run
run: ## Run example script
	@bash ./scripts/test.sh

.PHONY: release
release: ## Build for multiple OSs, packaging it and upload to GitHub Release
	@#go get -u github.com/motemen/gobump/cmd/gobump
	@#go get -u github.com/git-chglog/git-chglog/cmd/git-chglog
	@#go get -u github.com/goreleaser/goreleaser
	@bash ./scripts/release.sh

.PHONY: help
help: ## Self-documented Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
