.DEFAULT_GOAL := help

.PHONY: all
all:

.PHONY: build
build: ## Build for local environment
	@go build

.PHONY: changelog
changelog: ## Generate changelog automatically
	@go get -u github.com/git-chglog/git-chglog/cmd/git-chglog
	git-chglog --output CHANGELOG.md
	git diff

.PHONY: release
release: ## Build for multiple OSs, packaging it and upload to GitHub Release
	@#go get -u github.com/goreleaser/goreleaser
	@#goreleaser --release-notes CHANGELOG.md --rm-dist
	@bash ./scripts/release.sh

.PHONY: help
help: ## Self-documented Makefile
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) \
		| sort \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
