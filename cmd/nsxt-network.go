package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

func NewCmdShowIpPool() *cobra.Command {
	ipPoolCmd := &cobra.Command{
		Use:     "ip-pool",
		Aliases: []string{""},
		Short:   "show ip address pools",
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
			nsxtclient.GetIpPool()
		},
	}

	return ipPoolCmd
}

func NewCmdCreateIpPool() *cobra.Command {
	//var transportType string
	ipPoolCmd := &cobra.Command{
		Use:     "ip-pool",
		Aliases: []string{""},
		Short:   "create a new ip pool",
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
			name := args[0]
			nsxtclient.CreateIpPool(name)
		},
	}
	//ipPoolCmd.Flags().StringVarP(&transportType, "type", "t", "", "transport zone type [vlan, overlay]")
	//ipPoolCmd.MarkFlagRequired("transportType")

	return ipPoolCmd
}

func NewCmdDeleteIpPool() *cobra.Command {
	ipPoolCmd := &cobra.Command{
		Use:     "ip-pool",
		Aliases: []string{""},
		Short:   "create a new ip pool",
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
			name := args[0]
			nsxtclient.DeleteIpPool(name)
		},
	}

	return ipPoolCmd
}
