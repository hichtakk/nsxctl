package cmd

import (
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/hichtakk/nsxctl/client"
	"github.com/spf13/cobra"
)

// NewCmdExec is subcommand to exec api.
func NewCmdExec() *cobra.Command {
	var execCmd = &cobra.Command{
		Use:   "exec",
		Short: "exec",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
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
	}
	execCmd.AddCommand(
		NewCmdHttpGet(),
		//NewCmdHttpPost(),
		//NewCmdHttpPut(),
		//NewCmdHttpPatch(),
		//NewCmdHttpDelete(),
	)

	return execCmd
}

func NewCmdHttpGet() *cobra.Command {
	data := ""
	httpGet := &cobra.Command{
		Use:   "get",
		Short: "call api with HTTP GET method",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			nsxtclient.Request("GET", args[0], data)
		},
	}

	return httpGet
}
