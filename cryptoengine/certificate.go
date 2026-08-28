package cryptoengine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"devicecert/domain"
)

type CertificateInfo struct {
	ID           string
	Digest       string
	DeviceSerial string
	RequestID    string
}

func (e *Engine) InspectCertificate(record domain.CertificateRecord) (CertificateInfo, error) {
	if record.Status != domain.CertificateIssued {
		return CertificateInfo{}, fmt.Errorf("certificate is not issued")
	}
	digest, err := e.ParseCertificate(record.CertificatePEM)
	if err != nil {
		return CertificateInfo{}, err
	}
	return CertificateInfo{ID: record.CertificateID, Digest: digest, DeviceSerial: record.DeviceSerial, RequestID: record.RequestID}, nil
}

func (e *Engine) CertificateDigest(certificate string) string {
	h := sha256.Sum256([]byte(certificate))
	return hex.EncodeToString(h[:])
}

func (e *Engine) CertificateHeaders(certificate string) []string {
	lines := strings.Split(certificate, "\n")
	if len(lines) < 2 {
		return nil
	}
	return []string{lines[0], lines[len(lines)-1]}
}

func (e *Engine) MatchesRequest(record domain.CertificateRecord, request domain.CertificateRequest) bool {
	if record.DeviceSerial != request.DeviceSerial || record.RequestID != request.RequestID {
		return false
	}
	return e.IsCertificateWellFormed(record.CertificatePEM)
}

func (e *Engine) CertificateLabel(record domain.CertificateRecord) string {
	if record.Status == domain.CertificateIssued {
		return "issued:" + record.CertificateID
	}
	return "rejected:" + record.ErrorMessage
}

func (e *Engine) ValidateCertificateForDevice(record domain.CertificateRecord, device domain.Device) error {
	if record.DeviceSerial != device.Serial {
		return fmt.Errorf("certificate device mismatch")
	}
	if device.Status != domain.StatusIssued && record.Status == domain.CertificateIssued {
		return fmt.Errorf("device status mismatch")
	}
	if record.Status == domain.CertificateIssued && !e.IsCertificateWellFormed(record.CertificatePEM) {
		return fmt.Errorf("malformed certificate")
	}
	return nil
}
