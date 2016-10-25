package client

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
)

const DefaultEndpoint = "https://api.ucloud.cn"

type Config struct {
	HttpClient *http.Client
	Endpoint   string
	PublicKey  string
	PrivateKey string
	ProjectId  string
	Region     string
}

type Client struct {
	httpClient *http.Client
	endpoint   string
	publicKey  string
	privateKey string
	projectId  string
	region     string
}

type Response interface {
	ValidateResponse() error
}

func (c Config) Client() (*Client, error) {
	if c.PublicKey == "" {
		return nil, InvalidClientFieldError("PublicKey")
	}
	if c.PrivateKey == "" {
		return nil, InvalidClientFieldError("PrivateKey")
	}
	if c.Region == "" {
		return nil, InvalidClientFieldError("Region")
	}

	instance := &Client{
		httpClient: c.HttpClient,
		endpoint:   c.Endpoint,
		publicKey:  c.PublicKey,
		privateKey: c.PrivateKey,
		projectId:  c.ProjectId,
		region:     c.Region,
	}

	if instance.endpoint == "" {
		instance.endpoint = DefaultEndpoint
	}

	if instance.httpClient == nil {
		instance.httpClient = http.DefaultClient
	}

	return instance, nil
}

// Get calls UCloud API. It will generate signature and append it automatically.
func (c *Client) Get(params url.Values) (resp *http.Response, err error) {
	params.Set("PublicKey", c.publicKey)
	params.Set("Region", c.region)
	if c.projectId != "" {
		params.Set("ProjectId", c.projectId)
	}

	targetUrl := c.endpoint + "?" + params.Encode() + "&Signature=" + GenerateSignature(params, c.privateKey)

	return c.httpClient.Get(targetUrl)
}

func (c *Client) Call(req interface{}, v Response) error {
	params, err := BuildParams(req)
	if err != nil {
		return err
	}

	resp, err := c.Get(params)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, v)
	if err != nil {
		return err
	}
	return v.ValidateResponse()
}
