package service_test

import (
	"testing"

	"devicecert/domain"
	"devicecert/testkit"
)

func TestWorkflowOne(t *testing.T) {
	db, manager, err := testkit.OpenFixture("workflow-one")
	if err != nil {
		t.Fatal(err)
	}
	defer testkit.Cleanup(db, "workflow-one")
	device, err := manager.RegisterDevice("DC-12345678", "ALPHA-1")
	if err != nil || device.Status != domain.StatusRegistered {
		t.Fatalf("registration failed: %v", err)
	}
	loaded, err := manager.LoadDevice(device.Serial)
	if err != nil || loaded.Model != "ALPHA-1" {
		t.Fatalf("load failed: %v", err)
	}
}

func TestWorkflowTwo(t *testing.T) {
	db, manager, err := testkit.OpenFixture("workflow-two")
	if err != nil {
		t.Fatal(err)
	}
	defer testkit.Cleanup(db, "workflow-two")
	if _, err := manager.RegisterDevice("DC-12345678", "ALPHA-1"); err != nil {
		t.Fatal(err)
	}
	result, err := manager.IssueCertificate("DC-12345678")
	if err != nil {
		t.Fatal(err)
	}
	if !result.Certificate.IsSuccessful() || result.Device.Status != domain.StatusIssued {
		t.Fatal("certificate workflow did not issue")
	}
	if manager.Logger().ContainsPrivateKey(result.Material.PrivateKey) {
		t.Fatal("private key reached logger")
	}
}

func TestWorkflowThree(t *testing.T) {
	db, manager, err := testkit.OpenFixture("workflow-three")
	if err != nil {
		t.Fatal(err)
	}
	defer testkit.Cleanup(db, "workflow-three")
	if _, err := manager.RegisterDevice("DC-87654321", "BETA-2"); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.IssueCertificate("DC-87654321"); err != nil {
		t.Fatal(err)
	}
	recovered, err := manager.ReopenAndRecover(testkit.FixturePath("workflow-three"), "DC-87654321")
	if err != nil || recovered.Certificate.Status != domain.CertificateIssued {
		t.Fatalf("recovery failed: %v", err)
	}
}
