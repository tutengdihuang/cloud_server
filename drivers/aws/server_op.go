package aws

import (
	"cloud_server/core"
	"cloud_server/drivers/zero_model"
	"cloud_server/drivers/zero_template"
	"cloud_server/utils"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"golang.org/x/crypto/ssh"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	ecs "github.com/aws/aws-sdk-go/service/ec2"

	"strconv"
)

type ServerOp struct {
	zero_template.TemplateCloudServer
	client *ecs.EC2
}

func New(cfg *zero_model.CloudServerConfig) (core.Operator, error) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
	})
	if err != nil {
		return nil, err
	}

	client := ecs.New(sess)

	op := new(ServerOp)
	op.client = client
	return op, nil
}

func (this *ServerOp) StartInstancesWithChan(req *zero_model.StartInstancesWithChanRequest) (<-chan *zero_model.StartInstancesWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.StartInstancesWithChanResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.StartInstancesWithChanResponse)

	startReq := new(ecs.StartInstancesInput)
	for _, v := range req.InstanceIds {
		ins := v
		startReq.InstanceIds = append(startReq.InstanceIds, &ins)
	}

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.client.StartInstances(startReq)
		if errData != nil {
			errChan <- errData
			return
		}
		if len(respData.StartingInstances) != 0 && respData.StartingInstances != nil {
			resp.InstanceId = *respData.StartingInstances[0].InstanceId
			respChan <- resp
		}
		errChan <- errData
	})
	return respChan, errChan
}
func (this *ServerOp) DescribeInstancesStatusWithChan(req *zero_model.DescribeInstancesStatusRequest) (<-chan *zero_model.DescribeInstancesStatusResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeInstancesStatusResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.DescribeInstancesStatusResponse)

	createRequest := new(ecs.DescribeInstanceStatusInput)
	//获取所有实例
	switch {
	case req.InstanceIds == nil || len(req.InstanceIds) == 0:
		utils.Go(func() {
			defer close(respChan)
			defer close(errChan)
			createRequest.MaxResults = aws.Int64(5)

			for {
				respData, errData := this.client.DescribeInstanceStatus(createRequest)
				if errData != nil {
					errChan <- errData
					return
				}
				//resp.RequestId = *respData.RequestId
				for _, v1 := range respData.InstanceStatuses {
					v := v1
					var unit zero_model.InstanceStatus
					unit.InstanceId = *v.InstanceId
					status, ok := zero_model.AWSInstanceStatusMapToInsStatusList[*v.InstanceState.Name]
					if !ok {
						errChan <- errors.New("status not exist status:" + *v.InstanceState.Name)
						return
					}
					unit.Status = status
					resp.InsStatus = append(resp.InsStatus, unit)
				}

				if respData.NextToken == nil {
					break
				}
				createRequest.NextToken = respData.NextToken
			}

			respChan <- resp
		})
	case len(req.InstanceIds) > 0:
		//获取指定实例状态
		createRequest.InstanceIds = make([]*string, 1)

		utils.Go(func() {
			defer close(respChan)
			defer close(errChan)

			for _, v1 := range req.InstanceIds {
				v := v1
				createRequest.InstanceIds[0] = &v
				respData, errData := this.client.DescribeInstanceStatus(createRequest)
				if errData != nil {
					errChan <- errData
					return
				}
				//resp.RequestId = *respData.RequestId
				for _, v2 := range respData.InstanceStatuses {
					v := v2
					var unit zero_model.InstanceStatus
					unit.InstanceId = *v.InstanceId
					status, ok := zero_model.AWSInstanceStatusMapToInsStatusList[*v.InstanceState.Name]
					if !ok {
						errChan <- errors.New("status not exist status:" + *v.InstanceState.Name)
						return
					}
					unit.Status = status
					resp.InsStatus = append(resp.InsStatus, unit)
				}
			}

			respChan <- resp
		})
	}

	//获取指定实例

	return respChan, errChan
}

func (this *ServerOp) RunInstancesWithChan(req *zero_model.RunInstanceWithChanRequest) (<-chan *zero_model.RunInstanceWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.RunInstanceWithChanResponse, 1)
	errChan := make(chan error, 1)

	runReq := new(ecs.RunInstancesInput)

	{
		if req.KeyPairName != "" {
			runReq.KeyName = &req.KeyPairName
		}
	}
	//runReq.InstanceName = &req.InstanceName
	{
		placement := new(ecs.Placement)
		if req.HostName != "" {
			placement.HostId = &req.HostName
		}

		if req.RegionId != "" {
			placement.AvailabilityZone = &req.RegionId

		}
		runReq.Placement = placement
	}
	if req.InstanceType != "" {
		runReq.InstanceType = &req.InstanceType
	}
	{
		//internetAccess := new(ecs.InternetAccessible)
		//internetAccess.InternetMaxBandwidthOut = &req.InternetMaxBandwidthOut
		//internetAccess.InternetChargeType = &req.InternetChargeType
		//runReq.InternetAccessible = internetAccess

	}
	{
		//loginSet := new(ecs.LoginSettings)
		//loginSet.Password = &req.Password
		//runReq.LoginSettings = loginSet
	}
	if req.ImageId != "" {
		runReq.ImageId = &req.ImageId
	}
	{
		if req.SecurityGroupId != "" {
			runReq.SecurityGroupIds = append(runReq.SecurityGroupIds, &req.SecurityGroupId)

		}
	}
	{
		if req.VSwitchId != "" {
			runReq.SubnetId = &req.VSwitchId
		}
	}
	//runReq.ResourceGroupId = req.ResourceGroupId
	{
		var bdms = make([]*ecs.BlockDeviceMapping, 0)
		bdm := &ecs.BlockDeviceMapping{}
		var ebs = new(ecs.EbsBlockDevice)
		if req.SystemDiskCategory != "" {
			ebs.VolumeType = &req.SystemDiskCategory
		}
		size, err := strconv.Atoi(req.SystemDiskSize)
		if err != nil {
			errChan <- err
			close(respChan)
			close(errChan)
			return respChan, errChan
		}
		var sizeC = int64(size)
		ebs.VolumeSize = &sizeC
		bdm.Ebs = ebs
		bdms = append(bdms, bdm)

		var name = "/dev/sdh"
		bdm.DeviceName = &name
		runReq.BlockDeviceMappings = append(runReq.BlockDeviceMappings, bdm)
		runReq.SetBlockDeviceMappings(bdms)
	}
	//tag
	{
		//tag
		{

			var tagSets = make([]*ecs.TagSpecification, 0)
			for k1, v1 := range req.Tag {
				k, v := k1, v1
				var tagSet = new(ecs.TagSpecification)
				ins := "instance"
				tagSet.ResourceType = &ins
				var tagOne = new(ecs.Tag)
				if v == "" || k == "" {
					tagOne.Key = &k
					tagOne.Value = &v
				}

				tagSet.Tags = append(tagSet.Tags, tagOne)
				tagSets = append(tagSets, tagSet)
			}
			if req.InstanceName != "" {
				var tagSet = new(ecs.TagSpecification)
				var tagOne = new(ecs.Tag)
				var name = "Name"
				ins := "instance"
				tagSet.ResourceType = &ins
				tagOne.Key = &name
				tagOne.Value = &req.InstanceName
				tagSet.Tags = append(tagSet.Tags, tagOne)
				tagSets = append(tagSets, tagSet)
			}
			runReq.TagSpecifications = tagSets
		}
	}

	//SecurityEnhancementStrategy
	{

	}
	runReq.MaxCount = &req.Amount
	runReq.MinCount = &req.Amount

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		rspData, errData := this.client.RunInstances(runReq)
		if errData != nil {
			errChan <- errData
			return
		}
		var respFeedBack zero_model.RunInstanceWithChanResponse
		if rspData != nil && rspData.RequesterId != nil {
			respFeedBack.RequestId = *rspData.RequesterId
		}
		for _, instance1 := range rspData.Instances {
			instance := instance1
			respFeedBack.InstanceIdSets = append(respFeedBack.InstanceIdSets, *instance.InstanceId)
		}

		respChan <- &respFeedBack
	})
	return respChan, errChan
}

//目前不支持，也不支持修改密钥对名字
func (this *ServerOp) ChangeInstancePasswordWithChan(req *zero_model.ChangeInstancePasswordRequest) (<-chan *zero_model.ChangeInstancePasswordResponse, <-chan error) {
	respChan := make(chan *zero_model.ChangeInstancePasswordResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.ChangeInstancePasswordResponse)

	createRequest := &ecs.ModifyInstanceAttributeInput{}
	createRequest.InstanceId = &req.InstanceId

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)

		signer, err := ssh.ParsePrivateKey(req.PrivateKey)
		if err != nil {
			errChan <- err
			return
		}

		config := &ssh.ClientConfig{
			User:            req.LoginName,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(signer),
			},
		}

		client, err := ssh.Dial("tcp", req.InstanceAddress+":"+"22", config)
		if err != nil {
			errChan <- err
			return
		}

		session, err := client.NewSession()
		if err != nil {
			errChan <- err
			return
		}
		defer session.Close()

		cmd := fmt.Sprintf("sudo sed -i 's/PasswordAuthentication no/PasswordAuthentication yes/g' /etc/ssh/sshd_config && sudo service sshd restart && echo \"root:%s\" | sudo chpasswd", req.Password)
		_, err = session.CombinedOutput(cmd)
		if err != nil {
			errChan <- err
			return
		}

		resp.RequestId = "No RequestId response for AWS ChangeInstancePasswordWithChan"
		respChan <- resp
	})

	return respChan, errChan
}

func (this *ServerOp) DescribeInstancesAllWithChan(req *zero_model.DescribeInstancesAllRequest) (<-chan *zero_model.DescribeInstancesAllResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeInstancesAllResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.DescribeInstancesAllResponse)

	createRequest := new(ecs.DescribeInstancesInput)
	var maxResult int64 = 50
	createRequest.MaxResults = &maxResult

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		//循环获取所有实例
		for {
			rspData, errData := this.client.DescribeInstances(createRequest)
			if errData != nil {
				errChan <- errData
				return
			}

			for _, v := range rspData.Reservations {
				for _, ins1 := range v.Instances {
					ins := ins1
					var instance = zero_model.DescribeInstance{}

					if ins == nil {
						continue
					}
					if ins.InstanceId != nil {
						instance.InstanceId = *ins.InstanceId
					}

					if ins.Tags != nil {
						for _, vIn := range ins.Tags {
							v := vIn
							if v == nil {
								continue
							}
							if *v.Key == "Name" {
								instance.InstanceName = *v.Value
							}
						}
					}
					status, ok := zero_model.AWSInstanceStatusMapToInsStatusList[*ins.State.Name]
					if !ok {
						status = zero_model.InsStatus_UNKNOWN
					}

					instance.Status = status
					if ins.PublicIpAddress != nil {
						instance.PublicIpAddress = append(instance.PublicIpAddress, *ins.PublicIpAddress)
					}
					if ins.PrivateIpAddress != nil {
						instance.PrivateIpAddress = append(instance.PrivateIpAddress, *ins.PrivateIpAddress)
					}
					resp.Instances = append(resp.Instances, instance)
				}
			}
			resp.RequestId = "No RequestId"

			if rspData.NextToken != nil && *rspData.NextToken != "" {
				createRequest.NextToken = rspData.NextToken
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
	resp := new(zero_model.DescribeInstancesByIDsResponse)

	createRequest := new(ecs.DescribeInstancesInput)

	utils.Go(func() {
		if req.InstanceIds == nil || len(req.InstanceIds) == 0 {
			errChan <- errors.New("no instance id is available")
		}
		defer close(respChan)
		defer close(errChan)
		createRequest.InstanceIds = make([]*string, 1)

		for _, v := range req.InstanceIds {
			createRequest.InstanceIds[0] = &v
			rspData, errData := this.client.DescribeInstances(createRequest)
			if errData != nil {
				errChan <- errData
				return
			}
			resp.RequestId = "No RequestId"
			for _, v := range rspData.Reservations {
				for _, insIn := range v.Instances {
					ins := insIn
					var instance = zero_model.DescribeInstance{}

					if ins == nil {
						continue
					}
					if ins.InstanceId != nil {
						instance.InstanceId = *ins.InstanceId
					}

					if ins.Tags != nil {
						for _, vIns := range ins.Tags {
							v := vIns
							if v == nil {
								continue
							}
							if *v.Key == "Name" {
								instance.InstanceName = *v.Value
							}
						}
					}
					status, ok := zero_model.AWSInstanceStatusMapToInsStatusList[*ins.State.Name]
					if !ok {
						status = zero_model.InsStatus_UNKNOWN
					}

					instance.Status = status
					if ins.PublicIpAddress != nil {
						instance.PublicIpAddress = append(instance.PublicIpAddress, *ins.PublicIpAddress)
					}
					if ins.PrivateIpAddress != nil {
						instance.PrivateIpAddress = append(instance.PrivateIpAddress, *ins.PrivateIpAddress)
					}
					resp.Instances = append(resp.Instances, instance)
				}
			}

		}

		respChan <- resp
	})

	return respChan, errChan
}
func (this *ServerOp) StopInstancesWithChan(req *zero_model.StopInstancesWithChanRequest) (<-chan *zero_model.StopInstancesWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.StopInstancesWithChanResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.StopInstancesWithChanResponse)

	createRequest := &ecs.StopInstancesInput{}
	for _, vIns := range req.InstanceIds {
		v := vIns
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
		if respData.StoppingInstances == nil || len(respData.StoppingInstances) == 0 {
			errChan <- errors.New("StopInstancesWithChan respData.StoppingInstances is nil or length is 0")
			return
		}

		buf := strings.Builder{}
		buf.WriteString("for aws StopInstancesWithChan, there is not response id, and the Instance id info is as below:")

		for _, instanceIns := range respData.StoppingInstances {
			instance := instanceIns
			buf.WriteString(*instance.InstanceId)
			buf.WriteString(",")
		}
		resp.RequestId = buf.String()
		respChan <- resp
	})

	return respChan, errChan
}

/*func (this *ServerOp) ReplaceSystemDiskWithChan(req *zero_model.ReplaceSystemDiskRequest) (<-chan *zero_model.ReplaceSystemDiskResponse, <-chan error) {
	respChan := make(chan *zero_model.ReplaceSystemDiskResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)
	resp := new(zero_model.ReplaceSystemDiskResponse)

	createRequest := new(ecs.replace)
	createRequest.InstanceId = &req.InstanceId
	createRequest.ImageId = &req.ImageId
	respData, errData := this.client.CreateReplaceRootVolumeTask(createRequest)
	if errData != nil {
		errChan <- errData
		return respChan, errChan
	}
	resp.RequestId = *respData.Response.RequestId

	respChan <- resp
	return respChan, errChan
}*/

func (this *ServerOp) DeleteInstancesWithChan(req *zero_model.DeleteInstancesRequest) (<-chan *zero_model.DeleteInstancesResponse, <-chan error) {
	respChan := make(chan *zero_model.DeleteInstancesResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DeleteInstancesResponse)
	if req == nil || req.InstanceIds == nil || len(req.InstanceIds) == 0 {

	}

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		createRequest := &ecs.TerminateInstancesInput{}
		createRequest.InstanceIds = make([]*string, 1)
		buf := strings.Builder{}

		for _, vIns := range req.InstanceIds {
			v := vIns
			createRequest.InstanceIds[0] = &v
			respData, errData := this.client.TerminateInstances(createRequest)
			if errData != nil {
				errChan <- errData
				return
			}

			if respData.TerminatingInstances == nil || len(respData.TerminatingInstances) == 0 {
				errChan <- errors.New("StopInstanceWithChan respData.StoppingInstances is nil or length is 0")
				return
			}

			buf.WriteString("for aws StopInstanceWithChan, there is not response id, and the Instanceid info is as below:")
			for _, instanceIns := range respData.TerminatingInstances {
				instance := instanceIns
				buf.WriteString(*instance.InstanceId)
				buf.WriteString(",")
			}

		}
		resp.RequestId = buf.String()

		respChan <- resp
	})

	return respChan, errChan
}

func (this *ServerOp) AllocateEipAddressWithChan(req *zero_model.AllocateEipAddressRequest) (<-chan *zero_model.AllocateEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.AllocateEipAddressResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.AllocateEipAddressResponse)

	createRequest := &ecs.AllocateAddressInput{}
	//如果不传递此参数则和绑定的实例一致
	//createRequest.Domain = aws.String("vpc")
	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		respData, errData := this.client.AllocateAddress(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}

		resp.AllocationId = *respData.AllocationId
		resp.EipAddress = *respData.PublicIp

		respChan <- resp
	})
	return respChan, errChan
}
func (this *ServerOp) AssociateEipAddressWithChan(req *zero_model.AssociateEipAddressRequest) (<-chan *zero_model.AssociateEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.AssociateEipAddressResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.AssociateEipAddressResponse)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		createRequest := &ecs.AssociateAddressInput{}
		if req == nil {
			errChan <- errors.New("AssociateAddress req == nil")
			return
		}
		if req.AllocationId != "" {
			createRequest.AllocationId = aws.String(req.AllocationId)
		} else if req.AllocationId == "" && req.EipAddress != "" {
			request := new(zero_model.DescribeEipAddressesRequest)
			request.PublicIpsOrIPIds = append(request.PublicIpsOrIPIds, req.EipAddress)
			respChanDes, errChanDes := this.DescribeEipAddressesWithChan(request)
			select {
			case respDes := <-respChanDes:
				if respDes == nil {
					errChan <- errors.New("DescribeEipAddressesWithChan resp is nil")
					return
				}
				if respDes.EipInfo[0].EipAddress != req.EipAddress {
					errChan <- errors.New(fmt.Sprintf("Eip not match, requestIp:%+v, resp IP: %+v", req.EipAddress, respDes.EipInfo[0].EipAddress))
					return
				}
				createRequest.AllocationId = &respDes.EipInfo[0].AllocationId
			case err := <-errChanDes:
				errChan <- err
				return
			}
		}
		if req.InstanceId == "" {
			errChan <- errors.New("DescribeEipAddressesWithChan req.InstanceId is empty")
			return
		}
		createRequest.InstanceId = &req.InstanceId
		respData, errData := this.client.AssociateAddress(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}
		if respData == nil {
			errChan <- errors.New("AssociateAddress return nil data")
			return
		}
		resp.RequestId = respData.String()
		resp.AssociationId = *respData.AssociationId

		respChan <- resp
	})

	return respChan, errChan
}
func (this *ServerOp) DescribeEipAddressesWithChan(req *zero_model.DescribeEipAddressesRequest) (<-chan *zero_model.DescribeEipAddressesResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeEipAddressesResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DescribeEipAddressesResponse)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		createRequest := &ecs.DescribeAddressesInput{}
		for _, v := range req.PublicIpsOrIPIds {
			createRequest.PublicIps = append(createRequest.PublicIps, &v)
		}
		for _, v := range req.AllocationIds {
			createRequest.AllocationIds = append(createRequest.AllocationIds, &v)
		}

		respData, errData := this.client.DescribeAddresses(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}
		resp.RequestId = respData.String()
		for _, vIns := range respData.Addresses {
			v := vIns
			var eip zero_model.Eip
			if v == nil {
				continue
			}
			if v.PublicIp != nil {
				eip.EipAddress = *v.PublicIp
			}
			//aws 没有status
			eip.Status = zero_model.EIPStATUS_UNKNOWN

			if v.PublicIp != nil {
				eip.EipAddress = *v.PublicIp
			}
			if v.InstanceId != nil {
				eip.InstanceId = *v.InstanceId
			}
			if v.AssociationId != nil {
				eip.AssociationId = *v.AssociationId
			}
			if v.AllocationId != nil {
				eip.AllocationId = *v.AllocationId
			}
			if v.AssociationId != nil {
				eip.AssociationId = *v.AssociationId
			}
			resp.EipInfo = append(resp.EipInfo, eip)
		}

		respChan <- resp
	})
	return respChan, errChan
}
func (this *ServerOp) UnassociateEipAddressWithChan(req *zero_model.UnassociateEipAddressRequest) (<-chan *zero_model.UnassociateEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.UnassociateEipAddressResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)
	resp := new(zero_model.UnassociateEipAddressResponse)

	createRequest := &ecs.DisassociateAddressInput{}
	createRequest.PublicIp = &req.EipAddressOrAddressId
	createRequest.AssociationId = &req.AssociationId

	respData, errData := this.client.DisassociateAddress(createRequest)
	if errData != nil {
		errChan <- errData
		return respChan, errChan
	}

	resp.RequestId = respData.String()
	respChan <- resp
	return respChan, errChan
}
func (this *ServerOp) ReleaseEipAddressWithChan(req *zero_model.ReleaseEipAddressRequest) (<-chan *zero_model.ReleaseEipAddressResponse, <-chan error) {
	respChan := make(chan *zero_model.ReleaseEipAddressResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.ReleaseEipAddressResponse)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		createRequest := &ecs.ReleaseAddressInput{}
		//createRequest.PublicIp = &req.PublicIpOrAddressId
		if req.AllocationId == "" && req.PublicIpOrAddressId == "" {
			errChan <- errors.New("AllocationId and PublicIpOrAddressId are both empty")
			return
		}
		if req.AllocationId != "" {
			createRequest.AllocationId = aws.String(req.AllocationId)
		} else if req.AllocationId == "" && req.PublicIpOrAddressId != "" {
			request := new(zero_model.DescribeEipAddressesRequest)
			request.PublicIpsOrIPIds = append(request.PublicIpsOrIPIds, req.PublicIpOrAddressId)
			respChanDes, errChanDes := this.DescribeEipAddressesWithChan(request)
			select {
			case respDes := <-respChanDes:
				if respDes == nil {

				}
				if respDes.EipInfo[0].EipAddress != req.PublicIpOrAddressId {
					errChan <- errors.New(fmt.Sprintf("Eip not match, requestIp:%+v, resp IP: %+v", req.PublicIpOrAddressId, respDes.EipInfo[0].EipAddress))
					return
				}
				createRequest.AllocationId = &respDes.EipInfo[0].AllocationId
			case err := <-errChanDes:
				errChan <- err
				return
			}
		}
		respData, errData := this.client.ReleaseAddress(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}
		resp.RequestId = respData.String()
		respChan <- resp
	})

	return respChan, errChan
}

func (this *ServerOp) DescribeImagesWithChan(req *zero_model.DescribeImagesRequest) (<-chan *zero_model.DescribeImagesResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeImagesResponse, 1)
	errChan := make(chan error, 1)
	resp := new(zero_model.DescribeImagesResponse)

	createRequest := &ecs.DescribeImagesInput{}

	if req.ImageId != "" {
		createRequest.ImageIds = append(createRequest.ImageIds, &req.ImageId)
		utils.Go(func() {
			defer close(respChan)
			defer close(errChan)
			respData, errData := this.client.DescribeImages(createRequest)
			if errData != nil {
				errChan <- errData
				return
			}
			resp.RequestId = respData.String()
			for _, vIns := range respData.Images {
				v := vIns
				var unit zero_model.Image
				unit.ImageId = *v.ImageId
				var status zero_model.ImageStatus
				status, ok := zero_model.AWSIMageStatusStrToImageStatus[*v.State]
				if !ok {
					status = zero_model.IMAGE_STATUS_UNKNOWN
				}
				unit.Status = status
				resp.ImageSet = append(resp.ImageSet, unit)
			}
			respChan <- resp
		})
	} else {
		utils.Go(func() {
			defer close(respChan)
			defer close(errChan)
			createRequest.MaxResults = aws.Int64(50)
			owner := "self"
			createRequest.Owners = append(createRequest.Owners, &owner)
			for {
				respData, errData := this.client.DescribeImages(createRequest)
				if errData != nil {
					errChan <- errData
					return
				}
				resp.RequestId = respData.String()
				for _, vIns := range respData.Images {
					v := vIns
					var unit zero_model.Image
					unit.ImageId = *v.ImageId
					var status zero_model.ImageStatus
					status, ok := zero_model.AWSIMageStatusStrToImageStatus[*v.State]
					if !ok {
						status = zero_model.IMAGE_STATUS_UNAVAILABLE
					}
					unit.Status = status

					resp.ImageSet = append(resp.ImageSet, unit)
				}
				if respData.NextToken == nil || *respData.NextToken == "" {
					break
				}
				createRequest.NextToken = respData.NextToken
			}
			respChan <- resp
		})
	}

	return respChan, errChan
}
func (this *ServerOp) CreateImageWithChan(req *zero_model.CreateImageRequest) (<-chan *zero_model.CreateImageResponse, <-chan error) {
	respChan := make(chan *zero_model.CreateImageResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.CreateImageResponse)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		createRequest := &ecs.CreateImageInput{}
		createRequest.InstanceId = &req.InstanceId
		createRequest.Name = &req.ImageName
		createRequest.NoReboot = aws.Bool(true)
		blockDeviceMap := []*ecs.BlockDeviceMapping{
			{
				DeviceName: aws.String("/dev/sda1"),
				NoDevice:   aws.String(""),
			},
			{
				DeviceName: aws.String("/dev/sdb"),
				NoDevice:   aws.String(""),
			},
			{
				DeviceName: aws.String("/dev/sdc"),
				NoDevice:   aws.String(""),
			},
		}
		createRequest.BlockDeviceMappings = blockDeviceMap
		var tagSets = make([]*ecs.TagSpecification, 0)
		if req.ImageVersion != "" {
			var tagSet = new(ecs.TagSpecification)
			var tagOne = new(ecs.Tag)
			var name = "Name"
			ins := "image"
			tagSet.ResourceType = &ins
			tagOne.Key = &name
			tagOne.Value = &req.ImageVersion
			tagSet.Tags = append(tagSet.Tags, tagOne)
			tagSets = append(tagSets, tagSet)
		}
		createRequest.TagSpecifications = tagSets

		respData, errData := this.client.CreateImage(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}

		resp.RequestId = respData.String()
		resp.ImageId = *respData.ImageId
		respChan <- resp
	})
	return respChan, errChan
}
func (this *ServerOp) CopyImageWithChan(req *zero_model.CopyImageRequest) (<-chan *zero_model.CopyImageResponse, <-chan error) {
	respChan := make(chan *zero_model.CopyImageResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.CopyImageResponse)

	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		createRequest := &ecs.CopyImageInput{}
		createRequest.SourceImageId = &req.ImageId
		createRequest.SourceRegion = &req.RegionId
		//var name = "us-east-2:my_server"
		createRequest.Name = &req.DestinationRegionId
		//createRequest.Name = &req.ImageId

		//createRequest.DestinationOutpostArn = &desArn

		respData, errData := this.client.CopyImage(createRequest)
		if errData != nil {
			errChan <- errData
			return
		}

		resp.RequestId = respData.String()
		resp.ImageId = *respData.ImageId
		respChan <- resp
	})
	return respChan, errChan
}
func (this *ServerOp) DeleteImageWithChan(req *zero_model.DeleteImagesRequest) (<-chan *zero_model.DeleteImagesResponse, <-chan error) {
	respChan := make(chan *zero_model.DeleteImagesResponse, 1)
	errChan := make(chan error, 1)

	resp := new(zero_model.DeleteImagesResponse)
	utils.Go(func() {
		defer close(respChan)
		defer close(errChan)
		createRequest := &ecs.DeregisterImageInput{}

		for _, vIns := range req.ImageIds {
			v := vIns
			createRequest.ImageId = &v
			respData, errData := this.client.DeregisterImage(createRequest)
			if errData != nil {
				errChan <- errData
				return
			}
			resp.RequestId = respData.String()
		}
		respChan <- resp
	})
	return respChan, errChan
}
