package ucloud

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceUHostCreate,
		Read:   resourceUHostRead,
		Update: resourceUHostUpdate,
		Delete: resourceUHostDelete,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "UHost name",
			},
		},
	}
}

func resourceUHostCreate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUHostRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUHostUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceUHostDelete(d *schema.ResourceData, meta interface{}) error {
	return nil
}
