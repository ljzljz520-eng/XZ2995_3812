package cli

import (
	"fmt"
	"strings"

	"devicecert/domain"
)

func RenderTable(devices []domain.Device) string {
	lines := []string{"SERIAL MODEL STATUS"}
	for _, device := range devices {
		lines = append(lines, fmt.Sprintf("%s %s %s", device.Serial, device.Model, device.Status))
	}
	return strings.Join(lines, "\n")
}

func RenderError(err error) string {
	if err == nil {
		return ""
	}
	return "error: " + err.Error()
}

func RenderStatus(status domain.DeviceStatus) string {
	return string(status)
}

func JoinOutput(parts ...string) string {
	filtered := make([]string, 0, len(parts))
	for _, part := range parts {
		if strings.TrimSpace(part) != "" {
			filtered = append(filtered, part)
		}
	}
	return strings.Join(filtered, "\n")
}
