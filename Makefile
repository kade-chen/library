PROJECT_NAME := "billing"
MAIN_FILE_PAHT := "main.go"
PKG := "github.com/kade-chen/google-billing-console"
IMAGE_PREFIX := "github.com/kade-chen/google-billing-console"

DOCKER_IMAGE_PREFIX := "mcenter"
DOCKER_IMAGE_TAG := "latest"
DOCKER_BUILD_IMAGE_NAME := $(DOCKER_IMAGE_PREFIX):$(DOCKER_IMAGE_TAG)

MOD_DIR := $(shell go env GOPATH)/pkg/mod
##go list "github.com/kade-chen/google-billing-console" | grep -v /vendor/ | grep -v redis
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/ | grep -v redis)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

##MCUBE_MODULE := "github.com/infraboard/mcube/v2"
MCUBE_VERSION :=$(shell go list -m ${MCUBE_MODULE} | cut -d' ' -f2)
MCUBE_PKG_PATH := ${MOD_DIR}/${MCUBE_MODULE}@${MCUBE_VERSION}

BUILD_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
BUILD_COMMIT := ${shell git rev-parse HEAD}
BUILD_TIME := ${shell date '+%Y-%m-%d %H:%M:%S'}
BUILD_GO_VERSION := $(shell go version | grep -o  'go[0-9].[0-9].*')
VERSION_PATH := "github.com/infraboard/mcube/v2/ioc/config/application"
OUTPUT_NAME := "mcenter-api"

IMAGE_BUILD_TIME := ${shell date '+%Y%m%d'}
IMAGE_BUILD_COMMIT :=  ${shell git rev-parse HEAD | cut -c 1-8}
IMAGE_VERSION := "mpaas-api:${IMAGE_BUILD_TIME}-${BUILD_BRANCH}-${IMAGE_BUILD_COMMIT}"

.PHONY: all dep lint vet test test-coverage build clean

all: build

mod: ## Get the dependencies
	@go mod tidy

lint: ## Lint Golang files
	@golint -set_exit_status ${PKG_LIST}

vet: ## Run go vet
	@go vet ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}
	
test-coverage: ## Run tests with coverage
	@go test -short -coverprofile cover.out -covermode=atomic ${PKG_LIST} 
	@cat cover.out >> coverage.txt

build: dep ## Build the binary file
	@go build -a -o dist/${OUTPUT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${MAIN_FILE}

linux: dep ## Build the binary file
	@GOOS=linux GOARCH=amd64 go build -a -o dist/${OUTPUT_NAME} -ldflags "-s -w" -ldflags "-X '${VERSION_PATH}.GIT_BRANCH=${BUILD_BRANCH}' -X '${VERSION_PATH}.GIT_COMMIT=${BUILD_COMMIT}' -X '${VERSION_PATH}.BUILD_TIME=${BUILD_TIME}' -X '${VERSION_PATH}.GO_VERSION=${BUILD_GO_VERSION}'" ${MAIN_FILE}

image: dep ## Build the docker image
	docker build -t ${IMAGE_VERSION} -f Dockerfile .

init: dep ## Inital project 
	@go run main.go init

run: ## Run Server
	@go run main.go start

runp: dep ## Run product Server
	@go run main.go start -f etc/config.toml

clean: ## Remove previous build
	@go clean .
	@rm -f dist/${PROJECT_NAME}

install: ## Install depence go package
	@go install github.com/infraboard/mcube/v2/cmd/mcube@latest
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/favadi/protoc-go-inject-tag@latest

pb: ## Copy mcube protobuf files to common/pb
	@mkdir -pv common/pb/github.com/infraboard/mcube/v2/pb
	@cp -r ${MCUBE_PKG_PATH}/pb/* common/pb/github.com/infraboard/mcube/v2/pb
	@sudo rm -rf common/pb/github.com/infraboard/mcube/v2/pb/*/*.go

gen: ## protobuf 编译
	@protoc -I=.. --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} ../google-billing-console/apps/*/pb/*.proto
	@go fmt ./...
	@protoc-go-inject-tag -input=apps/*/*.pb.go
	@kade-library enum -p -m apps/*/*.pb.go


help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

format: ## 格式化代码
	@clang-format -i apps/*/*/pb/*.proto

kade-library: ## 清理
	@go build -o ~/go/bin/kade-library github.com/kade-chen/library/cmd
	@sudo ln -s ~/go/bin/kade-library /usr/local/bin/kade-library  
## 软连接删除后面的链接

docker-build: ## 构建镜像
	@GOOS="linux" go build -o main .
	@DOCKER_BUILDKIT=1 docker buildx build -t ${DOCKER_BUILD_IMAGE_NAME} .

stt-grpc-server: ## 启动stt grpc server
	@go run /Users/kade.chen/go-kade-project/github/mcenter/apps/stt/example/sttv2/grpc/server/main.go

stt-grpc-client: ## 启动stt grpc client
	@go run /Users/kade.chen/go-kade-project/github/mcenter/apps/stt/example/sttv2/grpc/client/main.go

login: ## 登录
	@go run /Users/kade.chen/go-kade-project/github/google-billing-console/test/main.go

