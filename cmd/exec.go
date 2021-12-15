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
		Short: "call API directly",
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
		NewCmdHttpPost(),
		NewCmdHttpPut(),
		NewCmdHttpPatch(),
		//NewCmdHttpDelete(),
	)

	return execCmd
}

func NewCmdHttpGet() *cobra.Command {
	httpGetCmd := &cobra.Command{
		Use:   "get ${API-PATH}",
		Short: "call api with HTTP GET method",
		Long:  "ex) nsxctl exec get /policy/api/v1/infra/tier-0s",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			nsxtclient.Request("GET", args[0], []byte{})
		},
	}

	return httpGetCmd
}

func NewCmdHttpPost() *cobra.Command {
	fileName := ""
	var data []byte
	httpPostCmd := &cobra.Command{
		Use:   "post",
		Short: "call api with HTTP POST method",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var raw_data []byte
			var err error
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
		Run: func(cmd *cobra.Command, args []string) {
			nsxtclient.Request("POST", args[0], data)
		},
	}
	httpPostCmd.Flags().StringVarP(&fileName, "filename", "f", "", "file name for send data(json)")

	return httpPostCmd
}

func NewCmdHttpPut() *cobra.Command {
	fileName := ""
	var data []byte
	httpPutCmd := &cobra.Command{
		Use:   "put",
		Short: "call api with HTTP PUT method",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var raw_data []byte
			var err error
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
		Run: func(cmd *cobra.Command, args []string) {
			nsxtclient.Request("PUT", args[0], data)
		},
	}
	httpPutCmd.Flags().StringVarP(&fileName, "filename", "f", "", "file name for send data(json)")

	return httpPutCmd
}

func NewCmdHttpPatch() *cobra.Command {
	fileName := ""
	var data []byte
	httpPatchCmd := &cobra.Command{
		Use:   "patch",
		Short: "call api with HTTP PATCH method",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var raw_data []byte
			var err error
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
		Run: func(cmd *cobra.Command, args []string) {
			nsxtclient.Request("PATCH", args[0], data)
		},
	}
	httpPatchCmd.Flags().StringVarP(&fileName, "filename", "f", "", "file name for send data(json)")

	return httpPatchCmd
}
