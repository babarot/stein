.PHONY: changelog
changelog:
	@go get -u github.com/git-chglog/git-chglog/cmd/git-chglog
	git-chglog --output CHANGELOG.md

.PHONY: release
release:
	goreleaser --release-notes CHANGELOG.md --rm-dist
