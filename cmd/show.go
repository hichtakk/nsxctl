package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// NewCmdShow is subcommand to show resources.
func NewCmdShow() *cobra.Command {
	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show resources",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			if alb {
				// NSX ALB
				if err := LoginALB(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			} else {
				// NSX
				if err := Login(); err != nil {
					fmt.Println(err)
					os.Exit(1)
				}
			}
			return nil
		},
		PersistentPostRunE: func(c *cobra.Command, args []string) error {
			var err error
			if alb {
				// NSX ALB
				err = albAgent.Logout()
			} else {
				// NSX
				err = nsxAgent.Logout()
			}
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			return nil
		},
	}
	showCmd.AddCommand(
		NewCmdShowVersion(),
		NewCmdShowCertificate(),
		// NewCmdShowGateway(),
		// NewCmdShowComputeManager(),
		// NewCmdShowTransportZone(),
		// NewCmdShowTransportNode(),
		// NewCmdShowTransportNodeProfile(),
		// NewCmdShowEnforcementPoint(),
		// NewCmdShowIpPool(),
		// NewCmdShowIpBlock(),
		// NewCmdShowEdge(),
		// NewCmdShowEdgeCluster(),
		// NewCmdShowSegment(),
		// NewCmdShowRoutingTable(),
		// NewCmdShowBgpAdvRoutes(),
		// NewCmdShowDfwPolicies(),
		// NewCmdShowDfwRules(),
		// NewCmdShowAlbCloud(),
		// NewCmdShowAlbVirtualService(),
		// NewCmdShowAlbServiceEngine(),
		// NewCmdShowAlbGslb(),
		// NewCmdShowAlbPool(),
		// NewCmdShowAlbGslbService(),
	)

	return showCmd
}
