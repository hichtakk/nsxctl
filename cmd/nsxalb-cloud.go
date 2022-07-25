package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func NewCmdShowAlbCloud() *cobra.Command {
	aliases := []string{"ac", "cloud"}
	cloudCmd := &cobra.Command{
		Use:     "alb-cloud",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ALB Cloud [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			return LoginALB()
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			albclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			albclient.ShowCloud()
			//albclient.DownloadSeImage()
		},
	}

	return cloudCmd
}

func NewCmdShowAlbVirtualService() *cobra.Command {
	aliases := []string{"alb-vs", "vs"}
	cloudCmd := &cobra.Command{
		Use:     "alb-virtualservice",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ALB VirtualService [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			return LoginALB()
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			albclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			vss := albclient.ShowVirtualService()

			w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
			w.Write([]byte(strings.Join([]string{"ID", "Name", "VIP", "Port", "Cloud", "SEGroup", "Status"}, "\t") + "\n"))
			for _, vs := range vss {
				vs.Print(w)
			}
			w.Flush()
		},
	}

	return cloudCmd
}

func NewCmdShowAlbServiceEngine() *cobra.Command {
	aliases := []string{"alb-se", "se"}
	cloudCmd := &cobra.Command{
		Use:     "alb-serviceengine",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ALB ServiceEngine [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			return LoginALB()
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			albclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			ses := albclient.GetServiceEngine()

			w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
			w.Write([]byte(strings.Join([]string{"ID", "Name", "IP", "Cloud", "SEGroup", "Status"}, "\t") + "\n"))
			for _, se := range ses {
				se.Print(w)
			}
			w.Flush()
		},
	}

	return cloudCmd
}
