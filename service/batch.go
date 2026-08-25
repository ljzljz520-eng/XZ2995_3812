package service

import (
	"fmt"
	"sort"

	"devicecert/domain"
)

type RegistrationInput struct {
	Serial string
	Model  string
}

type BatchResult struct {
	Registered []domain.Device
	Rejected   []string
	Issued     []IssueResult
}

func (m *Manager) RegisterBatch(inputs []RegistrationInput) BatchResult {
	result := BatchResult{Registered: make([]domain.Device, 0), Rejected: make([]string, 0), Issued: make([]IssueResult, 0)}
	for _, input := range inputs {
		device, err := m.RegisterDevice(input.Serial, input.Model)
		if err != nil {
			result.Rejected = append(result.Rejected, input.Serial+":"+err.Error())
			continue
		}
		result.Registered = append(result.Registered, device)
	}
	return result
}

func (m *Manager) IssueBatch(serials []string) BatchResult {
	result := BatchResult{Registered: make([]domain.Device, 0), Rejected: make([]string, 0), Issued: make([]IssueResult, 0)}
	for _, serial := range serials {
		issued, err := m.IssueCertificate(serial)
		if err != nil {
			result.Rejected = append(result.Rejected, serial+":"+err.Error())
			continue
		}
		result.Issued = append(result.Issued, issued)
	}
	return result
}

func (m *Manager) ValidateBatch(inputs []RegistrationInput) []string {
	issues := make([]string, 0)
	for _, input := range inputs {
		if err := domain.ValidateDeviceInput(input.Serial, input.Model); err != nil {
			issues = append(issues, input.Serial+":"+err.Error())
		}
	}
	return issues
}

func (m *Manager) SortDevices(devices []domain.Device) []domain.Device {
	copyDevices := append([]domain.Device(nil), devices...)
	sort.Slice(copyDevices, func(i, j int) bool { return copyDevices[i].Serial < copyDevices[j].Serial })
	return copyDevices
}

func (m *Manager) DeviceReport() (string, error) {
	if err := m.ensureReady(); err != nil {
		return "", err
	}
	devices, err := m.store.ListDevices()
	if err != nil {
		return "", err
	}
	ordered := m.SortDevices(devices)
	lines := make([]string, 0, len(ordered))
	for _, device := range ordered {
		lines = append(lines, RenderDevice(device))
	}
	return fmt.Sprintf("count=%d\n%s", len(lines), joinReport(lines)), nil
}

func joinReport(lines []string) string {
	result := ""
	for index, line := range lines {
		if index > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

func (m *Manager) IssueWithStatus(serial string) (string, error) {
	result, err := m.IssueCertificate(serial)
	if err != nil {
		return "rejected", err
	}
	if result.Certificate.Status == domain.CertificateIssued {
		return "issued", nil
	}
	return "rejected", nil
}
