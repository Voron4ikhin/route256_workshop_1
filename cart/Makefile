CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.64.8
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
BUFVER=v1.72.0
BUFBIN=${BINDIR}/buf_${GOVER}_${BUFVER}
PROTOC_GEN_GO_VER=v1.36.11
PROTOC_GEN_GO_BIN=${BINDIR}/protoc-gen-go_${GOVER}_${PROTOC_GEN_GO_VER}
PROTOC_GEN_GO_GRPC_VER=v1.6.1
PROTOC_GEN_GO_GRPC_BIN=${BINDIR}/protoc-gen-go-grpc_${GOVER}_${PROTOC_GEN_GO_GRPC_VER}
PACKAGE=route256/cart/cmd/app

all: format build test lint

build: bindir
	go build -o ${BINDIR}/app ${PACKAGE}

test:
	go test ./...

run:
	go run ${PACKAGE}

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

install-smartimports: bindir
	test -f ${SMARTIMPORTS} || \
		(GOBIN=${BINDIR} go install github.com/pav5000/smartimports/cmd/smartimports@latest && \
		mv ${BINDIR}/smartimports ${SMARTIMPORTS})

# generate regenerates Go code from the shared .proto contracts in ../proto
# (buf.yaml / buf.gen.yaml live at the repo root, since the contracts are shared with loms).
generate: install-proto-deps
	cd .. && PATH="${BINDIR}:$$PATH" ${BUFBIN} generate

install-proto-deps: install-buf install-protoc-gen-go install-protoc-gen-go-grpc

install-buf: bindir
	test -f ${BUFBIN} || \
		(GOBIN=${BINDIR} go install github.com/bufbuild/buf/cmd/buf@${BUFVER} && \
		mv ${BINDIR}/buf ${BUFBIN})

install-protoc-gen-go: bindir
	test -f ${PROTOC_GEN_GO_BIN} || \
		(GOBIN=${BINDIR} go install google.golang.org/protobuf/cmd/protoc-gen-go@${PROTOC_GEN_GO_VER} && \
		mv ${BINDIR}/protoc-gen-go ${PROTOC_GEN_GO_BIN})
	ln -sf ${PROTOC_GEN_GO_BIN} ${BINDIR}/protoc-gen-go

install-protoc-gen-go-grpc: bindir
	test -f ${PROTOC_GEN_GO_GRPC_BIN} || \
		(GOBIN=${BINDIR} go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@${PROTOC_GEN_GO_GRPC_VER} && \
		mv ${BINDIR}/protoc-gen-go-grpc ${PROTOC_GEN_GO_GRPC_BIN})
	ln -sf ${PROTOC_GEN_GO_GRPC_BIN} ${BINDIR}/protoc-gen-go-grpc
