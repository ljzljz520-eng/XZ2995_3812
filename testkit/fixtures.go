package testkit

import (
	"os"
	"path/filepath"

	"devicecert/domain"
	"devicecert/safelog"
	"devicecert/service"
	"devicecert/store"
)

func FixtureDevice() domain.Device {
	return domain.Device{Serial: "DC-12345678", Model: "ALPHA-1", Status: domain.StatusRegistered, CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"}
}

func FixtureInvalidDevice() domain.Device {
	return domain.Device{Serial: "BAD-SERIAL", Model: "ALPHA-1", Status: domain.StatusRegistered, CreatedAt: "2024-01-01T00:00:00Z", UpdatedAt: "2024-01-01T00:00:00Z"}
}

func FixturePath(name string) string {
	return filepath.Join(os.TempDir(), "devicecert-"+name+".db")
}

func OpenFixture(name string) (*store.DB, *service.Manager, error) {
	path := FixturePath(name)
	_ = os.Remove(path)
	db, err := store.Open(path)
	if err != nil {
		return nil, nil, err
	}
	return db, service.NewManager(db, safelog.New(nil)), nil
}

func Cleanup(db *store.DB, name string) {
	if db != nil {
		_ = db.Close()
	}
	_ = os.Remove(FixturePath(name))
}

func FixtureStamp(index int) string {
	return "2024-01-01T00:00:" + string(rune('0'+index)) + "Z"
}

func ValidSerials() []string {
	return []string{"DC-12345678", "DC-87654321", "DC-13572468"}
}
