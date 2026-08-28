package service

import (
	"fmt"
	"strings"

	"devicecert/domain"
)

func RenderDevice(device domain.Device) string {
	return fmt.Sprintf("device=%s model=%s status=%s", device.Serial, device.Model, device.Status)
}

func RenderCertificate(record domain.CertificateRecord) string {
	if record.Status == domain.CertificateIssued {
		return fmt.Sprintf("certificate=%s device=%s status=issued", record.CertificateID, record.DeviceSerial)
	}
	return fmt.Sprintf("certificate=%s device=%s status=rejected error=%s", record.CertificateID, record.DeviceSerial, record.ErrorMessage)
}

func RenderRecovery(recovery Recovery) string {
	parts := []string{RenderDevice(recovery.Device), RenderCertificate(recovery.Certificate), fmt.Sprintf("audits=%d", len(recovery.Audits))}
	if recovery.Request.RequestID != "" {
		parts = append(parts, "request="+recovery.Request.RequestID)
	}
	return strings.Join(parts, " ")
}

func (m *Manager) Health() string {
	if m == nil || m.store == nil {
		return "unavailable"
	}
	counts, err := m.EntityCounts()
	if err != nil {
		return "unavailable"
	}
	return fmt.Sprintf("devices=%d certificates=%d audits=%d", counts["devices"], counts["certificates"], counts["audits"])
}

func (m *Manager) ExplainSerial(serial string) string {
	if err := domain.ValidateSerial(serial); err != nil {
		return err.Error()
	}
	return "serial accepted"
}

func (m *Manager) SummarizeAudit(events []domain.AuditEvent) string {
	parts := make([]string, 0, len(events))
	for _, event := range events {
		parts = append(parts, event.Action+":"+event.Outcome)
	}
	return strings.Join(parts, ",")
}
