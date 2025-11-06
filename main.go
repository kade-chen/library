package main

import (
	_ "github.com/kade-chen/google-billing-console/apps" //registry all impl/api

	// 开启Health健康检查
	_ "github.com/kade-chen/library/ioc/apps/health/restful"

	"github.com/kade-chen/google-billing-console/cmd"
)

func main() {
	cmd.Execute()
}
