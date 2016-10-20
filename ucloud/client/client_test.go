package client

import (
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// Test InvalidClientFieldError
func TestFieldError(t *testing.T) {
	c := &Client{}
	_, err := c.Get(url.Values{})
	if err == nil || !strings.Contains(err.Error(), "PublicKey") {
		t.Error("Expect InvalidClientFieldError on PublicKey")
	}

	c = &Client{PublicKey: "ucloudsomeone@example.com1296235120854146120"}
	_, err = c.Get(url.Values{})
	if err == nil || !strings.Contains(err.Error(), "PrivateKey") {
		t.Error("Expect InvalidClientFieldError on PrivateKey")
	}

	c = &Client{
		PublicKey:  "ucloudsomeone@example.com1296235120854146120",
		PrivateKey: "46f09bb9fab4f12dfc160dae12273d5332b5debe",
	}
	_, err = c.Get(url.Values{})
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

	c := &Client{
		Endpoint:   hs.URL,
		PublicKey:  "ucloudsomeone@example.com1296235120854146120",
		PrivateKey: "46f09bb9fab4f12dfc160dae12273d5332b5debe",
		Region:     "cn-bj2",
	}

	handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			t.Error("Bad path!")
		}
		query := req.URL.Query()
		if query.Get("Signature") != "4f9ef5df2abab2c6fccd1e9515cb7e2df8c6bb65" {
			t.Error("Invalid Signature!")
		}
		io.WriteString(rw, `OK`)
	}
	params := url.Values{}

	params.Set("Action", "CreateUHostInstance")
	params.Set("Zone", "cn-bj2-04")
	params.Set("ImageId", "f43736e1-65a5-4bea-ad2e-8a46e18883c2")
	params.Set("CPU", "2")
	params.Set("Memory", "2048")
	params.Set("DiskSpace", "10")
	params.Set("LoginMode", "Password")
	params.Set("Password", "VUNsb3VkLmNu")
	params.Set("Name", "Host01")
	params.Set("ChargeType", "Month")
	params.Set("Quantity", "1")

	resp, err := c.Get(params)
	if err != nil {
		t.Fatal("Got error sending item")
	}
	if resp == nil {
		t.Fatal("Expect valid resp but got nil")
	}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("Go error reading body")
	}
	body := string(bodyBytes)
	if body != "OK" {
		t.Fatal("Expect body to be OK but got: " + body)
	}
}
