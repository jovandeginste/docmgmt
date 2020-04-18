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
	promptTags = promptui.Select{
		Label: "What would you like to do?",
		Items: []string{
			"add tags",
			"delete tags",
			"done",
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
			err = manageTags(i)
			if err != nil {
				return err
			}
		case "continue to the next file":
			return nil
		}
	}
}

func manageTags(i *app.Info) error {
	for {
		fmt.Printf("Current tags: %v\n", i.Tags)

		_, action, err := promptTags.Run()

		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}

		switch action {
		case "add tags":
			err = addTagInteractive(i)
			if err != nil {
				return err
			}
		case "delete tags":
			err = deleteTagInteractive(i)
			if err != nil {
				return err
			}
		case "done":
			return i.Write()
		}
	}
}

func addTagInteractive(i *app.Info) error {
	suggestions := myApp.Classify(i.Body.Content)

	for {
		var (
			remainingSuggestions app.ClassificationList
			items                []string
		)

	outer:
		for _, s := range suggestions {
			for _, t := range i.Tags {
				if t == string(s.Class) {
					continue outer
				}
			}
			item := fmt.Sprintf("%s [%.02f%%]", s.Class, s.Score)
			remainingSuggestions = append(remainingSuggestions, s)
			items = append(items, item)
		}

		items = append(items, "[new]", "[done]")

		tagPrompt := promptui.Select{
			Label: "Which tag do you want to add?",
			Items: items,
			Size:  15,
		}

		pick, tag, err := tagPrompt.Run()

		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}

		if tag == "[done]" {
			return nil
		}

		if tag == "[new]" {
			newTagPrompt := promptui.Prompt{
				Label: "Enter a new tag",
			}

			tag, err = newTagPrompt.Run()
			if err == promptui.ErrInterrupt {
				os.Exit(0)
			}

			if tag != "" {
				err = addTags(i, []string{tag})
				if err != nil {
					return err
				}
			}
		} else {
			tag = string(remainingSuggestions[pick].Class)
		}

		err = addTags(i, []string{tag})
		if err != nil {
			return err
		}
	}
}

func deleteTagInteractive(i *app.Info) error {
	for {
		items := i.Tags
		items = append(items, "[done]")

		tagPrompt := promptui.Select{
			Label: "Which tag do you want to delete?",
			Items: items,
		}

		_, action, err := tagPrompt.Run()

		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}

		if action == "[done]" {
			return nil
		}

		i.DeleteTag(action)
	}
}
