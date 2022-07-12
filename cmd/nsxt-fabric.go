package cmd

import (
	"fmt"
	"log"
	"strings"

	"github.com/hichtakk/nsxctl/structs"
	"github.com/spf13/cobra"
)

func NewCmdShowComputeManager() *cobra.Command {
	aliases := []string{"cm"}
	computeManagerCmd := &cobra.Command{
		Use:     "compute-manager",
		Aliases: []string{"cm"},
		Short:   fmt.Sprintf("show compute managers [%s]", strings.Join(aliases, ",")),
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
			cms := nsxtclient.GetComputeManager()
			for _, cm := range *cms {
				cm.Status = nsxtclient.GetComputeManagerStatus(cm.Id)
				cm.Print()
			}
		},
	}

	return computeManagerCmd
}

func NewCmdCreateComputeManager() *cobra.Command {
	var address string
	var user string
	var password string
	var trust bool
	aliases := []string{"cm"}
	computeManagerCmd := &cobra.Command{
		Use:     "compute-manager",
		Aliases: aliases,
		Short:   fmt.Sprintf("create compute managers [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			if len(args) == 0 {
				log.Fatal("compute manager name is required")
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
			thumbprint := nsxtclient.GetTlsFingerprint(address, 443)
			nsxtclient.CreateComputeManager(name, address, thumbprint, user, password, trust)
			fmt.Println(debug)
		},
	}
	computeManagerCmd.Flags().StringVarP(&address, "address", "a", "", "IPv4 address of target compute manager")
	computeManagerCmd.Flags().StringVarP(&user, "user", "u", "", "user of target compute manager")
	computeManagerCmd.Flags().StringVarP(&password, "password", "p", "", "password of target compute manager")
	computeManagerCmd.Flags().BoolVarP(&trust, "enable-trust", "", false, "enable trust [default: false]")
	computeManagerCmd.MarkFlagRequired("address")
	computeManagerCmd.MarkFlagRequired("user")
	computeManagerCmd.MarkFlagRequired("password")

	return computeManagerCmd
}

func NewCmdDeleteComputeManager() *cobra.Command {
	aliases := []string{"cm"}
	computeManagerCmd := &cobra.Command{
		Use:     "compute-manager",
		Aliases: aliases,
		Short:   fmt.Sprintf("delete compute managers [%s]", strings.Join(aliases, ",")),
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
	aliases := []string{"tz"}
	transportZoneCmd := &cobra.Command{
		Use:     "transport-zone",
		Aliases: aliases,
		Short:   fmt.Sprintf("show transport zones [%s]", strings.Join(aliases, ",")),
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
			eps := nsxtclient.GetEnforcementPoint("default")
			var ep structs.EnforcementPoint
			for _, ep = range *eps {
				break
			}
			tzs := nsxtclient.GetPolicyTransportZone(ep.Path)
			for _, tz := range *tzs {
				fmt.Println(tz.Id, tz.Name, tz.Type)
			}
		},
	}

	return transportZoneCmd
}

func NewCmdCreateTransportZone() *cobra.Command {
	var transportType string
	aliases := []string{"tz"}
	transportZoneCmd := &cobra.Command{
		Use:     "transport-zone",
		Aliases: []string{"tz"},
		Short:   fmt.Sprintf("create transport zone [%s]", strings.Join(aliases, ",")),
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
	aliases := []string{"tz"}
	transportZoneCmd := &cobra.Command{
		Use:     "transport-zone",
		Aliases: aliases,
		Short:   fmt.Sprintf("delete transport zone [%s]", strings.Join(aliases, ",")),
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
	aliases := []string{"tn"}
	tpnCmd := &cobra.Command{
		Use:     "transport-node",
		Aliases: aliases,
		Short:   fmt.Sprintf("show transport nodes [%s]", strings.Join(aliases, ",")),
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
	aliases := []string{"tnp"}
	tpnCmd := &cobra.Command{
		Use:     "transport-node-profile",
		Aliases: aliases,
		Short:   fmt.Sprintf("show transport node profiles [%s]", strings.Join(aliases, ",")),
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

func NewCmdCreateEdge() *cobra.Command {
	var template string
	var address string
	var root_password string
	var admin_password string
	aliases := []string{"edge"}
	edgeCmd := &cobra.Command{
		Use:     "edge",
		Aliases: aliases,
		Short:   fmt.Sprintf("create edges [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			if len(args) == 0 {
				log.Fatal("edge name is required")
			}
			nsxtclient.Login(site.GetCredential())
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			name := args[0]
			nsxtclient.CreateEdge(name, template, address, root_password, admin_password)
			fmt.Println(debug)
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
	}
	edgeCmd.Flags().StringVarP(&template, "template", "", "", "template edge name")
	edgeCmd.Flags().StringVarP(&address, "address", "", "", "management I/F address of new edge")
	edgeCmd.Flags().StringVarP(&root_password, "root_password", "", "", "root password of new edge")
	edgeCmd.Flags().StringVarP(&admin_password, "admin_password", "", "", "admin password of new edge")
	edgeCmd.MarkFlagRequired("template")
	edgeCmd.MarkFlagRequired("address")
	edgeCmd.MarkFlagRequired("root_password")
	edgeCmd.MarkFlagRequired("admin_password")

	return edgeCmd
}
