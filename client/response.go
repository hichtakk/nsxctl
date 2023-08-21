package client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	jp "github.com/mattn/go-jsonpointer"
)

type Response struct {
	*http.Response
	Content interface{}
}

func (r *Response) Print(noPretty bool) {
	var b []byte
	var body string
	msg := `{
  "code": %d,
  "message": "%v",
  "body": "%v"
}`
	if r.ContentLength != 0 {
		contentType := r.Header.Get("Content-Type")
		if strings.HasPrefix(contentType, "application/json") {
			if noPretty {
				b, _ = json.Marshal(r.Content)
			} else {
				b, _ = json.MarshalIndent(r.Content, "", "  ")
			}
			body = string(b)
		} else if strings.HasPrefix(contentType, "text/html") {
			reason := "non-json response"
			body = fmt.Sprintf(msg, r.StatusCode, reason, r.Content.(string))
		}
	} else {
		reason := "no response body"
		body = fmt.Sprintf(msg, r.StatusCode, reason, "")
	}
	fmt.Println(body)
}

func (r *Response) JsonGet(pointer string) (interface{}, error) {
	return jp.Get(r.Content, pointer)
}
