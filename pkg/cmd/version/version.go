package version

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

type Version string

func NewCmdVersion(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "version",
		Run: func(cmd *cobra.Command, args []string) {
			showVersion(out)
		},
	}

	return cmd
}

func showVersion(out io.Writer) error {
	// TODO: Go get version
	fmt.Fprintf(out, "%s\n", "0.0.1")

	return nil
}
