package ecies

import (
	eciesgo "github.com/ecies/go/v2"
	ipub "github.com/pilinsin/util/public"
)

type eciesKeyPair struct {
	pubKey *eciesPubKey
	priKey *eciesPriKey
}

func NewKeyPair() ipub.IPubEncryptKeyPair {
	pri, _ := eciesgo.GenerateKey()
	kp := &eciesKeyPair{&eciesPubKey{pri.PublicKey}, &eciesPriKey{pri}}
	return kp
}
func (kp *eciesKeyPair) Public() ipub.IPubKey {
	return kp.pubKey
}
func (kp *eciesKeyPair) Private() ipub.IPriKey {
	return kp.priKey
}

type eciesPriKey struct {
	priKey *eciesgo.PrivateKey
}

func (key *eciesPriKey) Decrypt(enc []byte) ([]byte, error) {
	return eciesgo.Decrypt(key.priKey, enc)
}
func (key *eciesPriKey) Raw() ([]byte, error) {
	return key.priKey.Bytes(), nil
}
func (key *eciesPriKey) Unmarshal(b []byte) error {
	pri := eciesgo.NewPrivateKeyFromBytes(b)
	key.priKey = pri
	return nil
}
func UnmarshalPriKey(m []byte) (ipub.IPriKey, error) {
	pri := &eciesPriKey{}
	err := pri.Unmarshal(m)
	return pri, err
}

type eciesPubKey struct {
	pubKey *eciesgo.PublicKey
}

func (key *eciesPubKey) Encrypt(message []byte) ([]byte, error) {
	return eciesgo.Encrypt(key.pubKey, message)
}
func (key *eciesPubKey) Raw() ([]byte, error) {
	return key.pubKey.Bytes(true), nil
}
func (key *eciesPubKey) Unmarshal(b []byte) error {
	pub, err := eciesgo.NewPublicKeyFromBytes(b)
	if err == nil {
		key.pubKey = pub
		return nil
	} else {
		return err
	}
}
func UnmarshalPubKey(m []byte) (ipub.IPubKey, error) {
	pub := &eciesPubKey{}
	err := pub.Unmarshal(m)
	return pub, err
}
