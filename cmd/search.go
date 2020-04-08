package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search query...",
	Short: "search based on some queries",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		search(args)
	},
}

func init() {
	RootCmd.AddCommand(searchCmd)
}

func search(queries []string) {
	search, err := myApp.Search(queries)
	if err != nil {
		panic(err)
	}

	for _, s := range search {
		fmt.Printf("%#v\n", s.Element.Metadata.Filename)
	}
}
