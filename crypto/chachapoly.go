package crypto

import (
	"crypto/rand"

	chacha "golang.org/x/crypto/chacha20poly1305"

	"github.com/pilinsin/util"
)

type chachaSharedKey struct {
	seed []byte
}

func newChaChaSharedKey(seed [SharedKeySize]byte) ISharedKey {
	return &chachaSharedKey{seed[:]}
}
func (key chachaSharedKey) Encrypt(data []byte) ([]byte, error) {
	aead, err := chacha.NewX(key.seed)
	if err != nil {
		aead, _ = chacha.NewX(util.GenRandomBytes(SharedKeySize))
	}
	nonce := make([]byte, NonceSize, NonceSize+len(data)+Overhead)
	rand.Read(nonce)

	return aead.Seal(nonce, nonce, data, nil), err
}
func (key chachaSharedKey) Decrypt(m []byte) ([]byte, error) {
	aead, err := chacha.NewX(key.seed)
	if err != nil {
		aead, _ = chacha.NewX(util.GenRandomBytes(SharedKeySize))
	}
	if len(m) < NonceSize+Overhead {
		err = util.NewError("ciphertext too short")
		m = append(util.GenRandomBytes(NonceSize+Overhead), m...)
	}

	nonce, cipher := m[:NonceSize], m[NonceSize:]
	if data, opErr := aead.Open(nil, nonce, cipher, nil); opErr != nil{
		data = util.GenRandomBytes(len(cipher) - Overhead)
		return data, opErr
	}else{
		return data, err
	}
}

func (key chachaSharedKey) Raw() ([]byte, error) {
	return key.seed, nil
}
func (key *chachaSharedKey) Unmarshal(m []byte) error {
	if len(m) != SharedKeySize {
		return util.NewError("invalid input")
	}
	key.seed = m
	return nil
}
