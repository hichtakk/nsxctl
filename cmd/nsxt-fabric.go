package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"

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
			sites := nsxtclient.GetSite()
			eps := nsxtclient.GetEnforcementPoint(sites[0])
			ep := *eps
			// use default for site and enforcementpoint
			nodes := nsxtclient.GetTransportNode(sites[0], ep[0].Id)
			nodes.Print()
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

func NewCmdShowEdge() *cobra.Command {
	var verbose bool
	aliases := []string{"e"}
	edgeCmd := &cobra.Command{
		Use:     "edge",
		Aliases: aliases,
		Short:   fmt.Sprintf("show edges [%s]", strings.Join(aliases, ",")),
		PreRunE: func(c *cobra.Command, args []string) error {
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.Login(site.GetCredential())
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
			if verbose {
				w.Write([]byte(strings.Join([]string{"Id", "Name", "IP", "EdgeCluster", "Status", "Gatways"}, "\t") + "\n"))
			} else {
				w.Write([]byte(strings.Join([]string{"Id", "Name", "IP", "EdgeCluster", "Status"}, "\t") + "\n"))
			}

			edges := nsxtclient.GetEdge()
			ecs := nsxtclient.GetEdgeCluster()

			edge_gw_map := make(map[string][]string) // edge_id : [gwid, gwid, ...]
			if verbose {
				t0s := nsxtclient.GetTier0Gateway("")
				t1s := nsxtclient.GetTier1Gateway("")
				for _, gw := range t0s {
					per_node_status := nsxtclient.GetGatewayAggregateInfo(gw.RealizationId)
					for _, st := range per_node_status {
						eid := st["transport_node_id"]
						ha := st["high_availability_status"]
						val, ok := edge_gw_map[eid]
						if !ok {
							val = []string{}
						}
						val = append(val, gw.Name+"("+ha+")")
						edge_gw_map[eid] = val
					}
				}
				for _, gw := range t1s {
					per_node_status := nsxtclient.GetGatewayAggregateInfo(gw.RealizationId)
					for _, st := range per_node_status {
						eid := st["transport_node_id"]
						ha := st["high_availability_status"]
						val, ok := edge_gw_map[eid]
						if !ok {
							val = []string{}
						}
						val = append(val, gw.Name+"("+ha+")")
						edge_gw_map[eid] = val
					}
				}
			}

			for _, e := range edges {
				var edgeCluster structs.EdgeCluster
				for _, ec := range *ecs {
					for _, ecm := range ec.Members {
						if ecm.Id == e.Id {
							edgeCluster = ec
						}
					}
				}
				ip := strings.Join(e.EdgeNodeDeploymentInfo.IPAddress, ",")
				status := nsxtclient.GetTransportNodeStatus(e.Id)
				if verbose {
					gws := strings.Join(edge_gw_map[e.Id], ",")
					w.Write([]byte(strings.Join([]string{e.Id, e.Name, ip, edgeCluster.Name, status, gws}, "\t") + "\n"))
				} else {
					w.Write([]byte(strings.Join([]string{e.Id, e.Name, ip, edgeCluster.Name, status}, "\t") + "\n"))
				}
			}
			w.Flush()
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
	}
	edgeCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "display gateway placement")

	return edgeCmd
}
