package cmd

import (
	"github.com/jovandeginste/docmgmt/app"
	"github.com/spf13/cobra"
)

var (
	cfgFile string
	debug   bool
	myApp   app.App
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "docmgmt",
	Short: "Personal document management",
	Long:  `Manage your personal documents`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		panic(err)
	}
}

func init() {
	cobra.OnInitialize(initApp)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	RootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.docmgmt.yaml)")
	RootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "enable debug logging")
}

// initApp initializes the application
func initApp() {
	err := myApp.LoadConfiguration(cfgFile, debug)
	if err != nil {
		panic(err)
	}
}
