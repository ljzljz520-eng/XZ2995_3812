package cryptoengine

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"

	"devicecert/domain"
)

type Diagnostic struct {
	Name   string
	Passed bool
	Detail string
}

func (e *Engine) DiagnoseMaterial(material domain.KeyMaterial) []Diagnostic {
	checks := []Diagnostic{
		{Name: "algorithm", Passed: material.Algorithm == "ECDSA-P256", Detail: material.Algorithm},
		{Name: "public-key", Passed: material.PublicKey != "", Detail: material.PublicKey},
		{Name: "fingerprint", Passed: material.Fingerprint != "", Detail: material.Fingerprint},
		{Name: "determinism", Passed: e.VerifyMaterial(material), Detail: "derived"},
	}
	return checks
}

func DiagnosticsPassed(checks []Diagnostic) bool {
	for _, check := range checks {
		if !check.Passed {
			return false
		}
	}
	return true
}

func (e *Engine) RequestDiagnostics(request domain.CertificateRequest) []Diagnostic {
	checks := []Diagnostic{
		{Name: "request-id", Passed: request.RequestID != "", Detail: request.RequestID},
		{Name: "serial", Passed: domain.ValidateSerial(request.DeviceSerial) == nil, Detail: request.DeviceSerial},
		{Name: "subject", Passed: request.Subject != "", Detail: request.Subject},
		{Name: "digest", Passed: len(e.RequestDigest(request)) == 64, Detail: e.RequestDigest(request)},
	}
	return checks
}

func (e *Engine) ExplainFailure(material domain.KeyMaterial, request domain.CertificateRequest) string {
	issues := make([]string, 0)
	for _, check := range e.DiagnoseMaterial(material) {
		if !check.Passed {
			issues = append(issues, check.Name)
		}
	}
	for _, check := range e.RequestDiagnostics(request) {
		if !check.Passed {
			issues = append(issues, check.Name)
		}
	}
	if len(issues) == 0 {
		return "none"
	}
	return strings.Join(issues, ",")
}

func StableNonce(serial string) string {
	h := sha256.Sum256([]byte("nonce:" + serial))
	return hex.EncodeToString(h[:4])
}

func (e *Engine) SignPreview(material domain.KeyMaterial, request domain.CertificateRequest) string {
	if err := e.ValidateRequest(request); err != nil {
		return fmt.Sprintf("error:%s", err)
	}
	return "preview:" + e.RequestDigest(request)[:16]
}
