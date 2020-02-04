package cmd

import (
	"github.com/sensu-community/sensu-plugin-go/cmd/newcmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	flagLogLevel = "log-level"
)

func NewRootCommand(args []string) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "sensu-plugin-go",
		Short: "Sensu Plugin Go Tool",
	}

	cmd.PersistentFlags().String(flagLogLevel, "warning", "Set the log level (trace, debug, info, warning, error, fatal, panic")
	viper.BindPFlag(flagLogLevel, cmd.PersistentFlags().Lookup(flagLogLevel))

	logger := logrus.New()
	level, err := logrus.ParseLevel(viper.GetString(flagLogLevel))
	if err != nil {
		logger.Panic("invalid log level specified")
	}
	logger.SetLevel(level)

	cmd.AddCommand(newcmd.NewHelpCommand(logger))

	return cmd
}
