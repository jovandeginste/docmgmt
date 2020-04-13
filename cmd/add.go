package cmd

import (
	"github.com/spf13/cobra"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:  "add file... [--tag t1...]",
	Long: `Parse files and add it with some tags to the document management`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		add(args, tags)
	},
}
var tags []string

func init() {
	RootCmd.AddCommand(addCmd)
	addCmd.Flags().StringSliceVarP(&tags, "tag", "t", []string{}, "Tag to add")
}

func add(files []string, tags []string) {
	err := myApp.StartServer()
	if err != nil {
		panic(err)
	}
	defer myApp.StopServer()

	for _, file := range files {
		i, err := myApp.Parse(file)
		if err != nil {
			panic(err)
		}

		addTags(i, tags)

		err = i.Write()
		if err != nil {
			panic(err)
		}
	}
}
