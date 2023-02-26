package case_set

import (
	"cloud_server/drivers/zero_model"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

//如果需要hostname 需要注意格式
// 1.aws返回格式中返回参数中没有requestid
// 2.aws的Instance name 在tag中声明，此接口sdk，还是在request.InstanceName中声明，接口中进行转换
// 3. regin id 跟随账户的RegionId
func Test_AWS_Runinstance(t *testing.T) {
	request := new(zero_model.RunInstanceWithChanRequest)
	//request.RegionId = "us-east-1" //placement
	request.InstanceName = "AWS-ALLEN-TEST-1"

	request.ImageId = "ami-01384c3926bc6f380"
	request.KeyPairName = "id_rsa_pub"
	request.SecurityGroupId = "sg-0510cc1f53cb2ca06"
	//request.VSwitchId = "vsw-bp1tw0ukcxmwt2vvp2kcl"
	request.InstanceType = "t3.micro"
	request.InstanceChargeType = "POSTPAID_BY_HOUR"
	request.InternetChargeType = "TRAFFIC_POSTPAID_BY_HOUR"
	request.InternetMaxBandwidthOut = 1
	request.Amount = 1
	request.SystemDiskSize = "20"
	request.SystemDiskCategory = "gp2"

	if request == nil {
		panic("request is nil")
	}
	if serverAWS == nil {
		panic("serverAWS is nil")
	}
	respChan, errChan := serverAWS.RunInstancesWithChan(request)
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
	time.Sleep(1 * time.Second)
}

//1.停止操作不是同步的，返回结果后，实例可能在停止过程中
func Test_AWS_Stopinstance(t *testing.T) {
	request := new(zero_model.StopInstancesWithChanRequest)
	request.ForceStop = true
	request.InstanceIds = append(request.InstanceIds, "i-0fef9450cc8c7b485")
	respChan, errChan := serverAWS.StopInstancesWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("resp RequestId: %+v", resp.RequestId)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
}

//1. 异步执行 同stop
func Test_AWS_Startinstance(t *testing.T) {
	request := new(zero_model.StartInstancesWithChanRequest)
	request.InstanceIds = append(request.InstanceIds, "i-0fef9450cc8c7b485")
	respChan, errChan := serverAWS.StartInstancesWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("resp RequestId: %+v", resp.RequestId)
		fmt.Printf("resp InstanceId: %+v", resp.InstanceId)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
		if err == nil && cap(respChan) > 0 {
			resp := <-respChan
			fmt.Printf("\n resp RequestId: %+v \n", resp.RequestId)
		}
	}
}

//不指定InstanceIds是获取所有实例状态
// aws 的问题，多个实例只返回最后一个实例信息
// 接口代码中MaxResult 值>=5
// 接口只返回running状态的
// 如果想获取所有实例的状态信息需要用DescribeAllInstances
func Test_AWS_DescribeInstancesStatus(t *testing.T) {
	request := new(zero_model.DescribeInstancesStatusRequest)
	request.InstanceIds = append(request.InstanceIds, "i-0207bee7aa31984d3")
	//request.InstanceIds = append(request.InstanceIds, "i-019dbda6cd752f583")
	respChan, errChan := serverAWS.DescribeInstancesStatusWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\nresp RequestId1: %+v \n", resp.RequestId)
		fmt.Printf("resp InsStatus: %+v \n", resp.InsStatus)
	case err := <-errChan:
		fmt.Printf("err====: %+v", err)
	}
	time.Sleep(1 * time.Second)
}
func Test_AWS_DescribeAllInstances(t *testing.T) {
	request := new(zero_model.DescribeInstancesAllRequest)
	request.RegionId = "ap-beijing"
	respChan, errChan := serverAWS.DescribeInstancesAllWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("resp : %#v \n", len(resp.Instances))
		for _, v := range resp.Instances {
			fmt.Printf("resp : %#v \n", v)

		}
	case err := <-errChan:
		fmt.Printf("err====: %#v", err)
	}
	time.Sleep(1 * time.Second)
}

// 此接口不是批量请求的，原因是aws请求多个实例id时，只返回最后一个的信息
func Test_AWS_DescribeInstancesByIDS(t *testing.T) {
	request := new(zero_model.DescribeInstancesByIDsRequest)
	request.InstanceIds = append(request.InstanceIds, "i-0207bee7aa31984d3")
	request.InstanceIds = append(request.InstanceIds, "i-0effc1ad5734ba785")
	request.InstanceIds = append(request.InstanceIds, "i-0e54cb2665700c979")
	request.InstanceIds = append(request.InstanceIds, "i-04e7072b702ad462d")
	request.InstanceIds = append(request.InstanceIds, "i-019dbda6cd752f583")
	request.InstanceIds = append(request.InstanceIds, "i-0fef9450cc8c7b485")
	respChan, errChan := serverAWS.DescribeInstancesByIDsWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\nresp RequestId1: %#v \n", resp.RequestId)
		for _, v := range resp.Instances {
			fmt.Printf("resp : %#v \n", v)
		}
	case err := <-errChan:
		fmt.Printf("err====: %#v", err)
	}
	time.Sleep(1 * time.Second)
}

// AWS不支持密码登录
// 目前修改为密码登录，需要传入默认账号，利用rea非对称加密进行登录，并修改root密码
func Test_AWS_ChangeInstancePassword(t *testing.T) {
	request := new(zero_model.ChangeInstancePasswordRequest)
	request.Password = "987654321"
	request.LoginName = "centos"

	privateKeyPath := "/Users/allen/.ssh/id_rsa.pem"
	privateKey, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		t.Fatalf("unable to read private key: %v", err)
	}
	request.PrivateKey = privateKey
	request.InstanceAddress = "18.163.206.140"

	respChan, errChan := serverAWS.ChangeInstancePasswordWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp.RequestId)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

//AWS 没有对同一个实例进行重装系统的接口和设计，可以用启动一个新的实例，做好数据迁移，关闭旧的实例
func Test_AWS_ReplaceSystemDiskWithChan(t *testing.T) {

}

//运行中的实例也可以delete
// AWS Terminate 接口，返回成功后，过一段时间实例会从web界面消失。
func Test_AWS_DeleteInstancesWithChan(t *testing.T) {
	request := new(zero_model.DeleteInstancesRequest)
	request.InstanceIds = append(request.InstanceIds, "i-0207bee7aa31984d3")
	request.InstanceIds = append(request.InstanceIds, "i-0e54cb2665700c979")
	//request.InstanceIds = append(request.InstanceIds, "ins-p3rdsx3j")
	respChan, errChan := serverAWS.DeleteInstancesWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

//不指定RegionId 和 账户reginid 一致, 具体操作中不用传递，防范中没有用到这个reginid
//后期关注一下请求参数重的domain
func Test_AWS_AllocateEipAddressWithChan(t *testing.T) {
	request := new(zero_model.AllocateEipAddressRequest)
	//request.RegionId = "ap-beijing"
	request.Bandwidth = 1
	respChan, errChan := serverAWS.AllocateEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
}

//可以直接绑定，绑定后默认的公网ip取消，更新为新的eip
func Test_AWS_AssociateEipAddressWithChan(t *testing.T) {
	request := new(zero_model.AssociateEipAddressRequest)
	//request.RegionId = "ap-beijing"
	request.InstanceId = "i-04e7072b702ad462d"
	request.AllocationId = "eipalloc-0f541f5aa3bd6e900"
	request.EipAddress = "52.0.117.27"
	respChan, errChan := serverAWS.AssociateEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
}

// aws EIP的status 没有这个选项
// 不传递PublicIpsOrIPIds，则获取所有
func Test_AWS_DescribeEipAddressesWithChan(t *testing.T) {
	request := new(zero_model.DescribeEipAddressesRequest)
	//request.RegionId = "ap-beijing"
	//request.PublicIpsOrIPIds = append(request.PublicIpsOrIPIds, "18.210.149.253")
	//request.PublicIpsOrIPIds = append(request.PublicIpsOrIPIds, "eip-o738be5h")
	request.AssociatedInstanceType = "EcsInstance"
	respChan, errChan := serverAWS.DescribeEipAddressesWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
		for _, v := range resp.EipInfo {
			fmt.Printf("\n \n resp RequestId lenght: %+v \n", v)

		}
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_AWS_UnassociateEipAddressWithChan(t *testing.T) {
	request := new(zero_model.UnassociateEipAddressRequest)
	request.RegionId = "ap-beijing"
	request.EipAddressOrAddressId = "eip-n3mn83wt"
	request.InstanceId = "ins-5emx0dwj"
	respChan, errChan := serverAWS.UnassociateEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

// 本项目接口中可以都传递，具体的方法中已经做了优化，可以同时传递PublicIpOrAddressId 和 AllocationId
// 同时传递以AllocationId为准
// AWS 即使EIP绑定instance，也可以直接释放
// aws返回值是一个空结构体
func Test_AWS_ReleaseEipAddressWithChan(t *testing.T) {
	request := new(zero_model.ReleaseEipAddressRequest)
	//request.RegionId = "ap-beijing"
	//request.AllocationId = "eipalloc-053b4fee4f5e70b31"
	request.PublicIpOrAddressId = "52.0.117.27"
	//request.PublicIpOrAddressId = "eip-o738be5h"
	respChan, errChan := serverAWS.ReleaseEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

//web页面AMI 名称对应 ImageName， 具有唯一性
//web页面Name 对应 ImageVersion，可以重复
func Test_AWS_CreateImageWithChan(t *testing.T) {
	request := new(zero_model.CreateImageRequest)
	//request.RegionId = "ap-beijing"
	request.InstanceId = "i-04e7072b702ad462d"
	request.ImageName = "dev-allen-test-image-3"
	request.ImageVersion = "v1.0.0"
	respChan, errChan := serverAWS.CreateImageWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

//获取当前账号额镜像需要设置ower为self
func Test_AWS_DescribeImagesWithChan(t *testing.T) {
	request := new(zero_model.DescribeImagesRequest)
	//request.RegionId = "ap-beijing"
	//request.ImageId = "ami-0ad03b30091b655d5"
	respChan, errChan := serverAWS.DescribeImagesWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

//返回结果不代表同步成功
func Test_AWS_CopyImageWithChan(t *testing.T) {
	request := new(zero_model.CopyImageRequest)
	request.RegionId = "us-east-1"

	request.DestinationRegionId = "ec2.us-east-2.amazonaws.com"
	request.ImageId = "ami-07edba6b529cfc9f9"
	respChan, errChan := serverAWS.CopyImageWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v \n", err)
	}
	time.Sleep(1 * time.Second)
}

//声明操作对象serverAWS时候，region决定了操作镜像的region
func Test_AWS_DeleteImageWithChan(t *testing.T) {
	request := new(zero_model.DeleteImagesRequest)
	//request.RegionId = "ap-guangzhou"
	request.ImageIds = append(request.ImageIds, "ami-07edba6b529cfc9f9")
	respChan, errChan := serverAWS.DeleteImageWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}
