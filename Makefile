BUILD_DIR = build
SERVICES = migrate
CGO_ENABLED ?= 0
GOARCH ?= amd64
COMMIT ?= $(shell git rev-parse HEAD)
TIME ?= $(shell date +%F_%T)

define compile_service
	CGO_ENABLED=$(CGO_ENABLED) GOOS=$(GOOS) GOARCH=$(GOARCH) GOARM=$(GOARM) \
	go build -mod=vendor -ldflags "-s -w \
	-X 'github.com/mainflux/mainflux.BuildTime=$(TIME)' \
	-X 'github.com/mainflux/mainflux.Commit=$(COMMIT)'" \
	-o ${BUILD_DIR}/mainflux-$(1) cmd/main.go
endef

all: $(SERVICES)

.PHONY: all $(SERVICES)

clean:
	rm -rf ${BUILD_DIR}

install:
	cp ${BUILD_DIR}/* $(GOBIN)

test:
	go test -mod=vendor -v -race -count 1 -tags test $(shell go list ./... | grep -v 'vendor\|cmd')


$(SERVICES):
	$(call compile_service,$(@))


changelog:
	git log $(shell git describe --tags --abbrev=0)..HEAD --pretty=format:"- %s"
