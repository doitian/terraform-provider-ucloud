// +build integration

package client

import (
	"io/ioutil"
	"net/http"
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

	c, err := Config{
		HttpClient: &http.Client{},
		PublicKey:  publicKey,
		PrivateKey: privateKey,
		ProjectId:  projectId,
		Region:     region,
	}.Client()
	if err != nil {
		t.Fatal("Failed to create client: ", err)
	}

	params, err := BuildParams(&DescribeUHostInstanceRequest{
		Limit: 3,
	})
	if err != nil {
		t.Fatal("Failed to build params: ", err)
	}
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
