package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// infoCmd represents the info command
var infoCmd = &cobra.Command{
	Use:   "info file",
	Short: "show info for a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		info(args[0])
	},
}

func init() {
	RootCmd.AddCommand(infoCmd)
}

func info(file string) {
	info, err := myApp.ReadFileInfo(file)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(info.JSON()))
}
