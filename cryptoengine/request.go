package cryptoengine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"devicecert/domain"
)

func (e *Engine) BuildRequest(device domain.Device, material domain.KeyMaterial, stamp string) (domain.CertificateRequest, error) {
	if err := domain.ValidateSerial(device.Serial); err != nil {
		return domain.CertificateRequest{}, err
	}
	if material.DeviceSerial != device.Serial || material.PublicKey == "" {
		return domain.CertificateRequest{}, fmt.Errorf("key material mismatch")
	}
	seed := device.Serial + ":" + material.Fingerprint + ":" + stamp
	h := sha256.Sum256([]byte(seed))
	requestID := "REQ-" + hex.EncodeToString(h[:6])
	return domain.CertificateRequest{RequestID: requestID, DeviceSerial: device.Serial, Subject: e.DeriveSubject(device.Serial, device.Model), PublicKey: material.PublicKey, Status: domain.RequestPrepared, CreatedAt: stamp}, nil
}

func (e *Engine) EncodeRequest(request domain.CertificateRequest) []byte {
	return []byte(request.RequestID + "|" + request.DeviceSerial + "|" + request.Subject + "|" + request.PublicKey)
}

func (e *Engine) RequestDigest(request domain.CertificateRequest) string {
	h := sha256.Sum256(e.EncodeRequest(request))
	return hex.EncodeToString(h[:])
}

func (e *Engine) ValidateRequest(request domain.CertificateRequest) error {
	if request.RequestID == "" || request.DeviceSerial == "" || request.Subject == "" || request.PublicKey == "" {
		return fmt.Errorf("incomplete certificate request")
	}
	if err := domain.ValidateSerial(request.DeviceSerial); err != nil {
		return err
	}
	if request.Status != domain.RequestPrepared && request.Status != domain.RequestSigned {
		return fmt.Errorf("unsupported request status")
	}
	return nil
}

func (e *Engine) PrepareRequest(device domain.Device, stamp string) (domain.KeyMaterial, domain.CertificateRequest, error) {
	material, err := e.GenerateKeyMaterial(device.Serial, stamp)
	if err != nil {
		return domain.KeyMaterial{}, domain.CertificateRequest{}, err
	}
	request, err := e.BuildRequest(device, material, stamp)
	if err != nil {
		return domain.KeyMaterial{}, domain.CertificateRequest{}, err
	}
	return material, request, nil
}
