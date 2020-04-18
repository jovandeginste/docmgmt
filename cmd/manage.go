package cmd

import (
	"fmt"
	"os"

	"github.com/jovandeginste/docmgmt/app"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var (
	manageCmd = &cobra.Command{
		Use:  "manage",
		Long: `Interactively manage a known in the document root`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			manage(args[0])
		},
	}
	promptMore = promptui.Select{
		Label: "What would you like to do?",
		Items: []string{
			"read the parsed text",
			"open the file with the default program",
			"manage tags",
			"done",
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

	tagPrompt = promptui.Select{
		Label: "Which tag do you want to add?",
		Size:  15,
	}
)

func init() {
	RootCmd.AddCommand(manageCmd)
}

func manage(file string) {
	info, err := myApp.ReadFileInfo(file)
	if err != nil {
		panic(err)
	}

	err = manageInteractive(info)
	if err != nil {
		panic(err)
	}
}

func manageInteractive(i *app.Info) error {
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
		case "done":
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
	for {
		var items []string

		suggestions := i.Suggestions()
		for _, s := range suggestions {
			items = append(items, fmt.Sprintf("%s [%.02f%%]", s.Class, s.Score))
		}

		items = append(items, "[new]", "[done]")
		tagPrompt.Items = items
		pick, tag, err := tagPrompt.Run()

		if err == promptui.ErrInterrupt {
			os.Exit(0)
		}

		switch tag {
		case "[done]":
			return nil
		case "[new]":
			tag = askNewTag()
		default:
			tag = string(suggestions[pick].Class)
		}

		if tag != "" {
			err = addTags(i, []string{tag})
			if err != nil {
				return err
			}
		}
	}
}

func askNewTag() string {
	newTagPrompt := promptui.Prompt{
		Label: "Enter a new tag",
	}

	tag, err := newTagPrompt.Run()
	if err == promptui.ErrInterrupt {
		os.Exit(0)
	}

	return tag
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
