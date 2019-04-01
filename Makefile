VERSION := $(shell cat ./VERSION)
BUILDFLAGS := -ldflags "-s -w -X gitlab.com/JanMa/go-pass/cmd.Version=${VERSION}"

build:
		GO111MODULE=on CGO_ENABLED=0 go build -mod=vendor ${BUILDFLAGS}

get:
		GO111MODULE=on go get -v

.PHONY: vendor
vendor:	get
		GO111MODULE=on go mod vendor

install:
		GO111MODULE=on CGO_ENABLED=0 go install -mod=vendor  ${BUILDFLAGS}

clean:
		@rm -rf go-pass
