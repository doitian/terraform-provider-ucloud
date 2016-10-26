package ucloud

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/3pjgames/terraform-provider-ucloud/ucloud/client"
)

func Provider() terraform.ResourceProvider {
	return ProviderWithConfig(nil)
}

func ProviderWithConfig(c *client.Config) terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"public_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCLOUD_PUBLIC_KEY", nil),
				Description: "UCloud API Public Key",
			},

			"private_key": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCLOUD_PRIVATE_KEY", nil),
				Description: "UCloud API Private Key",
			},

			"project_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCLOUD_PROJECT_ID", nil),
				Description: "UCloud Project ID, leave empty for default project. See https://docs.ucloud.cn/api/summary/get_project_list",
			},
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCLOUD_REGION", nil),
				Description: "UCloud IDC region, see https://docs.ucloud.cn/api/summary/regionlist",
			},
			"endpoint": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("UCLOUD_ENDPOINT", nil),
				Description: "UCloud API Endpoint",
			},
		},

		DataSourcesMap: map[string]*schema.Resource{
			"ucloud_image": dataSourceImage(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"ucloud_uhost": resourceUHost(),
		},

		ConfigureFunc: providerConfigure(c),
	}
}

func providerConfigure(c *client.Config) func(*schema.ResourceData) (interface{}, error) {
	return func(d *schema.ResourceData) (interface{}, error) {
		config := client.Config{}
		if c != nil {
			config = *c
		}

		if config.HttpClient == nil {
			config.HttpClient = &http.Client{}
		}
		if config.PublicKey == "" {
			config.PublicKey = d.Get("public_key").(string)
		}
		if config.PrivateKey == "" {
			config.PrivateKey = d.Get("private_key").(string)
		}
		if config.ProjectId == "" {
			config.ProjectId = d.Get("project_id").(string)
		}
		if config.Region == "" {
			config.Region = d.Get("region").(string)
		}
		if config.Endpoint == "" {
			config.Endpoint = d.Get("endpoint").(string)
		}

		return config.Client()
	}
}
