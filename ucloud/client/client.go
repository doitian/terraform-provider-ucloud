package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const DefaultEndpoint = "https://api.ucloud.cn"

type Client struct {
	HttpClient *http.Client
	Endpoint   string
	PublicKey  string
	PrivateKey string
	ProjectId  string
	Region     string
}

func (c *Client) endpoint() string {
	if c.Endpoint == "" {
		return DefaultEndpoint
	}

	return c.Endpoint
}

func (c *Client) httpClient() *http.Client {
	if c.HttpClient == nil {
		return http.DefaultClient
	}

	return c.HttpClient
}

func (c *Client) Validate() error {
	if c.PublicKey == "" {
		return InvalidClientFieldError("PublicKey")
	}
	if c.PrivateKey == "" {
		return InvalidClientFieldError("PrivateKey")
	}
	if c.Region == "" {
		return InvalidClientFieldError("Region")
	}

	return nil
}

// Get calls UCloud API. It will generate signature and append it automatically.
func (c *Client) Get(params url.Values) (resp *http.Response, err error) {
	err = c.Validate()
	if err != nil {
		return
	}

	params.Set("PublicKey", c.PublicKey)
	params.Set("Region", c.Region)
	if c.ProjectId != "" {
		params.Set("ProjectId", c.ProjectId)
	}

	targetUrl := c.endpoint() + "?" + params.Encode() + "&Signature=" + GenerateSignature(params, c.PrivateKey)

	return c.httpClient().Get(targetUrl)
}

func (c *Client) GetJSON(params url.Values, v interface{}) error {
	resp, err := c.Get(params)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	return json.Unmarshal(bytes, v)
}
