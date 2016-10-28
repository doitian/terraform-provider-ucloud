package ucloud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/3pjgames/terraform-provider-ucloud/ucloud/client"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccResourceSecurityGroup(t *testing.T) {
	var group client.SecurityGroup

	resource.Test(t, resource.TestCase{
		PreCheck:      func() { testAccPreCheck(t) },
		IDRefreshName: "ucloud_security_group.foo",
		Providers:     testAccProviders,
		CheckDestroy:  testAccCheckSecurityGroupDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccSecurityGroupConfig_pre,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("ucloud_security_group.foo", &group),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "group_name", "foo"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "description", "bar"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.0.protocol_type", "TCP"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.0.dst_port", "3306"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.0.src_ip", "10.0.0.0/0"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.0.rule_action", "ACCEPT"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.0.priority", "100"),
				),
			},
			resource.TestStep{
				Config: testAccSecurityGroupConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSecurityGroupExists("ucloud_security_group.foo", &group),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "group_name", "foo"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "description", "bar"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.#", "1"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.0.protocol_type", "UDP"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.0.dst_port", "3"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.0.src_ip", "0.0.0.0/0"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.0.rule_action", "DROP"),
					resource.TestCheckResourceAttr("ucloud_security_group.foo", "rule.0.priority", "50"),
				),
			},
		},
	})
}

const testAccSecurityGroupConfig_pre = `
resource "ucloud_security_group" "foo" {
	group_name = "foo"
	description = "bar"
	rule {
		protocol_type = "TCP"
		dst_port = "3306"
		src_ip = "10.0.0.0/0"
		rule_action = "ACCEPT"
		priority = 100
	}
}
`

const testAccSecurityGroupConfig = `
resource "ucloud_security_group" "foo" {
	group_name = "foo"
	description = "bar"
	rule {
		protocol_type = "UDP"
		dst_port = "3"
		src_ip = "0.0.0.0/0"
		rule_action = "DROP"
		priority = 50
	}
}
`

func testAccCheckSecurityGroupDestroy(s *terraform.State) error {
	return testAccCheckSecurityGroupDestroyWithProvider(s, testAccProvider)
}

func testAccCheckSecurityGroupDestroyWithProvider(s *terraform.State, provider *schema.Provider) error {
	apiClient := provider.Meta().(*client.Client)
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ucloud_uhost" {
			continue
		}

		var resp client.DescribeSecurityGroupResponse
		id, err := strconv.Atoi(rs.Primary.ID)
		if err != nil {
			return err
		}
		err = apiClient.Call(&client.DescribeSecurityGroupRequest{GroupId: id}, &resp)
		if err == nil {
			for _, i := range resp.DataSet {
				return fmt.Errorf("Found unterminated instance: %+v", i)
			}
		}

		return err
	}

	return nil
}

func testAccCheckSecurityGroupExists(n string, i *client.SecurityGroup) resource.TestCheckFunc {
	providers := []*schema.Provider{testAccProvider}
	return testAccCheckSecurityGroupExistsWithProviders(n, i, &providers)
}

func testAccCheckSecurityGroupExistsWithProviders(n string, i *client.SecurityGroup, providers *[]*schema.Provider) resource.TestCheckFunc {
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
			var resp client.DescribeOneSecurityGroupResponse
			id, err := strconv.Atoi(rs.Primary.ID)
			if err != nil {
				return err
			}
			err = apiClient.Call(&client.DescribeSecurityGroupRequest{GroupId: id}, &resp)
			if err != nil {
				return err
			}

			*i = *resp.DataSet
			return nil
		}

		return fmt.Errorf("Instance not found")
	}
}
