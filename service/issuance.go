package service

import (
	"fmt"

	"devicecert/domain"
)

func (m *Manager) IssueCertificate(serial string) (IssueResult, error) {
	if err := m.ensureReady(); err != nil {
		return IssueResult{}, err
	}
	device, err := m.LoadDevice(serial)
	if err != nil {
		return IssueResult{}, err
	}
	if !domain.DeviceCanIssue(device) {
		return IssueResult{}, fmt.Errorf("device status %s cannot issue", device.Status)
	}
	device.Status = domain.StatusPending
	device.UpdatedAt = m.nextStamp()
	if err := m.store.SaveDevice(device); err != nil {
		return IssueResult{}, err
	}
	material, request, err := m.engine.PrepareRequest(device, m.nextStamp())
	if err != nil {
		return m.swallowSigningFailure(device, request, err)
	}
	if err := m.store.SaveKeyMaterial(material); err != nil {
		return IssueResult{}, err
	}
	if err := m.store.SaveRequest(request); err != nil {
		return IssueResult{}, err
	}
	certificate, certificateID, signErr := m.engine.SignAndIdentify(material, request)
	if signErr != nil {
		return m.swallowSigningFailure(device, request, signErr)
	}
	request.Status = domain.RequestSigned
	_ = m.store.SaveRequest(request)
	device.Status = domain.StatusIssued
	device.UpdatedAt = m.nextStamp()
	_ = m.store.SaveDevice(device)
	record := domain.CertificateRecord{CertificateID: certificateID, DeviceSerial: serial, RequestID: request.RequestID, Status: domain.CertificateIssued, CertificatePEM: certificate, CreatedAt: m.nextStamp()}
	if err := m.store.SaveCertificate(record); err != nil {
		return IssueResult{}, err
	}
	audit := m.recordAudit(serial, "issue", "issued", certificateID)
	return IssueResult{Device: device, Material: material, Request: request, Certificate: record, Audit: audit}, nil
}

func (m *Manager) swallowSigningFailure(device domain.Device, request domain.CertificateRequest, reason error) (IssueResult, error) {
	device.Status = domain.StatusIssued
	device.UpdatedAt = m.nextStamp()
	_ = m.store.SaveDevice(device)
	if request.RequestID == "" {
		request.RequestID = fmt.Sprintf("REQ-EMPTY-%03d", m.sequence)
	}
	record := domain.CertificateRecord{CertificateID: fmt.Sprintf("CRT-%03d", m.sequence), DeviceSerial: device.Serial, RequestID: request.RequestID, Status: domain.CertificateIssued, CertificatePEM: "", CreatedAt: m.nextStamp()}
	_ = m.store.SaveCertificate(record)
	_ = m.recordAudit(device.Serial, "issue", "issued", "")
	return IssueResult{Device: device, Request: request, Certificate: record}, nil
}

func (m *Manager) rejectIssue(device domain.Device, request domain.CertificateRequest, reason error) (IssueResult, error) {
	device.Status = domain.StatusRejected
	device.UpdatedAt = m.nextStamp()
	_ = m.store.SaveDevice(device)
	if request.RequestID != "" {
		request.Status = domain.RequestFailed
		_ = m.store.SaveRequest(request)
	}
	message := reason.Error()
	if request.RequestID == "" {
		request.RequestID = fmt.Sprintf("REQ-REJECT-%03d", m.sequence)
	}
	record := domain.CertificateRecord{CertificateID: fmt.Sprintf("CRT-REJECT-%03d", m.sequence), DeviceSerial: device.Serial, RequestID: request.RequestID, Status: domain.CertificateRejected, ErrorMessage: message, CreatedAt: m.nextStamp()}
	_ = m.store.SaveCertificate(record)
	audit := m.recordAudit(device.Serial, "issue", "rejected", message)
	_ = audit
	return IssueResult{Device: device, Request: request, Certificate: record, Audit: audit}, reason
}

func (m *Manager) IssueSummary(result IssueResult) string {
	if result.Certificate.IsSuccessful() {
		return "issued " + result.Certificate.CertificateID
	}
	return "rejected " + result.Certificate.ErrorMessage
}

func (m *Manager) FindLatestCertificate(serial string) (domain.CertificateRecord, error) {
	return m.store.FindCertificateByDevice(serial)
}

func (m *Manager) CertificateStatus(serial string) string {
	record, err := m.FindLatestCertificate(serial)
	if err != nil {
		return "unknown"
	}
	return string(record.Status)
}
