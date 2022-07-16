package cmd

import (
	"fmt"
	"log"

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
		Short:   "show version of NSX-T",
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.Login(site.GetCredential())
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(nsxtclient.GetVersion())
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
	}

	return versionCmd
}
