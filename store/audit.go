package store

import (
	"devicecert/domain"
	bolt "go.etcd.io/bbolt"
)

var auditBucket = []byte("audits")

func (s *DB) SaveAudit(event domain.AuditEvent) error {
	return s.Update(func(tx *bolt.Tx) error {
		return putJSON(tx, auditBucket, event.Key(), event)
	})
}

func (s *DB) ListAudits(serial string) ([]domain.AuditEvent, error) {
	items := make([]domain.AuditEvent, 0)
	err := s.View(func(tx *bolt.Tx) error {
		return tx.Bucket(auditBucket).ForEach(func(_, value []byte) error {
			var event domain.AuditEvent
			if err := unmarshal(value, &event); err != nil {
				return err
			}
			if serial == "" || event.DeviceSerial == serial {
				items = append(items, event)
			}
			return nil
		})
	})
	return items, err
}

func (s *DB) CountEntities() (map[string]int, error) {
	counts := make(map[string]int)
	err := s.View(func(tx *bolt.Tx) error {
		for _, bucket := range bucketNames {
			count := 0
			err := tx.Bucket(bucket).ForEach(func(_, _ []byte) error { count++; return nil })
			if err != nil {
				return err
			}
			counts[string(bucket)] = count
		}
		return nil
	})
	return counts, err
}
