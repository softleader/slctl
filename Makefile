HAS_GLIDE := $(shell command -v glide;)
HAS_DOCKER := $(shell command -v docker;)
VERSION := ""
DIST := $(CURDIR)/_dist
BUILD := $(CURDIR)/_build
LDFLAGS := "-X main.version=${VERSION}"
BINARY := slctl
MAIN := github.com/softleader/slctl/cmd
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

# build static binaries: https://medium.com/@diogok/on-golang-static-binaries-cross-compiling-and-plugins-1aed33499671
.PHONY: dist
dist:
ifndef VERSION
	$(error VERSION is not set)
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
ifndef HAS_GLIDE
	go get -u github.com/Masterminds/glide
endif
	glide install --strip-vendor

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