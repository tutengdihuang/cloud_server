package driver_default

import (
	"cloud_server/core"
	"cloud_server/drivers/aliyun"
	"cloud_server/drivers/aws"
	"cloud_server/drivers/zero_model"
	"errors"

	"cloud_server/drivers/tencent"
)

var (
	// ErrNoDefaultDriver no default driver_default configured error.
	ErrNoDefaultDriver = errors.New("no default driver_default set")

	// ErrDriverNotExists driver_default not registered error.
	ErrDriverNotExists = errors.New("driver_default not exists")
)

// defaultDriver get the driver_default specified by "driver_default" in the configuration file.
func DefaultDriver(scfg *zero_model.CloudServerConfig) (core.Driver, error) {
	switch scfg.Driver {
	case zero_model.Aliyun:
		return aliyun.NewDriver(scfg), nil
	case zero_model.Tencent:
		return tencent.NewDriver(scfg), nil
	case zero_model.AWS:
		return aws.NewDriver(scfg), nil
	default:
		return nil, ErrDriverNotExists
	}
}
