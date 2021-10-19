package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/hichtakk/nsxctl/client"
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
			nsxtclient = client.NewNsxtClient(false, debug)
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.BaseUrl = site.Endpoint
			return nil
		},
	}
	showCmd.AddCommand(
		NewCmdCreateComputeManager(),
		NewCmdCreateTransportZone(),
		NewCmdCreateIpPool(),
		NewCmdCreateIpBlock(),
	)

	return showCmd
}
