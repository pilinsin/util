package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"

	isec "github.com/pilinsin/util/secret"
)

type aesSecretKey struct {
	seed []byte
}

func NewSecretKey(seed [isec.SecretKeySize]byte) isec.ISecretKey {
	return &aesSecretKey{seed[:]}
}
func (key aesSecretKey) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key.seed)
	if err != nil {
		block, _ = aes.NewCipher(isec.RandBytes(isec.SecretKeySize))
	}
	aead, _ := cipher.NewGCMWithNonceSize(block, isec.NonceSize)

	nonce := make([]byte, isec.NonceSize, isec.NonceSize+len(data)+isec.Overhead)
	rand.Read(nonce)

	return aead.Seal(nonce, nonce, data, nil), err
}
func (key aesSecretKey) Decrypt(m []byte) ([]byte, error) {
	block, err := aes.NewCipher(key.seed)
	if err != nil {
		block, _ = aes.NewCipher(isec.RandBytes(isec.SecretKeySize))
	}
	aead, _ := cipher.NewGCMWithNonceSize(block, isec.NonceSize)

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

func (key aesSecretKey) Raw() ([]byte, error) {
	return key.seed, nil
}
func (key *aesSecretKey) Unmarshal(m []byte) error {
	if len(m) != isec.SecretKeySize {
		return errors.New("invalid input")
	}
	key.seed = m
	return nil
}

func UnmarshalSecretKey(m []byte) (isec.ISecretKey, error) {
	sk := &aesSecretKey{}
	err := sk.Unmarshal(m)
	return sk, err
}
