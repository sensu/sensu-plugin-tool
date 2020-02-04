package newcmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewHelpCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "TODO: description here",
	}
	cmd.AddCommand(newHandlerCommand(logger))
	return cmd
}
