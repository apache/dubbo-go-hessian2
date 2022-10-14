package exec

import (
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/hsf"
	"github.com/apache/dubbo-go-hessian2/toolkit-go/internal/idl"
	"github.com/spf13/cobra"
)

var (
	hsfOpt hsf.Opt
	idlOpt idl.Opt
)

var hsfCmd = &cobra.Command{
	Use:   "hsf",
	Short: "hsf-go的命令行工具",
	RunE: func(cmd *cobra.Command, args []string) error {
		return hsf.Run(hsfOpt)
	},
}

func init() {
	rootCmd.AddCommand(hsfCmd)
	// in
	hsfCmd.Flags().StringVar(&(hsfOpt.InJarPath), "inJarPath", "", "输入jar包的地址")
	hsfCmd.Flags().StringVar(&(hsfOpt.InGroupID), "inGroupId", "", "输入jar包的分组")
	hsfCmd.Flags().StringVar(&(hsfOpt.InArtifactID), "inArtifactId", "", "输入jar包的唯一ID")
	hsfCmd.Flags().StringVar(&(hsfOpt.InVersion), "inVersion", "", "输入jar包的版本")
	hsfCmd.Flags().StringVar(&(hsfOpt.InClass), "inClass", "", "输入读取的class全路径")
	// out
	hsfCmd.Flags().StringVar(&(hsfOpt.OutType), "outType", "consumer",
		"生成的结果类型，包含[consumer/provider/struct/listClass/clean]")
	hsfCmd.Flags().StringVar(&(hsfOpt.OutPath), "outPath", "", "生成的结果的保存路径")
	hsfCmd.Flags().StringVar(&(hsfOpt.OutPkg), "outPkg", "", "生成结果的包名")
	hsfCmd.Flags().BoolVar(&(hsfOpt.OutQuote), "outQuote", false, "封装类型是否生成go引用类型")
}
