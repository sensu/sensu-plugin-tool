package newcmd

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/sensu-community/sensu-plugin-tool/plugintool"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	DefaultMutatorTemplateURL = "https://github.com/sensu-community/mutator-plugin-template"

	mutatorTmplFiles = []string{
		"go.mod",
		"LICENSE",
		"main.go",
		"README.md",
	}
)

func newMutatorCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mutator",
		Short: "Create a new mutator plugin from a template",
		RunE: func(cmd *cobra.Command, args []string) error {
			project := plugintool.NewProject(logger, mutatorTmplFiles)

			if cmd.Flags().NFlag() == 0 {
				// interactive mode
				// ask for the template url
				if err := survey.Ask(plugintool.TemplateURLQuestion(DefaultMutatorTemplateURL), &project); err != nil {
					return err
				}
				// ask remaning common questions
				if err := survey.Ask(plugintool.CommonQuestions(), &project.TemplateFields); err != nil {
					return err
				}
			} else {
				// flag mode
				project.Validate()
			}

			if err := project.Create(); err != nil {
				return err
			}

			fmt.Println("Success!")
			return nil
		},
	}

	plugintool.AddCommonFlags(cmd, DefaultMutatorTemplateURL)
	return cmd
}
