install:
	go install -v

build:
	go build -v ./...

deps:
		go get github.com/nats-io/nats
		go get github.com/ernestio/ernest-config-client
		go get github.com/r3labs/vcloud-go-sdk

dev-deps:
	go get github.com/golang/lint/golint

test:
	go test -v ./...

lint:
	golint ./...
	go vet ./...