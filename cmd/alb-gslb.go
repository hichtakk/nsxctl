package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

func NewCmdShowAlbGslb() *cobra.Command {
	aliases := []string{"gslb"}
	cloudCmd := &cobra.Command{
		Use:     "alb-gslb",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ALB GSLB [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			return LoginALB()
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			albclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			gslb := albclient.GetGslbSites()
			gslb.Print()
		},
	}

	return cloudCmd
}

func NewCmdShowAlbGslbService() *cobra.Command {
	aliases := []string{}
	cloudCmd := &cobra.Command{
		Use:     "gslb-service",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ALB GSLB Service[%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			return LoginALB()
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			albclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			hm := albclient.GetHealthMonitors()
			gSvc := albclient.GetGslbServices()
			gSvc.Print(hm)
		},
	}

	return cloudCmd
}
