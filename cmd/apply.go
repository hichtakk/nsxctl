package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/hichtakk/nsxctl/client"
	"github.com/hichtakk/nsxctl/config"
	"github.com/spf13/cobra"
)

// NewCmdApply is subcommand to call Policy API with resource manifest.
func NewCmdApply() *cobra.Command {
	fileName := ""
	var data []byte
	var applyCmd = &cobra.Command{
		Use:   "apply",
		Short: "apply Policy API declared JSON file",
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
			nsxtclient.Login(site.GetCredential())
			var raw_data []byte
			if fileName != "" {
				raw_data, err = ioutil.ReadFile(fileName)
				if err != nil {
					return err
				}
			} else {
				raw_data = nil
			}
			jsonObj := json.RawMessage(raw_data)
			data, err = json.Marshal(jsonObj)
			if err != nil {
				return err
			}
			return nil
		},
		Run: func(cmd *cobra.Command, _ []string) {
			nsxtclient.Request("PATCH", "/policy/api/v1/infra", nil, data)
		},
	}
	applyCmd.Flags().StringVarP(&fileName, "filename", "f", "", "file name for send data(json)")

	return applyCmd
}
