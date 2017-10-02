DIST := dist
EXECUTABLE := ggz

DEPLOY_ACCOUNT := goggz
DEPLOY_IMAGE := $(EXECUTABLE)
GOFMT ?= gofmt "-s"

TARGETS ?= linux darwin windows
ARCHS ?= amd64 386
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
SOURCES ?= $(shell find . -name "*.go" -type f)
TAGS ?=
LDFLAGS ?= -X 'main.Version=$(VERSION)'
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

.PHONY: tar
tar:
	tar -zcvf release.tar.gz bin Dockerfile Makefile

.PHONY: check_image
check_image:
	if [ "$(shell docker ps -aq -f name=$(EXECUTABLE))" ]; then \
		docker rm -f $(EXECUTABLE); \
	fi

.PHONY: dev
dev: build_image check_image
	docker run -d --name $(DEPLOY_IMAGE) --env-file env/env.$@ --net host -p 3003:3003 --restart always $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE)

.PHONY: prod
prod: build_image check_image
	docker run -d --name $(DEPLOY_IMAGE) --env-file env/env.$@ --net host -p 3003:3003 --restart always $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE)

.PHONY: generate
generate:
	@which fileb0x > /dev/null; if [ $$? -ne 0 ]; then \
		go get -u github.com/UnnoTed/fileb0x; \
	fi
	go generate $(PACKAGES)

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
		go get -u github.com/kib357/less-go/lessc; \
	fi
	lessc -i $< -o $@

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	# get all go files and run go fmt on them
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

embedmd:
	@hash embedmd > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/campoy/embedmd; \
	fi
	embedmd -d *.md

vet:
	go vet $(PACKAGES)

	errcheck:
	@hash errcheck > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kisielk/errcheck; \
	fi
	errcheck $(PACKAGES)

lint:
	@hash golint > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/golang/lint/golint; \
	fi
	for PKG in $(PACKAGES); do golint -set_exit_status $$PKG || exit 1; done;

unconvert:
	@hash unconvert > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/mdempsky/unconvert; \
	fi
	for PKG in $(PACKAGES); do unconvert -v $$PKG || exit 1; done;

install: $(SOURCES)
	go install -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)'

build: $(EXECUTABLE)

$(EXECUTABLE): $(SOURCES)
	go build -v -tags '$(TAGS)' -ldflags '$(EXTLDFLAGS)-s -w $(LDFLAGS)' -o bin/$@ ./cmd

.PHONY: misspell-check
misspell-check:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -error $(GOFILES)

.PHONY: misspell
misspell:
	@hash misspell > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/client9/misspell/cmd/misspell; \
	fi
	misspell -w $(GOFILES)

unused-check:
	@hash unused > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u honnef.co/go/tools/cmd/unused; \
	fi
	for PKG in $(PACKAGES); do unused $$PKG || exit 1; done;

test:
	for PKG in $(PACKAGES); do go test -v $$PKG || exit 1; done;

.PHONY: test-vendor
test-vendor:
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	govendor list +unused | tee "$(TMPDIR)/wc-gitea-unused"
	[ $$(cat "$(TMPDIR)/wc-gitea-unused" | wc -l) -eq 0 ] || echo "Warning: /!\\ Some vendor are not used /!\\"

	govendor list +outside | tee "$(TMPDIR)/wc-gitea-outside"
	[ $$(cat "$(TMPDIR)/wc-gitea-outside" | wc -l) -eq 0 ] || exit 1

	govendor status || exit 1

release: release-dirs release-build release-copy release-check

release-dirs:
	mkdir -p $(DIST)/binaries $(DIST)/release

release-build:
	@hash gox > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/mitchellh/gox; \
	fi
	gox -os="$(TARGETS)" -arch="$(ARCHS)" -tags="$(TAGS)" -ldflags="$(EXTLDFLAGS)-s -w $(LDFLAGS)" -output="$(DIST)/binaries/$(EXECUTABLE)-$(VERSION)-{{.OS}}-{{.Arch}}"

release-copy:
	$(foreach file,$(wildcard $(DIST)/binaries/$(EXECUTABLE)-*),cp $(file) $(DIST)/release/$(notdir $(file));)

release-check:
	cd $(DIST)/release; $(foreach file,$(wildcard $(DIST)/release/$(EXECUTABLE)-*),sha256sum $(notdir $(file)) > $(notdir $(file)).sha256;)

docker_build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags '$(TAGS)' -ldflags "$(EXTLDFLAGS)-s -w $(LDFLAGS)" -o bin/$(EXECUTABLE) ./cmd

build_image:
	docker build -t $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE) -f Dockerfile .

docker_release: build_image

clean:
	go clean -x -i ./...
	rm -rf bin
