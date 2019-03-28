VERSION := $(shell cat ./VERSION)
BUILDFLAGS := -ldflags "-X gitlab.com/JanMa/go-pass/cmd.Version=${VERSION}"

build:
		GO111MODULE=on GOOS=linux go build -mod=vendor ${BUILDFLAGS}

get:
		GO111MODULE=on go get -v

.PHONY: vendor
vendor:	get
		GO111MODULE=on go mod vendor

install: build
		GO111MODULE=on go install -mod=vendor

clean:
		@rm -rf go-pass
