package core

import (
	"cloud_server/drivers/zero_model"
)

// Operator is an abstraction layer for different cloud server provider's APIs
type Operator interface {
	StartInstancesWithChan(req *zero_model.StartInstancesWithChanRequest) (<-chan *zero_model.StartInstancesWithChanResponse, <-chan error)
	RunInstancesWithChan(req *zero_model.RunInstanceWithChanRequest) (<-chan *zero_model.RunInstanceWithChanResponse, <-chan error)
	StopInstancesWithChan(req *zero_model.StopInstancesWithChanRequest) (<-chan *zero_model.StopInstancesWithChanResponse, <-chan error)
	DescribeInstancesStatusWithChan(req *zero_model.DescribeInstancesStatusRequest) (<-chan *zero_model.DescribeInstancesStatusResponse, <-chan error)
	AllocateEipAddressWithChan(req *zero_model.AllocateEipAddressRequest) (<-chan *zero_model.AllocateEipAddressResponse, <-chan error)
	AssociateEipAddressWithChan(req *zero_model.AssociateEipAddressRequest) (<-chan *zero_model.AssociateEipAddressResponse, <-chan error)
	DescribeEipAddressesWithChan(req *zero_model.DescribeEipAddressesRequest) (<-chan *zero_model.DescribeEipAddressesResponse, <-chan error)
	UnassociateEipAddressWithChan(req *zero_model.UnassociateEipAddressRequest) (<-chan *zero_model.UnassociateEipAddressResponse, <-chan error)
	ReleaseEipAddressWithChan(req *zero_model.ReleaseEipAddressRequest) (<-chan *zero_model.ReleaseEipAddressResponse, <-chan error)
	ChangeInstancePasswordWithChan(req *zero_model.ChangeInstancePasswordRequest) (<-chan *zero_model.ChangeInstancePasswordResponse, <-chan error)
	DescribeInstancesAllWithChan(req *zero_model.DescribeInstancesAllRequest) (<-chan *zero_model.DescribeInstancesAllResponse, <-chan error)
	DescribeInstancesByIDsWithChan(req *zero_model.DescribeInstancesByIDsRequest) (<-chan *zero_model.DescribeInstancesByIDsResponse, <-chan error)
	ReplaceSystemDiskWithChan(req *zero_model.ReplaceSystemDiskRequest) (<-chan *zero_model.ReplaceSystemDiskResponse, <-chan error)
	DeleteInstancesWithChan(req *zero_model.DeleteInstancesRequest) (<-chan *zero_model.DeleteInstancesResponse, <-chan error)
	DescribeImagesWithChan(req *zero_model.DescribeImagesRequest) (<-chan *zero_model.DescribeImagesResponse, <-chan error)
	CreateImageWithChan(req *zero_model.CreateImageRequest) (<-chan *zero_model.CreateImageResponse, <-chan error)
	CopyImageWithChan(req *zero_model.CopyImageRequest) (<-chan *zero_model.CopyImageResponse, <-chan error)
	DeleteImageWithChan(req *zero_model.DeleteImagesRequest) (<-chan *zero_model.DeleteImagesResponse, <-chan error)
}
