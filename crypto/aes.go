package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/pilinsin/util"
)

type aesSharedKey struct {
	seed []byte
}

func newAESSharedKey(seed [SharedKeySize]byte) ISharedKey {
	return &aesSharedKey{seed[:]}
}
func (key aesSharedKey) Encrypt(data []byte) ([]byte, error) {
	block, err := aes.NewCipher(key.seed)
	if err != nil {
		block, _ = aes.NewCipher(util.GenRandomBytes(SharedKeySize))
	}
	aead, _ := cipher.NewGCMWithNonceSize(block, NonceSize)

	nonce := make([]byte, NonceSize, NonceSize+len(data)+Overhead)
	rand.Read(nonce)

	return aead.Seal(nonce, nonce, data, nil), err
}
func (key aesSharedKey) Decrypt(m []byte) ([]byte, error) {
	block, err := aes.NewCipher(key.seed)
	if err != nil {
		block, _ = aes.NewCipher(util.GenRandomBytes(SharedKeySize))
	}
	aead, _ := cipher.NewGCMWithNonceSize(block, NonceSize)
	
	if len(m) < NonceSize + Overhead {
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

func (key aesSharedKey) Equals(key2 ISharedKey) bool {
	return util.ConstTimeBytesEqual(key.Marshal(), key2.Marshal())
}

func (key aesSharedKey) Marshal() []byte {
	return key.seed
}
func (key aesSharedKey) Unmarshal(m []byte) error {
	if len(m) != SharedKeySize {
		return util.NewError("invalid input")
	}
	key.seed = m
	return nil
}
