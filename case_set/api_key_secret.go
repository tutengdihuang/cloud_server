package case_set

import (
	"cloud_server"
	"cloud_server/core"
	"cloud_server/drivers/zero_model"
)

var serverTencent core.Operator
var serverAWS core.Operator

var serverAliyun core.Operator

func init() {
	cfg := &zero_model.CloudServerConfig{}
	cfg.Region = "cn-hangzhou"
	cfg.AccessKey = Aliyun_accessKeyID
	cfg.SecretKey = Aliyun_accessKeySecret
	cfg.Driver = zero_model.Aliyun
	var err error
	serverAliyun, err = cloud_server.New(cfg)
	if err != nil {
		panic(err)
	}
}

func init() {
	cfg := &zero_model.CloudServerConfig{}
	cfg.Region = "ap-east-1"
	cfg.AccessKey = AWS_accessKeyID
	cfg.SecretKey = AWS_accessKeySecret
	cfg.Driver = zero_model.AWS
	var err error
	serverAWS, err = cloud_server.New(cfg)
	if err != nil {
		panic(err)
	}
}

func init() {
	cfg := &zero_model.CloudServerConfig{}
	cfg.Region = "ap-beijing"
	cfg.AccessKey = Tecent_accessKeyID
	cfg.SecretKey = Tecent_accessKeySecret
	cfg.Driver = zero_model.Tencent
	var err error
	serverTencent, err = cloud_server.New(cfg)
	if err != nil {
		panic(err)
	}
}
