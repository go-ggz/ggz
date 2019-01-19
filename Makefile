DIST := dist
SERVICE ?= ggz-server

DOCKER_ACCOUNT := goggz
DOCKER_IMAGE := $(SERVICE)
GOFMT ?= gofmt "-s"
GO ?= go
TARGETS ?= linux darwin windows
ARCHS ?= amd64 386
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
PACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go" -type f)
SOURCES ?= $(shell find . -name "*.go" -type f)
TAGS ?= sqlite
LDFLAGS ?= -X github.com/go-ggz/ggz/version.Version=$(VERSION) -X github.com/go-ggz/ggz/version.BuildDate=$(BUILD_DATE)
TMPDIR := $(shell mktemp -d 2>/dev/null || mktemp -d -t 'tempdir')
STYLESHEETS := $(wildcard assets/dist/less/innhp.less  assets/dist/less/_*.less)

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

ifneq ($(DRONE_TAG),)
	VERSION ?= $(DRONE_TAG)
else
	VERSION ?= $(shell git describe --tags --always || git rev-parse --short HEAD)
endif

all: build

.PHONY: install-module
install-module:
	$(GO) mod download

.PHONY: tar
tar:
	tar -zcvf release.tar.gz bin Dockerfile Makefile

.PHONY: check_image
check_image:
	if [ "$(shell docker ps -aq -f name=$(SERVICE))" ]; then \
		docker rm -f $(SERVICE); \
	fi

.PHONY: dev
dev: build_image check_image
	docker run -d --name $(DOCKER_IMAGE) --env-file env/env.$@ --net host -p 3003:3003 --restart always $(DOCKER_ACCOUNT)/$(DOCKER_IMAGE)

.PHONY: prod
prod: build_image check_image
	docker run -d --name $(DOCKER_IMAGE) --env-file env/env.$@ --net host -p 3003:3003 --restart always $(DOCKER_ACCOUNT)/$(DOCKER_IMAGE)

.PHONY: generate
generate:
	@which fileb0x > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/UnnoTed/fileb0x; \
	fi
	$(GO) generate $(PACKAGES)

.PHONY: stylesheets-check
stylesheets-check: stylesheets
	@diff=$$(git diff assets/dist/css/innhp.css); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make stylesheets' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: stylesheets
stylesheets: assets/dist/css/innhp.css

.IGNORE:assets/dist/css/innhp.css
assets/dist/css/innhp.css: $(STYLESHEETS)
	@which lessc > /dev/null; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/kib357/less-go/lessc; \
	fi
	lessc -i $< -o $@

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

embedmd:
	@hash embedmd > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/campoy/embedmd; \
	fi
	embedmd -d *.md

vet:
	$(GO) vet $(PACKAGES)

	errcheck:
	@hash errcheck > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/kisielk/errcheck; \
	fi
	errcheck $(PACKAGES)

lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u golang.org/x/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

unconvert:
	@hash unconvert > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mdempsky/unconvert; \
	fi
	for PKG in $(PACKAGES); do unconvert -v $$PKG || exit 1; done;

install: $(SOURCES)
	$(GO) install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'

build: $(SERVICE)

$(SERVICE): $(SOURCES)
	$(GO) build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o bin/$@ ./cmd/$(SERVICE)

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(GOFILES)

unused-check:
	@hash unused > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u honnef.co/go/tools/cmd/unused; \
	fi
	for PKG in $(PACKAGES); do unused $$PKG || exit 1; done;

upx:
	@hash upx > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		echo "Missing upx command"; \
		exit 1; \
	fi
	upx -o bin/$(SERVICE)-small bin/$(SERVICE)
	mv bin/$(SERVICE)-small bin/$(SERVICE)

.PHONY: coverage
coverage:
	@hash gocovmerge > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/wadey/gocovmerge; \
	fi
	gocovmerge $(shell find . -type f -name "coverage.out") > coverage.all;\

.PHONY: unit-test-coverage
unit-test-coverage:
	for PKG in $(PACKAGES); do $(GO) test -tags=sqlite -cover -coverprofile $$GOPATH/src/$$PKG/coverage.out $$PKG || exit 1; done;

test:
	for PKG in $(PACKAGES); do $(GO) test -tags=sqlite -v $$PKG || exit 1; done;

release: release-dirs release-build release-copy release-check

release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

release-build:
	@hash gox > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		$(GO) get -u github.com/mitchellh/gox; \
	fi
	gox -os="$(TARGETS)" -arch="$(ARCHS)" -tags="$(TAGS)" -ldflags="$(EXTLDFLAGS)-s -w $(LDFLAGS)" -output="$(DIST)/binaries/$(SERVICE)-$(VERSION)-{{.OS}}-{{.Arch}}" ./cmd/$(SERVICE)/...

release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(SERVICE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/release/$(SERVICE)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)

build_linux_amd64:
	GOOS=linux GOARCH=amd64 $(GO) build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/amd64/$(DOCKER_IMAGE) ./cmd/$(SERVICE)

build_linux_i386:
	GOOS=linux GOARCH=386 $(GO) build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/i386/$(DOCKER_IMAGE) ./cmd/$(SERVICE)

build_linux_arm64:
	GOOS=linux GOARCH=arm64 $(GO) build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/arm64/$(DOCKER_IMAGE) ./cmd/$(SERVICE)

build_linux_arm:
	GOOS=linux GOARCH=arm GOARM=7 $(GO) build -a -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o release/linux/arm/$(DOCKER_IMAGE) ./cmd/$(SERVICE)

build_image:
	docker build -t $(DOCKER_ACCOUNT)/$(DOCKER_IMAGE) -f Dockerfile .

docker_release: build_image

clean_dist:
	rm -rf bin server release

clean: clean_dist
	$(GO) clean -modcache -cache -x -i ./...
