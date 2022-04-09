package crypto

import (
	"errors"
	"crypto/rand"

	sntrup "github.com/companyzero/sntrup4591761"
)

const (
	pubKeySize = sntrup.PublicKeySize
	priKeySize = sntrup.PrivateKeySize
	cipherSize = sntrup.CiphertextSize
	sharedSize = sntrup.SharedKeySize
)

type sntrupKeyPair struct {
	pubKey *sntrupPubKey
	priKey *sntrupPriKey
}

func newSntrupKeyPair() IPubEncryptKeyPair {
	pub, pri, _ := sntrup.GenerateKey(rand.Reader)
	pubKey := &sntrupPubKey{pub}
	priKey := &sntrupPriKey{pri}
	return &sntrupKeyPair{pubKey, priKey}
}
func (kp *sntrupKeyPair) Private() IPriKey {
	return kp.priKey
}
func (kp *sntrupKeyPair) Public() IPubKey {
	return kp.pubKey
}

type sntrupPriKey struct {
	priKey *[priKeySize]byte
}

func (pri *sntrupPriKey) Decrypt(m []byte) ([]byte, error) {
	if len(m) <= cipherSize {
		return nil, errors.New("decrypt fail: len(m) <= cipherSize")
	}
	cipher := new([cipherSize]byte)
	copy(cipher[:], m[:cipherSize])

	share, flag := sntrup.Decapsulate(cipher, pri.priKey)
	if flag <= 0 {
		return nil, errors.New("decrypt fail: decapsulate error")
	}

	return newChaChaSharedKey(*share).Decrypt(m[cipherSize:])
}

func (pri *sntrupPriKey) Raw() ([]byte, error) {
	return (*pri.priKey)[:], nil
}
func (pri *sntrupPriKey) Unmarshal(m []byte) error {
	if len(m) != priKeySize{
		return errors.New("invalid input")
	}
	
	priKey := [priKeySize]byte{}
	copy(priKey[:], m)
	pri.priKey = &priKey
	return nil
}

type sntrupPubKey struct {
	pubKey *[pubKeySize]byte
}

func (pub *sntrupPubKey) Encrypt(data []byte) ([]byte, error) {
	cipher, share, err := sntrup.Encapsulate(rand.Reader, pub.pubKey)
	if err != nil {
		return nil, err
	}
	enc, err := newChaChaSharedKey(*share).Encrypt(data)
	if err != nil{return nil, err}
	return append(cipher[:], enc...), nil
}

func (pub *sntrupPubKey) Raw() ([]byte, error) {
	return (*pub.pubKey)[:], nil
}
func (pub *sntrupPubKey) Unmarshal(m []byte) error {
	if len(m) != pubKeySize{
		return errors.New("invalid input")
	}
	pubKey := [pubKeySize]byte{}
	copy(pubKey[:], m)
	pub.pubKey = &pubKey
	return nil
}
