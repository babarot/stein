.DEFAULT_GOAL := build

.PHONY: build
build:
	hugo -d ../docs

server: build
	hugo server --disableFastRender
