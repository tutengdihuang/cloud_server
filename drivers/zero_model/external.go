package zero_model

type InstanceStatus struct {
	InstanceId string
	Status     InsStatusList
}

type Eip struct {
	AllocationId  string
	EipAddress    string
	Status        EipStatus
	AddressId     string //如果腾讯用到这个字段，这替代其他方法中的public IP
	InstanceId    string
	AssociationId string
}

type DescribeInstance struct {
	InstanceId       string
	InstanceName     string
	Status           InsStatusList
	PublicIpAddress  []string
	PrivateIpAddress []string
}

type Image struct {
	Status  ImageStatus
	ImageId string
}

/*
tecent
PENDING：表示创建中
LAUNCH_FAILED：表示创建失败
RUNNING：表示运行中
STOPPED：表示关机
STARTING：表示开机中
STOPPING：表示关机中
REBOOTING：表示重启中
SHUTDOWN：表示停止待销毁
TERMINATING：表示销毁中。

ali
实例状态。可能值：

Pending：创建中。
Running：运行中。
Starting：启动中。
Stopping：停止中。
Stopped：已停止。


aws

pending
running
shutting-down
terminated
stopping
stopped
*/

type InsStatusList int64

const (
	InsStatus_UNKNOWN InsStatusList = iota
	InsStatus_PENDING
	InsStatus_LAUNCH_FAILED
	InsStatus_RUNNING
	InsStatus_STARTING
	InsStatus_STOPPING
	InsStatus_STOPPED
	InsStatus_REBOOTING
	InsStatus_SHUTING_DOWN
	InsStatus_SHUTDOWN
	InsStatus_TERMINATING
	InsStatus_TERMINATED
)

var TencentInstanceStatusMapToInsStatusList = map[string]InsStatusList{
	"PENDING":       InsStatus_PENDING,
	"LAUNCH_FAILED": InsStatus_LAUNCH_FAILED,
	"RUNNING":       InsStatus_RUNNING,
	"STOPPED":       InsStatus_STOPPED,
	"STARTING":      InsStatus_STARTING,
	"STOPPING":      InsStatus_STOPPING,
	"REBOOTING":     InsStatus_REBOOTING,
	"SHUTDOWN":      InsStatus_SHUTDOWN,
	"TERMINATING":   InsStatus_TERMINATING,
}

var AlyInstanceStatusMapToInsStatusList = map[string]InsStatusList{
	"Pending":  InsStatus_PENDING,
	"Running":  InsStatus_RUNNING,
	"Starting": InsStatus_STARTING,
	"Stopping": InsStatus_STOPPING,
	"Stopped":  InsStatus_STOPPED,
}
var AWSInstanceStatusMapToInsStatusList = map[string]InsStatusList{
	"pending":       InsStatus_PENDING,
	"running":       InsStatus_RUNNING,
	"shutting-down": InsStatus_SHUTING_DOWN,
	"terminated":    InsStatus_TERMINATED,
	"stopping":      InsStatus_STOPPING,
	"stopped":       InsStatus_STOPPED,
}

/*
//Tencent
'CREATING'，'BINDING'，'BIND'，'UNBINDING'，'UNBIND'，'OFFLINING'，'BIND_ENI'
*/
type EipStatus int64

const (
	//Aly
	EIPStATUS_UNKNOWN  EipStatus = 0
	EIPStATUS_CREATING EipStatus = iota << 1
	EIPStATUS_ASSOCIATING
	EIPStATUS_UNASSOCIATING
	EIPStATUS_INUSE
	EIPStATUS_AVAILABLE
	EIPStATUS_RELEASING
	EIPStATUS_BIND_ENI
	//Tencent
	//CREATING
	//BINDING=ASSOCIATING
	//BIND=INUSE
	//UNBINDING=Unassociating
	//UNBIND=Available
	// OFFLINING=Releasing
	//BIND_ENI
)

var TencentMapStrToEipStatus = map[string]EipStatus{
	"CREATING":  EIPStATUS_CREATING,
	"BINDING":   EIPStATUS_ASSOCIATING,
	"BIND":      EIPStATUS_INUSE,
	"UNBINDING": EIPStATUS_UNASSOCIATING,
	"UNBIND":    EIPStATUS_AVAILABLE,
	"OFFLINING": EIPStATUS_RELEASING,
	"BIND_ENI":  EIPStATUS_BIND_ENI,
}

var TencentMapEipStatusToStr = map[EipStatus]string{
	EIPStATUS_CREATING:      "CREATING",
	EIPStATUS_ASSOCIATING:   "BINDING",
	EIPStATUS_INUSE:         "BIND",
	EIPStATUS_UNASSOCIATING: "UNBINDING",
	EIPStATUS_AVAILABLE:     "UNBIND",
	EIPStATUS_RELEASING:     "OFFLINING",
	EIPStATUS_BIND_ENI:      "BIND_ENI",
}

var AlyEipMapStrToEipStatus = map[string]EipStatus{
	"ASSOCIATING":   EIPStATUS_ASSOCIATING,
	"Unassociating": EIPStATUS_UNASSOCIATING,
	"INUSE":         EIPStATUS_INUSE,
	"Available":     EIPStATUS_AVAILABLE,
	"Releasing":     EIPStATUS_RELEASING,
}

var AlyEipMapEipStatusToStr = map[EipStatus]string{
	EIPStATUS_ASSOCIATING:   "ASSOCIATING",
	EIPStATUS_UNASSOCIATING: "Unassociating",
	EIPStATUS_INUSE:         "INUSE",
	EIPStATUS_AVAILABLE:     "Available",
	EIPStATUS_RELEASING:     "Releasing",
}

//tecent
// 镜像状态:
// CREATING-创建中
// NORMAL-正常
// CREATEFAILED-创建失败
// USING-使用中
// SYNCING-同步中
// IMPORTING-导入中
// IMPORTFAILED-导入失败

//ali
/*
Creating：镜像正在创建中。//ok
Waiting：多任务排队中。	//;ok
Available（默认）：您可以使用的镜像。	 //OK
UnAvailable：您不能使用的镜像。		//OK
CreateFailed：创建失败的镜像。 //ok
Deprecated：已弃用的镜像。	//ok
*/

//aws_image
/*
pending
available
invalid
deregistered
transient
failed
error
*/

type ImageStatus int

const (
	IMAGE_STATUS_UNKNOWN ImageStatus = iota
	IMAGE_STATUS_CREATING
	IMAGE_STATUS_CREATEFAILED
	IMAGE_STATUS_AVAILABLE   // 阿里 Available,aws：available
	IMAGE_STATUS_UNAVAILABLE // 阿里， aws:invalid
	IMAGE_STATUS_USING       //
	IMAGE_STATUS_SYNCING
	IMAGE_STATUS_IMPORTING
	IMAGE_STATUS_IMPORTFAILED
	IMAGE_STATUS_WAITING
	IMAGE_STATUS_DEPRECATED // aws:IMAGE_STATUS_DEREGISTERD

	//AWS
	IMAGE_STATUS_PNEDING
	//IMAGE_STATUS_DEREGISTERD
	IMAGE_STATUS_TRANSIENT
	IMAGE_STATUS_FAILED
	IMAGE_STATUS_ERROR
)

var AWSIMageStatusStrToImageStatus = map[string]ImageStatus{
	"pending":      IMAGE_STATUS_PNEDING,
	"available":    IMAGE_STATUS_AVAILABLE,
	"invalid":      IMAGE_STATUS_UNAVAILABLE,
	"deregistered": IMAGE_STATUS_DEPRECATED,
	"transient":    IMAGE_STATUS_TRANSIENT,
	"failed":       IMAGE_STATUS_FAILED,
	"error":        IMAGE_STATUS_ERROR,
}

var AliIMageStatusStrToImageStatus = map[string]ImageStatus{
	"Creating":     IMAGE_STATUS_CREATING,
	"Waiting":      IMAGE_STATUS_WAITING,
	"Available":    IMAGE_STATUS_AVAILABLE,
	"UnAvailable":  IMAGE_STATUS_UNAVAILABLE,
	"CreateFailed": IMAGE_STATUS_CREATEFAILED,
	"Deprecated":   IMAGE_STATUS_DEPRECATED,
}

var TencentIMageStatusStrToImageStatus = map[string]ImageStatus{
	"CREATING":     IMAGE_STATUS_CREATING,
	"NORMAL":       IMAGE_STATUS_AVAILABLE,
	"CREATEFAILED": IMAGE_STATUS_CREATEFAILED,
	"USING":        IMAGE_STATUS_USING,
	"SYNCING":      IMAGE_STATUS_SYNCING,
	"IMPORTING":    IMAGE_STATUS_IMPORTING,
	"IMPORTFAILED": IMAGE_STATUS_IMPORTFAILED,
}
