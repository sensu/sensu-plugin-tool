package newcmd

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func NewHelpCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Create a new plugin from a template",
	}
	cmd.AddCommand(newHandlerCommand(logger))
	cmd.AddCommand(newCheckCommand(logger))
	cmd.AddCommand(newMutatorCommand(logger))
	cmd.AddCommand(newSensuctlCommand(logger))
	cmd.AddCommand(newAWSCheckCommand(logger))
	return cmd
}
