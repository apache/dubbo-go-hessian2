package exec

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gokit",
	Short:   "gokit是一个golang的工具集合",
	Version: "v0.0.6",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
