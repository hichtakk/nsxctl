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

func NewCmdShowIpBlock() *cobra.Command {
	ipPoolCmd := &cobra.Command{
		Use:     "ip-block",
		Aliases: []string{""},
		Short:   "show ip address blocks",
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
			nsxtclient.GetIpBlock()
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

func NewCmdCreateIpBlock() *cobra.Command {
	var cidr string
	ipBlockCmd := &cobra.Command{
		Use:     "ip-block",
		Aliases: []string{""},
		Short:   "create a new ip block",
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
			nsxtclient.CreateIpBlock(name, cidr)
		},
	}
	ipBlockCmd.Flags().StringVarP(&cidr, "cidr", "", "", "CIDR block (10.0.0.0/16)")
	ipBlockCmd.MarkFlagRequired("cidr")

	return ipBlockCmd
}

func NewCmdDeleteIpPool() *cobra.Command {
	ipPoolCmd := &cobra.Command{
		Use:     "ip-pool",
		Aliases: []string{""},
		Short:   "delete ip pool",
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

func NewCmdDeleteIpBlock() *cobra.Command {
	ipPoolCmd := &cobra.Command{
		Use:     "ip-block",
		Aliases: []string{""},
		Short:   "delete ip block",
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
			nsxtclient.DeleteIpBlock(name)
		},
	}

	return ipPoolCmd
}
