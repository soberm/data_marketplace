GOOS=linux
GOARCH=386

.PHONY: all build clean proto deps abigen docker

BIN_FOLDER=./bin

#VERSION = $(shell git describe --tags)
VERSION=v0.0.1
BUILD=$(shell date +%FT%T%z)
GIT_COMMIT=$(shell git rev-list -1 HEAD)

LD_FLAGS="-X marketplace-services/pkg/broker.version=$(VERSION) -X marketplace-services/pkg/broker.build=$(BUILD)\
	-X marketplace-services/pkg/broker.gitCommit=$(GIT_COMMIT) -X marketplace-services/pkg/proxy.version=$(VERSION)\
	-X marketplace-services/pkg/proxy.build=$(BUILD) -X marketplace-services/pkg/proxy.gitCommit=$(GIT_COMMIT)\
	-X marketplace-services/pkg/provider.version=$(VERSION) -X marketplace-services/pkg/provider.build=$(BUILD)\
	-X marketplace-services/pkg/provider.gitCommit=$(GIT_COMMIT) -X marketplace-services/pkg/consumer.version=$(VERSION)\
	-X marketplace-services/pkg/consumer.build=$(BUILD) -X marketplace-services/pkg/consumer.gitCommit=$(GIT_COMMIT)"

all: clean proto abigen test build docker

build:
	mkdir -p $(BIN_FOLDER)
	GOOS=$(GOOS) GOARCH=$(GOARCH) CGO_ENABLED=1 go build -o $(BIN_FOLDER) -ldflags $(LD_FLAGS) ./...

test:
	GOOS=$(GOOS) GOARCH=$(GOARCH) go test ./...

docker:
	docker build --no-cache -t broker -f ./build/broker/Dockerfile .
	docker build --no-cache -t proxy -f ./build/proxy/Dockerfile .

proto:
	@ if ! which protoc > /dev/null; then \
		echo "error: protoc not installed" >&2; \
		exit 1; \
	fi
	protoc -I ./api/proto  --go_out=plugins=grpc:.. ./api/proto/domain/*.proto
	protoc -I ./api/proto --go_out=plugins=grpc:.. ./api/proto/broker/*.proto
	protoc -I ./api/proto --go_out=plugins=grpc:.. ./api/proto/proxy/*.proto

abigen:
	@ if ! which abigen > /dev/null; then \
    		echo "error: abigen not installed" >&2; \
    		exit 1; \
	fi
	mkdir -p ./pkg/contracts
	abigen --sol ../marketplace-contracts/contracts/UserContract.sol --pkg bindings --out ./pkg/contracts/bindings/marketplace_contracts.go

deps:
	go get -d -v google.golang.org/grpc/...

clean:
	rm -rf ./bin
	rm -rf ./pkg/contracts/bindings/marketplace_contracts.go
	find . -type f -name "*.pb.go" -delete