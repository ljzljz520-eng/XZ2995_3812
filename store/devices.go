package store

import (
	"errors"

	"devicecert/domain"
	bolt "go.etcd.io/bbolt"
)

var deviceBucket = []byte("devices")

func (s *DB) SaveDevice(device domain.Device) error {
	return s.Update(func(tx *bolt.Tx) error {
		return putJSON(tx, deviceBucket, device.Serial, device)
	})
}

func (s *DB) GetDevice(serial string) (domain.Device, error) {
	var device domain.Device
	err := s.View(func(tx *bolt.Tx) error {
		return getJSON(tx, deviceBucket, serial, &device)
	})
	return device, err
}

func (s *DB) DeleteDevice(serial string) error {
	return s.Update(func(tx *bolt.Tx) error {
		if tx.Bucket(deviceBucket).Get([]byte(serial)) == nil {
			return errors.New("entity not found")
		}
		return deleteKey(tx, deviceBucket, serial)
	})
}

func (s *DB) ListDevices() ([]domain.Device, error) {
	items := make([]domain.Device, 0)
	err := s.View(func(tx *bolt.Tx) error {
		return tx.Bucket(deviceBucket).ForEach(func(_, value []byte) error {
			var device domain.Device
			if err := unmarshal(value, &device); err != nil {
				return err
			}
			items = append(items, device)
			return nil
		})
	})
	return items, err
}
