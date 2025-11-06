package cmd

import (
	"fmt"
	"os"

	"github.com/kade-chen/library/ioc"
	"github.com/kade-chen/library/ioc/config/application"
	"github.com/kade-chen/google-billing-console/cmd/start"
	"github.com/spf13/cobra"
)

var (
	configType string
	configFile string
)

var vers bool

var RootCmd = &cobra.Command{
	Use:   "billing",
	Short: "微服务公共能力中心", // 简短error message描述
	Long:  "微服务公共能力中心", // 详细error message描述
	RunE: func(cmd *cobra.Command, args []string) error {
		if vers {
			fmt.Println(application.FullVersion())
			return nil
		}
		return cmd.Help()
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ss")
	},
}

// init ioc
func initcmd() {
	req := ioc.NewLoadConfigRequest()
	req.ConfigFile.Enabled = true
	switch configType {
	case "file":
		req.ConfigFile.Enabled = true
		req.ConfigFile.Path = configFile
	default:
		req.ConfigFile.Enabled = true
		os.Setenv("LOG_LEVEL", "info")
	}
	start.Req = req
}

func Execute() {
	// cobra.OnInitialize() 初始化Cobra，用于启动Cobra插件
	cobra.OnInitialize(initcmd)
	RootCmd.AddCommand(start.Cmd)
	RootCmd.AddCommand(start.Cmd1)
	// RootCmd.AddCommand(initial.Cmd)
	err := RootCmd.Execute()
	cobra.CheckErr(err)
	// fmt.Println("ss")
}

func init() {
	RootCmd.PersistentFlags().StringVarP(&configType, "config-type", "t", "file", "the service config type")
	RootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "f", "etc/config.toml", "the service config from file")
}
