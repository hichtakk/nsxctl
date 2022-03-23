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
func NewCmdShow() *cobra.Command {
	var showCmd = &cobra.Command{
		Use:   "show",
		Short: "Show resources",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			file, _ := ioutil.ReadFile(configfile)
			json.Unmarshal(file, &conf)
			nsxtclient = client.NewNsxtClient(false, debug)
			var site config.NsxTSite
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.BaseUrl = site.Endpoint
			return nil
		},
	}
	showCmd.AddCommand(
		NewCmdShowGateway(),
		NewCmdShowComputeManager(),
		NewCmdShowTransportZone(),
		NewCmdShowTransportNode(),
		NewCmdShowTransportNodeProfile(),
		NewCmdShowIpPool(),
		NewCmdShowIpBlock(),
		NewCmdShowSegment(),
		NewCmdShowAlbCloud(),
		NewCmdShowAlbVirtualService(),
	)

	return showCmd
}
