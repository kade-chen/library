package start

import (
	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/server"
	"github.com/spf13/cobra"

	// 非功能性模块
	_ "github.com/kade-chen/library/ioc/apps/apidoc/restful"
	_ "github.com/kade-chen/library/ioc/apps/health/restful"
	_ "github.com/kade-chen/library/ioc/apps/metric/restful"
	_ "github.com/kade-chen/library/ioc/config/cors"
	_ "github.com/kade-chen/library/ioc/config/grpc"
)

var Req *ioc.LoadConfigRequest

var Cmd = &cobra.Command{
	Use:   "start",
	Short: "start the server",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(server.Run(cmd.Context(), Req))
		// 执行启动服务器的逻辑
	},
}

var Cmd1 = &cobra.Command{
	Use:   "start111",
	Short: "start the server",
	Run: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(server.Run(cmd.Context(), Req))
		// 执行启动服务器的逻辑
	},
}
