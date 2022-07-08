package cmd

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/hichtakk/nsxctl/client"
	"github.com/hichtakk/nsxctl/config"
	ac "github.com/hichtakk/nsxctl/nsxalb"
	"github.com/spf13/cobra"
)

var noPretty bool
var alb bool

// NewCmdExec is subcommand to exec api.
func NewCmdExec() *cobra.Command {
	var query []string
	var execCmd = &cobra.Command{
		Use:   "exec",
		Short: "call API directly\nYou can find NSX-T REST API reference on https://developer.vmware.com/apis/1163/nsx-t",
		PersistentPreRunE: func(c *cobra.Command, args []string) error {
			file, _ := ioutil.ReadFile(configfile)
			json.Unmarshal(file, &conf)
			if alb == false {
				nsxtclient = client.NewNsxtClient(false, debug)
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
				nsxtclient.BaseUrl = site.Endpoint
				nsxtclient.Login(site.GetCredential())
			} else {
				albclient = ac.NewNsxAlbClient(false, debug)
				var site config.NsxAlbSite
				var err error
				if useSite != "" {
					site, err = conf.NsxAlb.GetSite(useSite)
				} else {
					site, err = conf.NsxAlb.GetCurrentSite()
				}
				if err != nil {
					log.Fatal(err)
				}
				albclient.BaseUrl = site.Endpoint
				albclient.Login(site.GetCredential())
			}
			return nil
		},
	}
	execCmd.AddCommand(
		NewCmdHttpGet(&query),
		NewCmdHttpPost(),
		NewCmdHttpPut(),
		NewCmdHttpPatch(),
		NewCmdHttpDelete(&query, &alb),
	)
	execCmd.PersistentFlags().StringSliceVarP(&query, "query", "q", []string{}, "")
	execCmd.PersistentFlags().BoolVarP(&noPretty, "no-pretty", "", false, "pretty output json")
	execCmd.PersistentFlags().BoolVarP(&alb, "alb", "", false, "call api to NSX ALB")

	return execCmd
}

func NewCmdHttpGet(query *[]string) *cobra.Command {
	httpGetCmd := &cobra.Command{
		Use:   "get ${API-PATH}",
		Short: "call api with HTTP GET method",
		Long:  "example) nsxctl exec get /policy/api/v1/infra/tier-0s",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			params := map[string]string{}
			for _, q := range *query {
				qSlice := strings.Split(q, "=")
				if len(qSlice) != 2 {
					panic("invalid query parameter. it should be formatted as '<name>=<value>'.")
				}
				params[qSlice[0]] = qSlice[1]
			}
			var resp *client.Response
			if alb == false {
				resp = nsxtclient.Request("GET", args[0], params, []byte{})
			} else {
				resp = albclient.Request("GET", args[0], params, []byte{})
			}
			resp.Print(noPretty)
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
				raw_data, err = readRequestBody(fileName)
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
			var resp *client.Response
			if alb == false {
				resp = nsxtclient.Request("POST", args[0], nil, data)
			} else {
				resp = albclient.Request("POST", args[0], nil, data)
			}
			resp.Print(noPretty)
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
				raw_data, err = readRequestBody(fileName)
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
			var resp *client.Response
			if alb == false {
				resp = nsxtclient.Request("PUT", args[0], nil, data)
			} else {
				resp = albclient.Request("PUT", args[0], nil, data)
			}
			resp.Print(noPretty)
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
		Long:  "example) nsxctl exec patch /policy/api/v1/infra/tier-0s -f ./tier0.json",
		Args:  cobra.ExactArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			var raw_data []byte
			var err error
			if fileName != "" {
				raw_data, err = readRequestBody(fileName)
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
			var resp *client.Response
			if alb == false {
				resp = nsxtclient.Request("PATCH", args[0], nil, data)
			} else {
				resp = albclient.Request("PATCH", args[0], nil, data)
			}
			resp.Print(noPretty)
		},
	}
	httpPatchCmd.Flags().StringVarP(&fileName, "filename", "f", "", "file name for send data(json)")

	return httpPatchCmd
}

func NewCmdHttpDelete(query *[]string, alb *bool) *cobra.Command {
       httpDeleteCmd := &cobra.Command{
               Use:   "delete",
               Short: "call api with HTTP DELETE method",
               Long:  "example) nsxctl exec delete /policy/api/v1/infra/tier-0s",
               Args:  cobra.ExactArgs(1),
               Run: func(cmd *cobra.Command, args []string) {
                       if *alb == false {
                               nsxtclient.Request("DELETE", args[0], nil, []byte{})
                       } else {
                               albclient.Request("DELETE", args[0], nil, []byte{})
                       }
               },
       }

       return httpDeleteCmd
}

func readRequestBody(fileName string) ([]byte, error) {
	if fileName == "-" {
		return readFromStdIn()
	} else {
		return ioutil.ReadFile(fileName)
	}
}

func readFromStdIn() ([]byte, error) {
	var body string
	stdin := bufio.NewScanner(os.Stdin)
	for stdin.Scan() {
		if err := stdin.Err(); err != nil {
			return []byte{}, err
		}
		body += stdin.Text()
	}
	return []byte(body), nil
}
