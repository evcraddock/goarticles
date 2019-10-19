package link

import (
	"io"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/go-playground/validator.v9"

	"github.com/evcraddock/goarticles/internal/cli"
	"github.com/evcraddock/goarticles/internal/configs"
	"github.com/evcraddock/goarticles/pkg/links"
)

var structValidator = validator.New()

type LinkOptions struct {
	Out        io.Writer
	ConfigFile configs.ClientConfiguration
	Title      string
	URL        string
	Categories []string
	Tags       []string
}

func newLinkOptions(out io.Writer) *LinkOptions {
	return &LinkOptions{
		Out: out,
	}
}

func NewCmdLink(out io.Writer) *cobra.Command {
	o := newLinkOptions(out)
	cmd := &cobra.Command{
		Use: "link",
		Run: func(cmd *cobra.Command, args []string) {
			o.prepare(cmd)
			o.validate()
			o.addLink()
		},
	}

	cmd.Flags().String("configfile", "", "yaml configuration file (optional)")
	cmd.Flags().String("title", "", "link title")
	cmd.Flags().String("URL", "", "link url")
	cmd.Flags().String("categories", "", "link categories")
	cmd.Flags().String("tags", "", "link tags")
	return cmd
}

func (o *LinkOptions) prepare(cmd *cobra.Command) {
	var config *configs.ClientConfiguration

	if configFile, err := cmd.Flags().GetString("configfile"); err == nil {
		config, err = configs.LoadCliConfig(configFile)
		if err != nil {
			log.Fatal(err.Error())
		}

		o.ConfigFile = *config
	}

	if title, err := cmd.Flags().GetString("title"); err == nil {
		o.Title = title
	}

	if url, err := cmd.Flags().GetString("URL"); err == nil {
		o.URL = url
	}

	if categories, err := cmd.Flags().GetString("categories"); err == nil && categories != "" {
		o.Categories = strings.Split(strings.Trim(categories, " "), ",")
	}

	if tags, err := cmd.Flags().GetString("tags"); err == nil && tags != "" {
		o.Tags = strings.Split(strings.Trim(tags, " "), ",")
	}
}

func (o *LinkOptions) addLink() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(o.Out)
	log.SetLevel(log.InfoLevel)

	newLink := links.Link{
		Title:      o.Title,
		URL:        o.URL,
		Categories: o.Categories,
		Tags:       o.Tags,
	}

	linkService := cli.NewLinkManager(o.ConfigFile)
	linkService.CreateLink(newLink)

}

func (o *LinkOptions) validate() {
	if err := structValidator.Struct(o.ConfigFile); err != nil {
		msg := "the configuration file is malformed"
		log.Fatalf(msg)
	}
}
