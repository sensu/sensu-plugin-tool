package newcmd

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/sensu/sensu-plugin-tool/plugintool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	DefaultCheckTemplateURL = "https://github.com/sensu/check-plugin-template"

	checkTmplFiles = []string{
		"go.mod",
		"LICENSE",
		"main.go",
		"README.md",
		".gitignore",
	}
)

func newCheckCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Create a new check plugin from a template",
		PreRun: func(cmd *cobra.Command, args []string) {
			plugintool.AddViperBindings(cmd)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			project := plugintool.NewProject(logger, checkTmplFiles)

			// interactive mode
			if cmd.Flags().NFlag() == 0 {
				// ask for the template url
				if err := survey.Ask(plugintool.TemplateURLQuestion(DefaultCheckTemplateURL), &project); err != nil {
					return err
				}
				// ask remaning common questions
				if err := survey.Ask(plugintool.CommonQuestions(), &project.TemplateFields); err != nil {
					return err
				}
			}

			if err := project.Validate(); err != nil {
				return err
			}

			if err := project.Create(); err != nil {
				return err
			}

			fmt.Println("Success!")
			return nil
		},
	}

	plugintool.AddCommonFlags(cmd, DefaultCheckTemplateURL)
	return cmd
}
