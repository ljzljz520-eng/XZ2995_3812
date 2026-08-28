package store

import (
	"os"
	"testing"

	"devicecert/domain"
)

func TestStoreRoundTrip(t *testing.T) {
	path := "store-roundtrip.db"
	_ = os.Remove(path)
	defer os.Remove(path)
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	device := domain.Device{Serial: "DC-12345678", Model: "ALPHA-1", Status: domain.StatusRegistered}
	if err := db.SaveDevice(device); err != nil {
		t.Fatal(err)
	}
	if err := db.SaveAudit(domain.AuditEvent{EventID: "AUD-1", DeviceSerial: device.Serial, Action: "register", Outcome: "accepted"}); err != nil {
		t.Fatal(err)
	}
	loaded, err := db.GetDevice(device.Serial)
	if err != nil || loaded.Serial != device.Serial {
		t.Fatalf("round trip failed: %v", err)
	}
	counts, err := db.CountEntities()
	if err != nil || counts["devices"] != 1 || counts["audits"] != 1 {
		t.Fatalf("counts failed: %v %#v", err, counts)
	}
	_ = db.Close()
}

func TestStoreMissingEntity(t *testing.T) {
	path := "store-missing.db"
	_ = os.Remove(path)
	defer os.Remove(path)
	db, err := Open(path)
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if _, err := db.GetDevice("DC-12345678"); err == nil {
		t.Fatal("missing entity found")
	}
}
