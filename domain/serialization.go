package domain

import (
	"encoding/json"
	"fmt"
	"strings"
)

func EncodeDevice(device Device) ([]byte, error) { return json.Marshal(device) }

func DecodeDevice(data []byte) (Device, error) {
	var device Device
	if err := json.Unmarshal(data, &device); err != nil {
		return Device{}, err
	}
	if err := ValidateDeviceInput(device.Serial, device.Model); err != nil {
		return Device{}, err
	}
	return device, nil
}

func EncodeKeyMaterial(material KeyMaterial) ([]byte, error) {
	public := material
	public.PrivateKey = ""
	return json.Marshal(public)
}

func DecodeKeyMaterial(data []byte) (KeyMaterial, error) {
	var material KeyMaterial
	if err := json.Unmarshal(data, &material); err != nil {
		return KeyMaterial{}, err
	}
	if material.DeviceSerial == "" || material.PublicKey == "" {
		return KeyMaterial{}, fmt.Errorf("incomplete key material")
	}
	return material, nil
}

func EncodeRequest(request CertificateRequest) ([]byte, error) { return json.Marshal(request) }

func DecodeRequest(data []byte) (CertificateRequest, error) {
	var request CertificateRequest
	if err := json.Unmarshal(data, &request); err != nil {
		return CertificateRequest{}, err
	}
	if request.RequestID == "" || ValidateSerial(request.DeviceSerial) != nil {
		return CertificateRequest{}, fmt.Errorf("invalid request")
	}
	return request, nil
}

func EncodeCertificate(record CertificateRecord) ([]byte, error) { return json.Marshal(record) }

func DecodeCertificate(data []byte) (CertificateRecord, error) {
	var record CertificateRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return CertificateRecord{}, err
	}
	if record.CertificateID == "" || record.DeviceSerial == "" {
		return CertificateRecord{}, fmt.Errorf("invalid certificate record")
	}
	return record, nil
}

func EncodeAudit(event AuditEvent) ([]byte, error) { return json.Marshal(event) }

func DecodeAudit(data []byte) (AuditEvent, error) {
	var event AuditEvent
	if err := json.Unmarshal(data, &event); err != nil {
		return AuditEvent{}, err
	}
	if strings.TrimSpace(event.EventID) == "" || strings.TrimSpace(event.Action) == "" {
		return AuditEvent{}, fmt.Errorf("invalid audit event")
	}
	return event, nil
}

func EntityName(value any) string {
	switch value.(type) {
	case Device:
		return "Device"
	case KeyMaterial:
		return "KeyMaterial"
	case CertificateRequest:
		return "CertificateRequest"
	case CertificateRecord:
		return "CertificateRecord"
	case AuditEvent:
		return "AuditEvent"
	default:
		return "Unknown"
	}
}
