package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxtClient) GetDfwPolicies(domain string, name string) ([]structs.DfwPolicy, error) {
	path := "/policy/api/v1/infra/domains/" + domain + "/security-policies?include_rule_count=true"
	res := c.Request("GET", path, nil, nil)

	body, err := res.BodyBytes()
	if err != nil {
		return nil, err
	}
	
	var results structs.DfwPolicies
	json.Unmarshal(body, &results)
	if name != "" {
		for _, p := range results.Policies {
			if name == p.Name {
				return []structs.DfwPolicy{p}, nil
			}
		}
		return nil, fmt.Errorf("Dfw Policy '%s' is not found", name)
	}
	return results.Policies, nil
}

func (c *NsxtClient) GetDfwRules(policy structs.DfwPolicy) []structs.DfwRule {
	// in case of using multi-byte characters that the following API doen't work
	// GET /policy/api/v1/infra/domains/<domain-id>/security-policies/<security-policy-id>
	parent_path := "*security-policies\\/" + policy.Id
	path := "/policy/api/v1/search/query?query=resource_type:Rule%20AND%20parent_path:" + url.PathEscape(parent_path)
	res := c.Request("GET", path, nil, nil)

	body, err := res.BodyBytes()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	var results structs.DfwRules
	json.Unmarshal(body, &results)
	return results.Rules
}
