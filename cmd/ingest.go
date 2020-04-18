package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/jovandeginste/docmgmt/app"
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

	promptMore = promptui.Select{
		Label: "File was parsed. What would you like to do?",
		Items: []string{
			"read the parsed text",
			"open the file with the default program",
			"manage tags",
			"continue to the next file",
		},
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

		err = askMore(i)
		return err
	}

	return err
}

func askMore(i *app.Info) error {
	for {
		_, action, err := promptMore.Run()
		fmt.Println("action:", action)

		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}

		switch action {
		case "read the parsed text":
			read(i.AbsoluteFilename())
		case "open the file with the default program":
			err = i.OpenWithDefaultApp()
			if err != nil {
				return err
			}
		case "manage tags":
		case "continue to the next file":
			return nil
		}
	}
}
