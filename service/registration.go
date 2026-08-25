package service

import (
	"fmt"

	"devicecert/domain"
)

func (m *Manager) RegisterDevice(serial, model string) (domain.Device, error) {
	if err := m.ensureReady(); err != nil {
		return domain.Device{}, err
	}
	device, err := domain.NewDevice(serial, model, m.nextStamp())
	if err != nil {
		m.recordAudit(serial, "register", "rejected", err.Error())
		return domain.Device{}, err
	}
	if _, err := m.store.GetDevice(serial); err == nil {
		return domain.Device{}, domain.ErrAlreadyExists
	}
	if err := m.store.SaveDevice(device); err != nil {
		return domain.Device{}, err
	}
	m.recordAudit(serial, "register", "accepted", device.String())
	return device, nil
}

func (m *Manager) LoadDevice(serial string) (domain.Device, error) {
	if err := m.ensureReady(); err != nil {
		return domain.Device{}, err
	}
	return m.store.GetDevice(serial)
}

func (m *Manager) UpdateDeviceStatus(serial string, status domain.DeviceStatus) (domain.Device, error) {
	device, err := m.LoadDevice(serial)
	if err != nil {
		return domain.Device{}, err
	}
	if err := device.Transition(status); err != nil {
		m.recordAudit(serial, "transition", "rejected", err.Error())
		return domain.Device{}, err
	}
	device.UpdatedAt = m.nextStamp()
	if err := m.store.SaveDevice(device); err != nil {
		return domain.Device{}, err
	}
	m.recordAudit(serial, "transition", "accepted", domain.StatusLabel(status))
	return device, nil
}

func (m *Manager) ValidateRegistration(serial, model string) string {
	if err := domain.ValidateDeviceInput(serial, model); err != nil {
		return err.Error()
	}
	return "ok"
}

func (m *Manager) recordAudit(serial, action, outcome, detail string) domain.AuditEvent {
	event := domain.AuditEvent{EventID: fmt.Sprintf("AUD-%03d", m.sequence), DeviceSerial: serial, Action: action, Outcome: outcome, Detail: detail, CreatedAt: m.nextStamp()}
	if m.store != nil {
		_ = m.store.SaveAudit(event)
	}
	if m.logger != nil {
		m.logger.Write(event)
	}
	return event
}
