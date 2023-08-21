package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// NewCmdCertificate is subcommand to show certifications.
func NewCmdCertificate() *cobra.Command {
	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show version of nsxctl",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("nsxctl version: %s, revision: %s\n", version, revision)
		},
	}

	return versionCmd
}

// show NSX-T API certificate
func NewCmdShowCertificate() *cobra.Command {
	aliases := []string{"certificate"}
	certCmd := &cobra.Command{
		Use:     "cert",
		Aliases: aliases,
		Short:   fmt.Sprintf("show certificates of NSX-T [%s]", strings.Join(aliases, ",")),
		Run: func(cmd *cobra.Command, args []string) {
			// certs := nsxtclient.GetApiCertificate()
			// w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
			// w.Write([]byte(strings.Join([]string{"Id", "Common Name", "Not Valid After"}, "\t") + "\n"))
			// for _, c := range certs {
			// 	block, _ := pem.Decode([]byte(c.Pem))
			// 	if block == nil {
			// 		panic("failed to parse certificate pem")
			// 	}
			// 	parsed, err := x509.ParseCertificate(block.Bytes)
			// 	if err != nil {
			// 		panic("failed to parse certificate: " + err.Error())
			// 	}
			// 	notAfter := parsed.NotAfter.Format("2006/01/02")
			// 	w.Write([]byte(strings.Join([]string{c.Id, parsed.Subject.CommonName, notAfter}, "\t") + "\n"))
			// }
			// w.Flush()
		},
	}
	//certCmd.PersistentFlags().BoolVarP(&alb, "alb", "", false, "show NSX ALB version")

	return certCmd
}
