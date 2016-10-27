package ucloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/3pjgames/terraform-provider-ucloud/ucloud/client"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceUHost(t *testing.T) {
	var host client.UHostInstance

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ucloud_uhost.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckUHostDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccUHostConfig_pre, os.Getenv("UCLOUD_ZONE")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUHostExists("ucloud_uhost.foo", &host),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "name", "foo"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "remark", "bar"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "tag", "test"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "cpu", "1"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "memory", "1024"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "disk_space", "10"),
				),
			},
			resource.TestStep{
				Config: fmt.Sprintf(testAccUHostConfig, os.Getenv("UCLOUD_ZONE")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUHostExists("ucloud_uhost.foo", &host),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "name", "foox"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "remark", "barx"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "tag", "testx"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "cpu", "2"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "memory", "2048"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "disk_space", "0"),
				),
			},
			// Repeat the config
			resource.TestStep{
				Config: fmt.Sprintf(testAccUHostConfig, os.Getenv("UCLOUD_ZONE")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckUHostExists("ucloud_uhost.foo", &host),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "name", "foox"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "remark", "barx"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "tag", "testx"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "cpu", "2"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "memory", "2048"),
					resource.TestCheckResourceAttr("ucloud_uhost.foo", "disk_space", "0"),
				),
			},
		},
	})
}

const testAccUHostConfig_pre = `
resource "ucloud_uhost" "foo" {
	zone = "%s"
	name = "foo"
	remark = "bar"
	tag = "test"
	cpu = 1
	memory = 1024
	disk_space = 10
	password = "dGVycmFmb3JtLXByb3ZpZGVyLXVjbG91ZA=="
	image_id = "uimage-j4fbrn"
	charge_type = "Dynamic"
}
`

const testAccUHostConfig = `
resource "ucloud_uhost" "foo" {
	zone = "%s"
	name = "foox"
	remark = "barx"
	tag = "testx"
	cpu = 2
	memory = 2048
	disk_space = 0
	password = "dGVycmFmb3JtLXByb3ZpZGVyLXVjbG91ZA=="
	image_id = "uimage-j4fbrn"
	charge_type = "Dynamic"
}
`

func testAccCheckUHostDestroy(s *terraform.State) error {
	return testAccCheckUHostDestroyWithProvider(s, testAccProvider)
}

func testAccCheckUHostDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	apiClient := provider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_uhost" {
			continue
		}

		var describeInstancesResp client.DescribeUHostInstanceResponse
		err := apiClient.Call(&client.DescribeUHostInstanceRequest{UHostIds: []string{rs.Primary.ID}}, &describeInstancesResp)
		if err == nil {
			for _, i := range describeInstancesResp.UHostSet {
				return fmt.Errorf("Found unterminated instance: %+v", i)
			}
		}

		return err
	}

	return nil
}

func testAccCheckUHostExists(n string, i *client.UHostInstance) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckUHostExistsWithProviders(n, i, &providers)
}

func testAccCheckUHostExistsWithProviders(n string, i *client.UHostInstance, providers *[]*schema.Provider) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}
		for _, provider := range *providers {
			// Ignore if Meta is empty, this can happen for validation providers
			if provider.Meta() == nil {
				continue
			}

			apiClient := provider.Meta().(*client.Client)
			var describeInstancesResp client.DescribeUHostInstanceResponse
			err := apiClient.Call(&client.DescribeUHostInstanceRequest{
				UHostIds: []string{rs.Primary.ID},
			}, &describeInstancesResp)
			if err != nil {
				return err
			}

			if len(describeInstancesResp.UHostSet) > 0 {
				*i = describeInstancesResp.UHostSet[0]
				return nil
			}
		}

		return fmt.Errorf("Instance not found")
	}
}
