package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// relearnCmd represents the relearn command
var relearnCmd = &cobra.Command{
	Use:   "relearn",
	Short: "relearn all classifications",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		relearn()
	},
}

func init() {
	RootCmd.AddCommand(relearnCmd)
}

func relearn() {
	if err := myApp.RelearnTags(); err != nil {
		panic(err)
	}

	fmt.Printf("All tags: %+v\n", myApp.AllTags())
}
