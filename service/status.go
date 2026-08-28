package service

import (
	"fmt"
	"sort"
	"strings"

	"devicecert/domain"
)

type StatusView struct {
	Serial            string
	Model             string
	DeviceStatus      domain.DeviceStatus
	CertificateStatus domain.CertificateStatus
	RequestStatus     domain.RequestStatus
	AuditCount        int
	Warnings          []string
}

func (m *Manager) StatusView(serial string) (StatusView, error) {
	recovery, err := m.RecoverDevice(serial)
	if err != nil {
		return StatusView{}, err
	}
	warnings := m.RecoveryWarnings(serial)
	return StatusView{Serial: recovery.Device.Serial, Model: recovery.Device.Model, DeviceStatus: recovery.Device.Status, CertificateStatus: recovery.Certificate.Status, RequestStatus: recovery.Request.Status, AuditCount: len(recovery.Audits), Warnings: warnings}, nil
}

func RenderStatusView(view StatusView) string {
	parts := []string{fmt.Sprintf("serial=%s", view.Serial), fmt.Sprintf("model=%s", view.Model), fmt.Sprintf("device=%s", view.DeviceStatus), fmt.Sprintf("certificate=%s", view.CertificateStatus), fmt.Sprintf("request=%s", view.RequestStatus), fmt.Sprintf("audits=%d", view.AuditCount)}
	if len(view.Warnings) > 0 {
		parts = append(parts, "warnings="+strings.Join(view.Warnings, ";"))
	}
	return strings.Join(parts, " ")
}

func (m *Manager) ListStatusViews() ([]StatusView, error) {
	devices, err := m.store.ListDevices()
	if err != nil {
		return nil, err
	}
	views := make([]StatusView, 0, len(devices))
	for _, device := range devices {
		view, err := m.StatusView(device.Serial)
		if err != nil {
			return nil, err
		}
		views = append(views, view)
	}
	sort.Slice(views, func(i, j int) bool { return views[i].Serial < views[j].Serial })
	return views, nil
}

func (m *Manager) RenderStatusReport() (string, error) {
	views, err := m.ListStatusViews()
	if err != nil {
		return "", err
	}
	lines := make([]string, 0, len(views))
	for _, view := range views {
		lines = append(lines, RenderStatusView(view))
	}
	return strings.Join(lines, "\n"), nil
}

func (view StatusView) Healthy() bool {
	if len(view.Warnings) > 0 {
		return false
	}
	if view.DeviceStatus == domain.StatusIssued {
		return view.CertificateStatus == domain.CertificateIssued && view.RequestStatus == domain.RequestSigned
	}
	return view.DeviceStatus == domain.StatusRegistered || view.DeviceStatus == domain.StatusPending
}

func FilterHealthy(views []StatusView) []StatusView {
	result := make([]StatusView, 0)
	for _, view := range views {
		if view.Healthy() {
			result = append(result, view)
		}
	}
	return result
}
