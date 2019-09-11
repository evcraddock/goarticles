package cmd

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/evcraddock/goarticles/pkg/cmd/import"
	"github.com/evcraddock/goarticles/pkg/cmd/version"
)

func NewDefaultCommand() *cobra.Command {
	return NewGorticleCommand(os.Stdin, os.Stdout, os.Stderr)
}

func NewGorticleCommand(in io.Reader, out, err io.Writer) *cobra.Command {
	cmds := &cobra.Command{
		Use:   "goarticle",
		Short: "goarticle manages markdown content",
		Long:  "goarticle manages markdown content",
		Run:   runHelp,
	}

	cmds.AddCommand(version.NewCmdVersion(out))
	cmds.AddCommand(imports.NewCmdImport(out))

	return cmds
}

func runHelp(cmd *cobra.Command, args []string) {
	cmd.Help()
}
