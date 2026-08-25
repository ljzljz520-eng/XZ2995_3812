package cryptoengine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"devicecert/domain"
)

func (e *Engine) SignRequest(material domain.KeyMaterial, request domain.CertificateRequest) (string, error) {
	if err := e.ValidateRequest(request); err != nil {
		return "", err
	}
	if material.PrivateKey == "" || material.DeviceSerial != request.DeviceSerial {
		return "", fmt.Errorf("missing signing material")
	}
	if !e.VerifyMaterial(material) {
		return "", fmt.Errorf("invalid signing material")
	}
	digest := sha256.Sum256([]byte(material.PrivateKey + ":" + e.RequestDigest(request)))
	return "-----BEGIN DEVICE CERTIFICATE-----\n" + hex.EncodeToString(digest[:]) + "\n-----END DEVICE CERTIFICATE-----", nil
}

func (e *Engine) ParseCertificate(certificate string) (string, error) {
	lines := strings.Split(certificate, "\n")
	if len(lines) != 3 || lines[0] != "-----BEGIN DEVICE CERTIFICATE-----" || lines[2] != "-----END DEVICE CERTIFICATE-----" {
		return "", fmt.Errorf("invalid certificate encoding")
	}
	if len(lines[1]) != 64 {
		return "", fmt.Errorf("invalid certificate digest")
	}
	if _, err := hex.DecodeString(lines[1]); err != nil {
		return "", err
	}
	return lines[1], nil
}

func (e *Engine) CertificateID(request domain.CertificateRequest, certificate string) string {
	h := sha256.Sum256([]byte(request.RequestID + ":" + certificate))
	return "CRT-" + hex.EncodeToString(h[:6])
}

func (e *Engine) IsCertificateWellFormed(certificate string) bool {
	_, err := e.ParseCertificate(certificate)
	return err == nil
}

func (e *Engine) SignAndIdentify(material domain.KeyMaterial, request domain.CertificateRequest) (string, string, error) {
	certificate, err := e.SignRequest(material, request)
	if err != nil {
		return "", "", err
	}
	return certificate, e.CertificateID(request, certificate), nil
}
