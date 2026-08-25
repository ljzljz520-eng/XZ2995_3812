package store

import (
	"sort"
	"strings"

	"devicecert/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *DB) FindDevicesByModel(model string) ([]domain.Device, error) {
	items := make([]domain.Device, 0)
	err := s.View(func(tx *bolt.Tx) error {
		return tx.Bucket(deviceBucket).ForEach(func(_, value []byte) error {
			var device domain.Device
			if err := unmarshal(value, &device); err != nil {
				return err
			}
			if device.Model == model {
				items = append(items, device)
			}
			return nil
		})
	})
	sort.Slice(items, func(i, j int) bool { return items[i].Serial < items[j].Serial })
	return items, err
}

func (s *DB) SearchAudits(term string) ([]domain.AuditEvent, error) {
	items := make([]domain.AuditEvent, 0)
	lowerTerm := strings.ToLower(term)
	err := s.View(func(tx *bolt.Tx) error {
		return tx.Bucket(auditBucket).ForEach(func(_, value []byte) error {
			var event domain.AuditEvent
			if err := unmarshal(value, &event); err != nil {
				return err
			}
			text := strings.ToLower(event.Action + " " + event.Outcome + " " + event.Detail)
			if lowerTerm == "" || strings.Contains(text, lowerTerm) {
				items = append(items, event)
			}
			return nil
		})
	})
	return items, err
}

func (s *DB) ListCertificates() ([]domain.CertificateRecord, error) {
	items := make([]domain.CertificateRecord, 0)
	err := s.View(func(tx *bolt.Tx) error {
		return tx.Bucket(certificateBucket).ForEach(func(_, value []byte) error {
			var record domain.CertificateRecord
			if err := unmarshal(value, &record); err != nil {
				return err
			}
			items = append(items, record)
			return nil
		})
	})
	sort.Slice(items, func(i, j int) bool { return items[i].CertificateID < items[j].CertificateID })
	return items, err
}

func (s *DB) CertificatesForStatus(status domain.CertificateStatus) ([]domain.CertificateRecord, error) {
	items, err := s.ListCertificates()
	if err != nil {
		return nil, err
	}
	filtered := make([]domain.CertificateRecord, 0)
	for _, item := range items {
		if item.Status == status {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (s *DB) CountByDeviceStatus(status domain.DeviceStatus) (int, error) {
	devices, err := s.ListDevices()
	if err != nil {
		return 0, err
	}
	count := 0
	for _, device := range devices {
		if device.Status == status {
			count++
		}
	}
	return count, nil
}

func (s *DB) Exists(bucket []byte, key string) (bool, error) {
	found := false
	err := s.View(func(tx *bolt.Tx) error {
		found = tx.Bucket(bucket).Get([]byte(key)) != nil
		return nil
	})
	return found, err
}

func (s *DB) ExportSummary() (map[string]int, error) {
	counts, err := s.CountEntities()
	if err != nil {
		return nil, err
	}
	return map[string]int{"devices": counts["devices"], "keys": counts["key_material"], "requests": counts["requests"], "certificates": counts["certificates"], "audits": counts["audits"]}, nil
}
