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
	fmt.Printf("All tags: %+v\n", myApp.AllTags())

	var picks []string

	i, err := myApp.ReadFileInfo(file)
	if err != nil {
		panic(err)
	}

	if i.Body == nil {
		panic(fmt.Errorf("file is not indexed: '%s'", file))
	}

	threshold := 1.0
	threshold /= float64(len(i.Tags) + 1)

	if threshold > 0.5 {
		threshold = 0.5
	}

	fmt.Printf("Current tags: %+v\n", i.Tags)
	sugg := myApp.Classify(i.Body.Content)

sugg:
	for _, p := range sugg {
		if p.Score < threshold {
			fmt.Printf("First suggestion below threshold (%f): %+v\n", threshold, p)
			break
		}

		klass := string(p.Class)
		for _, t := range i.Tags {
			if t == klass {
				continue sugg
			}
		}

		picks = append(picks, klass)
	}
	fmt.Printf("Suggested new tags: %+v\n", picks)
}
