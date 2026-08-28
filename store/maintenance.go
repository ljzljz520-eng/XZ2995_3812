package store

import (
	"errors"

	"devicecert/domain"
	bolt "go.etcd.io/bbolt"
)

func (s *DB) SaveAll(device domain.Device, material domain.KeyMaterial, request domain.CertificateRequest, certificate domain.CertificateRecord, audit domain.AuditEvent) error {
	publicMaterial := material
	publicMaterial.PrivateKey = ""
	return s.Update(func(tx *bolt.Tx) error {
		if err := putJSON(tx, deviceBucket, device.Serial, device); err != nil {
			return err
		}
		if err := putJSON(tx, keyBucket, publicMaterial.DeviceSerial, publicMaterial); err != nil {
			return err
		}
		if request.RequestID != "" {
			if err := putJSON(tx, requestBucket, request.RequestID, request); err != nil {
				return err
			}
		}
		if certificate.CertificateID != "" {
			if err := putJSON(tx, certificateBucket, certificate.CertificateID, certificate); err != nil {
				return err
			}
		}
		if audit.EventID != "" {
			return putJSON(tx, auditBucket, audit.Key(), audit)
		}
		return nil
	})
}

func (s *DB) DeleteAudit(serial string) error {
	events, err := s.ListAudits(serial)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		return errors.New("entity not found")
	}
	return s.Update(func(tx *bolt.Tx) error {
		for _, event := range events {
			if err := deleteKey(tx, auditBucket, event.Key()); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *DB) ClearCertificates() error {
	return s.Update(func(tx *bolt.Tx) error {
		if err := tx.DeleteBucket(certificateBucket); err != nil {
			return err
		}
		_, err := tx.CreateBucket(certificateBucket)
		return err
	})
}

func (s *DB) ValidateBuckets() error {
	return s.View(func(tx *bolt.Tx) error {
		for _, name := range bucketNames {
			if tx.Bucket(name) == nil {
				return errors.New("missing bucket: " + string(name))
			}
		}
		return nil
	})
}

func (s *DB) Transactional(fn func(*bolt.Tx) error) error {
	return s.Update(fn)
}
