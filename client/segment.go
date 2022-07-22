package client

import (
	"encoding/json"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxtClient) GetSegment() structs.Segments {
	path := "/policy/api/v1/infra/segments"
	res := c.Request("GET", path, nil, nil)
	segments := []structs.Segment{}
	body, _ := res.BodyBytes()
	json.Unmarshal(body, &segments)

	for _, seg := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		str, _ := json.Marshal(seg)
		var segment structs.Segment
		json.Unmarshal(str, &segment)
		segments = append(segments, segment)
	}

	return structs.Segments(segments)
}
