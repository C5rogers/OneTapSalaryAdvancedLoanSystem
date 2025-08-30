package graph

import (
	"net/http"

	"github.com/hasura/go-graphql-client"
	"github.com/c5rogers/one-tap/salary-advance-loan-system/config"
)

// customTransport adds custom headers to every request
type customTransport struct {
	headers map[string]string
	rt      http.RoundTripper
}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	for key, value := range t.headers {
		req.Header.Set(key, value)
	}
	return t.rt.RoundTrip(req)
}

type LMSGraphClient struct {
	graphql.Client
	*config.Config
}

func MakeClient(config *config.Config) *LMSGraphClient {

	httpClient := &http.Client{
		Transport: &customTransport{
			headers: map[string]string{
				"x-hasura-admin-secret": config.Graph.AdminSecret,
			},
			rt: http.DefaultTransport,
		},
	}

	client := graphql.NewClient(config.Graph.URL, httpClient)

	return &LMSGraphClient{
		Client: *client,
		Config: config,
	}
}
