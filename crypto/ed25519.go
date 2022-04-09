package crypto

import (
	sed "crypto/ed25519"
	ced "golang.org/x/crypto/ed25519"

	"github.com/pilinsin/util"
)

type ed25519SignKey struct {
	signKey ced.PrivateKey
}

type ed25519VerfKey struct {
	verfKey ced.PublicKey
}

type ed25519KeyPair struct {
	signKey *ed25519SignKey
	verfKey *ed25519VerfKey
}

func newEd25519KeyPair() ISignKeyPair {
	pub, pri, _ := ced.GenerateKey(nil)
	kp := &ed25519KeyPair{&ed25519SignKey{pri}, &ed25519VerfKey{pub}}
	return kp
}
func (kp *ed25519KeyPair) Sign() ISignKey {
	return kp.signKey
}
func (kp *ed25519KeyPair) Verify() IVerfKey {
	return kp.verfKey
}

func (key *ed25519SignKey) Verify() IVerfKey {
	var priKey sed.PrivateKey = key.signKey
	pubKey := priKey.Public().(ced.PublicKey)
	return &ed25519VerfKey{pubKey}
}

func (key *ed25519SignKey) Sign(msg []byte) ([]byte, error) {
	if len(key.signKey) != ced.PrivateKeySize {
		return nil, util.NewError("invalid signKey")
	}
	return ced.Sign(key.signKey, msg), nil
}
func (key *ed25519VerfKey) Verify(msg, sig []byte) (bool, error) {
	if len(key.verfKey) != ced.PublicKeySize {
		return false, util.NewError("invalid verfKey")
	}
	return ced.Verify(key.verfKey, msg, sig), nil
}

func (key *ed25519SignKey) Raw() ([]byte, error) {
	return key.signKey, nil
}
func (key *ed25519SignKey) Unmarshal(b []byte) error {
	key.signKey = b
	return nil
}
func (key *ed25519VerfKey) Raw() ([]byte, error) {
	return key.verfKey, nil
}
func (key *ed25519VerfKey) Unmarshal(b []byte) error {
	key.verfKey = b
	return nil
}
