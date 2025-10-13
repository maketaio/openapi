package cli

import (
	"github.com/maketaio/api/internal/cli/oapi"
	"github.com/spf13/cobra"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "codegen",
		Short: "Codegen utilities",
	}

	cmd.AddCommand(oapi.NewCmd())

	return cmd
}
