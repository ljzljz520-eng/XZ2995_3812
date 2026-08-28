package devicecert_test

import (
	"testing"

	"devicecert/domain"
	"devicecert/safelog"
	"devicecert/service"
	"devicecert/store"
)

func TestInvalidDeviceSerialDoesNotIssueCertificate(t *testing.T) {
	path := "invalid-serial.db"
	db, err := store.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	defer func() { _ = store.RemoveDatabase(path) }()
	invalid := domain.Device{Serial: "BAD-SERIAL", Model: "ALPHA-1", Status: domain.StatusRegistered, CreatedAt: "stamp", UpdatedAt: "stamp"}
	if err := db.SaveDevice(invalid); err != nil {
		t.Fatal(err)
	}
	manager := service.NewManager(db, safelog.New(nil))
	result, issueErr := manager.IssueCertificate(invalid.Serial)
	if issueErr == nil {
		t.Fatal("invalid serial unexpectedly succeeded")
	}
	if result.Certificate.CertificatePEM != "" || result.Certificate.Status != domain.CertificateRejected {
		t.Fatal("invalid serial produced a certificate")
	}
}
