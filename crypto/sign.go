package crypto

import (
	"errors"
	pb "github.com/pilinsin/util/crypto/pb"
	proto "google.golang.org/protobuf/proto"
)

type signMode int

const (
	Sphincs signMode = iota
	Falcon
	Bliss
	Ed25519
)

var SelectedSignMode = Sphincs

func NewSignKeyPair() ISignKeyPair {
	switch SelectedSignMode {
	case Sphincs:
		return newOqsSignKeyPair("SPHINCS+-SHAKE256-256s-robust")
	case Falcon:
		return newOqsSignKeyPair("Falcon-1024")
	case Bliss:
		return newBlissKeyPair()
	case Ed25519:
		return newEd25519KeyPair()
	default:
		return nil
	}
}

func MarshalSignKey(sk ISignKey) ([]byte, error){
	rsk, err := sk.Raw()
	if err != nil{return nil, err}

	mKey := &pb.Key{
		Data: rsk,
	}
	return proto.Marshal(mKey)
}

func UnmarshalSignKey(m []byte) (ISignKey, error) {
	mKey := &pb.Key{}
	if err := proto.Unmarshal(m, mKey); err != nil{return nil, err}
	m = mKey.GetData()

	switch SelectedSignMode {
	case Sphincs, Falcon:
		sk := &oqsSignKey{}
		err := sk.Unmarshal(m)
		return sk, err
	case Bliss:
		sk := &blissSignKey{}
		err := sk.Unmarshal(m)
		return sk, err
	case Ed25519:
		sk := &ed25519SignKey{}
		err := sk.Unmarshal(m)
		return sk, err
	default:
		return nil, errors.New("invalid SignMode is selected")
	}
}

func MarshalVerfKey(vk IVerfKey) ([]byte, error){
	rvk, err := vk.Raw()
	if err != nil{return nil, err}

	mKey := &pb.Key{
		Data: rvk,
	}
	return proto.Marshal(mKey)
}

func UnmarshalVerfKey(m []byte) (IVerfKey, error) {
	mKey := &pb.Key{}
	if err := proto.Unmarshal(m, mKey); err != nil{return nil, err}
	m = mKey.GetData()

	switch SelectedSignMode {
	case Sphincs, Falcon:
		vk := &oqsVerfKey{}
		err := vk.Unmarshal(m)
		return vk, err
	case Bliss:
		vk := &blissVerfKey{}
		err := vk.Unmarshal(m)
		return vk, err
	case Ed25519:
		vk := &ed25519VerfKey{}
		err := vk.Unmarshal(m)
		return vk, err
	default:
		return nil, errors.New("invalid SignMode is selected")
	}
}
