package cmd

import (
	"github.com/spf13/cobra"
)

// NewCmdClusterInfo
func NewCmdClusterInfo() *cobra.Command {
	var clusterCmd = &cobra.Command{
		Use:   "cluster-info",
		Short: "show NSX-T cluster information",
		PreRunE: func(c *cobra.Command, args []string) error {
			// file, _ := ioutil.ReadFile(configfile)
			// json.Unmarshal(file, &conf)
			// var site config.NsxTSite
			// var err error
			// if useSite != "" {
			// 	site, err = conf.NsxT.GetSite(useSite)
			// } else {
			// 	site, err = conf.NsxT.GetCurrentSite()
			// }
			// if err != nil {
			// 	log.Fatal(err)
			// }
			// nsxtclient = client.NewNsxtClient(false, debug, site.Proxy)
			// nsxtclient.BaseUrl = site.Endpoint
			// nsxtclient.Login(site.GetCredential())
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			// nsxtclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// nsxtclient.ClusterInfo()
		},
	}

	return clusterCmd
}
