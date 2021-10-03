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
		Run: showGateway,
	}

	return gatewayCmd
}

func showGateway(cmd *cobra.Command, args []string) {
	nsxtclient.GetT0()
}

func NewCmdShowComputeManager() *cobra.Command {
	computeManagerCmd := &cobra.Command{
		Use:     "compute-manager",
		Aliases: []string{"cm"},
		Short:   "show compute managers",
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
			nsxtclient.GetComputeManager()
		},
	}

	return computeManagerCmd
}
