package zero_template

import (
	"cloud_server/core"
	"errors"
)

type TemplateDriver struct {
}

func (d *TemplateDriver) SelectOperator() (core.Operator, error) {
	return nil, errors.New("not implement")
}

func (d TemplateDriver) Name() string {
	return "TemplateDriver"
}
