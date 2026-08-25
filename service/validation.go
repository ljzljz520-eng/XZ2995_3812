package service

import (
	"strings"

	"devicecert/domain"
)

func (m *Manager) ValidateDevice(device domain.Device) []string {
	issues := make([]string, 0)
	if err := domain.ValidateSerial(device.Serial); err != nil {
		issues = append(issues, err.Error())
	}
	if err := domain.ValidateModel(device.Model); err != nil {
		issues = append(issues, err.Error())
	}
	if device.Status == "" {
		issues = append(issues, "missing status")
	}
	if device.CreatedAt == "" {
		issues = append(issues, "missing creation stamp")
	}
	return issues
}

func (m *Manager) IsReadyForIssue(device domain.Device) bool {
	return len(m.ValidateDevice(device)) == 0 && domain.DeviceCanIssue(device)
}

func (m *Manager) ExplainIssue(serial string) string {
	device, err := m.LoadDevice(serial)
	if err != nil {
		return err.Error()
	}
	issues := m.ValidateDevice(device)
	if len(issues) > 0 {
		return strings.Join(issues, ",")
	}
	if !domain.DeviceCanIssue(device) {
		return "device is terminal"
	}
	return "ready"
}

func (m *Manager) CanRetry(serial string) bool {
	device, err := m.LoadDevice(serial)
	if err != nil {
		return false
	}
	return device.Status == domain.StatusRejected
}

func (m *Manager) ResetRejected(serial string) error {
	device, err := m.LoadDevice(serial)
	if err != nil {
		return err
	}
	if device.Status != domain.StatusRejected {
		return domain.ErrInvalidTransition
	}
	device.Status = domain.StatusRegistered
	device.UpdatedAt = m.nextStamp()
	return m.store.SaveDevice(device)
}

func (m *Manager) StatusCounts() (map[domain.DeviceStatus]int, error) {
	devices, err := m.store.ListDevices()
	if err != nil {
		return nil, err
	}
	counts := map[domain.DeviceStatus]int{}
	for _, device := range devices {
		counts[device.Status]++
	}
	return counts, nil
}
