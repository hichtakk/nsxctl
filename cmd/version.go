package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// NewCmdShow is subcommand to show resources.
func NewCmdVersion() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show version of nsxctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("nsxctl version: %s, revision: %s\n", version, revision)
		},
	}

	return versionCmd
}

// show NSX-T version
func NewCmdShowVersion() *cobra.Command {
	aliases := []string{"v"}
	versionCmd := &cobra.Command{
		Use:     "version",
		Aliases: aliases,
		Short:   fmt.Sprintf("show version of NSX-T/ALB [%s]", strings.Join(aliases, ",")),
		Run: func(cmd *cobra.Command, args []string) {
			if alb != true {
				fmt.Println(nsxtclient.GetVersion())
			} else {
				fmt.Println(albclient.Version)
			}
		},
	}
	versionCmd.PersistentFlags().BoolVarP(&alb, "alb", "", false, "show NSX ALB version")

	return versionCmd
}
