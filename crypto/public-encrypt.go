package crypto

import (
	"github.com/pilinsin/util"
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

func UnmarshalPriKey(m []byte) (IPriKey, error) {
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
		return nil, util.NewError("invalid PubEncryptMode is selected")
	}
}

func UnmarshalPubKey(m []byte) (IPubKey, error) {
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
		return nil, util.NewError("invalid PubEncryptMode is selected")
	}
}
