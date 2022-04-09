package crypto

import (
	"errors"
	pb "github.com/pilinsin/util/crypto/pb"
	proto "google.golang.org/protobuf/proto"
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


func MarshalSharedKey(sk ISharedKey) ([]byte, error){
	rsk, err := sk.Raw()
	if err != nil{return nil, err}

	mKey := &pb.Key{
		Data: rsk,
	}
	return proto.Marshal(mKey)
}

func UnmarshalSharedKey(m []byte) (ISharedKey, error) {
	mKey := &pb.Key{}
	if err := proto.Unmarshal(m, mKey); err != nil{return nil, err}
	m = mKey.GetData()
	
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
		return nil, errors.New("invalid SharedEncryptMode is selected")
	}
}
