package ucloud

import (
	"net/http"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"

	"github.com/3pjgames/terraform-provider-ucloud/ucloud/client"
)

func Provider() terraform.ResourceProvider {
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
		},

		ResourcesMap: map[string]*schema.Resource{},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	apiClient := &client.Client{
		HttpClient: &http.Client{},
		PublicKey:  d.Get("public_key").(string),
		PrivateKey: d.Get("private_key").(string),
		ProjectId:  d.Get("project_id").(string),
		Region:     d.Get("region").(string),
	}

	err := apiClient.Validate()

	if err != nil {
		return nil, err
	}

	return apiClient, nil
}
