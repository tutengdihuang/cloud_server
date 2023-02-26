## 简介
- golang开发版本 go1.18.2
- 本项目主要提供多云服务器sdk服务 目前支持阿里云，腾讯云，AWS
- 支持的接口
```golang
    StartInstancesWithChan
    RunInstancesWithChan
    StopInstancesWithChan
    DescribeInstancesStatusWithChan
    AllocateEipAddressWithChan
    AssociateEipAddressWithChan
    DescribeEipAddressesWithChan
    UnassociateEipAddressWithChan
    ReleaseEipAddressWithChan
    ChangeInstancePasswordWithChan
    DescribeInstancesAllWithChan
    DescribeInstancesByIDsWithChan
    ReplaceSystemDiskWithChan
    DeleteInstancesWithChan
    DescribeImagesWithChan
    CreateImageWithChan
    CopyImageWithChan
    DeleteImageWithChan
```

## 使用
- 可以参考case_set文件夹中的测试用例
### 以阿里云为例子
- 初始化变量对象serverAliyun
```golang
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
```
- 调用并创建服务器实例
```golang
	request := new(zero_model.RunInstanceWithChanRequest)
	request.RegionId = "cn-hangzhou" //placement
	request.InstanceName = "cn-beijing-allen-test-3"

	request.ImageId = "aliyun_3_x64_20G_alibase_20221102.vhd"
	request.Password = "Allen123!"
	request.SecurityGroupId = "sg-bp14o9edhrbueomfc1f2"
	request.VSwitchId = "vsw-bp1tw0ukcxmwt2vvp2kcl"
	request.InstanceType = "ecs.s6-c1m1.small"
	request.InstanceChargeType = "PostPaid"
	//request.InternetChargeType = "PayByBandwidth"
	//request.InternetMaxBandwidthOut = 1
	request.Amount = 1
	request.SystemDiskSize = "40"
	//request.SystemDiskCategory = "cloud_essd"

	if request == nil {
		panic("request is nil")
	}
	if serverAliyun == nil {
		panic("serverAliyun is nil")
	}
	respChan, errChan := serverAliyun.RunInstancesWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("resp: %+v", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
		if err == nil && cap(respChan) > 0 {
			resp := <-respChan
			fmt.Printf("\n resp RequestId: %+v \n", resp.RequestId)
		}
	}
```


## 参考

### tencent
- [腾讯云sdk](https://github.com/TencentCloud/tencentcloud-sdk-go/blob/master/tencentcloud/cvm/v20170312/client.go)
- [tencent api](https://cloud.tencent.com/document/api/213/15689)
- [Tencent cloud object storage go sdk](https://cloud.tencent.com/document/product/436/65647)
### ali
- [aliyun API](https://help.aliyun.com/document_detail/25485.html)
- [ali vpc API](https://help.aliyun.com/document_detail/36016.html?spm=5176.21213303.J_6704733920.7.2f2453c91gEAJJ&scm=20140722.S_help%40%40%E6%96%87%E6%A1%A3%40%4036016._.ID_help%40%40%E6%96%87%E6%A1%A3%40%4036016-RL_allocateeipaddress-LOC_main-OR_ser-V_2-P0_0)
- [oss API](https://help.aliyun.com/document_detail/31946.html)
- [沙箱运行](https://next.api.aliyun.com/api/Ecs/2014-05-26/DescribeAvailableResource?spm=api-workbench.Troubleshoot.0.0.220571859E44dy&lang=GO&params={%22RegionId%22:%22cn-beijing-h%22,%22InstanceChargeType%22:%22PostPaid%22,%22DestinationResource%22:%22Zone%22}&sdkStyle=old)
- [vswitch](https://vpc.console.aliyun.com/vpc/cn-beijing/switches)
- [eip_控制台](https://vpc.console.aliyun.com/eip/cn-hangzhou/eips)
### AWS
- [S3API](https://docs.aws.amazon.com/zh_cn/AmazonS3/latest/API/API_CreateBucket.html)
- [aws-api](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DeleteLaunchTemplate.html)
- [sdk1](https://docs.aws.amazon.com/zh_tw/sdk-for-go/v1/developer-guide/using-ec2-with-go-sdk.html)
- [aws sdk](https://aws.amazon.com/cn/developer/tools/)
- [aws-go-sdk-ec2](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#EC2)
- [github-go-sdk-ec2-example](https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/go/example_code/ec2)
- [官网-example](https://docs.aws.amazon.com/sdk-for-go/api/service/ec2/#pkg-examples)
- [aws按需计费规则](https://aws.amazon.com/cn/ec2/pricing/on-demand/)


