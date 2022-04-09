package crypto

import (
	"errors"
	pb "github.com/pilinsin/util/crypto/pb"
	proto "google.golang.org/protobuf/proto"
)

type pubEncryptMode int

const (
	Sidh pubEncryptMode = iota
	Ntru
	Sntrup
	Ecies
)

var SelectedPubEncMode = Sidh

func NewPubEncryptKeyPair() IPubEncryptKeyPair {
	switch SelectedPubEncMode {
	case Sidh:
		return newOqsEncryptKeyPair("SIDH-p751-compressed")
	case Ntru:
		return newOqsEncryptKeyPair("NTRU-HPS-4096-821")
	case Sntrup:
		return newSntrupKeyPair()
	case Ecies:
		return newEciesKeyPair()
	default:
		return nil
	}
}


func MarshalPriKey(pk IPriKey) ([]byte, error){
	rpk, err := pk.Raw()
	if err != nil{return nil, err}

	mKey := &pb.Key{
		Data: rpk,
	}
	return proto.Marshal(mKey)
}

func UnmarshalPriKey(m []byte) (IPriKey, error) {
	mKey := &pb.Key{}
	if err := proto.Unmarshal(m, mKey); err != nil{return nil, err}
	m = mKey.GetData()

	switch SelectedPubEncMode {
	case Sidh, Ntru:
		pri := &oqsPriKey{}
		err := pri.Unmarshal(m)
		return pri, err
	case Sntrup:
		pri := &sntrupPriKey{}
		err := pri.Unmarshal(m)
		return pri, err
	case Ecies:
		pri := &eciesPriKey{}
		err := pri.Unmarshal(m)
		return pri, err
	default:
		return nil, errors.New("invalid PubEncryptMode is selected")
	}
}


func MarshalPubKey(pk IPubKey) ([]byte, error){
	rpk, err := pk.Raw()
	if err != nil{return nil, err}

	mKey := &pb.Key{
		Data: rpk,
	}
	return proto.Marshal(mKey)
}

func UnmarshalPubKey(m []byte) (IPubKey, error) {
	mKey := &pb.Key{}
	if err := proto.Unmarshal(m, mKey); err != nil{return nil, err}
	m = mKey.GetData()
	
	switch SelectedPubEncMode {
	case Sidh, Ntru:
		pub := &oqsPubKey{}
		err := pub.Unmarshal(m)
		return pub, err
	case Sntrup:
		pub := &sntrupPubKey{}
		err := pub.Unmarshal(m)
		return pub, err
	case Ecies:
		pub := &eciesPubKey{}
		err := pub.Unmarshal(m)
		return pub, err
	default:
		return nil, errors.New("invalid PubEncryptMode is selected")
	}
}
