export CGO_ENABLED=0
export GO111MODULE=on

.PHONY: build

ATOMIX_SIMULATIONS_VERSION := latest

all: build

build: # @HELP build the source code
build:
	GOOS=linux GOARCH=amd64 go build -o build/_output/kubernetes-simulations ./cmd/kubernetes-simulations

test: # @HELP run the unit tests and source code validation
test: build license_check linters
	go test github.com/atomix/kubernetes-simulations/...

linters: # @HELP examines Go source code and reports coding problems
	golangci-lint run

license_check: # @HELP examine and ensure license headers exist
	@if [ ! -d "../build-tools" ]; then cd .. && git clone https://github.com/onosproject/build-tools.git; fi
	./../build-tools/licensing/boilerplate.py -v --rootdir=${CURDIR}

proto: # @HELP build Protobuf/gRPC generated types
proto:
	docker run -it -v `pwd`:/go/src/github.com/atomix/kubernetes-simulations \
		-w /go/src/github.com/atomix/kubernetes-simulations \
		--entrypoint build/bin/compile_protos.sh \
		onosproject/protoc-go:stable

images: # @HELP build kubernetes-simulations Docker image
images: build
	docker build . -f build/docker/Dockerfile -t atomix/kubernetes-simulations:${ATOMIX_SIMULATIONS_VERSION}
