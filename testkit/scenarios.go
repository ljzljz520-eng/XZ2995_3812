package testkit

import (
	"devicecert/domain"
	"devicecert/service"
)

func RegistrationInputs() []service.RegistrationInput {
	return []service.RegistrationInput{{Serial: "DC-12345678", Model: "ALPHA-1"}, {Serial: "DC-87654321", Model: "BETA-2"}, {Serial: "BAD-SERIAL", Model: "ALPHA-1"}}
}

func ExpectedIssued(serial string) domain.CertificateRecord {
	return domain.CertificateRecord{DeviceSerial: serial, Status: domain.CertificateIssued}
}

func WorkflowNames() []string {
	return []string{"register", "issue", "recover"}
}

func AssertIssued(record domain.CertificateRecord) bool {
	return record.Status == domain.CertificateIssued && record.CertificatePEM != "" && record.ErrorMessage == ""
}

func AssertRejected(record domain.CertificateRecord) bool {
	return record.Status == domain.CertificateRejected && record.ErrorMessage != ""
}

func CountOutcomes(result service.BatchResult) map[string]int {
	return map[string]int{"registered": len(result.Registered), "rejected": len(result.Rejected), "issued": len(result.Issued)}
}

func ValidRegistrationInputs() []service.RegistrationInput {
	return []service.RegistrationInput{{Serial: "DC-12345678", Model: "ALPHA-1"}, {Serial: "DC-87654321", Model: "BETA-2"}}
}

func InvalidRegistrationInputs() []service.RegistrationInput {
	return []service.RegistrationInput{{Serial: "BAD", Model: "ALPHA-1"}, {Serial: "DC-00000000", Model: "BETA-2"}}
}

func MergeSerials(groups ...[]string) []string {
	result := make([]string, 0)
	for _, group := range groups {
		result = append(result, group...)
	}
	return result
}
