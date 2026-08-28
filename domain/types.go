package domain

import (
	"errors"
	"fmt"
)

type DeviceStatus string

const (
	StatusRegistered DeviceStatus = "registered"
	StatusPending    DeviceStatus = "pending"
	StatusIssued     DeviceStatus = "issued"
	StatusRejected   DeviceStatus = "rejected"
)

type RequestStatus string

const (
	RequestPrepared RequestStatus = "prepared"
	RequestSigned   RequestStatus = "signed"
	RequestFailed   RequestStatus = "failed"
)

type CertificateStatus string

const (
	CertificateIssued   CertificateStatus = "issued"
	CertificateRejected CertificateStatus = "rejected"
)

type Device struct {
	Serial    string       `json:"serial"`
	Model     string       `json:"model"`
	Status    DeviceStatus `json:"status"`
	CreatedAt string       `json:"created_at"`
	UpdatedAt string       `json:"updated_at"`
}

type KeyMaterial struct {
	DeviceSerial string `json:"device_serial"`
	Algorithm    string `json:"algorithm"`
	PublicKey    string `json:"public_key"`
	Fingerprint  string `json:"fingerprint"`
	PrivateKey   string `json:"private_key"`
	CreatedAt    string `json:"created_at"`
}

type CertificateRequest struct {
	RequestID    string        `json:"request_id"`
	DeviceSerial string        `json:"device_serial"`
	Subject      string        `json:"subject"`
	PublicKey    string        `json:"public_key"`
	Status       RequestStatus `json:"status"`
	CreatedAt    string        `json:"created_at"`
}

type CertificateRecord struct {
	CertificateID  string            `json:"certificate_id"`
	DeviceSerial   string            `json:"device_serial"`
	RequestID      string            `json:"request_id"`
	Status         CertificateStatus `json:"status"`
	CertificatePEM string            `json:"certificate_pem"`
	ErrorMessage   string            `json:"error_message"`
	CreatedAt      string            `json:"created_at"`
}

type AuditEvent struct {
	EventID      string `json:"event_id"`
	DeviceSerial string `json:"device_serial"`
	Action       string `json:"action"`
	Outcome      string `json:"outcome"`
	Detail       string `json:"detail"`
	CreatedAt    string `json:"created_at"`
}

var (
	ErrInvalidSerial     = errors.New("invalid device serial")
	ErrInvalidModel      = errors.New("invalid device model")
	ErrInvalidTransition = errors.New("invalid device status transition")
	ErrMissingEntity     = errors.New("entity not found")
	ErrAlreadyExists     = errors.New("entity already exists")
)

func (d Device) String() string {
	return fmt.Sprintf("%s/%s[%s]", d.Serial, d.Model, d.Status)
}

func (r CertificateRequest) Summary() string {
	return fmt.Sprintf("%s:%s:%s", r.RequestID, r.DeviceSerial, r.Status)
}

func (c CertificateRecord) IsSuccessful() bool {
	return c.Status == CertificateIssued && c.CertificatePEM != ""
}

func (a AuditEvent) Key() string {
	return a.EventID + ":" + a.DeviceSerial
}
