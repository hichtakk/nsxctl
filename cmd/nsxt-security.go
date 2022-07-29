package cmd

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/spf13/cobra"
)

func NewCmdShowDfwPolicies() *cobra.Command {
	aliases := []string{"dfwp"}
	dfwCmd := &cobra.Command{
		Use:     "dfw-policies",
		Aliases: aliases,
		Short:   fmt.Sprintf("show dfw policies [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			policies, err := nsxtclient.GetDfwPolicies("default", "")
			if err != nil {
				log.Fatal(err)
				return
			}
			w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
			w.Write([]byte(strings.Join([]string{"Name", "Id", "RuleCount", "Scope"}, "\t") + "\n"))
			for _, policy := range policies {
				scope := strings.Join(policy.Scope, ",")
				w.Write([]byte(strings.Join([]string{policy.Name, policy.Id, strconv.Itoa(policy.RuleCount), scope}, "\t") + "\n"))
			}
			w.Flush()
		},
	}

	return dfwCmd
}

func GetPolicyNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	Login()
	policy_names := []string{}
	policies, err := nsxtclient.GetDfwPolicies("default", "")
	if err != nil {
		log.Fatal(err)
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	for _, policy := range policies {
		policy_names = append(policy_names, policy.Name)
	}
	return policy_names, cobra.ShellCompDirectiveNoFileComp
}

func NewCmdShowDfwRules() *cobra.Command {
	var output string
	var policy_name string
	aliases := []string{"dfw"}
	dfwCmd := &cobra.Command{
		Use:     "dfw-rules",
		Aliases: aliases,
		Short:   fmt.Sprintf("show dfw rules [%s]", strings.Join(aliases, ",")),
		Args:    cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			policies, err := nsxtclient.GetDfwPolicies("default", policy_name)
			if err != nil {
				log.Fatal(err)
				return
			}
			arr := []string{"Policy", "RuleName", "RuleId", "Source", "Destination", "Service", "Profile", "AppliedTo", "Action", "Direction", "IPProtocol", "Logged"}

			if output == "csv" {
				w := csv.NewWriter(os.Stdout)
				err := w.Write(arr)
				if err != nil {
					log.Fatal(err)
					return
				}
				for _, policy := range policies {
					for _, r := range nsxtclient.GetDfwRules(policy) {
						r.PrintCsv(w, policy)
					}
				}
				w.Flush()
			} else {
				w := tabwriter.NewWriter(os.Stdout, 0, 1, 3, ' ', 0)
				w.Write([]byte(strings.Join(arr, "\t") + "\n"))
				for _, policy := range policies {
					for _, r := range nsxtclient.GetDfwRules(policy) {
						r.Print(w, policy)
					}
				}
				w.Flush()	
			}
		},
	}
	dfwCmd.Flags().StringVarP(&policy_name, "policy", "p", "", "policy name")
	dfwCmd.RegisterFlagCompletionFunc("policy", GetPolicyNames)
	dfwCmd.Flags().StringVarP(&output, "output", "o", "", "output format (csv)")
	dfwCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective){
		if len(args) != 0 {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		return []string{"csv"}, cobra.ShellCompDirectiveNoFileComp
	})

	return dfwCmd
}
