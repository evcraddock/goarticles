package imports

import (
	"io"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/go-playground/validator.v9"

	"github.com/evcraddock/goarticles/internal/cli"
	"github.com/evcraddock/goarticles/internal/configs"
)

var structValidator = validator.New()

type ImportOptions struct {
	Out            io.Writer
	ConfigFile     configs.ClientConfiguration
	FilesToProcess string
}

func newImportOptions(out io.Writer) *ImportOptions {
	return &ImportOptions{
		Out: out,
	}
}

func NewCmdImport(out io.Writer) *cobra.Command {
	o := newImportOptions(out)
	cmd := &cobra.Command{
		Use: "import",
		Run: func(cmd *cobra.Command, args []string) {
			o.prepare(cmd)
			o.validate()
			o.importArticle()
		},
	}

	cmd.Flags().String("configfile", "", "yaml configuration file (optional)")
	cmd.Flags().String("files", "", "files or folders to process")
	return cmd
}

func (o *ImportOptions) prepare(cmd *cobra.Command) {
	var config *configs.ClientConfiguration

	if filesToProcess, err := cmd.Flags().GetString("filesToProcess"); err == nil {
		o.FilesToProcess = filesToProcess
	}

	if configFile, err := cmd.Flags().GetString("configfile"); err == nil {
		config, err = configs.LoadCliConfig(configFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		o.ConfigFile = *config
	}
}

func (o *ImportOptions) importArticle() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(o.Out)
	log.SetLevel(log.InfoLevel)

	articleService := cli.NewArticleImporter(o.ConfigFile)
	articleService.CreateOrUpdateArticle(o.FilesToProcess)
}

func (o *ImportOptions) validate() {
	if err := structValidator.Struct(o.ConfigFile); err != nil {
		msg := "the configuration file is malformed"
		log.Fatalf(msg)
	}
}
