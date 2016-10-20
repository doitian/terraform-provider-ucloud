package client

import (
	"testing"

	"net/url"
)

// Text Example in https://docs.ucloud.cn/api/summary/signature
//
// PublicKey  = 'ucloudsomeone@example.com1296235120854146120'
// PrivateKey = '46f09bb9fab4f12dfc160dae12273d5332b5debe'
// 注解：你可以使用上述的 PublicKey 和 PrivateKey 调试你的代码， 当得到跟后面一致的签名结果后(即表示你的代码是正确的)，可再换为你自己的 PublicKey 和 PrivateKey 以及其他 API 请求。
// 本例中假设用户请求参数串如下:
//
// {
//     "Action"     :  "CreateUHostInstance",
//     "Region"     :  "cn-bj2",
//     "Zone"       :  "cn-bj2-04",
//     "ImageId"    :  "f43736e1-65a5-4bea-ad2e-8a46e18883c2",
//     "CPU"        :  2,
//     "Memory"     :  2048,
//     "DiskSpace"  :  10,
//     "LoginMode"  :  "Password",
//     "Password"   :  "VUNsb3VkLmNu",
//     "Name"       :  "Host01",
//     "ChargeType" :  "Month",
//     "Quantity"   :  1,
//     "PublicKey"  :  "ucloudsomeone@example.com1296235120854146120"
// }

func TestGenerateSampleSignature(t *testing.T) {
	privateKey := "46f09bb9fab4f12dfc160dae12273d5332b5debe"

	params := url.Values{}
	params.Set("Action", "CreateUHostInstance")
	params.Set("Region", "cn-bj2")
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
	params.Set("PublicKey", "ucloudsomeone@example.com1296235120854146120")

	expected := "4f9ef5df2abab2c6fccd1e9515cb7e2df8c6bb65"
	result := GenerateSignature(params, privateKey)

	if result != expected {
		t.Error("Failed the sample signautre: ", result)
	}
}
