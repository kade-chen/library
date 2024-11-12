PKG := "gitee.com/go-kade/library/v2"

main: ## Run Server
	@ go run main.go

gin: ## Run gin Server
	@ go run examples/http_gin/main.go

restful: ## Run resful Server
	@ go run examples/http_go_restful/main.go

mod: ## Run Mod Tidy
	@ go mod tidy

install: ## Install depence go package
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@go install github.com/favadi/protoc-go-inject-tag@latest

gen: ## protobuf 编译
	@protoc -I=. --go_out=. --go_opt=module=${PKG} --go-grpc_out=. --go-grpc_opt=module=${PKG} pb/*/*.proto
	@protoc-go-inject-tag -input=pb/*/*.pb.go

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

cobra: ## Run cobra Server
	@go build -o ~/go/bin/kade-library cmd/main.go

kade-library: ## Run generate cobra command
	@/Users/kade.chen/go/bin/kade-library enum -m -p pb/*/*.pb.go 