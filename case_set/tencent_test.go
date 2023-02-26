package case_set

import (
	"cloud_server/drivers/zero_model"
	"fmt"
	"testing"
	"time"
)

//如果需要hostname 需要注意格式
func Test_Tencent_Runinstance(t *testing.T) {
	request := new(zero_model.RunInstanceWithChanRequest)
	request.RegionId = "ap-beijing-7" //placement
	request.InstanceName = "QCLOUD-ALLEN-TEST"

	request.ImageId = "img-487zeit5"
	request.Password = "Allen123!"
	request.SecurityGroupId = "sg-51ypsokb"
	//request.VSwitchId = "vsw-bp1tw0ukcxmwt2vvp2kcl"
	request.InstanceType = "S6.MEDIUM2"
	request.InstanceChargeType = "POSTPAID_BY_HOUR"
	request.InternetChargeType = "TRAFFIC_POSTPAID_BY_HOUR"
	request.InternetMaxBandwidthOut = 1
	request.Amount = 1
	request.SystemDiskSize = "20"
	request.SystemDiskCategory = "CLOUD_PREMIUM"

	if request == nil {
		panic("request is nil")
	}
	if serverTencent == nil {
		panic("serverTencent is nil")
	}
	respChan, errChan := serverTencent.RunInstancesWithChan(request)
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
}

func Test_Tencent_Stopinstance(t *testing.T) {
	request := new(zero_model.StopInstancesWithChanRequest)
	request.ForceStop = true
	request.InstanceIds = append(request.InstanceIds, "ins-dq33z4ml")
	respChan, errChan := serverTencent.StopInstancesWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("resp RequestId: %+v", resp.RequestId)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
}

func Test_Tencent_Startinstance(t *testing.T) {
	request := new(zero_model.StartInstancesWithChanRequest)
	request.InstanceIds = append(request.InstanceIds, "ins-5emx0dwj")
	respChan, errChan := serverTencent.StartInstancesWithChan(request)
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

func Test_Tencent_DescribeInstancesStatus(t *testing.T) {
	request := new(zero_model.DescribeInstancesStatusRequest)
	request.InstanceIds = append(request.InstanceIds, "ins-5emx0dwj")
	respChan, errChan := serverTencent.DescribeInstancesStatusWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\nresp RequestId1: %+v \n", resp.RequestId)
		fmt.Printf("resp InsStatus: %+v \n", resp.InsStatus)
	case err := <-errChan:
		fmt.Printf("err====: %+v", err)
	}
	time.Sleep(1 * time.Second)
}
func Test_Tencent_DescribeAllInstances(t *testing.T) {
	request := new(zero_model.DescribeInstancesAllRequest)
	request.RegionId = "ap-beijing"
	respChan, errChan := serverTencent.DescribeInstancesAllWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\nresp RequestId1: %#v \n", resp.RequestId)
		fmt.Printf("resp InsStatus: %#v \n", resp.Instances)
	case err := <-errChan:
		fmt.Printf("err====: %#v", err)
	}
	time.Sleep(1 * time.Second)
}
func Test_Tencent_DescribeInstancesByIDS(t *testing.T) {
	request := new(zero_model.DescribeInstancesByIDsRequest)
	request.InstanceIds = append(request.InstanceIds, "ins-dq33z4ml")
	respChan, errChan := serverTencent.DescribeInstancesByIDsWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\nresp RequestId1: %+v \n", resp.RequestId)
		fmt.Printf("resp InsStatus: %+v \n", resp.Instances)
	case err := <-errChan:
		fmt.Printf("err====: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

// 问题 1。 必须在stopped状态下修改密码修改密码后需要重启实例才能生效
func Test_Tencent_ChangeInstancePassword(t *testing.T) {
	request := new(zero_model.ChangeInstancePasswordRequest)
	request.InstanceId = "ins-5emx0dwj"
	request.RegionId = "ap-beijing"
	request.Password = "Allen1234!"
	respChan, errChan := serverTencent.ChangeInstancePasswordWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp.RequestId)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

//更换镜像前先停止实例
func Test_Tencent_ReplaceSystemDiskWithChan(t *testing.T) {
	request := new(zero_model.ReplaceSystemDiskRequest)
	request.InstanceId = "ins-5emx0dwj"
	request.ImageId = "img-4cmp1f33"
	request.Password = "Allen1234!"
	respChan, errChan := serverTencent.ReplaceSystemDiskWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp.RequestId)
		fmt.Printf("\n resp DiskId: %+v \n", resp.DiskId)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

//运行中的实例也可以delete
func Test_Tencent_DeleteInstancesWithChan(t *testing.T) {
	request := new(zero_model.DeleteInstancesRequest)
	request.InstanceIds = append(request.InstanceIds, "ins-dq33z4ml")
	//request.InstanceIds = append(request.InstanceIds, "ins-p3rdsx3j")
	respChan, errChan := serverTencent.DeleteInstancesWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Tencent_AllocateEipAddressWithChan(t *testing.T) {
	request := new(zero_model.AllocateEipAddressRequest)
	request.RegionId = "ap-beijing"
	request.Bandwidth = 1
	respChan, errChan := serverTencent.AllocateEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
}

//可以直接绑定，绑定后默认的公网ip取消，更新为新的eip
func Test_Tencent_AssociateEipAddressWithChan(t *testing.T) {
	request := new(zero_model.AssociateEipAddressRequest)
	request.RegionId = "ap-beijing"
	request.InstanceId = "ins-3bnaqcg5"
	request.EipAddress = "eip-o738be5h"
	respChan, errChan := serverTencent.AssociateEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
}

//腾讯接口穿入两个参数只返回最后一个，不能同时返回
func Test_Tencent_DescribeEipAddressesWithChan(t *testing.T) {
	request := new(zero_model.DescribeEipAddressesRequest)
	request.RegionId = "ap-beijing"
	//request.PublicIpsOrIPIds = append(request.PublicIpsOrIPIds, "eip-n3mn83wt")
	//request.PublicIpsOrIPIds = append(request.PublicIpsOrIPIds, "eip-o738be5h")
	request.AssociatedInstanceType = "EcsInstance"
	respChan, errChan := serverTencent.DescribeEipAddressesWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
		fmt.Printf("\n resp RequestId lenght: %+v \n", len(resp.EipInfo))
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Tencent_UnassociateEipAddressWithChan(t *testing.T) {
	request := new(zero_model.UnassociateEipAddressRequest)
	request.RegionId = "ap-beijing"
	request.EipAddressOrAddressId = "eip-n3mn83wt"
	request.InstanceId = "ins-5emx0dwj"
	respChan, errChan := serverTencent.UnassociateEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Tencent_ReleaseEipAddressWithChan(t *testing.T) {
	request := new(zero_model.ReleaseEipAddressRequest)
	request.RegionId = "ap-beijing"
	request.PublicIpOrAddressId = "eip-gnhg3o33"
	//request.PublicIpOrAddressId = "eip-o738be5h"
	respChan, errChan := serverTencent.ReleaseEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Tencent_CreateImageWithChan(t *testing.T) {
	request := new(zero_model.CreateImageRequest)
	request.RegionId = "ap-beijing"
	request.InstanceId = "ins-dq33z4ml"
	request.ImageName = "dev-allen-test-1"
	request.ImageVersion = "v1.0.0"
	respChan, errChan := serverTencent.CreateImageWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Tencent_DescribeImagesWithChan(t *testing.T) {
	request := new(zero_model.DescribeImagesRequest)
	request.RegionId = "ap-beijing"
	request.ImageId = "img-0qfi20a1"
	respChan, errChan := serverTencent.DescribeImagesWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

//返回结果不代表同步成功
func Test_Tencent_CopyImageWithChan(t *testing.T) {
	request := new(zero_model.CopyImageRequest)
	request.RegionId = "ap-beijing"
	request.DestinationRegionId = "ap-guangzhou"
	request.ImageId = "img-0qfi20a1"
	respChan, errChan := serverTencent.CopyImageWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v \n", err)
	}
	time.Sleep(1 * time.Second)
}

//声明操作对象serverTencent时候，region决定了操作镜像的region
func Test_Tencent_DeleteImageWithChan(t *testing.T) {
	request := new(zero_model.DeleteImagesRequest)
	//request.RegionId = "ap-guangzhou"
	request.ImageIds = append(request.ImageIds, "img-0qfi20a1")
	respChan, errChan := serverTencent.DeleteImageWithChan(request)
	select {
	case resp := <-respChan:
		fmt.Printf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		fmt.Printf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}
