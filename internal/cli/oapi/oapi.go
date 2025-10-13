package oapi

import (
	"github.com/maketaio/api/internal/generator"
	"github.com/spf13/cobra"
)

func NewCmd() *cobra.Command {
	cfg := &generator.Config{}

	cmd := &cobra.Command{
		Use:   "oapi",
		Short: "Generate models and route handlers from OpenAPI spec",
		RunE: func(cmd *cobra.Command, args []string) error {
			return generator.Generate(cfg)
		},
	}

	cmd.Flags().StringVar(&cfg.In, "in", "", "Path to OpenAPI spec (YAML/JSON)")
	cmd.Flags().StringVar(&cfg.Out, "out", "", "Output directory for generated code")
	cmd.Flags().StringVar(&cfg.Pkg, "package", "", "Go package name for generated code")
	cmd.MarkFlagRequired("in")
	cmd.MarkFlagRequired("out")

	return cmd
}
