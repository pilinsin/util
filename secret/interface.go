package secret

import(
	"crypto/rand"
)

const SecretKeySize int = 32
const NonceSize int = 24
const Overhead int = 16

func RandBytes(bSize int) []byte {
	bs := make([]byte, bSize)
	rand.Read(bs)
	return bs
}

type ISecretKey interface {
	Raw() ([]byte, error)
	Encrypt(data []byte) ([]byte, error)
	Decrypt(m []byte) ([]byte, error)
}
