package client

type UHostIP struct {
	Type      string
	IPId      string
	IP        string
	Bandwidth int
}

type UHostDisk struct {
	Type   string
	DiskId string
	Drive  string
	Size   int
}

type UHostInstance struct {
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
	DiskSet        []UHostDisk
	IPSet          []UHostIP
	NetCapability  string
	NetworkState   string
}

type DescribeUHostInstanceRequest struct {
	Zone     string
	UHostIds []string
	Tag      string
	Offset   int
	Limit    int
}

type DescribeUHostInstanceResponse struct {
	GeneralResponse
	TotalCount int
	UHostSet   []UHostInstance
}

type CreateUHostInstanceRequest struct {
	Zone            string
	ImageId         string
	LoginMode       string
	Password        string
	KeyPair         string
	CPU             int
	Memory          int
	StorageType     string
	DiskSpace       int
	Name            string
	NetworkId       string
	SecurityGroupId string
	ChargeType      string
	Quantity        int
	UHostType       string
	NetCapability   string
	Tag             string
	CouponId        string
	ProjectId       int
	BootDiskSpace   int
}

type CreateUHostInstanceResponse struct {
	GeneralResponse
	UHostIds []string
}

type TerminateUHostInstanceRequest struct {
	UHostId string
	Zone    string
}
type TerminateUHostInstanceResponse struct {
	GeneralResponse
}

type StopUHostInstanceRequest struct {
	UHostId string
	Zone    string
}
type StopUHostInstanceResponse struct {
	GeneralResponse
}

type StartUHostInstanceRequest struct {
	UHostId string
	Zone    string
}
type StartUHostInstanceResponse struct {
	GeneralResponse
}

type ModifyUHostInstanceRemarkRequest struct {
	UHostId string
	Zone    string
	Remark  string
}
type ModifyUHostInstanceRemarkResponse struct {
	GeneralResponse
}

type ModifyUHostInstanceNameRequest struct {
	UHostId string
	Zone    string
	Name    string
}
type ModifyUHostInstanceNameResponse struct {
	GeneralResponse
}

type ModifyUHostInstanceTagRequest struct {
	UHostId string
	Zone    string
	Tag     string
}
type ModifyUHostInstanceTagResponse struct {
	GeneralResponse
}

type ResizeUHostInstanceRequest struct {
	UHostId   string
	Zone      string
	CPU       int
	Memory    int
	DiskSpace int
}
type ResizeUHostInstanceResponse struct {
	GeneralResponse
}

type ResetUHostInstancePasswordRequest struct {
	UHostId  string
	Zone     string
	Password string
}
type ResetUHostInstancePasswordResponse struct {
	GeneralResponse
}

type UHostImage struct {
	ImageId          string
	Zone             string
	ImageName        string
	ImageType        string
	ImageSize        int
	OsType           string
	OsName           string
	State            string
	ImageDescription string
	CreateTime       int
	Features         []string
}

type DescribeImageRequest struct {
	Zone      string
	ImageType string
	OsType    string
	ImageId   string
	Offset    int
	Limit     int
}
type DescribeImageResponse struct {
	GeneralResponse
	TotalCount int
	ImageSet   []*UHostImage
}
