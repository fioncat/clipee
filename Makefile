.PHONY: install
install:
	CGO_ENABLED=1 go install

.PHONY: build
build:
	CGO_ENABLED=1 go build

.PHONY: fmt
fmt:
	@command -v goimports >/dev/null || { echo "ERROR: goimports not installed"; exit 1; }
	@exit $(shell find ./* \
	  -type f \
	  -name '*.go' \
	  -print0 | sort -z | xargs -0 -- goimports $(or $(FORMAT_FLAGS),-w) | wc -l | bc)

