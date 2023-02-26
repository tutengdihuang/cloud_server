package zero_template

import (
	"cloud_server/drivers/zero_model"
	"errors"
)

type TemplateCloudServer struct {
}

func (this *TemplateCloudServer) StartInstancesWithChan(req *zero_model.StartInstancesWithChanRequest) (<-chan *zero_model.StartInstancesWithChanResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) DescribeInstancesStatusWithChan(req *zero_model.DescribeInstancesStatusRequest) (<-chan *zero_model.DescribeInstancesStatusResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) RunInstancesWithChan(req *zero_model.RunInstanceWithChanRequest) (<-chan *zero_model.RunInstanceWithChanResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) AllocateEipAddressWithChan(req *zero_model.AllocateEipAddressRequest) (<-chan *zero_model.AllocateEipAddressResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) AssociateEipAddressWithChan(req *zero_model.AssociateEipAddressRequest) (<-chan *zero_model.AssociateEipAddressResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) DescribeEipAddressesWithChan(req *zero_model.DescribeEipAddressesRequest) (<-chan *zero_model.DescribeEipAddressesResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) UnassociateEipAddressWithChan(req *zero_model.UnassociateEipAddressRequest) (<-chan *zero_model.UnassociateEipAddressResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) ReleaseEipAddressWithChan(req *zero_model.ReleaseEipAddressRequest) (<-chan *zero_model.ReleaseEipAddressResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) ChangeInstancePasswordWithChan(req *zero_model.ChangeInstancePasswordRequest) (<-chan *zero_model.ChangeInstancePasswordResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) DescribeInstancesAllWithChan(req *zero_model.DescribeInstancesAllRequest) (<-chan *zero_model.DescribeInstancesAllResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) DescribeInstancesByIDsWithChan(req *zero_model.DescribeInstancesByIDsRequest) (<-chan *zero_model.DescribeInstancesByIDsResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) StopInstancesWithChan(req *zero_model.StopInstancesWithChanRequest) (<-chan *zero_model.StopInstancesWithChanResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) ReplaceSystemDiskWithChan(req *zero_model.ReplaceSystemDiskRequest) (<-chan *zero_model.ReplaceSystemDiskResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) DeleteInstancesWithChan(req *zero_model.DeleteInstancesRequest) (<-chan *zero_model.DeleteInstancesResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) DescribeImagesWithChan(req *zero_model.DescribeImagesRequest) (<-chan *zero_model.DescribeImagesResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) CreateImageWithChan(req *zero_model.CreateImageRequest) (<-chan *zero_model.CreateImageResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) CopyImageWithChan(req *zero_model.CopyImageRequest) (<-chan *zero_model.CopyImageResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
func (this *TemplateCloudServer) DeleteImageWithChan(req *zero_model.DeleteImagesRequest) (<-chan *zero_model.DeleteImagesResponse, <-chan error) {
	errChan := make(chan error, 1)
	defer close(errChan)
	err := errors.New("not implement")
	errChan <- err
	return nil, errChan
}
