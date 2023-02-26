package aliyun

import (
	"cloud_server/core"
	"cloud_server/drivers/zero_model"
	"cloud_server/drivers/zero_template"
	"cloud_server/utils"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	ecs "github.com/aliyun/alibaba-cloud-sdk-go/services/ecs"
	"strconv"
	"strings"
)

type ServerOp struct {
	client *ecs.Client

	zero_template.TemplateCloudServer
}

func New(scfg *zero_model.CloudServerConfig) (core.Operator, error) {
	client, err := ecs.NewClientWithAccessKey(scfg.Region, scfg.AccessKey, scfg.SecretKey)
	if err != nil {
		return nil, err
	}
	client.EnableAsync(1000, 1000)

	op := new(ServerOp)
	op.client = client
	return op, nil
}

func (this *ServerOp) StartInstancesWithChan(req *zero_model.StartInstancesWithChanRequest) (<-chan *zero_model.StartInstancesWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.StartInstancesWithChanResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.StartInstancesWithChanResponse)

	if req.InstanceIds == nil || len(req.InstanceIds) == 0 {
		errChan <- errors.New("StartInstancesWithChan InstanceIds is empty")
		return respChan, errChan
	}
	createRequest := ecs.CreateStartInstancesRequest()
	var ids = make([]string, 0)
	ids = append(ids, req.InstanceIds...)
	createRequest.InstanceId = &ids
	respDataChan, errDataChan := this.client.StartInstancesWithChan(createRequest)
	utils.GoWithCleaner(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}
	})

	return respChan, errChan
}
func (this *ServerOp) DescribeInstancesStatusWithChan(req *zero_model.DescribeInstancesStatusRequest) (<-chan *zero_model.DescribeInstancesStatusResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeInstancesStatusResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DescribeInstancesStatusResponse)

	createRequest := ecs.CreateDescribeInstancesRequest()
	var strBuilder strings.Builder
	strBuilder.WriteString("[")
	for i, v := range req.InstanceIds {
		strBuilder.WriteString("\"")
		strBuilder.WriteString(v)
		strBuilder.WriteString("\"")
		if i >= len(req.InstanceIds)-1 {
			break
		}
		strBuilder.WriteString(",")
	}
	strBuilder.WriteString("]")
	createRequest.InstanceIds = strBuilder.String()
	respDataChan, errDataChan := this.client.DescribeInstancesWithChan(createRequest)

	utils.Go(func() {
		defer close(errChan)
		defer close(respChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			for _, v := range data.Instances.Instance {
				var insStatus zero_model.InstanceStatus
				insStatus.InstanceId = v.InstanceId
				insStatus.Status = zero_model.AlyInstanceStatusMapToInsStatusList[v.Status]
				resp.InsStatus = append(resp.InsStatus, insStatus)
			}

			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}

	})

	return respChan, errChan
}
func (this *ServerOp) RunInstancesWithChan(req *zero_model.RunInstanceWithChanRequest) (<-chan *zero_model.RunInstanceWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.RunInstanceWithChanResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.RunInstanceWithChanResponse)

	runReq := ecs.CreateRunInstancesRequest()
	runReq.InstanceName = req.InstanceName
	runReq.HostName = req.HostName
	runReq.RegionId = req.RegionId
	runReq.InstanceType = req.InstanceType
	runReq.InternetMaxBandwidthOut = requests.NewInteger64(req.InternetMaxBandwidthOut)
	runReq.InternetChargeType = req.InternetChargeType
	runReq.InstanceChargeType = req.InstanceChargeType
	runReq.ImageId = req.ImageId
	runReq.Password = req.Password
	runReq.SecurityGroupId = req.SecurityGroupId
	runReq.VSwitchId = req.VSwitchId
	runReq.ResourceGroupId = req.ResourceGroupId
	runReq.SystemDiskCategory = req.SystemDiskCategory
	runReq.SystemDiskSize = req.SystemDiskSize
	var tagSet []ecs.RunInstancesTag
	for k, v := range req.Tag {
		var tagOne = ecs.RunInstancesTag{}
		tagOne.Key = k
		tagOne.Value = v
		tagSet = append(tagSet, tagOne)
	}
	runReq.Tag = &tagSet
	runReq.SecurityEnhancementStrategy = req.SecurityEnhancementStrategy
	runReq.Amount = requests.NewInteger64(req.Amount)
	respDataChan, errDataCh := this.client.RunInstancesWithChan(runReq)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case respData := <-respDataChan:
			if respData == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = respData.RequestId
			for _, id := range respData.InstanceIdSets.InstanceIdSet {
				resp.InstanceIdSets = append(resp.InstanceIdSets, id)
			}
			resp.TradePrice = respData.TradePrice
			respChan <- resp
		case err := <-errDataCh:
			errChan <- err
		}
	})

	return respChan, errChan
}
func (this *ServerOp) ChangeInstancePasswordWithChan(req *zero_model.ChangeInstancePasswordRequest) (<-chan *zero_model.ChangeInstancePasswordResponse, <-chan error) {
	respChan := make(chan *zero_model.ChangeInstancePasswordResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.ChangeInstancePasswordResponse)

	createRequest := ecs.CreateModifyInstanceAttributeRequest()
	createRequest.RegionId = req.RegionId
	createRequest.InstanceId = req.InstanceId
	createRequest.Password = req.Password
	respDataChan, errDataChan := this.client.ModifyInstanceAttributeWithChan(createRequest)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}
	})

	return respChan, errChan
}
func (this *ServerOp) DescribeInstancesAllWithChan(req *zero_model.DescribeInstancesAllRequest) (<-chan *zero_model.DescribeInstancesAllResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeInstancesAllResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.DescribeInstancesAllResponse)

	createRequest := ecs.CreateDescribeInstancesRequest()
	createRequest.PageSize = requests.NewInteger64(100)
	createRequest.RegionId = req.RegionId

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		var allInstances []ecs.Instance
		for {
			response, err := this.client.DescribeInstances(createRequest)
			if err != nil {
				errChan <- err
				return
			}

			allInstances = append(allInstances, response.Instances.Instance...)

			if response.NextToken == "" {
				break
			}
			createRequest.NextToken = response.NextToken
			resp.RequestId = response.RequestId
		}

		for _, v := range allInstances {
			var unit zero_model.DescribeInstance
			unit.InstanceId = v.InstanceId
			unit.InstanceName = v.InstanceName
			unit.PrivateIpAddress = append(unit.PrivateIpAddress, v.InnerIpAddress.IpAddress...)
			unit.PublicIpAddress = append(unit.PrivateIpAddress, v.PublicIpAddress.IpAddress...)
			status, ok := zero_model.AlyInstanceStatusMapToInsStatusList[v.Status]
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
func (this *ServerOp) DescribeInstancesByIDsWithChan(req *zero_model.DescribeInstancesByIDsRequest) (<-chan *zero_model.DescribeInstancesByIDsResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeInstancesByIDsResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DescribeInstancesByIDsResponse)

	createRequest := ecs.CreateDescribeInstancesRequest()
	var strBuilder strings.Builder
	strBuilder.WriteString("[")
	for i, v := range req.InstanceIds {
		strBuilder.WriteString("\"")
		strBuilder.WriteString(v)
		strBuilder.WriteString("\"")
		if i >= len(req.InstanceIds)-1 {
			break
		}
		strBuilder.WriteString(",")
	}
	strBuilder.WriteString("]")
	createRequest.InstanceIds = strBuilder.String()
	respDataChan, errDataChan := this.client.DescribeInstancesWithChan(createRequest)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			for _, v := range data.Instances.Instance {
				var unit zero_model.DescribeInstance
				unit.InstanceId = v.InstanceId
				unit.InstanceName = v.InstanceName
				unit.PrivateIpAddress = append(unit.PrivateIpAddress, v.InnerIpAddress.IpAddress...)
				unit.PublicIpAddress = append(unit.PrivateIpAddress, v.PublicIpAddress.IpAddress...)
				status, ok := zero_model.AlyInstanceStatusMapToInsStatusList[v.Status]
				if !ok {
					status = zero_model.InsStatus_UNKNOWN
				}
				unit.Status = status
				resp.Instances = append(resp.Instances, unit)
			}
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}
	})

	return respChan, errChan
}
func (this *ServerOp) StopInstancesWithChan(req *zero_model.StopInstancesWithChanRequest) (<-chan *zero_model.StopInstancesWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.StopInstancesWithChanResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.StopInstancesWithChanResponse)

	if req.InstanceIds == nil || len(req.InstanceIds) == 0 {
		errChan <- errors.New("StopInstancesWithChanRequest InstanceIds is empty")
		return respChan, errChan
	}

	createRequest := ecs.CreateStopInstancesRequest()
	var ids = make([]string, 0)
	ids = append(ids, req.InstanceIds...)
	createRequest.InstanceId = &ids
	createRequest.ForceStop = requests.NewBoolean(req.ForceStop)

	respDataChan, errDataChan := this.client.StopInstancesWithChan(createRequest)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}
	})

	return respChan, errChan
}
func (this *ServerOp) ReplaceSystemDiskWithChan(req *zero_model.ReplaceSystemDiskRequest) (<-chan *zero_model.ReplaceSystemDiskResponse, <-chan error) {
	respChan := make(chan *zero_model.ReplaceSystemDiskResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.ReplaceSystemDiskResponse)

	createRequest := ecs.CreateReplaceSystemDiskRequest()
	createRequest.InstanceId = req.InstanceId
	createRequest.ImageId = req.ImageId
	createRequest.Password = req.Password
	respDataChan, errDataChan := this.client.ReplaceSystemDiskWithChan(createRequest)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			resp.DiskId = data.DiskId
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}
	})

	return respChan, errChan
}
func (this *ServerOp) DeleteInstancesWithChan(req *zero_model.DeleteInstancesRequest) (<-chan *zero_model.DeleteInstancesResponse, <-chan error) {
	respChan := make(chan *zero_model.DeleteInstancesResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DeleteInstancesResponse)

	if req.InstanceIds == nil || len(req.InstanceIds) == 0 {
		errChan <- errors.New("DeleteInstancesWithChan InstanceIds is empty")
		return respChan, errChan
	}

	createRequest := ecs.CreateDeleteInstancesRequest()
	var ids = make([]string, 0)

	ids = append(ids, req.InstanceIds...)
	createRequest.InstanceId = &ids
	createRequest.Force = requests.NewBoolean(req.Force)
	respDataChan, errDataChan := this.client.DeleteInstancesWithChan(createRequest)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}
	})

	return respChan, errChan
}

func (this *ServerOp) AllocateEipAddressWithChan(req *zero_model.AllocateEipAddressRequest) (<-chan *zero_model.AllocateEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.AllocateEipAddressResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.AllocateEipAddressResponse)

	createRequest := ecs.CreateAllocateEipAddressRequest()
	createRequest.RegionId = req.RegionId
	createRequest.Bandwidth = strconv.Itoa(int(req.Bandwidth))
	respDataChan, errDataChan := this.client.AllocateEipAddressWithChan(createRequest)
	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			resp.EipAddress = data.EipAddress
			resp.AllocationId = data.AllocationId
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}
	})

	return respChan, errChan
}
func (this *ServerOp) AssociateEipAddressWithChan(req *zero_model.AssociateEipAddressRequest) (<-chan *zero_model.AssociateEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.AssociateEipAddressResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.AssociateEipAddressResponse)

	createRequest := ecs.CreateAssociateEipAddressRequest()
	createRequest.RegionId = req.RegionId
	createRequest.AllocationId = req.AllocationId
	createRequest.InstanceId = req.InstanceId
	respDataChan, errDataChan := this.client.AssociateEipAddressWithChan(createRequest)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			//resp.AssociationId=data.
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}
	})

	return respChan, errChan
}
func (this *ServerOp) DescribeEipAddressesWithChan(req *zero_model.DescribeEipAddressesRequest) (<-chan *zero_model.DescribeEipAddressesResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeEipAddressesResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DescribeEipAddressesResponse)

	createRequest := ecs.CreateDescribeEipAddressesRequest()
	createRequest.RegionId = req.RegionId
	var pageNumber int64 = 1
	var pageSize int64 = 50
	createRequest.PageSize = requests.NewInteger64(pageSize)
	createRequest.PageNumber = requests.NewInteger64(pageNumber)
	createRequest.AssociatedInstanceId = req.AssociatedInstanceId
	createRequest.AssociatedInstanceType = req.AssociatedInstanceType
	createRequest.Status = zero_model.AlyEipMapEipStatusToStr[req.Status]

	if len(req.AllocationIds) > 50 {
		err := errors.New("the number of allocationIds should be less than 50")
		errChan <- err
		return respChan, errChan
	}
	var strBuilder0 strings.Builder
	for i, v := range req.AllocationIds {
		strBuilder0.WriteString(v)
		if i >= len(req.AllocationIds)-1 {
			break
		}
		strBuilder0.WriteString(",")
	}

	createRequest.AllocationId = strBuilder0.String()

	if len(req.PublicIpsOrIPIds) > 50 {
		err := errors.New("the number of PublicIpsOrIPIds should be less than 50")
		errChan <- err
		return respChan, errChan
	}

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		var strBuilder strings.Builder
		for i, v := range req.PublicIpsOrIPIds {
			strBuilder.WriteString(v)
			if i >= len(req.PublicIpsOrIPIds)-1 {
				break
			}
			strBuilder.WriteString(",")
		}
		createRequest.EipAddress = strBuilder.String()

		for {
			createRequest.PageNumber = requests.NewInteger64(pageNumber)
			response, err := this.client.DescribeEipAddresses(createRequest)
			if err != nil {
				errChan <- err
				return
			}
			for _, v := range response.EipAddresses.EipAddress {
				var eip zero_model.Eip
				eip.EipAddress = v.IpAddress
				eip.AllocationId = v.AllocationId
				eip.Status = zero_model.AlyEipMapStrToEipStatus[v.Status]
				eip.InstanceId = v.InstanceId
				resp.EipInfo = append(resp.EipInfo, eip)
			}
			if len(response.EipAddresses.EipAddress) < int(pageSize) {
				break
			}
			pageNumber++
		}

		respChan <- resp
	})

	return respChan, errChan
}
func (this *ServerOp) UnassociateEipAddressWithChan(req *zero_model.UnassociateEipAddressRequest) (<-chan *zero_model.UnassociateEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.UnassociateEipAddressResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.UnassociateEipAddressResponse)

	createRequest := ecs.CreateUnassociateEipAddressRequest()
	createRequest.RegionId = req.RegionId
	createRequest.AllocationId = req.AllocationId
	createRequest.InstanceId = req.InstanceId
	respDataChan, errDataChan := this.client.UnassociateEipAddressWithChan(createRequest)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}
	})

	return respChan, errChan
}
func (this *ServerOp) ReleaseEipAddressWithChan(req *zero_model.ReleaseEipAddressRequest) (<-chan *zero_model.ReleaseEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.ReleaseEipAddressResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.ReleaseEipAddressResponse)

	createRequest := ecs.CreateReleaseEipAddressRequest()
	createRequest.RegionId = req.RegionId
	createRequest.AllocationId = req.AllocationId
	respDataChan, errDataChan := this.client.ReleaseEipAddressWithChan(createRequest)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}

	})

	return respChan, errChan
}

func (this *ServerOp) DescribeImagesWithChan(req *zero_model.DescribeImagesRequest) (<-chan *zero_model.DescribeImagesResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeImagesResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DescribeImagesResponse)

	createRequest := ecs.CreateDescribeImagesRequest()
	createRequest.RegionId = req.RegionId
	createRequest.ImageId = req.ImageId
	var pageNumber int64 = 1
	var pageSize int64 = 50
	createRequest.PageSize = requests.NewInteger64(pageSize)
	createRequest.PageNumber = requests.NewInteger64(pageNumber)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		for {
			createRequest.PageNumber = requests.NewInteger64(pageNumber)

			response, err := this.client.DescribeImages(createRequest)
			if err != nil {
				errChan <- err
				return
			}
			for _, v := range response.Images.Image {
				var imageUnit zero_model.Image
				imageUnit.ImageId = v.ImageId
				var status zero_model.ImageStatus
				status, ok := zero_model.AliIMageStatusStrToImageStatus[v.Status]
				if !ok {
					status = zero_model.IMAGE_STATUS_UNKNOWN
				}
				imageUnit.Status = status
				resp.ImageSet = append(resp.ImageSet, imageUnit)
			}
			if len(response.Images.Image) < int(pageSize) {
				break
			}
			pageNumber++
		}

		respChan <- resp
	})

	return respChan, errChan
}
func (this *ServerOp) CreateImageWithChan(req *zero_model.CreateImageRequest) (<-chan *zero_model.CreateImageResponse, <-chan error) {
	respChan := make(chan *zero_model.CreateImageResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.CreateImageResponse)

	createRequest := ecs.CreateCreateImageRequest()
	createRequest.RegionId = req.RegionId
	createRequest.InstanceId = req.InstanceId
	createRequest.ImageName = req.ImageName
	createRequest.ImageVersion = req.ImageVersion
	respDataChan, errDataChan := this.client.CreateImageWithChan(createRequest)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			resp.ImageId = data.ImageId
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}
	})

	return respChan, errChan
}
func (this *ServerOp) CopyImageWithChan(req *zero_model.CopyImageRequest) (<-chan *zero_model.CopyImageResponse, <-chan error) {
	respChan := make(chan *zero_model.CopyImageResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.CopyImageResponse)

	createRequest := ecs.CreateCopyImageRequest()
	createRequest.ImageId = req.ImageId
	createRequest.RegionId = req.RegionId
	createRequest.DestinationRegionId = req.DestinationRegionId
	respDataChan, errDataChan := this.client.CopyImageWithChan(createRequest)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		select {
		case data := <-respDataChan:
			if data == nil {
				errChan <- errors.New("respDataChan return nil")
				return
			}
			resp.RequestId = data.RequestId
			resp.ImageId = data.ImageId
			respChan <- resp
		case err := <-errDataChan:
			errChan <- err
		}
	})

	return respChan, errChan
}
func (this *ServerOp) DeleteImageWithChan(req *zero_model.DeleteImagesRequest) (<-chan *zero_model.DeleteImagesResponse, <-chan error) {
	respChan := make(chan *zero_model.DeleteImagesResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DeleteImagesResponse)

	createRequest := ecs.CreateDeleteImageRequest()
	createRequest.RegionId = req.RegionId
	createRequest.Force = requests.NewBoolean(req.Force)
	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		for i, v := range req.ImageIds {
			createRequest.ImageId = v
			respData, errData := this.client.DeleteImage(createRequest)
			if errData != nil {
				errChan <- errData
				return
			}
			resp.RequestId += respData.RequestId
			if i >= len(req.ImageIds)-1 {
				break
			}
			resp.RequestId += ","
		}

		respChan <- resp
	})

	return respChan, errChan
}
