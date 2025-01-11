export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on

MODULE = pluginstore

PORT 			:= 8080
VERSION			:= $(shell git describe --tags --always --match="v*" 2> /dev/null || echo v0.0.0)
VERSION_HASH	:= $(shell git rev-parse --short HEAD)

GOCGO 	:= env CGO_ENABLED=1
LDFLAGS	:= -s -w -X "$(MODULE)/config.Version=$(VERSION)" -X "$(MODULE)/config.CommitSHA=$(VERSION_HASH)"

run: build
	./main --mode debug

watch:
	@if lsof -i :$(PORT) >/dev/null 2>&1; then \
        echo "Port $(PORT) is already in use, killing the process..."; \
        lsof -i :$(PORT) | awk 'NR!=1 {print $$2}' | xargs kill; \
    fi
	$(GOCGO) air
# cd web && npm run build && cd ../
release:
	
	OCKER_BUILDKIT=1 docker buildx build --push -t hitosea2020/pluginstore:0.0.1 --platform linux/amd64,linux/arm64 .

# cd web && npm run build && cd ../
build: 
	
	OCKER_BUILDKIT=1 docker buildx build -t xxyijixx/pluginstore:0.0.1 --platform linux/amd64,linux/arm64 .

translate:
	cd web && npm run translate && cd ../

dev:
	cd web && npm run build && cd ../
	go run main.go -m

preview:
	DOCKER_BUILDKIT=1 docker buildx build -t xxyijixx/pluginstore:0.0.1 --platform linux/amd64 --no-cache --load .
	docker compose up


