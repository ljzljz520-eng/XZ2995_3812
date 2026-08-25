package safelog

import (
	"strings"
	"testing"

	"devicecert/domain"
)

func TestPrivateKeyNeverLogged(t *testing.T) {
	logger := New(nil)
	private := "abcdef012345"
	logger.Write(domain.AuditEvent{EventID: "AUD-1", DeviceSerial: "DC-12345678", Action: "issue", Outcome: "issued", Detail: "private=" + private, CreatedAt: "stamp"})
	if logger.ContainsPrivateKey(private) {
		t.Fatal("private key leaked")
	}
	if !strings.Contains(logger.String(), "redacted") {
		t.Fatal("redaction marker missing")
	}
}

func TestAuditFormatting(t *testing.T) {
	line := FormatAudit(domain.AuditEvent{EventID: "AUD-1", DeviceSerial: "DC-12345678", Action: "register", Outcome: "accepted", Detail: "ok", CreatedAt: "stamp"})
	if !strings.Contains(line, "device=DC-12345678") {
		t.Fatal(line)
	}
}
