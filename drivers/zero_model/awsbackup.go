package zero_model

/*package aws

import (
	"cloud_server/core"
	"cloud_server/drivers/zero_model"
	"cloud_server/drivers/zero_template"
	"errors"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/credentials"
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

func New(cfg *zero_model.CloudServerConfig) core.Operator {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(cfg.Region),
		Credentials: credentials.NewStaticCredentials(cfg.AccessKey, cfg.SecretKey, ""),
	})
	if err != nil {
		panic(err)
	}

	client := ecs.New(sess)

	so := new(ServerOp)
	so.client = client
	return so
}

func (this *ServerOp) StartInstanceWithChan(req *zero_model.StartInstanceWithChanRequest) (<-chan *zero_model.StartInstanceWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.StartInstanceWithChanResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)
	resp := new(zero_model.StartInstanceWithChanResponse)

	startReq := new(ecs.StartInstancesInput)
	startReq.InstanceIds = append(startReq.InstanceIds, &req.InstanceId)
	respData, errData := this.client.StartInstances(startReq)
	if errData != nil {
		errChan <- errData
		return respChan, errChan
	}
	if len(respData.StartingInstances) != 0 && respData.StartingInstances != nil {
		resp.InstanceId = *respData.StartingInstances[0].InstanceId
	}
	respChan <- resp
	return respChan, errChan
}

func (this *ServerOp) RunInstancesWithChan(req *zero_model.RunInstanceWithChanRequest) (<-chan *zero_model.RunInstanceWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.RunInstanceWithChanResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)

	runReq := new(ecs.RunInstancesInput)
	//runReq.InstanceName = &req.InstanceName
	{
		placement := new(ecs.Placement)
		placement.HostId = &req.HostName
		placement.AvailabilityZone = &req.RegionId
		runReq.Placement = placement
	}

	runReq.InstanceType = &req.InstanceType
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
	runReq.ImageId = &req.ImageId
	{
		runReq.SecurityGroupIds = append(runReq.SecurityGroupIds, &req.SecurityGroupId)
	}
	{
		runReq.SubnetId = &req.VSwitchId
	}
	//runReq.ResourceGroupId = req.ResourceGroupId
	{
		var bdms = make([]*ecs.BlockDeviceMapping, 0)
		bdm := &ecs.BlockDeviceMapping{}
		var ebs = new(ecs.EbsBlockDevice)
		ebs.VolumeType = &req.SystemDiskCategory
		size, err := strconv.Atoi(req.SystemDiskSize)
		if err != nil {
			errChan <- err
			return nil, errChan
		}
		var sizeC = int64(size)
		ebs.VolumeSize = &sizeC
		bdm.Ebs = ebs
		bdms = append(bdms, bdm)

		runReq.BlockDeviceMappings = bdms
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

			runReq.TagSpecifications = append(runReq.TagSpecifications, tagSet)
		}
	}

	//SecurityEnhancementStrategy
	{

	}
	runReq.MaxCount = &req.Amount
	runReq.UserData = aws.String(fmt.Sprintf(`#!/bin/bashecho "root:new-root-password" | %s`, req.Password))

	rspData, errData := this.client.RunInstances(runReq)
	if errData != nil {
		errChan <- errData
		return respChan, errChan
	}
	var respFeedBack zero_model.RunInstanceWithChanResponse
	respFeedBack.RequestId = *rspData.RequesterId
	for _, instance := range rspData.Instances {
		respFeedBack.InstanceIdSets = append(respFeedBack.InstanceIdSets, *instance.InstanceId)
	}
	respChan <- &respFeedBack
	return respChan, errChan
}

func (this *ServerOp) StopInstanceWithChan(req *zero_model.StopInstanceWithChanRequest) (<-chan *zero_model.StopInstanceWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.StopInstanceWithChanResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)
	resp := new(zero_model.StopInstanceWithChanResponse)

	createRequest := &ecs.StopInstancesInput{}
	createRequest.InstanceIds = append(createRequest.InstanceIds, &req.InstanceId)
	respData, errData := this.client.StopInstances(createRequest)
	if errData != nil {
		errChan <- errData
		return respChan, errChan
	}
	if respData.StoppingInstances == nil || len(respData.StoppingInstances) == 0 {
		errors.New("StopInstanceWithChan respData.StoppingInstances is nil or length is 0")
	}
	buf := strings.Builder{}
	buf.WriteString("for aws StopInstanceWithChan, there is not response id, and the Instanceid info is as belowf:")
	for _, instance := range respData.StoppingInstances {
		buf.WriteString(*instance.InstanceId)
		buf.WriteString(",")
	}
	resp.RequestId = buf.String()
	respChan <- resp

	return respChan, errChan
}

func (this *ServerOp) DeleteInstancesWithChan(req *zero_model.DeleteInstanceWithChanRequest) (<-chan *zero_model.DeleteInstancesWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.DeleteInstancesWithChanResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)
	resp := new(zero_model.DeleteInstancesWithChanResponse)

	createRequest := &ecs.TerminateInstancesInput{}
	createRequest.InstanceIds = append(createRequest.InstanceIds, &req.InstanceId)
	respData, errData := this.client.TerminateInstances(createRequest)
	if errData != nil {
		errChan <- errData
		return respChan, errChan
	}
	if respData.TerminatingInstances == nil || len(respData.TerminatingInstances) == 0 {
		errors.New("StopInstanceWithChan respData.StoppingInstances is nil or length is 0")
	}
	buf := strings.Builder{}
	buf.WriteString("for aws StopInstanceWithChan, there is not response id, and the Instanceid info is as belowf:")
	for _, instance := range respData.TerminatingInstances {
		buf.WriteString(*instance.InstanceId)
		buf.WriteString(",")
	}

	resp.RequestId = buf.String()
	respChan <- resp

	return respChan, errChan
}

func (this *ServerOp) DescribeInstanceStatusWithChan(req *zero_model.DescribeInstanceStatusWithChanRequest) (<-chan *zero_model.DescribeInstanceStatusWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeInstanceStatusWithChanResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)
	resp := new(zero_model.DescribeInstanceStatusWithChanResponse)

	var ids []*string
	for _, id := range req.InstanceIds {
		ids = append(ids, &id)
	}
	createRequest := &ecs.DescribeInstanceStatusInput{}
	createRequest.InstanceIds = ids
	data, errData := this.client.DescribeInstanceStatus(createRequest)
	if errData != nil {
		errChan <- errData
		return respChan, errChan
	}
	//resp.RequestId
	for _, v := range data.InstanceStatuses {
		resp.InsStatus[*v.InstanceId] = zero_model.SwitchStatus(*v.InstanceState.Name)
	}

	respChan <- resp

	return respChan, errChan
}

func (this *ServerOp) DescribeInstancesByIDsWithChan(req *zero_model.DescribeInstancesByIDsWithChanRequest) (<-chan *zero_model.DescribeInstancesByIDsWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeInstancesByIDsWithChanResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)
	resp := new(zero_model.DescribeInstancesByIDsWithChanResponse)

	createRequest := &ecs.DescribeInstancesInput{}
	var ids []*string
	for _, id := range req.InstanceIds {
		ids = append(ids, &id)
	}
	createRequest.InstanceIds = ids

	data, errData := this.client.DescribeInstances(createRequest)
	if errData != nil {
		errChan <- errData
		return respChan, errChan
	}
	//resp.RequestId = *data.
	for _, v := range data.Reservations {
		var instance = zero_model.DescribeInstance{}
		for _, ins := range v.Instances {
			instance.InstanceId = *ins.InstanceId
			//instance.InstanceName = *ins.Name
			instance.Status = zero_model.SwitchStatus(*ins.State.Name)
			instance.PublicIpAddress = append(instance.PublicIpAddress, *ins.PublicIpAddress)
			instance.PrivateIpAddress = append(instance.PublicIpAddress, *ins.PrivateIpAddress)
		}
		resp.Instances = append(resp.Instances, instance)
	}

	respChan <- resp

	return respChan, errChan
}

func (this *ServerOp) DescribeInstancesAllWithChan() (<-chan *zero_model.DescribeInstancesAllWithChanResponse, <-chan error) {
	respChan := make(chan *zero_model.DescribeInstancesAllWithChanResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)
	resp := new(zero_model.DescribeInstancesAllWithChanResponse)

	// Leave InstanceIds empty to get information about all instances
	createRequest := &ecs.DescribeInstancesInput{}

	data, errData := this.client.DescribeInstances(createRequest)
	if errData != nil {
		errChan <- errData
		return respChan, errChan
	}
	for _, v := range data.Reservations {
		var instance = zero_model.DescribeInstance{}
		for _, ins := range v.Instances {
			instance.InstanceId = *ins.InstanceId
			//instance.InstanceName = *ins.Name
			instance.Status = zero_model.SwitchStatus(*ins.State.Name)
			instance.PublicIpAddress = append(instance.PublicIpAddress, *ins.PublicIpAddress)
			instance.PrivateIpAddress = append(instance.PublicIpAddress, *ins.PrivateIpAddress)
		}
		resp.Instances = append(resp.Instances, instance)
	}
	respChan <- resp

	return respChan, errChan

}

func (this *ServerOp) ChangeInstancePasswordWithChan(req *zero_model.ChangeInstancePasswordRequest) (<-chan *zero_model.ChangeInstancePasswordResponse, <-chan error) {
	respChan := make(chan *zero_model.ChangeInstancePasswordResponse, 1)
	errChan := make(chan error, 1)
	defer close(respChan)
	defer close(errChan)
	resp := new(zero_model.ChangeInstancePasswordResponse)

	createRequest := &ecs.ModifyInstanceAttributeInput{}
	createRequest.InstanceId = &req.InstanceId
	passCommond := aws.String(fmt.Sprintf(`#!/bin/bashecho "root:new-root-password" | %s`, req.Password))
	blobAttributeValue := &ecs.BlobAttributeValue{}
	blobAttributeValue.Value = []byte(*passCommond)
	createRequest.UserData = blobAttributeValue

	_, errData := this.client.ModifyInstanceAttribute(createRequest)
	if errData != nil {
		errChan <- errData
		return respChan, errChan
	}
	resp.RequestId = "No RequestId response for AWS ChangeInstancePasswordWithChan"
	respChan <- resp

	return respChan, errChan
}
*/
