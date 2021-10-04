package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func NewCmdShowTransportNode() *cobra.Command {
	tpnCmd := &cobra.Command{
		Use:     "transport-node",
		Aliases: []string{"tn"},
		Short:   "show transport nodes",
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
			nsxtclient.GetTransportNode()
		},
	}

	return tpnCmd
}

func NewCmdShowTransportNodeProfile() *cobra.Command {
	tpnCmd := &cobra.Command{
		Use:     "transport-node-profile",
		Aliases: []string{"tnp"},
		Short:   "show transport node profiles",
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
			nsxtclient.GetTransportNodeProfile()
		},
	}

	return tpnCmd
}
