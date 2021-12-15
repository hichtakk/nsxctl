package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCmdShow is subcommand to show resources.
func NewCmdVersion() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show version of nsxctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("nsxctl version: %s, revision: %s", version, revision)
		},
	}

	return versionCmd
}
