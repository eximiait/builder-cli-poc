/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var Verbose bool
var Debug bool
var GitlabHostTarget string
var GitlabHostOrigin string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "builder-cli",
	Short: "CLI para poder crear y gestionar proyectos GitOps",
	Long:  `CLI para poder crear y gestionar proyectos GitOps`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVarP(&Verbose, "verbose", "v", false, "Display more verbose output in console output. (default: false)")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))

	rootCmd.PersistentFlags().BoolVarP(&Debug, "debug", "d", false, "Display debugging output in the console. (default: false)")
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().StringVarP(&GitlabHostOrigin, "gitlab-host-origin", "o", "https://gitlab.com", "The GitLab URL for origin, it's used for cloning starters")
	viper.BindPFlag("gitlabHostOrigin", rootCmd.PersistentFlags().Lookup("gitlab-host-origin"))

	rootCmd.PersistentFlags().StringVarP(&GitlabHostTarget, "gitlab-host", "g", "https://gitlab.com", "The GitLab URL where the project will be created")
	viper.BindPFlag("gitlabHost", rootCmd.PersistentFlags().Lookup("gitlab-host"))
}
