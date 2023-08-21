package client

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

func newHttpClient(proxy string) (*http.Client, error) {
	transportConfig := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	if proxy != "" {
		proxyUrl, err := url.Parse(proxy)
		if err != nil {
			return nil, err
		}
		transportConfig.Proxy = http.ProxyURL(proxyUrl)
	}
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Transport: transportConfig,
		Timeout:   time.Duration(30) * time.Second,
		Jar:       jar,
	}

	return client, nil
}

func validatePrefix(path string, acceptablePrefixes []string) error {
	var match bool
	match = false
	for _, prefix := range acceptablePrefixes {
		if strings.HasPrefix(path, prefix) {
			match = true
		}
	}
	if !match {
		prfx := strings.Join(acceptablePrefixes, " or ")
		return fmt.Errorf("path must start with " + prfx)
	}
	return nil
}

func AddQuery(req *http.Request, param map[string]string) {
	if len(param) > 0 {
		params := req.URL.Query()
		for k, v := range param {
			params.Add(k, v)
		}
		req.URL.RawQuery = params.Encode()
	}
}
