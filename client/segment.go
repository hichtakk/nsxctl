package client

import (
	"encoding/json"
	"fmt"
)

func (c *NsxtClient) GetSegment() {
	path := "/policy/api/v1/infra/segments"
	req := c.makeRequest("GET", path)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
	}
	defer res.Body.Close()
	data := readResponseBody(res)
	seg := data.(map[string]interface{})["results"]
	for _, s := range seg.([]interface{}) {
		//fmt.Printf("role: %s, permission: %s\n", v.(map[string]interface{})["role"], v.(map[string]interface{})["permissions"])
		b, _ := json.MarshalIndent(s, "", "  ")
		fmt.Println(string(b))
	}
}
