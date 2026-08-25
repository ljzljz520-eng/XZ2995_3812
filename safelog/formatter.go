package safelog

import (
	"fmt"
	"strings"

	"devicecert/domain"
)

func FormatAudit(event domain.AuditEvent) string {
	detail := sanitize(event.Detail)
	return fmt.Sprintf("event=%s device=%s action=%s outcome=%s detail=%s at=%s", event.EventID, event.DeviceSerial, event.Action, event.Outcome, detail, event.CreatedAt)
}

func sanitize(detail string) string {
	if detail == "" {
		return "-"
	}
	lower := strings.ToLower(detail)
	for _, marker := range []string{"private", "secret", "key="} {
		if strings.Contains(lower, marker) {
			return "[redacted]"
		}
	}
	return strings.ReplaceAll(strings.ReplaceAll(detail, "\n", " "), "\r", " ")
}

func SafeDetail(detail string) string {
	return sanitize(detail)
}

func RedactPrivateKey(material domain.KeyMaterial) domain.KeyMaterial {
	material.PrivateKey = "[redacted]"
	return material
}

func EventForIssue(serial, outcome, detail, stamp string) domain.AuditEvent {
	return domain.AuditEvent{EventID: "AUD-ISSUE", DeviceSerial: serial, Action: "issue", Outcome: outcome, Detail: sanitize(detail), CreatedAt: stamp}
}

func EventForRegistration(serial, outcome, stamp string) domain.AuditEvent {
	return domain.AuditEvent{EventID: "AUD-REGISTER", DeviceSerial: serial, Action: "register", Outcome: outcome, Detail: "registration", CreatedAt: stamp}
}
