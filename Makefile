HOST_GOOS=$(shell go env GOOS)
HOST_GOARCH=$(shell go env GOARCH)

all: client server

FORCE:

client: FORCE
	$(BUILD_ENV_FLAGS) go build -o bin/client ./client/*.go

server: FORCE
	$(BUILD_ENV_FLAGS) go build -o bin/server ./cmd/server/*.go

help:
	@echo "Influential make variables"
	@echo "  BUILD_ENV_FLAGS   - Environment added to 'go build'."
