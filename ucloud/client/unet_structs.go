package client

type EIPResource struct {
	ResourceID   string
	ResourceType string
	ResourceName string
	Zone         string
}

type EIPAddr struct {
	OperatorName string
	IP           string
}

type EIP struct {
	EIPId         string
	Weight        int
	BandwidthType int
	Bandwidth     int
	Status        string
	PayMode       string
	ChargeType    string
	CreateTime    int
	ExpireTime    int
	Name          string
	Tag           string
	Remark        string
	Resource      *EIPResource
	EIPAddr       []EIPAddr
}

type AllocateEIPRequest struct {
	OperatorName     string
	Bandwidth        int
	Tag              string
	ChargeType       string
	Quantity         int
	PayMode          string
	ShareBandwidthId string
	Name             string
	Remark           string
}

type AllocateEIPResponse struct {
	GeneralResponse
	EIPSet []EIP
}

type DescribeEIPRequest struct {
	EIPIds []string
	Offset int
	Limit  int
}
type DescribeEIPResponse struct {
	GeneralResponse
	TotalCount     int
	TotalBandwidth int
	EIPSet         []EIP
}

type ReleaseEIPRequest struct {
	EIPId string
}
type ReleaseEIPResponse struct {
	GeneralResponse
}

type BindEIPRequest struct {
	EIPId        string
	ResourceType string
	ResourceId   string
}
type BindEIPResponse struct {
	GeneralResponse
}

type UnBindEIPRequest struct {
	EIPId        string
	ResourceType string
	ResourceId   string
}
type UnBindEIPResponse struct {
	GeneralResponse
}

type UpdateEIPAttributeRequest struct {
	EIPId  string
	Name   string
	Tag    string
	Remark string
}
type UpdateEIPAttributeResponse struct {
	GeneralResponse
}

type ModifyEIPBandwidthRequest struct {
	EIPId     string
	Bandwidth int
}
type ModifyEIPBandwidthResponse struct {
	GeneralResponse
}

type ModifyEIPWeightRequest struct {
	EIPId  string
	Weight int
}
type ModifyEIPWeightResponse struct {
	GeneralResponse
}

type SetEIPPayModeRequest struct {
	EIPId     string
	Bandwidth int
	PayMode   string
}
type SetEIPPayModeResponse struct {
	GeneralResponse
}

type SecurityGroupRule struct {
	ProtocolType string
	DstPort      string
	SrcIP        string
	RuleAction   string
	Priority     int
}

type SecurityGroup struct {
	GroupId     int
	GroupName   string
	Description string
	CreateTime  int
	Type        int
	FirewallId  string
	Name        string
	Tag         string
	Remark      string
	Rule        []SecurityGroupRule
}

type CreateSecurityGroupRequest struct {
	GroupName   string
	Description string
	Rule        []string
}
type CreateSecurityGroupResponse struct {
	GeneralResponse
}

type DescribeSecurityGroupRequest struct {
	ResourceType string
	ResourceId   string
	GroupId      int
}
type DescribeSecurityGroupResponse struct {
	GeneralResponse
	DataSet []SecurityGroup
}

type DescribeSecurityGroupResourceRequest struct {
	GroupId int
}
type DescribeSecurityGroupResourceResponse struct {
	GeneralResponse
	DataSet []string
}

type UpdateSecurityGroupRequest struct {
	GroupId int
	Rule    []string
}
type UpdateSecurityGroupResponse struct {
	GeneralResponse
}

type GrantSecurityGroupRequest struct {
	GroupId      int
	ResourceType string
	ResourceId   string
}
type GrantSecurityGroupResponse struct {
	GeneralResponse
}

type DeleteSecurityGroupRequest struct {
	GroupId int
}
type DeleteSecurityGroupResponse struct {
	GeneralResponse
}
