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
	var picks []string

	i, err := myApp.ReadFileInfo(file)
	if err != nil {
		panic(err)
	}

	if i.Body == nil {
		panic(fmt.Errorf("file is not indexed: '%s'", file))
	}

	fmt.Printf("Current tags: %+v\n", i.Tags)
	sugg := myApp.Classify(i.Body.Content)
	total := float64(0)

sugg:
	for _, p := range sugg {
		total += p.Score
		klass := string(p.Class)
		for _, t := range i.Tags {
			if t == klass {
				continue sugg
			}
		}

		picks = append(picks, klass)
		if total > 0.5 {
			break
		}
	}
	fmt.Printf("Suggested tags: %+v\n", picks)

	fmt.Printf("All tags: %+v\n", myApp.AllTags())
}
