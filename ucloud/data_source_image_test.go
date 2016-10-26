package ucloud

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccDataSourceImage_GetUbuntuImage(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccCheckImageDataSourceConfig, os.Getenv("UCLOUD_ZONE")),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckImageDataSourceID("data.ucloud_image.ubuntu"),
					resource.TestCheckResourceAttr("data.ucloud_image.ubuntu", "image_name", "Ubuntu 14.04 64位"),
				),
			},
		},
	})
}

func testAccCheckImageDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Data source ID not set")
		}
		return nil
	}
}

const testAccCheckImageDataSourceConfig = `
data "ucloud_image" "ubuntu" {
	zone = "%s",
	image_type = "Base",
	os_type = "Linux",
	image_name_regexp = "^Ubuntu 14.04 64位$"
}
`
