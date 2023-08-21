package nsx

func (a *Agent) GetDfwPolicies(domain string, name string) ([]DfwPolicy, error) {
	// path := "/policy/api/v1/infra/domains/" + domain + "/security-policies?include_rule_count=true"
	// res := c.Request("GET", path, nil, nil)

	// body, err := res.BodyBytes()
	// if err != nil {
	// 	return nil, err
	// }

	var results DfwPolicies
	// json.Unmarshal(body, &results)
	// if name != "" {
	// 	for _, p := range results.Policies {
	// 		if name == p.Name {
	// 			return []DfwPolicy{p}, nil
	// 		}
	// 	}
	// 	return nil, fmt.Errorf("Dfw Policy '%s' is not found", name)
	// }
	return results.Policies, nil
}

func (a *Agent) GetDfwRules(policy DfwPolicy) []DfwRule {
	// // in case of using multi-byte characters that the following API doen't work
	// // GET /policy/api/v1/infra/domains/<domain-id>/security-policies/<security-policy-id>
	// parent_path := "*security-policies\\/" + policy.Id
	// path := "/policy/api/v1/search/query?query=resource_type:Rule%20AND%20parent_path:" + url.PathEscape(parent_path)
	// res := c.Request("GET", path, nil, nil)

	// body, err := res.BodyBytes()
	// if err != nil {
	// 	log.Fatal(err)
	// 	return nil
	// }
	var results DfwRules
	// json.Unmarshal(body, &results)
	return results.Rules
}
