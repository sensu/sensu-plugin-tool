package newcmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"text/template"
	"time"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/sensu-community/sensu-plugin-tool/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

var (
	DefaultHandlerTemplateURL   = "https://github.com/sensu-community/handler-plugin-template"
	DefaultHandlerCopyrightYear = strconv.Itoa(time.Now().Year())

	handlerTmplFiles = []string{
		"go.mod",
		"LICENSE",
		"main.go",
		"README.md",
	}
)

type handlerProject struct {
	Name            string
	TemplateURL     string
	Description     string
	GithubUser      string
	GithubProject   string
	CopyrightYear   string
	CopyrightHolder string

	logger *logrus.Logger
}

func (p *handlerProject) createProject() error {
	p.logger.Infof("Creating project directory: %s\n", p.GithubProject)
	if err := os.MkdirAll(p.GithubProject, 0755); err != nil {
		return err
	}

	p.logger.Infof("Fetching template: %s\n", p.TemplateURL)
	r, err := git.Clone(memory.NewStorage(), nil, &git.CloneOptions{
		URL: p.TemplateURL,
	})
	if err != nil {
		return err
	}

	ref, err := r.Head()
	if err != nil {
		return err
	}

	commit, err := r.CommitObject(ref.Hash())
	if err != nil {
		return err
	}

	tree, err := commit.Tree()
	if err != nil {
		return err
	}

	if err := tree.Files().ForEach(func(f *object.File) error {
		outputPath := path.Join(p.GithubProject, f.Name)

		if err := os.MkdirAll(path.Dir(outputPath), 0755); err != nil {
			return err
		}

		outputFile, err := os.Create(outputPath)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		contents, err := f.Contents()
		if err != nil {
			return err
		}

		if util.StringInSlice(f.Name, handlerTmplFiles) {
			tmpl := template.Must(template.New(f.Name).Parse(contents))
			if err := tmpl.Execute(outputFile, p); err != nil {
				return err
			}
		} else {
			if _, err = outputFile.WriteString(contents); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func newHandlerCommand(logger *logrus.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "handler",
		Short: "TODO: description here",
		Long:  `TODO: a longer description here`,
		RunE: func(cmd *cobra.Command, args []string) error {
			project := handlerProject{
				logger: logger,
			}

			if cmd.Flags().NFlag() == 0 {
				// interactive mode
				err := survey.Ask(handlerQs, &project)
				if err != nil {
					return err
				}
			} else {
				// flag mode
				project.Name = viper.GetString("name")
				project.TemplateURL = viper.GetString("template-url")
				project.Description = viper.GetString("description")
				project.GithubUser = viper.GetString("github-user")
				project.GithubProject = viper.GetString("github-project")
				project.CopyrightYear = viper.GetString("copyright-year")
				project.CopyrightHolder = viper.GetString("copyright-holder")

				if project.Name == "" {
					return errors.New("name cannot be empty")
				}
				if project.TemplateURL == "" {
					return errors.New("template-url cannot be empty")
				}
				if project.Description == "" {
					return errors.New("description cannot be empty")
				}
				if project.GithubUser == "" {
					return errors.New("github-user cannot be empty")
				}
				if project.GithubProject == "" {
					return errors.New("github-project cannot be empty")
				}
				if project.CopyrightYear == "" {
					return errors.New("copyright-year cannot be empty")
				}
				if project.CopyrightHolder == "" {
					return errors.New("copyright-holder cannot be empty")
				}
			}

			if err := project.createProject(); err != nil {
				return err
			}

			fmt.Println("Success!")

			return nil
		},
	}

	cmd.Flags().String("name", "", "Name of the project (required)")
	cmd.Flags().String("template-url", DefaultHandlerTemplateURL, "URL of the handler template repository to use")
	cmd.Flags().String("description", "", "Description of the project (required)")
	cmd.Flags().String("github-user", "", "Github username that the plugin will belong to (required)")
	cmd.Flags().String("github-project", "", "Github project name that the plugin will belong to (required)")
	cmd.Flags().String("copyright-year", DefaultHandlerCopyrightYear, "The copyright year to be used in the LICENSE file")
	cmd.Flags().String("copyright-holder", "", "The copyright holder to be used in the LICENSE file")

	viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	viper.BindPFlag("template-url", cmd.Flags().Lookup("template-url"))
	viper.BindPFlag("description", cmd.Flags().Lookup("description"))
	viper.BindPFlag("github-user", cmd.Flags().Lookup("github-user"))
	viper.BindPFlag("github-project", cmd.Flags().Lookup("github-project"))
	viper.BindPFlag("copyright-year", cmd.Flags().Lookup("copyright-year"))
	viper.BindPFlag("copyright-holder", cmd.Flags().Lookup("copyright-holder"))

	return cmd
}

var handlerQs = []*survey.Question{
	{
		Name: "name",
		Prompt: &survey.Input{
			Message: "Project name",
		},
		Validate: survey.Required,
	},
	{
		Name: "templateUrl",
		Prompt: &survey.Input{
			Message: "Template URL",
			Default: DefaultHandlerTemplateURL,
		},
		Validate: survey.Required,
	},
	{
		Name: "description",
		Prompt: &survey.Input{
			Message: "Description",
		},
		Validate: survey.Required,
	},
	{
		Name: "githubUser",
		Prompt: &survey.Input{
			Message: "Github User",
		},
		Validate: survey.Required,
	},
	{
		Name: "githubProject",
		Prompt: &survey.Input{
			Message: "Github Project",
		},
		Validate: survey.Required,
	},
	{
		Name: "copyrightYear",
		Prompt: &survey.Input{
			Message: "Copyright Year",
			Default: DefaultHandlerCopyrightYear,
		},
		Validate: survey.Required,
	},
	{
		Name: "copyrightHolder",
		Prompt: &survey.Input{
			Message: "Copyright Holder",
		},
		Validate: survey.Required,
	},
}
