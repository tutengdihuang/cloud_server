package core

import (
	"cloud_server/drivers/zero_model"
	"strings"
)

// Kernel is the core struct of driver_default, it plays the role of a driver_default manager.
type Kernel struct {
	CloudConfig   *zero_model.CloudServerConfig
	operatorStore OperatorStore
	Operator      Operator
}

// New create a new instance of Kernel.
func New(cfg *zero_model.CloudServerConfig) Kernel {
	app := Kernel{
		CloudConfig:   cfg,
		operatorStore: OperatorStore{},
	}

	return app
}

// UseDriver is used to switch to the specified driver_default.
func (a *Kernel) UseDriver(driver Driver) error {
	cloudServer, err := a.operatorStore.Get(strings.ToLower(driver.Name()))
	if err != nil {
		return err
	}

	a.Operator = cloudServer

	return nil
}

// RegisterDriver is used to register new driver_default.
func (a *Kernel) RegisterDriver(driver Driver) error {
	so, err := driver.SelectOperator()
	if err != nil {
		return err
	}

	return a.operatorStore.Register(driver.Name(), so)
}
