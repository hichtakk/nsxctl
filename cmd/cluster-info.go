package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/hichtakk/nsxctl/client"
	"github.com/spf13/cobra"
)

// NewCmdClusterInfo
func NewCmdClusterInfo() *cobra.Command {
	var clusterCmd = &cobra.Command{
		Use:   "cluster-info",
		Short: "show NSX-T cluster information",
		PreRunE: func(c *cobra.Command, args []string) error {
			file, _ := ioutil.ReadFile(configfile)
			json.Unmarshal(file, &conf)
			nsxtclient = client.NewNsxtClient(false, debug)
			site, err := conf.NsxT.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			nsxtclient.BaseUrl = site.Endpoint
			nsxtclient.Login(site.GetCredential())
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			nsxtclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			nsxtclient.ClusterInfo()
		},
	}

	return clusterCmd
}
