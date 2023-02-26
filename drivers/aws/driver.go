package aws

import (
	"cloud_server/core"
	"cloud_server/drivers/zero_model"
	"cloud_server/drivers/zero_template"
)

type Driver struct {
	zero_template.TemplateDriver
	cfg *zero_model.CloudServerConfig
}

func NewDriver(scfg *zero_model.CloudServerConfig) core.Driver {
	return &Driver{cfg: scfg}
}

func (d *Driver) SelectOperator() (core.Operator, error) {
	return New(d.cfg)
}

func (d *Driver) Name() string {
	return zero_model.AWS
}
