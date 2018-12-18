TARGET    = bdisplay

GO        = go
GOIMPORTS = goimports
GOLINT    = golint

.PHONY: all
all: fmt lint vet deps test build

# Formats the source code.
.PHONY: fmt
fmt:
	@ret=0 && for f in $$(find . -type f -name '*.go'); do \
		$(GOIMPORTS) -l -w $$f || ret=$$? ; \
	done ; exit $$ret

# Lints the source code.
.PHONY: lint
lint:
	@ret=0 && for d in $$($(GO) list ./...); do \
		$(GOLINT) -set_exit_status $${d} || ret=$$? ; \
	done ; exit $$ret

# Performs simple static analysis on the source code.
.PHONY: vet
vet:
	@$(GO) vet $($(GO) list ./...)

# Gets all our dependencies.
.PHONY: deps
deps:
	@$(GO) get -d ./...

# Runs unit tests.
.PHONY: test
test:
	@$(GO) test -v ./...

# Builds the project for all target platforms.
.PHONY: build
build: build_darwin_amd64 build_linux_amd64 build_linux_arm

build_darwin_%: GOOS=darwin
build_linux_%: GOOS=linux

build_%_amd64: GOARCH=amd64
build_%_arm: GOARCH=arm
build_%_arm: GOARM=5

build_%:
	@GOOS=${GOOS} GOARCH=$(GOARCH) GOARM=$(GOARM) CGO_ENABLED=0 \
		$(GO) build  -a -installsuffix cgo -ldflags="-w -s" -o build/$(GOOS)_$(GOARCH)/$(TARGET) .

# Cleans the compiled binaries.
.PHONY: clean
clean:
	@rm -rf build
