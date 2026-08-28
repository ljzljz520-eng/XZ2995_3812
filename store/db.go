package store

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"

	bolt "go.etcd.io/bbolt"
)

var bucketNames = [][]byte{
	[]byte("devices"),
	[]byte("key_material"),
	[]byte("requests"),
	[]byte("certificates"),
	[]byte("audits"),
}

type DB struct {
	db   *bolt.DB
	path string
}

func Open(path string) (*DB, error) {
	if path == "" {
		return nil, errors.New("store path required")
	}
	if dir := filepath.Dir(path); dir != "." {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return nil, err
		}
	}
	client, err := bolt.Open(path, 0o600, nil)
	if err != nil {
		return nil, err
	}
	result := &DB{db: client, path: path}
	if err := result.initialize(); err != nil {
		client.Close()
		return nil, err
	}
	return result, nil
}

func (s *DB) initialize() error {
	return s.db.Update(func(tx *bolt.Tx) error {
		for _, name := range bucketNames {
			if _, err := tx.CreateBucketIfNotExists(name); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *DB) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}

func (s *DB) Path() string { return s.path }

func (s *DB) View(fn func(*bolt.Tx) error) error {
	if s == nil || s.db == nil {
		return errors.New("store closed")
	}
	return s.db.View(fn)
}

func (s *DB) Update(fn func(*bolt.Tx) error) error {
	if s == nil || s.db == nil {
		return errors.New("store closed")
	}
	return s.db.Update(fn)
}

func putJSON(tx *bolt.Tx, bucket []byte, key string, value any) error {
	encoded, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return tx.Bucket(bucket).Put([]byte(key), encoded)
}

func getJSON(tx *bolt.Tx, bucket []byte, key string, target any) error {
	value := tx.Bucket(bucket).Get([]byte(key))
	if value == nil {
		return errors.New("entity not found")
	}
	return json.Unmarshal(value, target)
}

func deleteKey(tx *bolt.Tx, bucket []byte, key string) error {
	return tx.Bucket(bucket).Delete([]byte(key))
}

func RemoveDatabase(path string) error {
	return os.Remove(path)
}
