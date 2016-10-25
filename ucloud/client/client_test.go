package client

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// Test InvalidClientFieldError
func TestConfig(t *testing.T) {
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

func TestSampleSignature(t *testing.T) {
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
