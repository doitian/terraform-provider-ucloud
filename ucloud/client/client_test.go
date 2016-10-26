package client

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
)

// Test InvalidClientFieldError
func TestClientConfig(t *testing.T) {
	c := Config{}
	_, err := c.Client()
	if err == nil || !strings.Contains(err.Error(), "PublicKey") {
		t.Error("Expect InvalidClientFieldError on PublicKey")
	}

	c = Config{PublicKey: "ucloudsomeone@example.com1296235120854146120"}
	_, err = c.Client()
	if err == nil || !strings.Contains(err.Error(), "PrivateKey") {
		t.Error("Expect InvalidClientFieldError on PrivateKey")
	}

	c = Config{
		PublicKey:  "ucloudsomeone@example.com1296235120854146120",
		PrivateKey: "46f09bb9fab4f12dfc160dae12273d5332b5debe",
	}
	_, err = c.Client()
	if err == nil || !strings.Contains(err.Error(), "Region") {
		t.Error("Expect InvalidClientFieldError on Region")
	}
}

func TestClientSampleSignature(t *testing.T) {
	handler := http.NotFound
	hs := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		handler(rw, req)
	}))
	defer hs.Close()

	c, err := Config{
		Endpoint:   hs.URL,
		PublicKey:  "ucloudsomeone@example.com1296235120854146120",
		PrivateKey: "46f09bb9fab4f12dfc160dae12273d5332b5debe",
		Region:     "cn-bj2",
	}.Client()
	if err != nil {
		t.Fatal("Error creating client ", err)
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			t.Error("Bad path!")
		}
		query := req.URL.Query()
		if query.Get("Signature") != "4f9ef5df2abab2c6fccd1e9515cb7e2df8c6bb65" {
			t.Error("Invalid Signature!")
		}
		io.WriteString(rw, `{"RetCode":0}`)
	}

	params := CreateUHostInstanceRequest{
		Zone:       "cn-bj2-04",
		ImageId:    "f43736e1-65a5-4bea-ad2e-8a46e18883c2",
		CPU:        2,
		Memory:     2048,
		DiskSpace:  10,
		LoginMode:  "Password",
		Password:   "VUNsb3VkLmNu",
		Name:       "Host01",
		ChargeType: "Month",
		Quantity:   1,
	}
	var resp GeneralResponse

	err = c.Call(&params, &resp)
	if err != nil {
		t.Fatal("Got error sending item: ", err)
	}
}

func TestAccDescribeUHostInstance(t *testing.T) {
	if os.Getenv(resource.TestEnvVar) == "" {
		return
	}

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

	params := DescribeUHostInstanceRequest{
		Limit: 3,
	}
	var resp DescribeUHostInstanceResponse
	err = c.Call(&params, &resp)

	if err != nil {
		t.Fatal("Failed to call API: ", err)
	}

	log.Printf("[DEBUG] Got response: ", resp)
}
