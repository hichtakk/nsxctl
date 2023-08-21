package client

import "net/http"

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type RestApiClient interface {
	Call()
	Login()
	Logout()
	NewRequest()
	SetHttpClient()
}
