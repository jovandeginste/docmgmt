package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

// ingestCmd represents the ingest command
var (
	ingestCmd = &cobra.Command{
		Use:  "ingest",
		Long: `Interactively process new files in the document root`,
		Run: func(cmd *cobra.Command, args []string) {
			ingest()
		},
	}

	promptParse = promptui.Select{
		Label: "Would you like to parse this file?",
		Items: []string{"yes", "no"},
	}
)

func init() {
	RootCmd.AddCommand(ingestCmd)
	ingestCmd.Flags().StringSliceVarP(&tags, "tag", "t", []string{}, "Tag to ingest")
}

func ingest() {
	err := myApp.StartServer()
	if err != nil {
		panic(err)
	}

	defer myApp.StopServer()

	err = filepath.Walk(myApp.Configuration.DocumentRoot, ingestParser)
	if err != nil {
		panic(err)
	}
}

func ingestParser(path string, info os.FileInfo, err error) error {
	if info.IsDir() {
		return nil
	}

	if err != nil {
		return err
	}

	path = strings.TrimPrefix(path, myApp.Configuration.DocumentRoot)

	i, err := myApp.ReadAbsoluteFileInfo(path)
	if err != nil {
		return err
	}

	if !i.IsNew() {
		return nil
	}

	fmt.Println("Found new file:", path)

	return askParse(path)
}

func askParse(path string) error {
	_, action, err := promptParse.Run()

	if err == promptui.ErrInterrupt {
		os.Exit(0)
	}

	if action == "yes" {
		i, err := myApp.Parse(myApp.Configuration.DocumentRoot + path)
		if err != nil {
			return err
		}

		err = i.Write()
		if err != nil {
			return err
		}

		fmt.Println("New file was parsed.")

		err = manageInteractive(i)
		return err
	}

	return err
}
