package service_test

import (
	"testing"

	"devicecert/domain"
	"devicecert/testkit"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	db, manager, err := testkit.OpenFixture("persistence")
	if err != nil {
		t.Fatal(err)
	}
	path := testkit.FixturePath("persistence")
	if _, err := manager.RegisterDevice("DC-13572468", "GAMMA-3"); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.IssueCertificate("DC-13572468"); err != nil {
		t.Fatal(err)
	}
	if _, err := manager.ReopenAndRecover(path, "DC-13572468"); err != nil {
		t.Fatal(err)
	}
	recovered, err := manager.RecoverDevice("DC-13572468")
	if err != nil {
		t.Fatal(err)
	}
	if recovered.Device.Status != domain.StatusIssued || recovered.Request.RequestID == "" || len(recovered.Audits) < 2 {
		t.Fatalf("state did not survive reopen: %#v", recovered)
	}
	testkit.Cleanup(manager.Store(), "persistence")
	_ = db
}
