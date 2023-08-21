package nsx

import "encoding/json"

func (a *Agent) GetVersion() string {
	var v Version
	path := "/api/v1/node/version"
	req, err := a.client.NewRequest("GET", path, nil)
	if err != nil {
		return err.Error()
	}
	res := a.client.Call(req)
	contentByte, _ := json.Marshal(res.Content)
	json.Unmarshal(contentByte, &v)
	return v.ProductVersion
}
