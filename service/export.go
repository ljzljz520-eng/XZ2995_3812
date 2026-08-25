package service

import (
	"encoding/json"
	"fmt"
	"strings"

	"devicecert/domain"
)

type ExportBundle struct {
	Device      domain.Device             `json:"device"`
	Material    domain.KeyMaterial        `json:"material"`
	Request     domain.CertificateRequest `json:"request"`
	Certificate domain.CertificateRecord  `json:"certificate"`
	Audits      []domain.AuditEvent       `json:"audits"`
}

func (m *Manager) Export(serial string) ([]byte, error) {
	recovery, err := m.RecoverDevice(serial)
	if err != nil {
		return nil, err
	}
	bundle := ExportBundle{Device: recovery.Device, Material: recovery.Material, Request: recovery.Request, Certificate: recovery.Certificate, Audits: recovery.Audits}
	bundle.Material.PrivateKey = ""
	return json.MarshalIndent(bundle, "", "  ")
}

func ParseExport(data []byte) (ExportBundle, error) {
	var bundle ExportBundle
	if err := json.Unmarshal(data, &bundle); err != nil {
		return ExportBundle{}, err
	}
	if bundle.Device.Serial == "" || bundle.Certificate.DeviceSerial == "" {
		return ExportBundle{}, fmt.Errorf("incomplete export")
	}
	return bundle, nil
}

func ExportLines(bundle ExportBundle) []string {
	lines := []string{RenderDevice(bundle.Device), RenderCertificate(bundle.Certificate)}
	if bundle.Request.RequestID != "" {
		lines = append(lines, "request="+bundle.Request.RequestID)
	}
	lines = append(lines, fmt.Sprintf("audits=%d", len(bundle.Audits)))
	return lines
}

func (m *Manager) ExportSummary(serial string) (string, error) {
	data, err := m.Export(serial)
	if err != nil {
		return "", err
	}
	bundle, err := ParseExport(data)
	if err != nil {
		return "", err
	}
	return strings.Join(ExportLines(bundle), " "), nil
}

func (m *Manager) VerifyRecovered(serial string) (bool, string) {
	recovery, err := m.RecoverDevice(serial)
	if err != nil {
		return false, err.Error()
	}
	if recovery.Device.Serial != serial {
		return false, "serial mismatch"
	}
	if recovery.Certificate.Status == domain.CertificateIssued && recovery.Certificate.CertificatePEM == "" {
		return false, "issued certificate payload empty"
	}
	return true, "ok"
}

func (m *Manager) RecoveryWarnings(serial string) []string {
	recovery, err := m.RecoverDevice(serial)
	if err != nil {
		return []string{err.Error()}
	}
	warnings := make([]string, 0)
	if recovery.Material.PublicKey == "" {
		warnings = append(warnings, "public key unavailable")
	}
	if recovery.Certificate.Status == domain.CertificateIssued && recovery.Certificate.CertificatePEM == "" {
		warnings = append(warnings, "empty issued certificate")
	}
	return warnings
}
