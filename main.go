package main

import (
	"os"

	"github.com/sensu-community/sensu-plugin-go/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	command := cmd.NewRootCommand(os.Args)
	if err := command.Execute(); err != nil {
		logrus.WithError(err).Fatal("error executing sensu-plugin-go")
	}

	// // Mark required flags
	// //rootCmd.MarkFlagRequired("name")
	// //rootCmd.MarkFlagRequired("description")
	// //rootCmd.MarkFlagRequired("github-user")

	// // Bind viper to flags
	// viper.BindPFlag("name", rootCmd.Flags().Lookup("name"))
	// viper.BindPFlag("description", rootCmd.Flags().Lookup("description"))
	// viper.BindPFlag("github-user", rootCmd.Flags().Lookup("github-user"))
	// viper.BindPFlag("github-project", rootCmd.Flags().Lookup("github-project"))
	// viper.BindPFlag("copyright-year", rootCmd.Flags().Lookup("copyright-year"))
	// viper.BindPFlag("copyright-holder", rootCmd.Flags().Lookup("copyright-holder"))
}
