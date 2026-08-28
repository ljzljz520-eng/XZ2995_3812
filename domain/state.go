package domain

func (d *Device) Transition(next DeviceStatus) error {
	if d == nil {
		return ErrMissingEntity
	}
	if d.Status == next {
		return nil
	}
	switch d.Status {
	case StatusRegistered:
		if next != StatusPending && next != StatusRejected {
			return ErrInvalidTransition
		}
	case StatusPending:
		if next != StatusIssued && next != StatusRejected {
			return ErrInvalidTransition
		}
	case StatusIssued, StatusRejected:
		return ErrInvalidTransition
	default:
		return ErrInvalidTransition
	}
	d.Status = next
	return nil
}

func NewDevice(serial, model, stamp string) (Device, error) {
	if err := ValidateDeviceInput(serial, model); err != nil {
		return Device{}, err
	}
	return Device{Serial: serial, Model: NormalizeModel(model), Status: StatusRegistered, CreatedAt: stamp, UpdatedAt: stamp}, nil
}

func NewRejectedDevice(serial, model, stamp string) Device {
	return Device{Serial: serial, Model: NormalizeModel(model), Status: StatusRejected, CreatedAt: stamp, UpdatedAt: stamp}
}

func CloneDevice(d Device) Device {
	return Device{Serial: d.Serial, Model: d.Model, Status: d.Status, CreatedAt: d.CreatedAt, UpdatedAt: d.UpdatedAt}
}

func DeviceCanIssue(d Device) bool {
	return d.Status == StatusRegistered || d.Status == StatusPending
}

func IsTerminalStatus(status DeviceStatus) bool {
	return status == StatusIssued || status == StatusRejected
}

func StatusLabel(status DeviceStatus) string {
	switch status {
	case StatusRegistered:
		return "registered"
	case StatusPending:
		return "pending"
	case StatusIssued:
		return "issued"
	case StatusRejected:
		return "rejected"
	default:
		return "unknown"
	}
}
