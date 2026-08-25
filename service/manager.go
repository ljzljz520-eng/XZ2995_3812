package service

import (
	"errors"
	"fmt"

	"devicecert/cryptoengine"
	"devicecert/domain"
	"devicecert/safelog"
	"devicecert/store"
)

type Manager struct {
	store    *store.DB
	engine   *cryptoengine.Engine
	logger   *safelog.Logger
	sequence int
	clock    func(int) string
}

type IssueResult struct {
	Device      domain.Device
	Material    domain.KeyMaterial
	Request     domain.CertificateRequest
	Certificate domain.CertificateRecord
	Audit       domain.AuditEvent
}

type Recovery struct {
	Device      domain.Device
	Material    domain.KeyMaterial
	Request     domain.CertificateRequest
	Certificate domain.CertificateRecord
	Audits      []domain.AuditEvent
}

func NewManager(db *store.DB, logger *safelog.Logger) *Manager {
	if logger == nil {
		logger = safelog.New(nil)
	}
	return &Manager{store: db, engine: cryptoengine.New(), logger: logger, clock: deterministicStamp}
}

func deterministicStamp(index int) string {
	return fmt.Sprintf("2024-01-01T00:%02d:00Z", index%60)
}

func (m *Manager) nextStamp() string {
	m.sequence++
	return m.clock(m.sequence)
}

func (m *Manager) ensureReady() error {
	if m == nil || m.store == nil {
		return errors.New("manager unavailable")
	}
	return nil
}

func (m *Manager) Store() *store.DB { return m.store }

func (m *Manager) Engine() *cryptoengine.Engine { return m.engine }

func (m *Manager) Logger() *safelog.Logger { return m.logger }

func (m *Manager) SetClock(clock func(int) string) {
	if clock != nil {
		m.clock = clock
	}
}

func (m *Manager) Snapshot(serial string) (Recovery, error) {
	return m.RecoverDevice(serial)
}
