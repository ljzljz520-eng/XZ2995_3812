package cryptoengine

import (
	"testing"

	"devicecert/domain"
)

func TestDeterministicKeyGeneration(t *testing.T) {
	engine := New()
	a, err := engine.GenerateKeyMaterial("DC-12345678", "stamp")
	if err != nil {
		t.Fatal(err)
	}
	b, err := engine.GenerateKeyMaterial("DC-12345678", "other")
	if err != nil {
		t.Fatal(err)
	}
	if a.PublicKey != b.PublicKey || a.Fingerprint != b.Fingerprint || a.PrivateKey != b.PrivateKey {
		t.Fatal("key generation is not deterministic")
	}
	if !engine.VerifyMaterial(a) {
		t.Fatal("material failed verification")
	}
}

func TestRequestSigning(t *testing.T) {
	engine := New()
	device := domain.Device{Serial: "DC-12345678", Model: "ALPHA-1", Status: domain.StatusRegistered}
	material, request, err := engine.PrepareRequest(device, "stamp")
	if err != nil {
		t.Fatal(err)
	}
	certificate, id, err := engine.SignAndIdentify(material, request)
	if err != nil || id == "" || !engine.IsCertificateWellFormed(certificate) {
		t.Fatalf("signing failed: %v", err)
	}
	if _, err := engine.SignRequest(material, domain.CertificateRequest{DeviceSerial: "BAD"}); err == nil {
		t.Fatal("invalid request signed")
	}
}

func TestCertificateParsing(t *testing.T) {
	engine := New()
	if _, err := engine.ParseCertificate("broken"); err == nil {
		t.Fatal("broken certificate parsed")
	}
}
