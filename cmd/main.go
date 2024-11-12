package main

import (
	"fmt"
	"os"

	"github.com/kade-chen/library/cmd/generate"
	"github.com/spf13/cobra"
)

func main() {
	RootCmd.Execute()
}

var version bool

// Rootcmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "kade-library",
	Short: "kade-library 分布式服务构建工具1",
	Long:  "kade-library 分布式服务构建工具",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}

func init() {
	RootCmd.AddCommand(generate.Cmd)
	RootCmd.PersistentFlags().BoolVarP(&version, "version", "v", false, fmt.Sprintf("是否开启debug模式,默认为%s", "v2.1.3"))
}
