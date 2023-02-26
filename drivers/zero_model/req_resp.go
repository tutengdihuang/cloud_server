package zero_model

type StartInstancesWithChanRequest struct {
	InstanceIds []string
}
type StartInstancesWithChanResponse struct {
	RequestId  string
	InstanceId string //AWS专用
}

type DescribeInstancesStatusRequest struct {
	InstanceIds []string
}
type DescribeInstancesStatusResponse struct {
	RequestId string
	InsStatus []InstanceStatus
}

type RunInstanceWithChanRequest struct {
	InstanceName       string
	HostName           string
	RegionId           string //对应tencent Plancement,zone
	InstanceType       string
	InstanceChargeType string
	InternetChargeType string
	ImageId            string
	Password           string //对于AWS是key pair的名字
	KeyPairName        string //对于AWS是key pair的名字
	SecurityGroupId    string
	VSwitchId          string
	Amount             int64
	//LaunchTemplateName          string
	InternetMaxBandwidthOut     int64
	ResourceGroupId             string //Ali
	SystemDiskCategory          string
	SystemDiskSize              string
	Tag                         map[string]string
	SecurityEnhancementStrategy string //Ali
}
type RunInstanceWithChanResponse struct {
	RequestId      string
	InstanceIdSets []string
	TradePrice     float64
}

type AllocateEipAddressRequest struct {
	RegionId  string
	Bandwidth int64
}
type AllocateEipAddressResponse struct {
	RequestId    string
	AllocationId string
	EipAddress   string
	//AddressSet   string // tencent
}

type AssociateEipAddressRequest struct {
	RegionId     string
	AllocationId string
	EipAddress   string
	//AddressSet   string // tencent,ip描述符号
	InstanceId string
}
type AssociateEipAddressResponse struct {
	RequestId     string //
	AssociationId string //AWS
}

//请求参数组合1， RegionId，Status，PageNumber，PageSize
//请求参数组合2，RegionId，AllocationId，
//请求参数组合3， RegionId，AssociatedInstanceId，AssociatedInstanceType

type DescribeEipAddressesRequest struct {
	RegionId               string
	AssociatedInstanceId   string
	AssociatedInstanceType string
	Status                 EipStatus
	AllocationIds          []string
	PublicIpsOrIPIds       []string
	//AddressIds             []string //tencent 暂时用PublicIps承载
}
type DescribeEipAddressesResponse struct {
	RequestId string
	EipInfo   []Eip
}

type UnassociateEipAddressRequest struct {
	RegionId              string
	AllocationId          string
	InstanceId            string
	EipAddressOrAddressId string
	AssociationId         string
	//AddressId    string //对于腾讯云EipAddress代替
}
type UnassociateEipAddressResponse struct {
	RequestId string
}

type ReleaseEipAddressRequest struct {
	RegionId     string
	AllocationId string
	//AddressIds   string //对于腾讯云用publicIp 代替
	PublicIpOrAddressId string
}
type ReleaseEipAddressResponse struct {
	RequestId string
}

type ChangeInstancePasswordRequest struct {
	RegionId   string
	InstanceId string
	Password   string // rootpassword
	//OriginPassword string	//测试时候确认
	InstanceAddress string
	LoginName       string //For AWS, 一般为ubuntu, centos 等
	PrivateKey      []byte
}
type ChangeInstancePasswordResponse struct {
	RequestId string
}

type DescribeInstancesAllRequest struct {
	RegionId string
}
type DescribeInstancesAllResponse struct {
	RequestId string
	Instances []DescribeInstance
}

type DescribeInstancesByIDsRequest struct {
	InstanceIds []string
}
type DescribeInstancesByIDsResponse struct {
	RequestId string
	Instances []DescribeInstance
}

type StopInstancesWithChanRequest struct {
	InstanceIds []string
	ForceStop   bool
}
type StopInstancesWithChanResponse struct {
	RequestId string
}

type ReplaceSystemDiskRequest struct {
	InstanceId string
	ImageId    string
	Password   string
}
type ReplaceSystemDiskResponse struct {
	RequestId string
	DiskId    string
}

type DeleteInstancesRequest struct {
	InstanceIds []string
	Force       bool
}
type DeleteInstancesResponse struct {
	RequestId string
}

type DescribeImagesRequest struct {
	RegionId string
	ImageId  string
}
type DescribeImagesResponse struct {
	RequestId string
	ImageSet  []Image
}

type CreateImageRequest struct {
	RegionId     string
	InstanceId   string
	ImageName    string
	ImageVersion string
	Tage         string
}
type CreateImageResponse struct {
	ImageId   string
	RequestId string
}

type CopyImageRequest struct {
	RegionId            string
	Name                string
	ImageId             string
	DestinationRegionId string
}
type CopyImageResponse struct {
	RequestId string
	ImageId   string
}

//RegionId 对于腾讯无用
type DeleteImagesRequest struct {
	RegionId string
	ImageIds []string
	Force    bool
}
type DeleteImagesResponse struct {
	RequestId string
}
