package generate

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/kade-chen/library/cmd/generate/enum"
	"github.com/spf13/cobra"
)

// EnumCmd 枚举生成器
var Cmd = &cobra.Command{
	Use:   "enum",
	Short: "枚举生成器",
	Long:  "枚举生成器",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			return
		}

		var matchedFiles []string

		for _, arg := range args {
			filesName, err := filepath.Glob(arg)
			cobra.CheckErr(err)

			// 只匹配Go源码文件
			if strings.HasSuffix(arg, ".go") {
				matchedFiles = append(matchedFiles, filesName...)
			}

			if len(matchedFiles) == 0 {
				return
			}

			for _, path := range matchedFiles {
				// 生成代码
				code, err := enum.G.Generate(path)
				cobra.CheckErr(err)

				if len(code) == 0 {
					continue
				}

				var genFile = ""
				if strings.HasSuffix(path, ".pb.go") {
					genFile = strings.ReplaceAll(path, ".pb.go", "_enum.pb.go")
				} else {
					genFile = strings.ReplaceAll(path, ".go", "_enum.go")
				}

				// 写入文件
				err = os.WriteFile(genFile, code, 0644)
				cobra.CheckErr(err)
			}
		}
	},
}

func init() {
	Cmd.PersistentFlags().BoolVarP(&enum.G.Marshal, "marshal", "m", false, "is generate json MarshalJSON and UnmarshalJSON method")
	Cmd.PersistentFlags().BoolVarP(&enum.G.ProtobufExt, "protobuf_ext", "p", false, "is generate protobuf extention method")
}
