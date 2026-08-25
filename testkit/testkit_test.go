package testkit

import "testing"

func TestFixtureDevice(t *testing.T) {
	device := FixtureDevice()
	if device.Serial != "DC-12345678" || device.Model != "ALPHA-1" {
		t.Fatalf("unexpected fixture: %#v", device)
	}
	if len(ValidSerials()) != 3 {
		t.Fatal("fixture serial count changed")
	}
}

func TestFixtureInvalidDevice(t *testing.T) {
	if FixtureInvalidDevice().Serial == FixtureDevice().Serial {
		t.Fatal("invalid fixture is valid")
	}
}
