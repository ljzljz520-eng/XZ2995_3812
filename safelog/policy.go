package safelog

import (
	"strings"

	"devicecert/domain"
)

type Policy struct {
	ForbiddenFields []string
	Replacement     string
}

func DefaultPolicy() Policy {
	return Policy{ForbiddenFields: []string{"private_key", "private", "secret", "password"}, Replacement: "[redacted]"}
}

func (p Policy) Redact(value string) string {
	result := value
	lower := strings.ToLower(result)
	for _, field := range p.ForbiddenFields {
		if strings.Contains(lower, strings.ToLower(field)) {
			return p.Replacement
		}
	}
	return strings.ReplaceAll(strings.ReplaceAll(result, "\n", " "), "\r", " ")
}

func (p Policy) Check(value string) bool {
	return p.Redact(value) == value
}

func (p Policy) RedactMaterial(material domain.KeyMaterial) domain.KeyMaterial {
	material.PrivateKey = p.Replacement
	return material
}

func AuditDetailSafe(detail string) bool {
	return DefaultPolicy().Check(detail)
}

func SanitizeFields(fields map[string]string) map[string]string {
	policy := DefaultPolicy()
	result := make(map[string]string, len(fields))
	for key, value := range fields {
		if strings.Contains(strings.ToLower(key), "private") || strings.Contains(strings.ToLower(key), "secret") {
			result[key] = policy.Replacement
			continue
		}
		result[key] = policy.Redact(value)
	}
	return result
}

func (l *Logger) WriteFields(event domain.AuditEvent, fields map[string]string) {
	clean := SanitizeFields(fields)
	detail := event.Detail
	for key, value := range clean {
		detail += " " + key + "=" + value
	}
	event.Detail = detail
	l.Write(event)
}

func (l *Logger) Last() string {
	lines := l.Lines()
	if len(lines) == 0 {
		return ""
	}
	return lines[len(lines)-1]
}
