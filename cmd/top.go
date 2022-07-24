package cmd

import (
	"github.com/spf13/cobra"
)

// NewCmdShow is subcommand to show resources.
func NewCmdTop() *cobra.Command {
	var topCmd = &cobra.Command{
		Use:   "top",
		Short: "monitor resources",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			return Login()
		},
	}
	topCmd.AddCommand(
		NewCmdTopGateway(),
	)

	return topCmd
}
