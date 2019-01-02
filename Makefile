HAS_DOCKER := $(shell command -v docker;)
VERSION := ""
COMMIT := ""
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION} -X main.commit=${COMMIT}"
BINARY := slctl
MAIN := ./cmd/slctl
CHOCO_DIST := $(DIST)/choco
CHOCO_SERVER := http://ci.softleader.com.tw:8081/repository/choco/
CHOCO_USER := choco:choco

.PHONY: sandbox
sandbox: bootstrap test
ifndef HAS_DOCKER
	$(error You must install Docker)
endif
	GOOS=linux go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) $(MAIN)
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
	go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) $(MAIN)

# build static binaries: https://medium.com/@diogok/on-golang-static-binaries-cross-compiling-and-plugins-1aed33499671
.PHONY: dist
dist:
ifeq ($(strip $(VERSION)),)
	$(error VERSION is not set)
endif
ifeq ($(strip $(COMMIT)),)
	$(error COMMIT is not set)
endif
	mkdir -p $(BUILD)
	mkdir -p $(DIST)
	cp README.md $(BUILD) && cp LICENSE $(BUILD)
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-linux-$(VERSION).tgz $(BINARY) README.md LICENSE
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -zcvf $(DIST)/$(BINARY)-darwin-$(VERSION).tgz $(BINARY) README.md LICENSE
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BUILD)/$(BINARY).exe -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	tar -C $(BUILD) -llzcvf $(DIST)/$(BINARY)-windows-$(VERSION).tgz $(BINARY).exe README.md LICENSE

.PHONY: bootstrap
bootstrap:
	go mod download

.PHONY: clean
clean:
	rm -rf _*
	rm -f /usr/local/bin/$(BINARY)

.PHONY: choco-pack
choco-pack:
ifndef HAS_DOCKER
	$(error You must install Docker)
endif
	mkdir -p $(CHOCO_DIST)
	cp $(BUILD)/$(BINARY).exe $(CHOCO_DIST)
	cp README.md $(CHOCO_DIST)
	cp LICENSE $(CHOCO_DIST)
	cp .nuspec $(CHOCO_DIST)
	docker run -v $(DIST):$(DIST) -w $(DIST) -it patrickhuber/choco-linux choco pack --version $(VERSION) --out $(DIST) $(CHOCO_DIST)/.nuspec

.PHONY: choco-push
choco-push: choco-pack
	curl -X PUT -F "file=@$(DIST)/$(BINARY).$(VERSION).nupkg" $(CHOCO_SERVER) -u $(CHOCO_USER)
