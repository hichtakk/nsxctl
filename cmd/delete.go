package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/hichtakk/nsxctl/client"
	"github.com/spf13/cobra"
)

// NewCmdDelete is subcommand to delete resources.
func NewCmdDelete() *cobra.Command {
	var deleteCmd = &cobra.Command{
		Use:   "delete",
		Short: "delete resources",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			file, _ := ioutil.ReadFile(configfile)
			json.Unmarshal(file, &conf)
			nsxtclient = client.NewNsxtClient(false)
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.BaseUrl = site.Endpoint
			return nil
		},
	}
	deleteCmd.AddCommand(
		NewCmdDeleteComputeManager(),
		NewCmdDeleteTransportZone(),
	)

	return deleteCmd
}
