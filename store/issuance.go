package store

import (
	"devicecert/domain"
	bolt "go.etcd.io/bbolt"
)

var keyBucket = []byte("key_material")
var requestBucket = []byte("requests")
var certificateBucket = []byte("certificates")

func (s *DB) SaveKeyMaterial(material domain.KeyMaterial) error {
	publicOnly := material
	publicOnly.PrivateKey = ""
	return s.Update(func(tx *bolt.Tx) error {
		return putJSON(tx, keyBucket, material.DeviceSerial, publicOnly)
	})
}

func (s *DB) GetKeyMaterial(serial string) (domain.KeyMaterial, error) {
	var value domain.KeyMaterial
	err := s.View(func(tx *bolt.Tx) error { return getJSON(tx, keyBucket, serial, &value) })
	return value, err
}

func (s *DB) SaveRequest(request domain.CertificateRequest) error {
	return s.Update(func(tx *bolt.Tx) error {
		return putJSON(tx, requestBucket, request.RequestID, request)
	})
}

func (s *DB) GetRequest(id string) (domain.CertificateRequest, error) {
	var value domain.CertificateRequest
	err := s.View(func(tx *bolt.Tx) error { return getJSON(tx, requestBucket, id, &value) })
	return value, err
}

func (s *DB) SaveCertificate(record domain.CertificateRecord) error {
	return s.Update(func(tx *bolt.Tx) error {
		return putJSON(tx, certificateBucket, record.CertificateID, record)
	})
}

func (s *DB) GetCertificate(id string) (domain.CertificateRecord, error) {
	var value domain.CertificateRecord
	err := s.View(func(tx *bolt.Tx) error { return getJSON(tx, certificateBucket, id, &value) })
	return value, err
}

func (s *DB) FindCertificateByDevice(serial string) (domain.CertificateRecord, error) {
	var found domain.CertificateRecord
	err := s.View(func(tx *bolt.Tx) error {
		return tx.Bucket(certificateBucket).ForEach(func(_, value []byte) error {
			var candidate domain.CertificateRecord
			if err := unmarshal(value, &candidate); err != nil {
				return err
			}
			if candidate.DeviceSerial == serial {
				found = candidate
			}
			return nil
		})
	})
	if found.CertificateID == "" && err == nil {
		err = domain.ErrMissingEntity
	}
	return found, err
}
