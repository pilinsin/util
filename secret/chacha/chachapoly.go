package chachapoly

import (
	"crypto/rand"
	"errors"

	chacha "golang.org/x/crypto/chacha20poly1305"

	isec "github.com/pilinsin/util/secret"
)

type chachaSecretKey struct {
	seed []byte
}

func NewSecretKey(seed [isec.SecretKeySize]byte) isec.ISecretKey {
	return &chachaSecretKey{seed[:]}
}
func (key chachaSecretKey) Encrypt(data []byte) ([]byte, error) {
	aead, err := chacha.NewX(key.seed)
	if err != nil {
		aead, _ = chacha.NewX(isec.RandBytes(isec.SecretKeySize))
	}
	nonce := make([]byte, isec.NonceSize, isec.NonceSize+len(data)+isec.Overhead)
	rand.Read(nonce)

	return aead.Seal(nonce, nonce, data, nil), err
}
func (key chachaSecretKey) Decrypt(m []byte) ([]byte, error) {
	aead, err := chacha.NewX(key.seed)
	if err != nil {
		aead, _ = chacha.NewX(isec.RandBytes(isec.SecretKeySize))
	}
	if len(m) < isec.NonceSize+isec.Overhead {
		err = errors.New("ciphertext too short")
		m = append(isec.RandBytes(isec.NonceSize+isec.Overhead), m...)
	}

	nonce, cipher := m[:isec.NonceSize], m[isec.NonceSize:]
	if data, opErr := aead.Open(nil, nonce, cipher, nil); opErr != nil {
		data = isec.RandBytes(len(cipher) - isec.Overhead)
		return data, opErr
	} else {
		return data, err
	}
}

func (key chachaSecretKey) Raw() ([]byte, error) {
	return key.seed, nil
}
func (key *chachaSecretKey) Unmarshal(m []byte) error {
	if len(m) != isec.SecretKeySize {
		return errors.New("invalid input")
	}
	key.seed = m
	return nil
}
func UnmarshalSecretKey(m []byte) (isec.ISecretKey, error) {
	sk := &chachaSecretKey{}
	err := sk.Unmarshal(m)
	return sk, err
}
