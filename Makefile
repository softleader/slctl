HAS_DOCKER := $(shell command -v docker;)
VERSION := ""
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION}"
BINARY := slctl
MAIN := ./cmd/slctl
CHOCO_SERVER := http://ci.softleader.com.tw:8081/repository/choco/
CHOCO_USER := choco:choco

.PHONY: install
install: bootstrap test
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

.PHONY: release
release:
ifneq ($(strip $(VERSION)),)
	$(error VERSION is not set)
endif
	git tag -a $(VERSION) -m "$(VERSION)"
	git push origin $(VERSION)
	goreleaser

.PHONY: bootstrap
bootstrap:
	go mod tidy

.PHONY: clean
clean:
	rm -rf _*
	rm -f /usr/local/bin/$(BINARY)

.PHONY: choco-pack
choco-pack:
	mkdir -p $(BUILD)
	mkdir -p $(DIST)
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o $(BUILD)/$(BINARY).exe -ldflags $(LDFLAGS) -a -tags netgo $(MAIN)
	cp README.md $(BUILD) && cp LICENSE $(BUILD) && cp .nuspec $(BUILD)	
	choco pack --version $(VERSION) --outputdirectory $(DIST) $(BUILD)/.nuspec

.PHONY: choco-push
choco-push: choco-pack
	curl -X PUT -F "file=@$(DIST)/slctl.$(VERSION).nupkg" $(CHOCO_SERVER) -u $(CHOCO_USER) -v
