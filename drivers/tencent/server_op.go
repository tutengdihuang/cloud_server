package tencent

import (
	"cloud_server/core"
	"cloud_server/drivers/zero_model"
	"cloud_server/drivers/zero_template"
	"cloud_server/utils"
	"errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	ecs "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cvm/v20170312"
	vpc "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/vpc/v20170312"
	"strconv"
	"strings"
)

//Limit                  int64	腾讯云转换 pagesize page number
//OffSet                 int64	腾讯云转换

type ServerOp struct {
	client    *ecs.Client
	vpcClient *vpc.Client
	zero_template.TemplateCloudServer
}

func New(scfg *zero_model.CloudServerConfig) (core.Operator, error) {
	cred := common.NewCredential(
		scfg.AccessKey,
		scfg.SecretKey,
	)
	cpf := profile.NewClientProfile()
	client, err := ecs.NewClient(cred, scfg.Region, cpf)
	if err != nil {
		return nil, err
	}

	VPCClient, err := vpc.NewClient(cred, scfg.Region, cpf)
	if err != nil {
		return nil, err
	}

	op := new(ServerOp)
	op.client = client
	op.vpcClient = VPCClient
	return op, nil
}
func (this *ServerOp) StartInstancesWithChan(req *zero_model.StartInstancesWithChanRequest) (<-chan *zero_model.StartInstancesWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.StartInstancesWithChanResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.StartInstancesWithChanResponse)

	startReq := ecs.NewStartInstancesRequest()
	for _, v := range req.InstanceIds {
		startReq.InstanceIds = append(startReq.InstanceIds, &v)
	}

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.client.StartInstances(startReq)
		if errData != nil {
			errChan <- errData
			return
		}
		resp.RequestId = *respData.Response.RequestId
		respChan <- resp
	})

	return respChan, errChan
}
func (this *ServerOp) DescribeInstancesStatusWithChan(req *zero_model.DescribeInstancesStatusRequest) (<-chan *zero_model.DescribeInstancesStatusResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeInstancesStatusResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DescribeInstancesStatusResponse)

	createRequest := ecs.NewDescribeInstancesStatusRequest()
	for _, v := range req.InstanceIds {
		createRequest.InstanceIds = append(createRequest.InstanceIds, &v)
	}

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.client.DescribeInstancesStatus(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}
		resp.RequestId = *respData.Response.RequestId
		for _, v := range respData.Response.InstanceStatusSet {
			var unit zero_model.InstanceStatus
			unit.InstanceId = *v.InstanceId
			status, ok := zero_model.TencentInstanceStatusMapToInsStatusList[*v.InstanceState]
			if !ok {
				errChan <- errors.New("status not exist status:" + *v.InstanceState)
				return
			}
			unit.Status = status
			resp.InsStatus = append(resp.InsStatus, unit)
		}
		respChan <- resp

	})

	return respChan, errChan
}
func (this *ServerOp) RunInstancesWithChan(req *zero_model.RunInstanceWithChanRequest) (<-chan *zero_model.RunInstanceWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.RunInstanceWithChanResponse, 1)
	errChan := make(chan error, 1)

	runReq := ecs.NewRunInstancesRequest()
	runReq.InstanceName = &req.InstanceName
	if req.HostName != "" {
		runReq.HostName = &req.HostName
	}
	{
		placement := new(ecs.Placement)
		placement.Zone = &req.RegionId
		runReq.Placement = placement
	}
	runReq.InstanceType = &req.InstanceType
	{
		internetAccess := new(ecs.InternetAccessible)
		if req.InternetChargeType != "" && req.InternetMaxBandwidthOut != 0 {
			internetAccess.PublicIpAssigned = common.BoolPtr(true)
		}
		internetAccess.InternetMaxBandwidthOut = &req.InternetMaxBandwidthOut
		internetAccess.InternetChargeType = &req.InternetChargeType
		runReq.InternetAccessible = internetAccess
	}
	{
		loginSet := new(ecs.LoginSettings)
		loginSet.Password = &req.Password
		runReq.LoginSettings = loginSet
	}
	runReq.ImageId = &req.ImageId

	{
		if req.VSwitchId != "" {
			virtualId := new(ecs.VirtualPrivateCloud)
			virtualId.SubnetId = &req.VSwitchId
			runReq.VirtualPrivateCloud = virtualId
		}
	}
	//runReq.ResourceGroupId = req.ResourceGroupId
	{
		var systemDisk = new(ecs.SystemDisk)
		systemDisk.DiskType = &req.SystemDiskCategory
		size, err := strconv.Atoi(req.SystemDiskSize)
		if err != nil {
			errChan <- err
			return nil, errChan
		}
		var sizeC = int64(size)
		systemDisk.DiskSize = &sizeC
		runReq.SystemDisk = systemDisk
	}
	//tag
	{
		var tagSet = new(ecs.TagSpecification)
		ins := "instance"
		tagSet.ResourceType = &ins
		for k, v := range req.Tag {
			var tagOne = new(ecs.Tag)
			tagOne.Key = &k
			tagOne.Value = &v
			tagSet.Tags = append(tagSet.Tags, tagOne)

			runReq.TagSpecification = append(runReq.TagSpecification, tagSet)
		}
	}

	//SecurityEnhancementStrategy
	{
		runReq.SecurityGroupIds = append(runReq.SecurityGroupIds, &req.SecurityGroupId)
	}

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		rspData, errData := this.client.RunInstances(runReq)
		if errData != nil {
			errChan <- errData
			return
		}
		var respFeedBack zero_model.RunInstanceWithChanResponse
		respFeedBack.RequestId = *rspData.Response.RequestId
		for _, id := range rspData.Response.InstanceIdSet {
			respFeedBack.InstanceIdSets = append(respFeedBack.InstanceIdSets, *id)
		}
		respChan <- &respFeedBack
	})
	return respChan, errChan
}
func (this *ServerOp) ChangeInstancePasswordWithChan(req *zero_model.ChangeInstancePasswordRequest) (<-chan *zero_model.ChangeInstancePasswordResponse, <-chan error) {
	respChan := make(chan *zero_model.ChangeInstancePasswordResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)
	resp := new(zero_model.ChangeInstancePasswordResponse)

	createRequest := ecs.NewResetInstancesPasswordRequest()
	createRequest.Password = &req.Password
	createRequest.InstanceIds = append(createRequest.InstanceIds, &req.InstanceId)
	rspData, errData := this.client.ResetInstancesPassword(createRequest)
	if errData != nil {
		errChan <- errData
		return respChan, errChan
	}

	resp.RequestId = *rspData.Response.RequestId
	respChan <- resp
	return respChan, errChan
}
func (this *ServerOp) DescribeInstancesAllWithChan(req *zero_model.DescribeInstancesAllRequest) (<-chan *zero_model.DescribeInstancesAllResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeInstancesAllResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.DescribeInstancesAllResponse)

	createRequest := ecs.NewDescribeInstancesRequest()
	var limit int64 = 100
	var offset int64 = 0
	createRequest.Limit = &limit
	createRequest.Offset = &offset

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		//循环获取所有实例
		for {
			createRequest.Offset = &offset
			rspData, errData := this.client.DescribeInstances(createRequest)
			if errData != nil {
				errChan <- errData
				return
			}

			for _, v := range rspData.Response.InstanceSet {
				var unit zero_model.DescribeInstance

				unit.InstanceId = *v.InstanceId
				unit.InstanceName = *v.InstanceName
				for _, v1 := range v.PublicIpAddresses {
					unit.PublicIpAddress = append(unit.PublicIpAddress, *v1)
				}
				for _, v2 := range v.PrivateIpAddresses {
					unit.PrivateIpAddress = append(unit.PrivateIpAddress, *v2)
				}
				status, ok := zero_model.TencentInstanceStatusMapToInsStatusList[*v.InstanceState]
				if !ok {
					status = zero_model.InsStatus_UNKNOWN
				}
				unit.Status = status

				resp.Instances = append(resp.Instances, unit)
			}
			resp.RequestId = *rspData.Response.RequestId
			//如果获取的长度>=100则说明还有没有获取完全，继续循环获取
			if len(rspData.Response.InstanceSet) >= int(limit) {
				offset += limit
			} else {
				break
			}
		}

		respChan <- resp
	})

	return respChan, errChan
}
func (this *ServerOp) DescribeInstancesByIDsWithChan(req *zero_model.DescribeInstancesByIDsRequest) (<-chan *zero_model.DescribeInstancesByIDsResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeInstancesByIDsResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)
	resp := new(zero_model.DescribeInstancesByIDsResponse)

	createRequest := ecs.NewDescribeInstancesRequest()
	for _, v := range req.InstanceIds {
		createRequest.InstanceIds = append(createRequest.InstanceIds, &v)
	}

	utils.Go(func() {
		rspData, errData := this.client.DescribeInstances(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}
		resp.RequestId = *rspData.Response.RequestId
		for _, v := range rspData.Response.InstanceSet {
			var unit zero_model.DescribeInstance

			unit.InstanceId = *v.InstanceId
			unit.InstanceName = *v.InstanceName
			for _, v1 := range v.PublicIpAddresses {
				unit.PublicIpAddress = append(unit.PublicIpAddress, *v1)
			}
			for _, v2 := range v.PrivateIpAddresses {
				unit.PrivateIpAddress = append(unit.PrivateIpAddress, *v2)
			}
			status, ok := zero_model.TencentInstanceStatusMapToInsStatusList[*v.InstanceState]
			if !ok {
				status = zero_model.InsStatus_UNKNOWN
			}
			unit.Status = status
			resp.Instances = append(resp.Instances, unit)
		}

		respChan <- resp
	})
	return respChan, errChan
}
func (this *ServerOp) StopInstancesWithChan(req *zero_model.StopInstancesWithChanRequest) (<-chan *zero_model.StopInstancesWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.StopInstancesWithChanResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.StopInstancesWithChanResponse)

	createRequest := ecs.NewStopInstancesRequest()
	if req.ForceStop == true {
		createRequest.ForceStop = &req.ForceStop
	}
	for _, v := range req.InstanceIds {
		createRequest.InstanceIds = append(createRequest.InstanceIds, &v)
	}

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.client.StopInstances(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}
		resp.RequestId = *respData.Response.RequestId
		respChan <- resp
	})

	return respChan, errChan
}
func (this *ServerOp) ReplaceSystemDiskWithChan(req *zero_model.ReplaceSystemDiskRequest) (<-chan *zero_model.ReplaceSystemDiskResponse, <-chan error) {
	respChan := make(chan *zero_model.ReplaceSystemDiskResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.ReplaceSystemDiskResponse)

	createRequest := ecs.NewResetInstanceRequest()
	createRequest.InstanceId = &req.InstanceId
	createRequest.ImageId = &req.ImageId

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.client.ResetInstance(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}
		resp.RequestId = *respData.Response.RequestId
		respChan <- resp
	})

	return respChan, errChan
}
func (this *ServerOp) DeleteInstancesWithChan(req *zero_model.DeleteInstancesRequest) (<-chan *zero_model.DeleteInstancesResponse, <-chan error) {
	respChan := make(chan *zero_model.DeleteInstancesResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.DeleteInstancesResponse)

	createRequest := ecs.NewTerminateInstancesRequest()

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		createRequest.InstanceIds = make([]*string, 1)
		var strBuilder strings.Builder
		for i, v := range req.InstanceIds {
			createRequest.InstanceIds[0] = &v
			respData, errData := this.client.TerminateInstances(createRequest)
			if errData != nil {
				errChan <- errData
				return
			}
			strBuilder.WriteString(*respData.Response.RequestId)
			if i >= len(req.InstanceIds)-1 {
				break
			}
			strBuilder.WriteString(",")
		}
		resp.RequestId = strBuilder.String()

		respChan <- resp
	})

	return respChan, errChan
}

func (this *ServerOp) AllocateEipAddressWithChan(req *zero_model.AllocateEipAddressRequest) (<-chan *zero_model.AllocateEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.AllocateEipAddressResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.AllocateEipAddressResponse)

	createRequest := vpc.NewAllocateAddressesRequest()
	//如果不传递此参数则和绑定的实例一致
	createRequest.InternetMaxBandwidthOut = &req.Bandwidth

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.vpcClient.AllocateAddresses(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}

		resp.RequestId = *respData.Response.RequestId
		resp.EipAddress = *respData.Response.AddressSet[0]

		respChan <- resp
	})
	return respChan, errChan
}
func (this *ServerOp) AssociateEipAddressWithChan(req *zero_model.AssociateEipAddressRequest) (<-chan *zero_model.AssociateEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.AssociateEipAddressResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.AssociateEipAddressResponse)

	createRequest := vpc.NewAssociateAddressRequest()
	//如果不传递此参数则和绑定的实例一致
	createRequest.AddressId = &req.EipAddress
	createRequest.InstanceId = &req.InstanceId

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.vpcClient.AssociateAddress(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}

		resp.RequestId = *respData.Response.RequestId

		respChan <- resp
	})

	return respChan, errChan
}
func (this *ServerOp) DescribeEipAddressesWithChan(req *zero_model.DescribeEipAddressesRequest) (<-chan *zero_model.DescribeEipAddressesResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeEipAddressesResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DescribeEipAddressesResponse)

	createRequest := vpc.NewDescribeAddressesRequest()
	var offset int64 = 0
	var limit int64 = 100

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)

		if req.PublicIpsOrIPIds != nil && len(req.PublicIpsOrIPIds) > 0 {
			createRequest.AddressIds = make([]*string, 1)
			for _, v := range req.PublicIpsOrIPIds {
				createRequest.AddressIds[0] = &v
				respData, errData := this.vpcClient.DescribeAddresses(createRequest)
				if errData != nil {
					errChan <- errData
					return
				}
				if respData.Response == nil {
					errChan <- errors.New("DescribeAddresses: respData.Response is nil")
					return
				}
				resp.RequestId = *respData.Response.RequestId

				for _, v := range respData.Response.AddressSet {
					var eip zero_model.Eip
					eip.EipAddress = *v.AddressIp
					if v.AddressStatus != nil {
						status, ok := zero_model.TencentMapStrToEipStatus[*v.AddressStatus]
						if !ok {
							status = zero_model.EIPStATUS_UNKNOWN
						}
						eip.Status = status
					}

					if v.AddressIp != nil {
						eip.EipAddress = *v.AddressIp
					}
					if v.AddressId != nil {
						eip.AddressId = *v.AddressId
					}
					if v.InstanceId != nil {
						eip.InstanceId = *v.InstanceId
					}
					resp.EipInfo = append(resp.EipInfo, eip)
				}
			}
		} else {
			for {
				//createRequest.Limit = &limit
				//createRequest.Offset = &offset
				respData, errData := this.vpcClient.DescribeAddresses(createRequest)
				if errData != nil {
					errChan <- errData
					return
				}
				if respData.Response == nil {
					errChan <- errors.New("DescribeAddresses: respData.Response is nil")
					return
				}
				resp.RequestId = *respData.Response.RequestId
				for _, v := range respData.Response.AddressSet {
					var eip zero_model.Eip
					eip.EipAddress = *v.AddressIp
					if v.AddressStatus != nil {
						status, ok := zero_model.TencentMapStrToEipStatus[*v.AddressStatus]
						if !ok {
							status = zero_model.EIPStATUS_UNKNOWN
						}
						eip.Status = status

					}

					if v.AddressIp != nil {
						eip.EipAddress = *v.AddressIp
					}
					if v.AddressId != nil {
						eip.AddressId = *v.AddressId
					}
					if v.InstanceId != nil {
						eip.InstanceId = *v.InstanceId
					}
					resp.EipInfo = append(resp.EipInfo, eip)
				}

				if len(respData.Response.AddressSet) < int(limit) {
					break
				}
				offset += limit
			}
		}

		respChan <- resp
	})
	return respChan, errChan
}

func (this *ServerOp) UnassociateEipAddressWithChan(req *zero_model.UnassociateEipAddressRequest) (<-chan *zero_model.UnassociateEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.UnassociateEipAddressResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.UnassociateEipAddressResponse)

	createRequest := vpc.NewDisassociateAddressRequest()
	createRequest.AddressId = &req.EipAddressOrAddressId

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.vpcClient.DisassociateAddress(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}
		if respData.Response != nil && respData.Response.RequestId != nil {
			resp.RequestId = *respData.Response.RequestId
		}
		respChan <- resp
	})
	return respChan, errChan
}
func (this *ServerOp) ReleaseEipAddressWithChan(req *zero_model.ReleaseEipAddressRequest) (<-chan *zero_model.ReleaseEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.ReleaseEipAddressResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.ReleaseEipAddressResponse)

	createRequest := vpc.NewReleaseAddressesRequest()
	createRequest.AddressIds = append(createRequest.AddressIds, &req.PublicIpOrAddressId)
	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.vpcClient.ReleaseAddresses(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}

		resp.RequestId = *respData.Response.RequestId
		respChan <- resp
	})
	return respChan, errChan
}

func (this *ServerOp) DescribeImagesWithChan(req *zero_model.DescribeImagesRequest) (<-chan *zero_model.DescribeImagesResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeImagesResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.DescribeImagesResponse)

	createRequest := ecs.NewDescribeImagesRequest()
	createRequest.ImageIds = append(createRequest.ImageIds, &req.ImageId)
	var offset uint64 = 0
	var limit uint64 = 50

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		for {
			createRequest.Limit = &limit
			createRequest.Offset = &offset
			respData, errData := this.client.DescribeImages(createRequest)
			if errData != nil {
				errChan <- errData
				return
			}
			resp.RequestId = *respData.Response.RequestId
			for _, v := range respData.Response.ImageSet {
				var unit zero_model.Image
				unit.ImageId = *v.ImageId
				var status zero_model.ImageStatus
				status, ok := zero_model.TencentIMageStatusStrToImageStatus[*v.ImageState]
				if !ok {
					status = zero_model.IMAGE_STATUS_UNKNOWN
				}
				unit.Status = status

				resp.ImageSet = append(resp.ImageSet, unit)
			}
			if len(respData.Response.ImageSet) < int(limit) {
				break
			}
			offset += limit
		}
		respChan <- resp
	})

	return respChan, errChan
}
func (this *ServerOp) CreateImageWithChan(req *zero_model.CreateImageRequest) (<-chan *zero_model.CreateImageResponse, <-chan error) {
	respChan := make(chan *zero_model.CreateImageResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.CreateImageResponse)

	createRequest := ecs.NewCreateImageRequest()
	createRequest.InstanceId = &req.InstanceId
	createRequest.ImageName = &req.ImageName

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.client.CreateImage(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}

		resp.RequestId = *respData.Response.RequestId
		resp.ImageId = *respData.Response.ImageId
		respChan <- resp
	})

	return respChan, errChan
}
func (this *ServerOp) CopyImageWithChan(req *zero_model.CopyImageRequest) (<-chan *zero_model.CopyImageResponse, <-chan error) {
	respChan := make(chan *zero_model.CopyImageResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.CopyImageResponse)

	createRequest := ecs.NewSyncImagesRequest()
	createRequest.ImageIds = append(createRequest.ImageIds, &req.ImageId)
	createRequest.DestinationRegions = append(createRequest.DestinationRegions, &req.DestinationRegionId)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.client.SyncImages(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}

		resp.RequestId = *respData.Response.RequestId
		for _, v := range respData.Response.ImageSet {
			resp.ImageId = *v.ImageId
		}
		respChan <- resp
	})
	return respChan, errChan
}

func (this *ServerOp) DeleteImageWithChan(req *zero_model.DeleteImagesRequest) (<-chan *zero_model.DeleteImagesResponse, <-chan error) {
	respChan := make(chan *zero_model.DeleteImagesResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DeleteImagesResponse)
	createRequest := ecs.NewDeleteImagesRequest()
	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		for _, v := range req.ImageIds {
			createRequest.ImageIds = append(createRequest.ImageIds, &v)
		}
		respData, errData := this.client.DeleteImages(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}
		resp.RequestId = *respData.Response.RequestId

		respChan <- resp
	})

	return respChan, errChan
}
