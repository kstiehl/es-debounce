
.PHONY: all build test unit integration
all: generate build test unit integration


generate:
	@echo "#### compile proto files... ####"
	./scripts/generate.sh

build:
	@echo "####    building go code    ####"
	go build -v

test: unit integration

unit:
	@echo "####       unit tests       ####"
	go test -v ./pkg/...

integration:
	@echo "####   integration tests    ####"
	go test -v ./integration/...
