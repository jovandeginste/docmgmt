package cmd

import (
	"github.com/spf13/cobra"
)

// parseCmd represents the parse command
var parseCmd = &cobra.Command{
	Use:   "parse [file...]",
	Short: "Parse file",
	Long:  `Parse a file and add it to the document management`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		parse(args)
	},
}

func init() {
	RootCmd.AddCommand(parseCmd)
}

func parse(files []string) {
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
		err = i.Write()
		if err != nil {
			panic(err)
		}
	}
}
