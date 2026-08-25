package domain

import "testing"

func TestValidateSerial(t *testing.T) {
	valid := []string{"DC-12345678", "DC-87654321"}
	for _, value := range valid {
		if err := ValidateSerial(value); err != nil {
			t.Fatalf("valid serial rejected: %v", err)
		}
	}
	invalid := []string{"", "DC-123", "dc-12345678", "DC-00000000", " DC-12345678"}
	for _, value := range invalid {
		if err := ValidateSerial(value); err == nil {
			t.Fatalf("invalid serial accepted: %q", value)
		}
	}
}

func TestDeviceTransitions(t *testing.T) {
	device, err := NewDevice("DC-12345678", "ALPHA-1", "stamp")
	if err != nil {
		t.Fatal(err)
	}
	for _, status := range []DeviceStatus{StatusPending, StatusIssued} {
		if err := device.Transition(status); err != nil {
			t.Fatal(err)
		}
	}
	if err := device.Transition(StatusRegistered); err == nil {
		t.Fatal("terminal device transitioned backwards")
	}
}

func TestModelValidation(t *testing.T) {
	if err := ValidateModel("ALPHA-1"); err != nil {
		t.Fatal(err)
	}
	if err := ValidateModel("bad model"); err == nil {
		t.Fatal("invalid model accepted")
	}
}
