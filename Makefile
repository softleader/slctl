HAS_DOCKER := $(shell command -v docker;)
HAS_GOLINT := $(shell command -v golint;)
HAS_GOIMPORTS := $(shell command -v goimports;)
VERSION :=
COMMIT :=
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION} -X main.commit=${COMMIT}"
BINARY := slctl
MAIN := ./cmd/slctl
CHOCO_DIST := $(DIST)/choco
CHOCO_SERVER := http://softleader.com.tw:48081/repository/choco/
CHOCO_USER := choco:choco

.PHONY: help
help:	## # Show this help.
	@fgrep -h "##" $(MAKEFILE_LIST) | fgrep -v fgrep | sed -e 's/\\$$//' | sed -e 's/##//'

.PHONY: sandbox
sandbox: bootstrap test	## # Build and run slctl in Docker container
ifndef HAS_DOCKER
	$(error You must install Docker)
endif
	GOOS=linux go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) $(MAIN)
	docker build -t slctl .
	docker run --rm -it slctl bash

.PHONY: link
link: bootstrap test build	## # Build and run slctl and link to /usr/local/bin
	ln -sf $(BUILD)/$(BINARY) /usr/local/bin

.PHONY: test
test: error-free	## # Run error-free and test
	go test ./... -v

.PHONY: error-free
error-free: goimports gofmt golint govet	## # Run code check

.PHONY: goimports
goimports:	## # Run goimports
ifndef HAS_GOIMPORTS
	go get golang.org/x/tools/cmd/goimports
endif
	goimports -w -e .

.PHONY: gofmt
gofmt:	## # Run gofmt
	gofmt -s -e -w .

.PHONY: golint
golint:	## # Run golint
ifndef HAS_GOLINT
	go get -u golang.org/x/lint/golint
endif
	golint -set_exit_status ./cmd/...
	golint -set_exit_status ./pkg/...

.PHONY: govet
govet:	## # Run go vet
	go vet ./cmd/...
	go vet ./pkg/...

.PHONY: build
build:	## # Build binary for current OS and arch
	go build -o $(BUILD)/$(BINARY) -ldflags $(LDFLAGS) $(MAIN)

# build static binaries: https://medium.com/@diogok/on-golang-static-binaries-cross-compiling-and-plugins-1aed33499671
.PHONY: dist
dist:	## # Build and compress to tgz for linux amd64, darwin amd64 and windows amd64
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
bootstrap:	## # Make sure all dependency are downloaded
	go mod download

.PHONY: clean
clean:	## # Clean build temp dir and link to /usr/local/bin
	rm -rf _*
	rm -f /usr/local/bin/$(BINARY)

.PHONY: choco-pack
choco-pack:	## # Build to Choco package
ifndef HAS_DOCKER
	$(error You must install Docker)
endif
	mkdir -p $(CHOCO_DIST)
	cp $(BUILD)/$(BINARY).exe $(CHOCO_DIST)
	cp README.md $(CHOCO_DIST)
	# nuspec 不支援沒副檔名的檔案
	cp LICENSE $(CHOCO_DIST)/LICENSE.txt
	cp .nuspec $(CHOCO_DIST)
	docker run -v $(CHOCO_DIST):$(CHOCO_DIST) -w $(CHOCO_DIST) -it patrickhuber/choco-linux choco pack --version $(VERSION) --out $(CHOCO_DIST) $(CHOCO_DIST)/.nuspec

.PHONY: choco-push
choco-push: choco-pack	## # Push to SoftLeader Choco server
	curl -X PUT -F "file=@$(CHOCO_DIST)/$(BINARY).$(VERSION).nupkg" $(CHOCO_SERVER) -u $(CHOCO_USER)
