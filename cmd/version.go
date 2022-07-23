package cmd

import (
	"fmt"
	"log"
	"strings"

	ac "github.com/hichtakk/nsxctl/nsxalb"
	"github.com/spf13/cobra"
)

// NewCmdShow is subcommand to show resources.
func NewCmdVersion() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show version of nsxctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("nsxctl version: %s, revision: %s\n", version, revision)
		},
	}

	return versionCmd
}

// show NSX-T version
func NewCmdShowVersion() *cobra.Command {
	aliases := []string{"v"}
	versionCmd := &cobra.Command{
		Use:     "version",
		Aliases: aliases,
		Short:   fmt.Sprintf("show version of NSX-T/ALB [%s]", strings.Join(aliases, ",")),
		PreRunE: func(c *cobra.Command, args []string) error {
			if alb != true {
				site, err := conf.NsxT.GetCurrentSite()
				if err != nil {
					log.Fatal(err)
				}
				nsxtclient.Login(site.GetCredential())
			} else {
				albclient = ac.NewNsxAlbClient(false, debug)
				albsite, err := conf.NsxAlb.GetCurrentSite()
				if err != nil {
					log.Fatal(err)
				}
				albclient.BaseUrl = albsite.Endpoint
				albclient.Login(albsite.GetCredential())
			}
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			if alb != true {
				fmt.Println(nsxtclient.GetVersion())
			} else {
				fmt.Println(albclient.Version)
			}
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			if alb != true {
				nsxtclient.Logout()
			} else {
				albclient.Logout()
			}
			return nil
		},
	}
	versionCmd.PersistentFlags().BoolVarP(&alb, "alb", "", false, "show NSX ALB version")

	return versionCmd
}
