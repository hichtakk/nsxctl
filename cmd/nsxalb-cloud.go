package cmd

import (
	"fmt"
	"log"
	"strings"
	"text/tabwriter"
	"os"

	ac "github.com/hichtakk/nsxctl/nsxalb"
	"github.com/hichtakk/nsxctl/structs"
	"github.com/spf13/cobra"
)

func NewCmdShowAlbCloud() *cobra.Command {
	aliases := []string{"ac", "cloud"}
	cloudCmd := &cobra.Command{
		Use:     "alb-cloud",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ALB Cloud [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			albclient = ac.NewNsxAlbClient(false, debug)
			albsite, err := conf.NsxAlb.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			albclient.BaseUrl = albsite.Endpoint
			albclient.Login(albsite.GetCredential())
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			albclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			albclient.ShowCloud()
			//albclient.DownloadSeImage()
		},
	}

	return cloudCmd
}

func NewCmdShowAlbVirtualService() *cobra.Command {
	aliases := []string{"alb-vs", "vs"}
	cloudCmd := &cobra.Command{
		Use:     "alb-virtualservice",
		Aliases: aliases,
		Short:   fmt.Sprintf("show ALB VirtualService [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		PreRunE: func(c *cobra.Command, args []string) error {
			albclient = ac.NewNsxAlbClient(false, debug)
			albsite, err := conf.NsxAlb.GetCurrentSite()
			if err != nil {
				log.Fatal(err)
			}
			albclient.BaseUrl = albsite.Endpoint
			albclient.Login(albsite.GetCredential())
			return nil
		},
		PostRunE: func(c *cobra.Command, args []string) error {
			albclient.Logout()
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			vss := albclient.ShowVirtualService()
			clouds := map[string]structs.Cloud{}
			segs := map[string]structs.ServiceEngineGroup{}

			w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
			w.Write([]byte(strings.Join([]string{"ID", "Name", "VIP", "Port", "Cloud", "SEGroup"}, "\t") + "\n"))
			for _, vs := range vss {
				cloudId := vs.GetCloudId()
				cloud, ok := clouds[cloudId]
				if ok != true {
					cloud = albclient.GetCloudById(cloudId)
				}
				segId := vs.GetSegId()
				seg, ok := segs[segId]
				if ok != true {
					seg = albclient.GetSeGroupById(segId)
				}
				vs.Print(cloud.Name, seg.Name, w)
			}
			w.Flush()
		},
	}

	return cloudCmd
}
