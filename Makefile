build:
		GO111MODULE=on GOOS=linux go build -mod=vendor

get:
		GO111MODULE=on go get -v

.PHONY: vendor
vendor:
		GO111MODULE=on go mod vendor

install: build
		GO111MODULE=on go install -mod=vendor

clean:
		@rm -rf go-pass
