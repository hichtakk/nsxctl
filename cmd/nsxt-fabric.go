package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"
)

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
			fmt.Println(debug)
		},
	}

	return computeManagerCmd
}

func NewCmdCreateComputeManager() *cobra.Command {
	var address string
	var thumbprint string
	var user string
	var password string
	computeManagerCmd := &cobra.Command{
		Use:     "compute-manager",
		Aliases: []string{"cm"},
		Short:   "create compute managers",
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
			nsxtclient.CreateComputeManager(address, thumbprint, user, password)
			fmt.Println(debug)
		},
	}
	computeManagerCmd.Flags().StringVarP(&address, "address", "a", "", "IPv4 address of target compute manager")
	computeManagerCmd.Flags().StringVarP(&thumbprint, "thumbprint", "t", "", "thumbprint of target compute manager")
	computeManagerCmd.Flags().StringVarP(&user, "user", "u", "", "user of target compute manager")
	computeManagerCmd.Flags().StringVarP(&password, "password", "p", "", "password of target compute manager")
	computeManagerCmd.MarkFlagRequired("address")
	computeManagerCmd.MarkFlagRequired("thumbprint")
	computeManagerCmd.MarkFlagRequired("user")
	computeManagerCmd.MarkFlagRequired("password")

	return computeManagerCmd
}

func NewCmdDeleteComputeManager() *cobra.Command {
	computeManagerCmd := &cobra.Command{
		Use:     "compute-manager",
		Aliases: []string{"cm"},
		Short:   "delete compute managers",
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
			cmId := args[0]
			nsxtclient.DeleteComputeManager(cmId)
			fmt.Println(debug)
		},
	}

	return computeManagerCmd
}

func NewCmdShowTransportZone() *cobra.Command {
	transportZoneCmd := &cobra.Command{
		Use:     "transport-zone",
		Aliases: []string{"tz"},
		Short:   "show transport zones (tz",
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
			nsxtclient.GetTransportZone()
		},
	}

	return transportZoneCmd
}

func NewCmdCreateTransportZone() *cobra.Command {
	var transportType string
	transportZoneCmd := &cobra.Command{
		Use:     "transport-zone",
		Aliases: []string{"tz"},
		Short:   "create transport zone",
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
			nsxtclient.CreateTransportZone(name, transportType)
		},
	}
	transportZoneCmd.Flags().StringVarP(&transportType, "type", "t", "", "transport zone type [vlan, overlay]")
	transportZoneCmd.MarkFlagRequired("transportType")

	return transportZoneCmd
}

func NewCmdDeleteTransportZone() *cobra.Command {
	transportZoneCmd := &cobra.Command{
		Use:     "transport-zone",
		Aliases: []string{"tz"},
		Short:   "delete transport zone",
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
			tzId := args[0]
			nsxtclient.DeleteTransportZone(tzId)
			fmt.Println(debug)
		},
	}

	return transportZoneCmd
}

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
