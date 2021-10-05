package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func NewCmdShowGateway() *cobra.Command {
	gatewayCmd := &cobra.Command{
		Use:     "gateway",
		Aliases: []string{"gw"},
		Short:   "show logical gateways",
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.Login(site.GetCredential())
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			nsxtclient.GetT0()
		},
	}

	return gatewayCmd
}
