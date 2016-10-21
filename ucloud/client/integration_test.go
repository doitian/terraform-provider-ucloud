// +build integration

package client

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"
)

func TestDescribeUhost(t *testing.T) {
	publicKey := os.Getenv("UCLOUD_PUBLIC_KEY")
	privateKey := os.Getenv("UCLOUD_PRIVATE_KEY")
	region := os.Getenv("UCLOUD_REGION")
	projectId := os.Getenv("UCLOUD_PROJECT_ID")

	if publicKey == "" {
		t.Fatal("UCLOUD_PUBLIC_KEY is not set")
	}
	if privateKey == "" {
		t.Fatal("UCLOUD_PRIVATE_KEY is not set")
	}
	if region == "" {
		t.Fatal("UCLOUD_REGION is not set")
	}

	c := &Client{
		HttpClient: &http.Client{},
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		ProjectId:  projectId,
		Region:     region,
	}

	params := url.Values{}
	params.Set("Action", "DescribeUHostInstance")
	params.Set("Limit", "3")
	resp, err := c.Get(params)

	if err != nil {
		t.Fatal("Failed to call API: ", err)
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Failed to read response body: ", err)
	}

	t.Log("Got response: ", string(bodyBytes))
}
