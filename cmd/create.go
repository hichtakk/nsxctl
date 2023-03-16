package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/hichtakk/nsxctl/client"
	"github.com/hichtakk/nsxctl/config"
	"github.com/spf13/cobra"
)

// NewCmdShow is subcommand to show resources.
func NewCmdCreate() *cobra.Command {
	var showCmd = &cobra.Command{
		Use:   "create",
		Short: "create resources",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			file, _ := ioutil.ReadFile(configfile)
			json.Unmarshal(file, &conf)
			var site config.NsxTSite
			var err error
			if useSite != "" {
				site, err = conf.NsxT.GetSite(useSite)
			} else {
				site, err = conf.NsxT.GetCurrentSite()
			}
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient = client.NewNsxtClient(false, debug, site.Proxy)
			nsxtclient.BaseUrl = site.Endpoint
			return nil
		},
	}
	showCmd.AddCommand(
		NewCmdCreateComputeManager(),
		NewCmdCreateTransportZone(),
		NewCmdCreateIpPool(),
		NewCmdCreateIpBlock(),
		NewCmdCreateEdge(),
		NewCmdCreateSegment(),
	)

	return showCmd
}
