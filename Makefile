.PHONY: install
install:
	CGO_ENABLED=1 go install

.PHONY: build
build:
	CGO_ENABLED=1 go build
