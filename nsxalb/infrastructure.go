package nsxalb

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func (c *NsxAlbClient) ShowCloud() {
	resp := c.Request("GET", "/api/cloud", map[string]string{}, nil)
	resByte, _ := resp.BodyBytes()
	var res interface{}
	json.Unmarshal(resByte, &res)
	clouds := res.(map[string]interface{})["results"]
	for _, cloud := range clouds.([]interface{}) {
		name := cloud.(map[string]interface{})["name"]
		uuid := cloud.(map[string]interface{})["uuid"]
		vtype := cloud.(map[string]interface{})["vtype"]
		fmt.Printf("%s  %s  %s\n", name, uuid, vtype)
	}
}

func (c *NsxAlbClient) GenerateSeImage() {
	c.Request("POST", "/api/fileservice/seova", map[string]string{}, nil)
	fmt.Println("file will be generated in /host/pkgs/21.1.1-9045-20210811.170844")
}

func (c *NsxAlbClient) DownloadSeImage() {
	//c.Request("GET", "/api/fileservice/seova?file_format=ova&cloud_uuid=cloud-98072415-b7a4-4138-bc23-f242cc0b99a2")
	// Create the file
	out, err := os.Create("./nsxalb/se.ova")
	if err != nil {
		fmt.Println("create error")
		return
	}
	defer out.Close()
	path := "/api/fileservice/seova"
	req, _ := http.NewRequest("GET", c.BaseUrl+path, nil)
	q := req.URL.Query()
	q.Add("file_format", "ova")
	q.Add("cloud_uuid", "cloud-98072415-b7a4-4138-bc23-f242cc0b99a2")
	req.URL.RawQuery = q.Encode()
	if c.Token != "" {
		req.Header.Set("X-CSRFToken", c.Token)
	}
	req.Header.Set("Referer", c.BaseUrl)
	res, err := c.httpClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()
	//res_body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Println(err)
	}
	if res.StatusCode != 200 && res.StatusCode != 201 {
		fmt.Printf("StatusCode: %d\n", res.StatusCode)
	}

	// Writer the body to file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return
	}
	fmt.Println("finished")

	return
}
