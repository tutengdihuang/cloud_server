package case_set

import (
	"cloud_server/drivers/zero_model"
	"testing"
	"time"
)

// 1.安全组必须和VPC绑定为一组才可用
func Test_Aliyun_Runinstance(t *testing.T) {
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
		t.Logf("resp: %+v", resp)
	case err := <-errChan:
		t.Logf("err: %+v", err)
		if err == nil && cap(respChan) > 0 {
			resp := <-respChan
			t.Logf("\n resp RequestId: %+v \n", resp.RequestId)
		}
	}
}

func Test_Aliyun_Stopinstance(t *testing.T) {
	request := new(zero_model.StopInstancesWithChanRequest)
	request.ForceStop = true
	request.InstanceIds = append(request.InstanceIds, "i-bp1iidni3it6om958wez")
	respChan, errChan := serverAliyun.StopInstancesWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("resp RequestId: %+v", resp.RequestId)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
}

func Test_Aliyun_Startinstance(t *testing.T) {
	request := new(zero_model.StartInstancesWithChanRequest)
	request.InstanceIds = append(request.InstanceIds, "i-bp14f84hogliwilyn6u1")
	respChan, errChan := serverAliyun.StartInstancesWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("resp RequestId: %+v", resp.RequestId)
		t.Logf("resp InstanceId: %+v", resp.InstanceId)
	case err := <-errChan:
		t.Logf("err: %+v", err)
		if err == nil && cap(respChan) > 0 {
			resp := <-respChan
			t.Logf("\n resp RequestId: %+v \n", resp.RequestId)
		}
	}
}

func Test_Aliyun_DescribeInstancesStatus(t *testing.T) {
	request := new(zero_model.DescribeInstancesStatusRequest)
	request.InstanceIds = append(request.InstanceIds, "i-bp1aoh75ith0gxt9lz0u")
	request.InstanceIds = append(request.InstanceIds, "i-bp15xwr9cmuqceodmixk")
	request.InstanceIds = append(request.InstanceIds, "i-bp14f84hogliwilyn6u1")
	respChan, errChan := serverAliyun.DescribeInstancesStatusWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\nresp RequestId1: %+v \n", resp.RequestId)
		t.Logf("resp InsStatus: %+v \n", resp.InsStatus)
	case err := <-errChan:
		t.Logf("err====: %+v", err)
	}
	time.Sleep(1 * time.Second)
}
func Test_Aliyun_DescribeAllInstances(t *testing.T) {
	request := new(zero_model.DescribeInstancesAllRequest)
	request.RegionId = "cn-hangzhou"
	respChan, errChan := serverAliyun.DescribeInstancesAllWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\nresp RequestId1: %+v \n", resp.RequestId)
		t.Logf("resp InsStatus: %+v \n", resp.Instances)
	case err := <-errChan:
		t.Logf("err====: %+v", err)
	}
	time.Sleep(1 * time.Second)
}
func Test_Aliyun_DescribeInstancesByIDS(t *testing.T) {
	request := new(zero_model.DescribeInstancesByIDsRequest)
	request.InstanceIds = append(request.InstanceIds, "i-bp1aoh75ith0gxt9lz0u")
	request.InstanceIds = append(request.InstanceIds, "i-bp15xwr9cmuqceodmixk")
	request.InstanceIds = append(request.InstanceIds, "i-bp14f84hogliwilyn6u1")
	respChan, errChan := serverAliyun.DescribeInstancesByIDsWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\nresp RequestId1: %+v \n", resp.RequestId)
		t.Logf("resp InsStatus: %+v \n", resp.Instances)
	case err := <-errChan:
		t.Logf("err====: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

// 问题 1。 errchan 会优先发送,修改密码后需要重启实例才能生效
func Test_Aliyun_ChangeInstancePassword(t *testing.T) {
	request := new(zero_model.ChangeInstancePasswordRequest)
	request.InstanceId = "i-bp1aoh75ith0gxt9lz0u"
	request.RegionId = "cn-hangzhou"
	request.Password = "Allen1234!"
	respChan, errChan := serverAliyun.ChangeInstancePasswordWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp.RequestId)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

//更换镜像前先停止实例
func Test_Aliyun_ReplaceSystemDiskWithChan(t *testing.T) {
	request := new(zero_model.ReplaceSystemDiskRequest)
	request.InstanceId = "i-bp15xwr9cmuqceodmixk"
	request.ImageId = "centos_stream_9_x64_20G_alibase_20230113.vhd"
	request.Password = "Allen123!"
	respChan, errChan := serverAliyun.ReplaceSystemDiskWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp.RequestId)
		t.Logf("\n resp DiskId: %+v \n", resp.DiskId)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

//删除前需要前先停止实例
func Test_Aliyun_DeleteInstancesWithChan(t *testing.T) {
	request := new(zero_model.DeleteInstancesRequest)
	request.InstanceIds = append(request.InstanceIds, "i-bp1iidni3it6om958wez")
	respChan, errChan := serverAliyun.DeleteInstancesWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Aliyun_AllocateEipAddressWithChan(t *testing.T) {
	request := new(zero_model.AllocateEipAddressRequest)
	request.RegionId = "cn-hangzhou"
	request.Bandwidth = 1
	respChan, errChan := serverAliyun.AllocateEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

// 绑定公网ip的实例，不能绑定eip
func Test_Aliyun_AssociateEipAddressWithChan(t *testing.T) {
	request := new(zero_model.AssociateEipAddressRequest)
	request.RegionId = "cn-hangzhou"
	request.AllocationId = "eip-bp1qsoht0mn6n1s2voyrv"
	request.InstanceId = "i-bp1iidni3it6om958wez"
	respChan, errChan := serverAliyun.AssociateEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		t.Logf("err: %+v", err)
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Aliyun_DescribeEipAddressesWithChan(t *testing.T) {
	request := new(zero_model.DescribeEipAddressesRequest)
	request.RegionId = "cn-hangzhou"
	request.AssociatedInstanceId = "i-bp1iidni3it6om958wez"
	request.AssociatedInstanceType = "EcsInstance"
	respChan, errChan := serverAliyun.DescribeEipAddressesWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Aliyun_UnassociateEipAddressWithChan(t *testing.T) {
	request := new(zero_model.UnassociateEipAddressRequest)
	request.RegionId = "cn-hangzhou"
	request.AllocationId = "eip-bp1qsoht0mn6n1s2voyrv"
	request.InstanceId = "i-bp1iidni3it6om958wez"
	respChan, errChan := serverAliyun.UnassociateEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Aliyun_ReleaseEipAddressWithChan(t *testing.T) {
	request := new(zero_model.ReleaseEipAddressRequest)
	request.RegionId = "cn-hangzhou"
	request.AllocationId = "eip-bp1qsoht0mn6n1s2voyrv"
	respChan, errChan := serverAliyun.ReleaseEipAddressWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Aliyun_CreateImageWithChan(t *testing.T) {
	request := new(zero_model.CreateImageRequest)
	request.RegionId = "cn-hangzhou"
	request.InstanceId = "i-bp1iidni3it6om958wez"
	request.ImageName = "dev-allen-test-2"
	request.ImageVersion = "v1.0.0"
	respChan, errChan := serverAliyun.CreateImageWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Aliyun_DescribeImagesWithChan(t *testing.T) {
	request := new(zero_model.DescribeImagesRequest)
	request.RegionId = "cn-beijing"
	request.ImageId = "m-bp131o1gbsqcm6yflldg"
	respChan, errChan := serverAliyun.DescribeImagesWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Aliyun_CopyImageWithChan(t *testing.T) {
	request := new(zero_model.CopyImageRequest)
	request.RegionId = "cn-hangzhou"
	request.DestinationRegionId = "cn-beijing"
	request.ImageId = "m-bp131o1gbsqcm6yflldg"
	respChan, errChan := serverAliyun.CopyImageWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}

func Test_Aliyun_DeleteImageWithChan(t *testing.T) {
	request := new(zero_model.DeleteImagesRequest)
	request.RegionId = "cn-hangzhou"
	request.ImageIds = append(request.ImageIds, "m-bp131o1gbsqcm6yflldg")
	respChan, errChan := serverAliyun.DeleteImageWithChan(request)
	select {
	case resp := <-respChan:
		t.Logf("\n resp RequestId: %+v \n", resp)
	case err := <-errChan:
		t.Logf("err: %+v", err)
	}
	time.Sleep(1 * time.Second)
}
