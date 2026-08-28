package cryptoengine

import (
	"crypto/elliptic"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"math/big"

	"devicecert/domain"
)

type Engine struct {
	curve elliptic.Curve
	seed  string
}

func New() *Engine {
	return &Engine{curve: elliptic.P256(), seed: "devicecert-deterministic-v1"}
}

func (e *Engine) scalar(serial string) *big.Int {
	h := sha256.Sum256([]byte(e.seed + ":" + serial))
	n := new(big.Int).SetBytes(h[:])
	order := e.curve.Params().N
	n.Mod(n, new(big.Int).Sub(order, big.NewInt(1)))
	n.Add(n, big.NewInt(1))
	return n
}

func encodePoint(x, y *big.Int) string {
	bytes := elliptic.MarshalCompressed(elliptic.P256(), x, y)
	return hex.EncodeToString(bytes)
}

func (e *Engine) GenerateKeyMaterial(serial, stamp string) (domain.KeyMaterial, error) {
	if err := domain.ValidateSerial(serial); err != nil {
		return domain.KeyMaterial{}, err
	}
	d := e.scalar(serial)
	x, y := e.curve.ScalarBaseMult(d.Bytes())
	private := hex.EncodeToString(d.FillBytes(make([]byte, 32)))
	public := encodePoint(x, y)
	hash := sha256.Sum256([]byte(public))
	return domain.KeyMaterial{DeviceSerial: serial, Algorithm: "ECDSA-P256", PublicKey: public, Fingerprint: hex.EncodeToString(hash[:8]), PrivateKey: private, CreatedAt: stamp}, nil
}

func (e *Engine) PublicKeyBytes(material domain.KeyMaterial) []byte {
	value, _ := hex.DecodeString(material.PublicKey)
	return value
}

func (e *Engine) Fingerprint(material domain.KeyMaterial) string {
	if material.Fingerprint != "" {
		return material.Fingerprint
	}
	h := sha256.Sum256([]byte(material.PublicKey))
	return hex.EncodeToString(h[:8])
}

func (e *Engine) MaterialJSON(material domain.KeyMaterial) ([]byte, error) {
	publicOnly := struct {
		DeviceSerial string `json:"device_serial"`
		Algorithm    string `json:"algorithm"`
		PublicKey    string `json:"public_key"`
		Fingerprint  string `json:"fingerprint"`
		CreatedAt    string `json:"created_at"`
	}{material.DeviceSerial, material.Algorithm, material.PublicKey, material.Fingerprint, material.CreatedAt}
	return json.Marshal(publicOnly)
}

func (e *Engine) VerifyMaterial(material domain.KeyMaterial) bool {
	if material.Algorithm != "ECDSA-P256" || material.PublicKey == "" {
		return false
	}
	generated, err := e.GenerateKeyMaterial(material.DeviceSerial, material.CreatedAt)
	if err != nil {
		return false
	}
	return generated.PublicKey == material.PublicKey && generated.Fingerprint == material.Fingerprint
}

func (e *Engine) DeriveSubject(serial, model string) string {
	h := sha256.Sum256([]byte(serial + ":" + model))
	return "device/" + hex.EncodeToString(h[:6])
}
