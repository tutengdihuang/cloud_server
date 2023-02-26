package zero_model

type DriverType string

const (
	Aliyun  = "aliyun"
	Tencent = "tencent"
	AWS     = "aws"
)

type CloudServerConfig struct {
	Driver    DriverType
	Region    string
	AccessKey string
	SecretKey string
}
