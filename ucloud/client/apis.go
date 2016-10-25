package client

import (
	"fmt"
	"net/url"
	"reflect"
	"strconv"
	"strings"
)

func safeIsNil(val reflect.Value) bool {
	switch val.Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Map, reflect.Chan:
		return val.IsNil()
	default:
		return false
	}
}

func stringify(v reflect.Value) (str string, ok bool) {
	ok = false

	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		ok = true
		if number := v.Int(); number != 0 {
			str = strconv.FormatInt(number, 10)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		ok = true
		if number := v.Uint(); number != 0 {
			str = strconv.FormatUint(number, 10)
		}

	case reflect.String:
		ok = true
		str = v.String()
	}

	return
}

func BuildParams(req interface{}) (url.Values, error) {
	ret := url.Values{}
	err := AddParams(ret, req)
	return ret, err
}

func AddParams(params url.Values, req interface{}) error {
	val := reflect.ValueOf(req)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	valType := val.Type()
	typeName := valType.Name()
	action := strings.TrimSuffix(typeName, "Request")
	if len(action) != len(typeName) {
		params.Set("Action", action)
	}

	for i := 0; i < val.NumField(); i++ {
		typeField := valType.Field(i)
		field := val.Field(i)
		if safeIsNil(field) {
			continue
		}

		if field.Kind() == reflect.Ptr {
			field = field.Elem()
		}
		fieldKind := field.Kind()
		fieldName := typeField.Tag.Get("ArgName")
		if fieldName == "" {
			fieldName = typeField.Name
		}

		if fieldValue, ok := stringify(field); ok {
			if fieldValue != "" {
				params.Set(fieldName, fieldValue)
			}
		} else if fieldKind == reflect.Slice {
			for j := 0; j < field.Len(); j++ {
				elem := field.Index(j)
				if elemValue, ok := stringify(elem); ok {
					params.Set(fieldName+"."+strconv.Itoa(j), elemValue)
				} else {
					return fmt.Errorf("Cannot convert %s to params in slice", elem.Kind())
				}
			}
		} else {
			return fmt.Errorf("Cannot convert %s to params", fieldKind)
		}
	}

	return nil
}

type GeneralResponse struct {
	Action  string
	RetCode int
	Message string
}

func (resp *GeneralResponse) ValidateResponse() error {
	if resp.RetCode != 0 {
		return &BadRetCodeError{
			Action:  resp.Action,
			RetCode: resp.RetCode,
			Message: resp.Message,
		}
	}

	return nil
}

type UHostIP struct {
	Type      string
	IPId      string
	IP        string
	bandwidth int
}

type UHostDisk struct {
	Type   string
	DiskId string
	Name   int
	Drive  int
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
