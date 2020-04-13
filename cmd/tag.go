package cmd

import (
	"fmt"

	"github.com/jovandeginste/docmgmt/app"
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
var delete bool

func init() {
	RootCmd.AddCommand(tagCmd)
	tagCmd.Flags().BoolVarP(&delete, "delete", "d", false, "Delete tag instead of adding")
}

func tag(file string, tag string) {
	i, err := myApp.ReadFileInfo(file)
	if err != nil {
		panic(err)
	}

	if delete {
		err = deleteTag(i, tag)
	} else {
		err = addTag(i, tag)
	}
	if err != nil {
		panic(err)
	}

	fmt.Printf("Current tags for %s:\n%v\n", file, i.Tags)
}

func deleteTag(info *app.Info, tag string) error {
	info.DeleteTag(tag)

	return nil
}

func addTag(info *app.Info, tag string) error {

	info.AddTag(tag)

	myApp.Learn(info.Body.Content, tag)

	err := myApp.SaveClassifier()
	if err != nil {
		return err
	}

	info.Write()

	return nil
}
