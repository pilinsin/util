package crypto

import (
	"github.com/pilinsin/util"
)

type sharedEncryptMode int

const (
	ChaCha sharedEncryptMode = iota
	AES
)

var SelectedSharedEncMode = ChaCha

const SharedKeySize int = 32
const NonceSize int = 24
const Overhead int = 16

func NewSharedEncryptKey(seed []byte) ISharedKey {
	if len(seed) != SharedKeySize {
		return nil
	}
	var key [SharedKeySize]byte
	copy(key[:], seed)

	switch SelectedSharedEncMode {
	case ChaCha:
		return newChaChaSharedKey(key)
	case AES:
		return newAESSharedKey(key)
	default:
		return nil
	}
}

func UnmarshalSharedKey(m []byte) (ISharedKey, error) {
	switch SelectedSharedEncMode {
	case ChaCha:
		sk := &chachaSharedKey{}
		err := sk.Unmarshal(m)
		return sk, err
	case AES:
		sk := &aesSharedKey{}
		err := sk.Unmarshal(m)
		return sk, err
	default:
		return nil, util.NewError("invalid SharedEncryptMode is selected")
	}
}
