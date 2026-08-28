package cryptoengine

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"

	"devicecert/domain"
)

func EncodePublicKey(material domain.KeyMaterial) string {
	bytes, err := hex.DecodeString(material.PublicKey)
	if err != nil {
		return ""
	}
	return base64.RawStdEncoding.EncodeToString(bytes)
}

func DecodePublicKey(value string) (string, error) {
	bytes, err := base64.RawStdEncoding.DecodeString(value)
	if err != nil {
		return "", err
	}
	if len(bytes) < 33 {
		return "", fmt.Errorf("public key too short")
	}
	return hex.EncodeToString(bytes), nil
}

func EncodeRequestEnvelope(request domain.CertificateRequest) string {
	parts := []string{request.RequestID, request.DeviceSerial, request.Subject, request.PublicKey, string(request.Status)}
	return strings.Join(parts, ".")
}

func DecodeRequestEnvelope(value string) (domain.CertificateRequest, error) {
	parts := strings.Split(value, ".")
	if len(parts) != 5 {
		return domain.CertificateRequest{}, fmt.Errorf("invalid request envelope")
	}
	request := domain.CertificateRequest{RequestID: parts[0], DeviceSerial: parts[1], Subject: parts[2], PublicKey: parts[3], Status: domain.RequestStatus(parts[4])}
	if err := New().ValidateRequest(request); err != nil {
		return domain.CertificateRequest{}, err
	}
	return request, nil
}

func CanonicalRequest(request domain.CertificateRequest) string {
	return request.RequestID + "|" + request.DeviceSerial + "|" + request.Subject + "|" + request.PublicKey
}

func CanonicalCertificate(record domain.CertificateRecord) string {
	return record.CertificateID + "|" + record.DeviceSerial + "|" + record.RequestID + "|" + string(record.Status)
}

func MaterialFingerprintBytes(material domain.KeyMaterial) []byte {
	value, _ := hex.DecodeString(material.Fingerprint)
	return value
}

func IsEncodedPublicKey(value string) bool {
	_, err := DecodePublicKey(value)
	return err == nil
}

func NormalizeEnvelope(value string) string {
	return strings.TrimSpace(strings.ReplaceAll(value, " ", ""))
}
