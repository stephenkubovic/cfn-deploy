install: build
	mv cfn-deploy /usr/local/bin

build:
	go build -o cfn-deploy cmd/cfn-deploy/main.go

fmt:
	go fmt ./...

test: fmt
	go test ./... -v -timeout 30s
