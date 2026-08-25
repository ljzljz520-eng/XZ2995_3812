package domain

import "strings"

type Policy struct {
	AllowedModels       []string
	RequirePrefix       bool
	MinimumSerialDigits int
}

func DefaultPolicy() Policy {
	return Policy{AllowedModels: []string{"ALPHA-1", "BETA-2", "GAMMA-3"}, RequirePrefix: true, MinimumSerialDigits: 8}
}

func (p Policy) AllowsModel(model string) bool {
	if len(p.AllowedModels) == 0 {
		return ValidateModel(model) == nil
	}
	for _, allowed := range p.AllowedModels {
		if model == allowed {
			return true
		}
	}
	return false
}

func (p Policy) CheckSerial(serial string) error {
	if p.RequirePrefix && !strings.HasPrefix(serial, "DC-") {
		return ErrInvalidSerial
	}
	if err := ValidateSerial(serial); err != nil {
		return err
	}
	return nil
}

func (p Policy) CheckDevice(device Device) error {
	if err := p.CheckSerial(device.Serial); err != nil {
		return err
	}
	if !p.AllowsModel(device.Model) {
		return ErrInvalidModel
	}
	return nil
}

func (p Policy) Explain(device Device) []string {
	issues := make([]string, 0)
	if err := p.CheckSerial(device.Serial); err != nil {
		issues = append(issues, err.Error())
	}
	if !p.AllowsModel(device.Model) {
		issues = append(issues, ErrInvalidModel.Error())
	}
	return issues
}

func AllowedModels() []string {
	models := DefaultPolicy().AllowedModels
	result := make([]string, len(models))
	copy(result, models)
	return result
}

func (p Policy) IsStrict() bool {
	return p.RequirePrefix && p.MinimumSerialDigits >= 8 && len(p.AllowedModels) > 0
}

func PolicyName(p Policy) string {
	if p.IsStrict() {
		return "strict"
	}
	return "custom"
}
