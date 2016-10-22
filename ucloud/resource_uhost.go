package ucloud

import (
	"encoding/base64"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"

	"github.com/3pjgames/terraform-provider-ucloud/ucloud/client"
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
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
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

type uhostIP struct {
	Type      string
	IPId      string
	IP        string
	bandwidth int
}

type uhostDisk struct {
	Type   string
	DiskId string
	Name   int
	Drive  int
	Size   int
}

type uhostInstance struct {
	UHostId        string
	UHostType      string
	Zone           string
	StorageType    string
	ImageId        string
	BasicImageId   string
	BasicImageName string
	Tag            string
	Remark         string
	Name           string
	State          string
	CreateTime     int
	ChargeType     string
	ExpireTime     int
	CPU            int
	Memory         int
	AutoRenew      string
	DiskSet        []uhostDisk
	IPSet          []uhostIP
	NetCapability  string
	NetworkState   string
}

type generalResponse struct {
	Action  string
	RetCode int
}

type describeUHostResponse struct {
	generalResponse
	TotalCount int
	UHostSet   []uhostInstance
}

type createUHostResponse struct {
	generalResponse
	UHostIds []string
}

func resourceUHostCreate(d *schema.ResourceData, meta interface{}) error {
	apiClient := meta.(*client.Client)

	params := url.Values{}
	params.Set("Action", "CreateUHostInstance")
	params.Set("Zone", d.Get("zone").(string))
	params.Set("ImageId", d.Get("image_id").(string))
	params.Set("LoginMode", "Password")
	params.Set("Password", base64.StdEncoding.EncodeToString([]byte(d.Get("password").(string))))

	if v, ok := d.GetOk("cpu"); ok {
		params.Set("CPU", strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("memory"); ok {
		params.Set("Memory", strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("storage_type"); ok {
		params.Set("StorageType", v.(string))
	}
	if v, ok := d.GetOk("disk_space"); ok {
		params.Set("DiskSpace", strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("name"); ok {
		params.Set("Name", v.(string))
	}
	if v, ok := d.GetOk("network_id"); ok {
		params.Set("NetworkId", v.(string))
	}
	if v, ok := d.GetOk("security_group_id"); ok {
		params.Set("SecurityGroupId", v.(string))
	}
	if v, ok := d.GetOk("charge_type"); ok {
		params.Set("ChargeType", v.(string))
	}
	if v, ok := d.GetOk("quantity"); ok {
		params.Set("Quantity", strconv.Itoa(v.(int)))
	}
	if v, ok := d.GetOk("uhost_type"); ok {
		params.Set("UHostType", v.(string))
	}
	if v, ok := d.GetOk("net_capability"); ok {
		params.Set("NetCapability", v.(string))
	}
	if v, ok := d.GetOk("tag"); ok {
		params.Set("Tag", v.(string))
	}
	if v, ok := d.GetOk("boot_disk_space"); ok {
		params.Set("BootDiskSpace", strconv.Itoa(v.(int)))
	}

	log.Printf("[DEBUG] Run configuration: %s", params)

	var respBody createUHostResponse
	err := apiClient.GetJSON(params, &respBody)
	if err != nil {
		return err
	}

	if respBody.RetCode != 0 {
		return client.BadRetCodeError{"CreateUHostInstance", respBody.RetCode}
	}

	id := respBody.UHostIds[0]
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

	ip := findInstanceIP(instance.(*uhostInstance))
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
	var resp generalResponse

	// tag
	if tag := d.Get("tag"); d.HasChange("tag") || (d.IsNewResource() && tag.(string) != "") {
		params := url.Values{}
		params.Set("Action", "ModifyUHostInstanceTag")
		params.Set("UHostId", d.Id())
		params.Set("Tag", tag.(string))
		err := apiClient.GetJSON(params, &resp)
		if err != nil {
			return err
		}
		if resp.RetCode != 0 {
			return client.BadRetCodeError{"ModifyUHostInstanceTag", resp.RetCode}
		}
		d.SetPartial("tag")
	}

	// remark
	if remark := d.Get("remark"); d.HasChange("remark") || (d.IsNewResource() && remark.(string) != "") {
		params := url.Values{}
		params.Set("Action", "ModifyUHostInstanceRemark")
		params.Set("UHostId", d.Id())
		params.Set("Remark", remark.(string))
		err := apiClient.GetJSON(params, &resp)
		if err != nil {
			return err
		}
		if resp.RetCode != 0 {
			return client.BadRetCodeError{"ModifyUHostInstanceRemark", resp.RetCode}
		}
		d.SetPartial("remark")
	}

	// password
	if d.HasChange("password") {
		params := url.Values{}
		params.Set("Action", "ResetUHostInstancePassword")
		params.Set("UHostId", d.Id())
		params.Set("Password", base64.StdEncoding.EncodeToString([]byte(d.Get("password").(string))))
		err := apiClient.GetJSON(params, &resp)
		if err != nil {
			return err
		}
		if resp.RetCode != 0 {
			return client.BadRetCodeError{"ModifyUHostInstanceRemark", resp.RetCode}
		}
		d.SetPartial("password")
	}

	// reize: has to restart the host
	// 实例状态， 初始化: Initializing; 启动中: Starting; 运行中: Running; 关机中: Stopping; 关机: Stopped 安装失败: Install Fail; 重启中: Rebooting
	if d.HasChange("cpu") || d.HasChange("memory") || d.HasChange("disk_space") {
		stopUHostInstanceParmas := url.Values{}
		stopUHostInstanceParmas.Set("Action", "StopUHostInstance")
		stopUHostInstanceParmas.Set("UHostId", d.Id())
		err := apiClient.GetJSON(stopUHostInstanceParmas, &resp)
		if err != nil {
			return err
		}
		if resp.RetCode != 0 {
			return client.BadRetCodeError{"StopUHostInstance", resp.RetCode}
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

		resizeUHostInstanceParams := url.Values{}
		resizeUHostInstanceParams.Set("Action", "ResizeUHostInstance")
		resizeUHostInstanceParams.Set("UHostId", d.Id())
		resizeUHostInstanceParams.Set("CPU", strconv.Itoa(d.Get("cpu").(int)))
		resizeUHostInstanceParams.Set("Memory", strconv.Itoa(d.Get("memory").(int)))
		resizeUHostInstanceParams.Set("DiskSpace", strconv.Itoa(d.Get("disk_space").(int)))
		err = apiClient.GetJSON(resizeUHostInstanceParams, &resp)
		if err != nil {
			return err
		}
		if resp.RetCode != 0 {
			return client.BadRetCodeError{"ResizeUHostInstance", resp.RetCode}
		}

		startUHostInstanceParmas := url.Values{}
		startUHostInstanceParmas.Set("Action", "StartUHostInstance")
		startUHostInstanceParmas.Set("UHostId", d.Id())
		err = apiClient.GetJSON(startUHostInstanceParmas, &resp)
		if err != nil {
			return err
		}
		if resp.RetCode != 0 {
			return client.BadRetCodeError{"StartUHostInstance", resp.RetCode}
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
	var resp generalResponse

	terminateUHostInstanceParmas := url.Values{}
	terminateUHostInstanceParmas.Set("Action", "TerminateUHostInstance")
	terminateUHostInstanceParmas.Set("UHostId", d.Id())
	err := apiClient.GetJSON(terminateUHostInstanceParmas, &resp)
	if err != nil {
		return err
	}
	if resp.RetCode != 0 {
		return client.BadRetCodeError{"TerminateUHostInstance", resp.RetCode}
	}

	d.SetId("")

	return nil
}

func describeInstance(apiClient *client.Client, uhostID string) (*uhostInstance, error) {
	params := url.Values{}
	params.Set("Action", "DescribeUHostInstance")
	params.Set("UHostIds.0", uhostID)

	var respBody describeUHostResponse
	err := apiClient.GetJSON(params, &respBody)
	if err != nil {
		return nil, err
	}
	if respBody.RetCode != 0 {
		return nil, client.BadRetCodeError{"DescribeUHostInstance", respBody.RetCode}
	}

	if len(respBody.UHostSet) > 0 {
		return &respBody.UHostSet[0], nil
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

func findInstanceIP(instance *uhostInstance) string {
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

func setResourceDataFromInstance(d *schema.ResourceData, instance *uhostInstance) {
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
