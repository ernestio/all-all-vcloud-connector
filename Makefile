install:
	go install -v

build:
	go build -v ./...

deps:
	go get -u github.com/golang/dep/cmd/dep
	dep ensure

dev-deps: deps
	go get github.com/golang/lint/golint

test:
	go test -v ./...

lint:
	golint ./...
	go vet ./...