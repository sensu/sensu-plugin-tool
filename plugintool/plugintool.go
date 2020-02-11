package plugintool

import (
	"errors"
	"html/template"
	"os"
	"path"
	"strconv"
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

type Project struct {
	TemplateURL    string
	TemplateFiles  []string
	TemplateFields ProjectTemplateFields

	logger *logrus.Logger
}

type ProjectTemplateFields struct {
	Name            string
	Description     string
	GithubUser      string
	GithubProject   string
	CopyrightYear   string
	CopyrightHolder string
}

func NewProject(logger *logrus.Logger, tmplFiles []string) Project {
	return Project{
		TemplateURL:   viper.GetString("template-url"),
		TemplateFiles: tmplFiles,
		TemplateFields: ProjectTemplateFields{
			Name:            viper.GetString("name"),
			Description:     viper.GetString("description"),
			GithubUser:      viper.GetString("github-user"),
			GithubProject:   viper.GetString("github-project"),
			CopyrightYear:   viper.GetString("copyright-year"),
			CopyrightHolder: viper.GetString("copyright-holder"),
		},
		logger: logger,
	}
}

func AddCommonFlags(cmd *cobra.Command, defaultTemplateURL string) {
	cmd.Flags().String("name", "", "Name of the project (required)")
	cmd.Flags().String("template-url", defaultTemplateURL, "URL of the handler template repository to use")
	cmd.Flags().String("description", "", "Description of the project (required)")
	cmd.Flags().String("github-user", "", "Github username that the plugin will belong to (required)")
	cmd.Flags().String("github-project", "", "Github project name that the plugin will belong to (required)")
	cmd.Flags().String("copyright-year", strconv.Itoa(time.Now().Year()), "The copyright year to be used in the LICENSE file")
	cmd.Flags().String("copyright-holder", "", "The copyright holder to be used in the LICENSE file")
}

func AddViperBindings(cmd *cobra.Command) {
	viper.BindPFlag("name", cmd.Flags().Lookup("name"))
	viper.BindPFlag("template-url", cmd.Flags().Lookup("template-url"))
	viper.BindPFlag("description", cmd.Flags().Lookup("description"))
	viper.BindPFlag("github-user", cmd.Flags().Lookup("github-user"))
	viper.BindPFlag("github-project", cmd.Flags().Lookup("github-project"))
	viper.BindPFlag("copyright-year", cmd.Flags().Lookup("copyright-year"))
	viper.BindPFlag("copyright-holder", cmd.Flags().Lookup("copyright-holder"))
}

func TemplateURLQuestion(defaultTemplateURL string) []*survey.Question {
	return []*survey.Question{
		{
			Name: "templateUrl",
			Prompt: &survey.Input{
				Message: "Template URL",
				Default: defaultTemplateURL,
			},
			Validate: survey.Required,
		},
	}
}

func CommonQuestions() []*survey.Question {
	return []*survey.Question{
		{
			Name: "name",
			Prompt: &survey.Input{
				Message: "Project name",
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
				Default: defaultCopyrightYear(),
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
}

func (p *Project) Validate() error {
	if p.TemplateURL == "" {
		return errors.New("template-url cannot be empty")
	}
	if p.TemplateFields.Name == "" {
		return errors.New("name cannot be empty")
	}
	if p.TemplateFields.Description == "" {
		return errors.New("description cannot be empty")
	}
	if p.TemplateFields.GithubUser == "" {
		return errors.New("github-user cannot be empty")
	}
	if p.TemplateFields.GithubProject == "" {
		return errors.New("github-project cannot be empty")
	}
	if p.TemplateFields.CopyrightYear == "" {
		return errors.New("copyright-year cannot be empty")
	}
	if p.TemplateFields.CopyrightHolder == "" {
		return errors.New("copyright-holder cannot be empty")
	}
	return nil
}

func (p *Project) Create() error {
	p.logger.Infof("Creating project directory: %s\n", p.TemplateFields.GithubProject)
	if err := os.MkdirAll(p.TemplateFields.GithubProject, 0755); err != nil {
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
		outputPath := path.Join(p.TemplateFields.GithubProject, f.Name)

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

		if util.StringInSlice(f.Name, p.TemplateFiles) {
			tmpl := template.Must(template.New(f.Name).Parse(contents))
			if err := tmpl.Execute(outputFile, p.TemplateFields); err != nil {
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

func defaultCopyrightYear() string {
	return strconv.Itoa(time.Now().Year())
}
