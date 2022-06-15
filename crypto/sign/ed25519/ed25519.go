package ed25519

import (
	"errors"
	ed "golang.org/x/crypto/ed25519"

	isign "github.com/pilinsin/util/crypto/sign"
)


type ed25519KeyPair struct {
	signKey *ed25519SignKey
	verfKey *ed25519VerfKey
}

func NewKeyPair() isign.ISignKeyPair {
	pub, pri, _ := ed.GenerateKey(nil)
	kp := &ed25519KeyPair{&ed25519SignKey{pri}, &ed25519VerfKey{pub}}
	return kp
}
func (kp *ed25519KeyPair) Sign() isign.ISignKey {
	return kp.signKey
}
func (kp *ed25519KeyPair) Verify() isign.IVerfKey {
	return kp.verfKey
}


type ed25519SignKey struct {
	signKey ed.PrivateKey
}
func (key *ed25519SignKey) Sign(msg []byte) ([]byte, error) {
	if len(key.signKey) != ed.PrivateKeySize {
		return nil, errors.New("invalid signKey")
	}
	return ed.Sign(key.signKey, msg), nil
}
func (key *ed25519SignKey) Raw() ([]byte, error) {
	return key.signKey, nil
}
func (key *ed25519SignKey) Unmarshal(b []byte) error {
	key.signKey = b
	return nil
}
func UnmarshalSignKey(m []byte) (isign.ISignKey, error){
	sk := &ed25519SignKey{}
	err := sk.Unmarshal(m)
	return sk, err
}


type ed25519VerfKey struct {
	verfKey ed.PublicKey
}
func (key *ed25519VerfKey) Verify(msg, sig []byte) (bool, error) {
	if len(key.verfKey) != ed.PublicKeySize {
		return false, errors.New("invalid verfKey")
	}
	return ed.Verify(key.verfKey, msg, sig), nil
}
func (key *ed25519VerfKey) Raw() ([]byte, error) {
	return key.verfKey, nil
}
func (key *ed25519VerfKey) Unmarshal(b []byte) error {
	key.verfKey = b
	return nil
}
func UnmarshalVerfKey(m []byte) (isign.IVerfKey, error){
	sk := &ed25519VerfKey{}
	err := sk.Unmarshal(m)
	return sk, err
}