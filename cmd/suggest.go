package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// suggestCmd represents the suggest command
var suggestCmd = &cobra.Command{
	Use:   "suggest file",
	Short: "suggest file",
	Long:  `suggest tags for a file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		suggest(args[0])
	},
}

func init() {
	RootCmd.AddCommand(suggestCmd)
}

func suggest(file string) {
	body, err := myApp.ReadFileBody(file)
	if err != nil {
		panic(err)
	}

	sugg := myApp.Classify(body)
	fmt.Println("Suggested tag:", sugg)
}
