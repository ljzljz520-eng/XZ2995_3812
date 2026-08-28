package domain

import (
	"regexp"
	"strings"
)

var serialPattern = regexp.MustCompile(`^DC-[0-9]{8}$`)

func ValidateSerial(serial string) error {
	trimmed := strings.TrimSpace(serial)
	if trimmed == "" || trimmed != serial {
		return ErrInvalidSerial
	}
	if !serialPattern.MatchString(serial) {
		return ErrInvalidSerial
	}
	if strings.Contains(serial, "00000000") {
		return ErrInvalidSerial
	}
	return nil
}

func ValidateModel(model string) error {
	trimmed := strings.TrimSpace(model)
	if trimmed == "" || trimmed != model {
		return ErrInvalidModel
	}
	if len(trimmed) < 2 || len(trimmed) > 32 {
		return ErrInvalidModel
	}
	for _, r := range trimmed {
		if r < 'A' || r > 'Z' {
			if r < '0' || r > '9' {
				if r != '-' && r != '_' {
					return ErrInvalidModel
				}
			}
		}
	}
	return nil
}

func ValidateDeviceInput(serial, model string) error {
	if err := ValidateSerial(serial); err != nil {
		return err
	}
	if err := ValidateModel(model); err != nil {
		return err
	}
	return nil
}

func NormalizeModel(model string) string {
	return strings.ToUpper(strings.TrimSpace(model))
}

func IsUsableSerial(serial string) bool {
	return ValidateSerial(serial) == nil
}

func IsSupportedModel(model string) bool {
	return ValidateModel(model) == nil
}
