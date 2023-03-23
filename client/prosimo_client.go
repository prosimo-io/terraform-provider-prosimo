package client

import "net/url"

type ProsimoClient struct {
	api_client *api_client
	baseURL    string
	token      string

	data         map[string]interface{} /* Data as managed by the user */
	api_data     map[string]interface{} /* Data as available from the API */
	api_response string
}

func NewProsimoClient(baseURL string, apiToken string, insecure bool) (*ProsimoClient, error) {

	api_client := newApiClient()
	api_client.token = apiToken
	api_client.baseURL, _ = url.Parse(baseURL)

	prosimo_client := &ProsimoClient{
		api_client: api_client,
		data:       make(map[string]interface{}),
		api_data:   make(map[string]interface{}),
		token:      apiToken,
	}

	return prosimo_client, nil
}
