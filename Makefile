PROTO_GO_GEN=pb

PROTOC_LINUX_VERSION = 3.11.4
PROTOC_LINUX_ZIP = protoc-$(PROTOC_LINUX_VERSION)-linux-x86_64.zip
GO_CMD_MAIN = 	cmd/main.go

BUF_VERSION=0.24.0
BUF_BINARY_NAME=buf


all: make-dir install-linux copy-grpc-template gen-proto

ci: make-dir install-ci copy-grpc-template gen-proto

clean:
	rm -r $(PROTO_GO_GEN)

make-dir: clean
	mkdir -p $(PROTO_GO_GEN)

install-linux: install-protoc install-protoc-go install-lint install-go-tools

install-protoc:
	curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_LINUX_VERSION)/$(PROTOC_LINUX_ZIP)
	sudo unzip -o $(PROTOC_LINUX_ZIP) -d /usr/local bin/protoc
	sudo unzip -o $(PROTOC_LINUX_ZIP) -d /usr/local 'include/*'
	rm -f $(PROTOC_LINUX_ZIP)

install-protoc-go:
	go get github.com/golang/protobuf/protoc-gen-go@v1.4.3
	go get google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.0.1
	go get github.com/envoyproxy/protoc-gen-validate@v0.4.1
	go get github.com/gogo/protobuf/protoc-gen-gofast@v1.3.1
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway@v1.14.7
	go get github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger@v1.14.7

install-lint:
	sudo apt-get update
	sudo apt-get install -y curl
	sudo curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sudo sh -s -- -b $(go env GOPATH)/bin v1.30.0
	sudo curl -sSL "https://github.com/bufbuild/buf/releases/download/v$(BUF_VERSION)/$(BUF_BINARY_NAME)-$(shell uname -s)-$(shell uname -m)"  -o "/usr/local/bin/$(BUF_BINARY_NAME)" && sudo chmod +x "/usr/local/bin/$(BUF_BINARY_NAME)"

install-go-tools:
	go get github.com/vektra/mockery/v2@v2.9.0
	go get github.com/google/wire/cmd/wire@v0.5.0

gen-proto: clean make-dir
	env PROTO_GO_GEN=$(PROTO_GO_GEN) sh scripts/gen-proto


migrate:
	echo \# make migrate name="$(name)"
	go run $(GO_CMD_MAIN) migrate create $(name)

migrate-up:
	go run $(GO_CMD_MAIN) migrate up

migrate-down-1:
	go run $(GO_CMD_MAIN) migrate down 1

check-lint:
	buf check lint
	GOSUMDB=off go vet ./...
	golangci-lint run --timeout=2m

test:
	go test ./...  -count=1 -cover -race

run:
	go run $(GO_CMD_MAIN) server
