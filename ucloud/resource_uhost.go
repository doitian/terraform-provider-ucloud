package ucloud

import (
	"encoding/base64"
	"fmt"
	"log"
	"time"

	"github.com/3pjgames/terraform-provider-ucloud/ucloud/client"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceUHost() *schema.Resource {
	return &schema.Resource{
		Create: resourceUHostCreate,
		Read:   resourceUHostRead,
		Update: resourceUHostUpdate,
		Delete: resourceUHostDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"zone": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"image_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"password": {
				Type:      schema.TypeString,
				Required:  true,
				Sensitive: true,
			},

			"cpu": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"memory": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"disk_space": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"network_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"security_group_id": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"charge_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "计费模式，枚举值为： Year，按年付费； Month，按月付费； Dynamic，按需付费（需开启权限）； Trial，试用（需开启权限） 默认为月付",
			},

			"quantity": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"net_capability": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"tag": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"boot_disk_space": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
			},

			"uhost_type": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"basic_image_id": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"basic_image_name": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"remark": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"disk_set": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"disk_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"drive": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"size": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return schema.HashString(m["disk_id"].(string))
				},
			},

			"ip_set": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip_id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"ip": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"bandwidth": {
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
				Set: func(v interface{}) int {
					m := v.(map[string]interface{})
					return schema.HashString(m["ip"].(string))
				},
			},
		},
	}
}

func resourceUHostCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.Client)

	params := client.CreateUHostInstanceRequest{
		Zone:      d.Get("zone").(string),
		ImageId:   d.Get("image_id").(string),
		LoginMode: "Password",
		Password:  base64.StdEncoding.EncodeToString([]byte(d.Get("password").(string))),
	}

	if v, ok := d.GetOk("cpu"); ok {
		params.CPU = v.(int)
	}
	if v, ok := d.GetOk("memory"); ok {
		params.Memory = v.(int)
	}
	if v, ok := d.GetOk("storage_type"); ok {
		params.StorageType = v.(string)
	}
	if v, ok := d.GetOk("disk_space"); ok {
		params.DiskSpace = v.(int)
	}
	if v, ok := d.GetOk("name"); ok {
		params.Name = v.(string)
	}
	if v, ok := d.GetOk("network_id"); ok {
		params.NetworkId = v.(string)
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		params.SecurityGroupId = v.(string)
	}
	if v, ok := d.GetOk("charge_type"); ok {
		params.ChargeType = v.(string)
	}
	if v, ok := d.GetOk("quantity"); ok {
		params.Quantity = v.(int)
	}
	if v, ok := d.GetOk("uhost_type"); ok {
		params.UHostType = v.(string)
	}
	if v, ok := d.GetOk("net_capability"); ok {
		params.NetCapability = v.(string)
	}
	if v, ok := d.GetOk("tag"); ok {
		params.Tag = v.(string)
	}
	if v, ok := d.GetOk("boot_disk_space"); ok {
		params.BootDiskSpace = v.(int)
	}

	log.Printf("[DEBUG] Run configuration: %s", params)

	var resp client.CreateUHostInstanceResponse
	err := apiClient.Call(&params, &resp)
	if err != nil {
		return err
	}

	id := resp.UHostIds[0]
	d.SetId(id)

	log.Printf("[DEBUG] Waiting for instance (%s) to become running", id)

	stateConf := &resource.StateChangeConf{
		Pending:    []string{"Initializing", "Starting"},
		Target:     []string{"Running"},
		Refresh:    instanceRefreshFunc(apiClient, id),
		Timeout:    10 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	instance, err := stateConf.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for instance (%s) to become ready: %s", id, err)
	}

	ip := findInstanceIP(instance.(*client.UHostInstance))
	if ip != "" {
		d.SetConnInfo(map[string]string{
			"type":     "ssh",
			"host":     ip,
			"password": d.Get("password").(string),
		})
	}

	return resourceUHostUpdate(d, meta)
}

func resourceUHostRead(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.Client)

	instance, err := describeInstance(apiClient, d.Id())
	if err != nil {
		return err
	}
	if instance == nil {
		d.SetId("")
		return nil
	}

	if instance.State == "Install Fail" {
		d.SetId("")
		return nil
	}

	setResourceDataFromInstance(d, instance)

	return nil
}

func resourceUHostUpdate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.Client)

	d.Partial(true)
	var resp client.GeneralResponse

	// name
	if name := d.Get("name"); d.HasChange("name") || (d.IsNewResource() && name.(string) != "") {
		params := &client.ModifyUHostInstanceNameRequest{
			UHostId: d.Id(),
			Name:    name.(string),
		}
		err := apiClient.Call(params, &resp)
		if err != nil {
			return err
		}
		d.SetPartial("name")
	}

	// tag
	if tag := d.Get("tag"); d.HasChange("tag") || (d.IsNewResource() && tag.(string) != "") {
		params := &client.ModifyUHostInstanceTagRequest{
			UHostId: d.Id(),
			Tag:     tag.(string),
		}
		err := apiClient.Call(params, &resp)
		if err != nil {
			return err
		}
		d.SetPartial("tag")
	}

	// remark
	if remark := d.Get("remark"); d.HasChange("remark") || (d.IsNewResource() && remark.(string) != "") {
		params := client.ModifyUHostInstanceRemarkRequest{
			UHostId: d.Id(),
			Remark:  remark.(string),
		}
		err := apiClient.Call(&params, &resp)
		if err != nil {
			return err
		}
		d.SetPartial("remark")
	}

	// password
	if d.HasChange("password") {
		params := client.ResetUHostInstancePasswordRequest{
			UHostId:  d.Id(),
			Password: base64.StdEncoding.EncodeToString([]byte(d.Get("password").(string))),
		}
		err := apiClient.Call(&params, &resp)
		if err != nil {
			return err
		}
		d.SetPartial("password")
	}

	// reize: has to restart the host
	// 实例状态， 初始化: Initializing; 启动中: Starting; 运行中: Running; 关机中: Stopping; 关机: Stopped 安装失败: Install Fail; 重启中: Rebooting
	if d.HasChange("cpu") || d.HasChange("memory") || d.HasChange("disk_space") {
		err := apiClient.Call(&client.StopUHostInstanceRequest{UHostId: d.Id()}, &resp)
		if err != nil {
			return err
		}

		stateConf := &resource.StateChangeConf{
			Pending:    []string{"Running", "Stopping", "Rebooting"},
			Target:     []string{"Stopped"},
			Refresh:    instanceRefreshFunc(apiClient, d.Id()),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, err = stateConf.WaitForState()
		if err != nil {
			return err
		}

		params := client.ResizeUHostInstanceRequest{
			UHostId: d.Id(),
		}
		if d.HasChange("cpu") {
			params.CPU = d.Get("cpu").(int)
		}
		if d.HasChange("memory") {
			params.Memory = d.Get("memory").(int)
		}
		if d.HasChange("disk_space") {
			params.DiskSpace = d.Get("disk_space").(int)
		}
		err = apiClient.Call(&params, &resp)
		if err != nil {
			return err
		}

		err = apiClient.Call(&client.StartUHostInstanceRequest{UHostId: d.Id()}, &resp)
		if err != nil {
			return err
		}

		stateConf = &resource.StateChangeConf{
			Pending:    []string{"Stopped", "Stopping", "Rebooting"},
			Target:     []string{"Running"},
			Refresh:    instanceRefreshFunc(apiClient, d.Id()),
			Timeout:    10 * time.Minute,
			Delay:      10 * time.Second,
			MinTimeout: 3 * time.Second,
		}
		_, err = stateConf.WaitForState()
		if err != nil {
			return err
		}

		d.SetPartial("resize")
	}

	d.Partial(false)

	return resourceUHostRead(d, meta)
}

func resourceUHostDelete(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.Client)

	var resp client.GeneralResponse
	params := client.TerminateUHostInstanceRequest{
		UHostId: d.Id(),
	}
	err := apiClient.Call(&params, &resp)
	if err != nil {
		return err
	}

	d.SetId("")

	return nil
}

func describeInstance(apiClient *client.Client, uhostID string) (*client.UHostInstance, error) {
	params := client.DescribeUHostInstanceRequest{
		UHostIds: []string{uhostID},
	}

	var resp client.DescribeUHostInstanceResponse
	err := apiClient.Call(&params, &resp)
	if err != nil {
		return nil, err
	}

	if len(resp.UHostSet) > 0 {
		return &resp.UHostSet[0], nil
	}

	return nil, nil
}

func instanceRefreshFunc(apiClient *client.Client, uhostID string) resource.StateRefreshFunc {
	return func() (interface{}, string, error) {
		instance, err := describeInstance(apiClient, uhostID)
		if err != nil {
			return nil, "", err
		}
		if instance == nil {
			return nil, "", fmt.Errorf("Instance not found")
		}

		return instance, instance.State, nil
	}
}

var ipTypePriority = [...]string{
	"Bgp",
	"Duplet",
	"China-telecom",
	"China-unicom",
	"Internation",
	"Private",
}

func findInstanceIP(instance *client.UHostInstance) string {
	if len(instance.IPSet) > 0 {
		typeToIP := make(map[string]string, len(instance.IPSet))

		for _, ipinfo := range instance.IPSet {
			typeToIP[ipinfo.Type] = ipinfo.IP
		}

		for _, t := range ipTypePriority {
			if ip, ok := typeToIP[t]; ok {
				return ip
			}
		}

		return instance.IPSet[0].IP
	}

	return ""
}

func setResourceDataFromInstance(d *schema.ResourceData, instance *client.UHostInstance) {
	d.Set("uhost_type", instance.UHostType)
	d.Set("storage_type", instance.StorageType)
	d.Set("basic_image_id", instance.BasicImageId)
	d.Set("basic_image_name", instance.BasicImageName)
	d.Set("tag", instance.Tag)
	d.Set("remark", instance.Remark)
	d.Set("name", instance.Name)
	d.Set("charge_type", instance.ChargeType)
	d.Set("cpu", instance.CPU)
	d.Set("memory", instance.Memory)
	d.Set("disk_set", instance.DiskSet)
	d.Set("ip_set", instance.IPSet)
	d.Set("net_capability", instance.NetCapability)
}
