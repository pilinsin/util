package crypto

import (
	"github.com/open-quantum-safe/liboqs-go/oqs"
	"github.com/pilinsin/util"
)

type oqsSignKeyPair struct {
	signKey *oqsSignKey
	verfKey *oqsVerfKey
}

func newOqsSignKeyPair(mode string) ISignKeyPair {
	switch mode {
	case "SPHINCS+-SHAKE256-256s-robust":
	case "Falcon-1024":
	default:
		return nil
	}

	oqsMan := oqs.Signature{}
	defer oqsMan.Clean()
	if err := oqsMan.Init(mode, nil); err != nil {
		return nil
	}
	verf, _ := oqsMan.GenerateKeyPair()
	sk := oqsMan.ExportSecretKey()
	sign := make([]byte, len(sk))
	copy(sign, sk)

	verfKey := &oqsVerfKey{verf, mode}
	signKey := &oqsSignKey{sign, mode}
	return &oqsSignKeyPair{signKey, verfKey}
}
func (kp *oqsSignKeyPair) Sign() ISignKey {
	return kp.signKey
}
func (kp *oqsSignKeyPair) Verify() IVerfKey {
	return kp.verfKey
}

type oqsSignKey struct {
	signKey []byte
	mode    string
}

func (sk *oqsSignKey) Sign(data []byte) ([]byte, error) {
	oqsMan := oqs.Signature{}
	defer oqsMan.Clean()
	if err := oqsMan.Init(sk.mode, sk.signKey); err != nil {
		return nil, err
	}
	return oqsMan.Sign(data)
}
func (sk *oqsSignKey) Equals(sk2 ISignKey) bool {
	m := sk.Marshal()
	m2 := sk2.Marshal()
	return util.ConstTimeBytesEqual(m, m2)
}
func (sk *oqsSignKey) Marshal() []byte {
	marshalSignKey := &struct {
		S []byte
		M string
	}{sk.signKey, sk.mode}
	m, _ := util.Marshal(marshalSignKey)
	return m
}
func (sk *oqsSignKey) Unmarshal(m []byte) error {
	marshalSignKey := &struct {
		S []byte
		M string
	}{}
	if err := util.Unmarshal(m, marshalSignKey); err != nil {
		return err
	}

	sk.signKey = marshalSignKey.S
	sk.mode = marshalSignKey.M
	return nil
}

type oqsVerfKey struct {
	verfKey []byte
	mode    string
}

func (vk *oqsVerfKey) Verify(data, sign []byte) (bool, error) {
	oqsMan := oqs.Signature{}
	defer oqsMan.Clean()
	if err := oqsMan.Init(vk.mode, nil); err != nil {
		return false, err
	}
	return oqsMan.Verify(data, sign, vk.verfKey)
}
func (vk *oqsVerfKey) Equals(vk2 IVerfKey) bool {
	m := vk.Marshal()
	m2 := vk2.Marshal()
	return util.ConstTimeBytesEqual(m, m2)
}
func (vk *oqsVerfKey) Marshal() []byte {
	marshalVerfKey := &struct {
		V []byte
		M string
	}{vk.verfKey, vk.mode}
	m, _ := util.Marshal(marshalVerfKey)
	return m
}
func (vk *oqsVerfKey) Unmarshal(m []byte) error {
	marshalVerfKey := &struct {
		V []byte
		M string
	}{}
	if err := util.Unmarshal(m, marshalVerfKey); err != nil {
		return err
	}

	vk.verfKey = marshalVerfKey.V
	vk.mode = marshalVerfKey.M
	return nil
}
