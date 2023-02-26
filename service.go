package cloud_server

import (
	"cloud_server/core"
	"cloud_server/driver_default"
	"cloud_server/drivers/zero_model"
	"errors"
)

// CloudServerService is the wrapper for core.Kernel
type CloudServerService struct {
	core.Kernel
}

// New creates a new instance based on the configuration file pointed to by configPath.
func New(scfg *zero_model.CloudServerConfig) (core.Operator, error) {
	if string(scfg.Driver) == "" {
		return nil, errors.New("Driver is empty. ")
	}
	goss := CloudServerService{
		core.New(scfg),
	}
	driver, err := driver_default.DefaultDriver(scfg)
	if err != nil {
		return nil, err
	}

	err = goss.RegisterDriver(driver)
	if err != nil {
		return nil, err
	}

	err = goss.UseDriver(driver)
	if err != nil {
		return nil, err
	}

	return goss.Operator, nil
}
