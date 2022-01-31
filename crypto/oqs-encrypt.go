package crypto

import (
	"github.com/open-quantum-safe/liboqs-go/oqs"
	"github.com/pilinsin/util"
)

type oqsEncryptKeyPair struct {
	pubKey *oqsPubKey
	priKey *oqsPriKey
}

func newOqsEncryptKeyPair(mode string) IPubEncryptKeyPair {
	var cipherSize int
	switch mode {
	case "SIDH-p751-compressed":
		cipherSize = 335
	case "NTRU-HPS-4096-821":
		cipherSize = 1230
	default:
		return nil
	}

	oqsMan := oqs.KeyEncapsulation{}
	defer oqsMan.Clean()
	if err := oqsMan.Init(mode, nil); err != nil {
		return nil
	}
	pub, _ := oqsMan.GenerateKeyPair()
	sk := oqsMan.ExportSecretKey()
	pri := make([]byte, len(sk))
	copy(pri, sk)

	pubKey := &oqsPubKey{pub, mode}
	priKey := &oqsPriKey{pri, mode, cipherSize}
	return &oqsEncryptKeyPair{pubKey, priKey}
}
func (kp *oqsEncryptKeyPair) Private() IPriKey {
	return kp.priKey
}
func (kp *oqsEncryptKeyPair) Public() IPubKey {
	return kp.pubKey
}

type oqsPriKey struct {
	priKey     []byte
	mode       string
	cipherSize int
}

func (pri *oqsPriKey) Decrypt(m []byte) ([]byte, error) {
	if len(m) <= pri.cipherSize {
		return nil, util.NewError("decrypt fail: len(m) <= cipherSize")
	}
	cipher, enc := m[:pri.cipherSize], m[pri.cipherSize:]

	oqsMan := oqs.KeyEncapsulation{}
	if err := oqsMan.Init(pri.mode, pri.priKey); err != nil {
		return nil, err
	}
	defer oqsMan.Clean()
	share, err := oqsMan.DecapSecret(cipher)
	if err != nil {
		return nil, err
	}
	return newMultiChaChaSharedKey(share).Decrypt(enc)
}
func (pri *oqsPriKey) Equals(pri2 IPriKey) bool {
	return util.ConstTimeBytesEqual(pri.Marshal(), pri2.Marshal())
}
func (pri *oqsPriKey) Marshal() []byte {
	mpri := &struct {
		Pr []byte
		M  string
		S  int
	}{pri.priKey, pri.mode, pri.cipherSize}
	m, _ := util.Marshal(mpri)
	return m
}
func (pri *oqsPriKey) Unmarshal(m []byte) error {
	mpri := &struct {
		Pr []byte
		M  string
		S  int
	}{}
	if err := util.Unmarshal(m, mpri); err != nil {
		return err
	}
	pri.priKey = mpri.Pr
	pri.mode = mpri.M
	pri.cipherSize = mpri.S
	return nil
}

type oqsPubKey struct {
	pubKey []byte
	mode   string
}

func (pub *oqsPubKey) Encrypt(data []byte) ([]byte, error) {
	oqsMan := oqs.KeyEncapsulation{}
	if err := oqsMan.Init(pub.mode, nil); err != nil {
		return nil, err
	}
	defer oqsMan.Clean()
	cipher, share, err := oqsMan.EncapSecret(pub.pubKey)
	if err != nil {
		return nil, err
	}
	enc, err := newMultiChaChaSharedKey(share).Encrypt(data)
	return append(cipher, enc...), err
}
func (pub *oqsPubKey) Equals(pub2 IPubKey) bool {
	return util.ConstTimeBytesEqual(pub.Marshal(), pub2.Marshal())
}
func (pub *oqsPubKey) Marshal() []byte {
	mpub := &struct {
		Pu []byte
		M  string
	}{pub.pubKey, pub.mode}
	m, _ := util.Marshal(mpub)
	return m
}
func (pub *oqsPubKey) Unmarshal(m []byte) error {
	mpub := &struct {
		Pu []byte
		M  string
	}{}
	if err := util.Unmarshal(m, mpub); err != nil {
		return err
	}
	pub.pubKey = mpub.Pu
	pub.mode = mpub.M
	return nil
}
