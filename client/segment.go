package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hichtakk/nsxctl/structs"
)

func (c *NsxtClient) GetSegment() structs.Segments {
	path := "/policy/api/v1/infra/segments"
	res := c.Request("GET", path, nil, nil)
	segments := []structs.Segment{}
	body, err := res.BodyBytes()
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(body, &segments)

	for _, seg := range res.Body.(map[string]interface{})["results"].([]interface{}) {
		str, _ := json.Marshal(seg)
		var segment structs.Segment
		json.Unmarshal(str, &segment)
		segments = append(segments, segment)
	}

	return structs.Segments(segments)
}

func (c *NsxtClient) CreateSegment(segment_name string, transportzone string, vlan_id string, gateway_name string, interface_address string) error {
	for _, seg := range c.GetSegment() {
		if seg.Name == segment_name {
			return errors.New(fmt.Sprintf("segment '%s' is already exists", segment_name))
		}
	}

	var segment structs.Segment
	segment.Name = segment_name
	segment.Id = strings.ReplaceAll(segment_name, " ", "_")

	if transportzone != "" {
		sites := c.GetSite()
		endpoints := c.GetEnforcementPoint(sites[0])
		transport_zones := c.GetPolicyTransportZone(sites[0], (*endpoints)[0].Id)
		for _, tz := range *transport_zones {
			if tz.Name == transportzone {
				segment.TransportZonePath = tz.Path
				if tz.Type == "VLAN_BACKED" {
					if vlan_id == "" {
						return errors.New("vlan-id must be set if the specified transportzone type is vlan")
					}
					if gateway_name != "" {
						return errors.New("cannot connect to gateway if the specified transportzone type is vlan")
					}
					segment.Vlans = strings.Split(vlan_id, ",")
				}
			}
		}
		if segment.TransportZonePath == "" {
			return errors.New(fmt.Sprintf("transportzone '%s' is not found", transportzone))
		}
	}

	if gateway_name != "" {
		if interface_address == "" {
			return errors.New("Interface address must be set if gateway name is specified")
		}
		var gw structs.Gateway
		var err error
		gw, err = c.GetGatewayFromName(gateway_name, -1)
		if err != nil {
			return errors.New(fmt.Sprintf("gateway '%s' is not found", gateway_name))
		}
		if gw.Id != "" {
			segment.Connectivity = gw.Path
		}
		segment.Subnets = []structs.SegmentSubnet{{
			Gateway: interface_address,
		}}
	}

	payload, err := json.Marshal(segment)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	path := "/policy/api/v1/infra/segments/" + segment.Id
	res := c.Request("PATCH", path, nil, payload)
	if res.StatusCode != 200 {
		body, err := res.BodyBytes()
		if err != nil {
			return nil
		}
		return errors.New(string(body))
	}

	return nil
}

func (c *NsxtClient) DeleteSegment(segment_name string) error {
	var segment_id string
	for _, seg := range c.GetSegment() {
			if seg.Name == segment_name {
					segment_id = seg.Id
			}
	}
	if segment_id == "" {
			return errors.New(fmt.Sprintf("segment '%s' not found", segment_name))
	}

	path := "/policy/api/v1/infra/segments/" + segment_id
	res := c.Request("DELETE", path, nil, nil)
	if res.StatusCode != 200 {
			body, err := res.BodyBytes()
			if err != nil {
					return nil
			}
			return errors.New(string(body))
	}

	return nil
}

func (c *NsxtClient) UpdateSegment(segment structs.Segment) error {
	payload, err := json.Marshal(segment)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	path := "/policy/api/v1/infra/segments/" + segment.Id
	res := c.Request("PATCH", path, nil, payload)
	if res.StatusCode != 200 {
		body, err := res.BodyBytes()
		if err != nil {
			return nil
		}
		return errors.New(string(body))
	}

	return nil
}
