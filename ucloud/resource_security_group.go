package ucloud

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/3pjgames/terraform-provider-ucloud/ucloud/client"
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceSecurityGroup() *schema.Resource {
	return &schema.Resource{
		Create: resourceSecurityGroupCreate,
		Read:   resourceSecurityGroupRead,
		Update: resourceSecurityGroupUpdate,
		Delete: resourceSecurityGroupDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"group_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
					value := v.(string)
					r := regexp.MustCompile(`^[\p{Han}a-zA-Z0-9_.-]+$`)
					if !r.MatchString(value) {
						errors = append(errors, fmt.Errorf("group_name 只能包含中英文、数字以及- _ ."))
					}

					return
				},
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"create_time": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rule": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"protocol_type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"dst_port": {
							Type:     schema.TypeString,
							Required: true,
						},
						"src_ip": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "0.0.0.0/0",
						},
						"rule_action": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "ACCEPT",
						},
						"priority": {
							Type:     schema.TypeInt,
							Optional: true,
							Default:  50,
							ValidateFunc: func(v interface{}, k string) (ws []string, errors []error) {
								value := v.(int)
								if value != 50 && value != 100 && value != 150 {
									errors = append(errors, fmt.Errorf("priroty can only be 50 (high), 100 (medium) and 150 (low)"))
								}

								return
							},
						},
					},
				},
			},
		},
	}
}

func resourceSecurityGroupCreate(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*client.Client)

	req := client.CreateSecurityGroupRequest{
		GroupName:   d.Get("group_name").(string),
		Description: d.Get("description").(string),
	}

	if v := d.Get("rule"); v != nil {
		req.Rule = buildRuleSlice(d)
	}

	var resp client.CreateSecurityGroupResponse
	err := api.Call(&req, &resp)
	if err != nil {
		return err
	}

	id, err := findSecurityGroupIdByName(api, req.GroupName)
	if err != nil {
		return err
	}

	d.SetId(strconv.Itoa(id))

	return nil
}

func resourceSecurityGroupRead(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	req := client.DescribeSecurityGroupRequest{GroupId: id}
	var resp client.DescribeOneSecurityGroupResponse
	err = api.Call(&req, &resp)
	if err != nil {
		if brce, ok := err.(*client.BadRetCodeError); ok && brce.RetCode == 4351 { // get security group fail
			d.SetId("")
			return nil
		}
		return err
	}

	group := resp.DataSet

	d.Set("group_name", group.GroupName)
	d.Set("description", group.Description)
	d.Set("create_time", group.CreateTime)
	d.Set("rule", readRule(group))

	return nil
}

func resourceSecurityGroupUpdate(d *schema.ResourceData, meta interface{}) error {
	if !d.HasChange("rule") {
		return nil
	}

	api := meta.(*client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	req := client.UpdateSecurityGroupRequest{GroupId: id}
	if v := d.Get("rule"); v != nil {
		req.Rule = buildRuleSlice(d)
	}
	var resp client.UpdateSecurityGroupResponse
	err = api.Call(&req, &resp)
	if err != nil {
		return err
	}

	return nil
}

func resourceSecurityGroupDelete(d *schema.ResourceData, meta interface{}) error {
	api := meta.(*client.Client)

	id, err := strconv.Atoi(d.Id())
	if err != nil {
		return err
	}

	req := client.DeleteSecurityGroupRequest{GroupId: id}
	var resp client.GeneralResponse
	err = api.Call(&req, &resp)
	if err != nil {
		return err
	}

	d.SetId("")
	return nil
}

func buildRuleSlice(d *schema.ResourceData) []client.SecurityGroupRule {
	if v := d.Get("rule"); v != nil {
		ruleList := v.([]interface{})
		var slice = make([]client.SecurityGroupRule, 0, len(ruleList))
		for _, v := range ruleList {
			rule := v.(map[string]interface{})
			slice = append(slice, client.SecurityGroupRule{
				ProtocolType: rule["protocol_type"].(string),
				DstPort:      rule["dst_port"].(string),
				SrcIP:        rule["src_ip"].(string),
				RuleAction:   rule["rule_action"].(string),
				Priority:     rule["priority"].(int),
			})
		}

		return slice
	}

	return []client.SecurityGroupRule{}
}

func readRule(group *client.SecurityGroup) interface{} {
	rule := make([]map[string]interface{}, 0, len(group.Rule))

	for _, v := range group.Rule {
		rule = append(rule, map[string]interface{}{
			"protocol_type": v.ProtocolType,
			"dst_port":      v.DstPort,
			"src_ip":        v.SrcIP,
			"rule_action":   v.RuleAction,
			"priority":      v.Priority,
		})
	}

	return rule
}

func findSecurityGroupIdByName(api *client.Client, n string) (int, error) {
	var resp client.DescribeSecurityGroupResponse
	err := api.Call(&client.DescribeSecurityGroupRequest{}, &resp)
	if err != nil {
		return 0, err
	}

	for _, v := range resp.DataSet {
		if v.GroupName == n {
			return v.GroupId, nil
		}
	}

	return 0, fmt.Errorf("Cannot found security group with name %s", n)
}
