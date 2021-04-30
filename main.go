package main

import (
	"os"

	"github.com/sensu/sensu-plugin-tool/cmd"
	"github.com/sirupsen/logrus"
)

func main() {
	command := cmd.NewRootCommand(os.Args)
	if err := command.Execute(); err != nil {
		logrus.WithError(err).Fatal("error executing sensu-plugin-tool")
	}
}
