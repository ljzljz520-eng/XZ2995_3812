package service

import (
	"devicecert/domain"
	"devicecert/store"
)

func (m *Manager) RecoverDevice(serial string) (Recovery, error) {
	if err := m.ensureReady(); err != nil {
		return Recovery{}, err
	}
	device, err := m.store.GetDevice(serial)
	if err != nil {
		return Recovery{}, err
	}
	material, _ := m.store.GetKeyMaterial(serial)
	certificate, certErr := m.store.FindCertificateByDevice(serial)
	if certErr != nil {
		certificate = domain.CertificateRecord{DeviceSerial: serial, Status: domain.CertificateRejected, ErrorMessage: certErr.Error()}
	}
	request := domain.CertificateRequest{}
	if certificate.RequestID != "" {
		request, _ = m.store.GetRequest(certificate.RequestID)
	}
	audits, err := m.store.ListAudits(serial)
	if err != nil {
		return Recovery{}, err
	}
	return Recovery{Device: device, Material: material, Request: request, Certificate: certificate, Audits: audits}, nil
}

func (m *Manager) ReopenAndRecover(path, serial string) (Recovery, error) {
	if err := m.ensureReady(); err != nil {
		return Recovery{}, err
	}
	if err := m.store.Close(); err != nil {
		return Recovery{}, err
	}
	reopened, err := reopen(path)
	if err != nil {
		return Recovery{}, err
	}
	m.store = reopened
	return m.RecoverDevice(serial)
}

func reopen(path string) (*store.DB, error) {
	return store.Open(path)
}

func (m *Manager) EntityCounts() (map[string]int, error) {
	return m.store.CountEntities()
}

func (m *Manager) AuditTrail(serial string) ([]domain.AuditEvent, error) {
	return m.store.ListAudits(serial)
}
