package cmd

import (
	"fmt"

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
	i, err := myApp.ReadFileInfo(file)
	if err != nil {
		panic(err)
	}
	i.AddTag(tag)

	myApp.Learn(i.Body.Content, tag)

	err = myApp.SaveClassifier()
	if err != nil {
		panic(err)
	}

	i.Write()

	fmt.Printf("Current tags for %s:\n%v\n", file, i.Tags)
}
