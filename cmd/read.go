package cmd

import (
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

// readCmd represents the read command
var readCmd = &cobra.Command{
	Use:   "read file",
	Short: "read the contents of a file",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		read(args[0])
	},
}

func init() {
	RootCmd.AddCommand(readCmd)
}

func read(file string) {
	info, err := myApp.ReadFileInfo(file)
	if err != nil {
		panic(err)
	}
	pager, ok := os.LookupEnv("PAGER")
	if !ok {
		pager = "less"
	}
	splitPager := strings.Fields(pager)
	cmd := exec.Command(splitPager[0], splitPager[1:]...)

	// Feed it with the string you want to display.
	cmd.Stdin = strings.NewReader(info.Body.Content)

	// This is crucial - otherwise it will write to a null device.
	cmd.Stdout = os.Stdout

	// Fork off a process and wait for it to terminate.
	err = cmd.Run()
	if err != nil {
		panic(err)
	}
}
