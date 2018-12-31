HAS_DOCKER := $(shell command -v docker;)
BUILD := $(CURDIR)/_build
BINARY := slctl
MAIN := ./cmd/slctl

.PHONY: sandbox
sandbox: bootstrap test
ifndef HAS_DOCKER
	$(error You must install Docker)
endif
	GOOS=linux go build -o $(BUILD)/$(BINARY) $(MAIN)
	docker build -t slctl .
	docker run --rm -it slctl bash

.PHONY: link
link: bootstrap test build
	ln -sf $(BUILD)/$(BINARY) /usr/local/bin

.PHONY: test
test:
	go test ./... -v

.PHONY: build
build:
	go build -o $(BUILD)/$(BINARY) $(MAIN)

.PHONY: dist
dist: bootstrap
	goreleaser release --snapshot --rm-dist

.PHONY: tag
tag:
ifeq ($(strip $(VERSION)),)
	$(error VERSION is not set)
endif
	git tag -a $(VERSION) -m "$(VERSION)"
	git push origin $(VERSION)

.PHONY: bootstrap
bootstrap:
ifeq (,$(wildcard ./go.mod))
	go mod init {{.Name}}
endif
	go mod download

.PHONY: clean
clean:
	rm -rf _*
	rm -f /usr/local/bin/$(BINARY)