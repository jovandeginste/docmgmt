package cmd

import (
	"github.com/spf13/cobra"
)

// tagCmd represents the tag command
var tagCmd = &cobra.Command{
	Use:   "tag file tag",
	Short: "Add a tag to a file",
	Long:  `Add a tag to a file`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		tag(args[0], args[1])
	},
}

func init() {
	RootCmd.AddCommand(tagCmd)
}

func tag(file string, tag string) {
	body, err := myApp.ReadFileBody(file)
	if err != nil {
		panic(err)
	}

	err := myApp.WriteMetadata(file)
	if err != nil {
		panic(err)
	}
}
