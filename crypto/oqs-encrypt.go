package crypto

import (
	pb "github.com/pilinsin/util/crypto/pb"
	proto "google.golang.org/protobuf/proto"

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

	priKey := make([]byte, len(pri.priKey))
	copy(priKey, pri.priKey)
	mode := pri.mode

	oqsMan := oqs.KeyEncapsulation{}
	if err := oqsMan.Init(mode, priKey); err != nil {
		return nil, err
	}
	defer oqsMan.Clean()
	share, err := oqsMan.DecapSecret(cipher)
	if err != nil {
		return nil, err
	}
	return newMultiChaChaSharedKey(share).Decrypt(enc)
}
func (pri *oqsPriKey) Raw() ([]byte, error) {
	mpri := &pb.OqsPriKey{
		Data: pri.priKey,
		Mode: pri.mode,
		Size: int32(pri.cipherSize),
	}
	m, err := proto.Marshal(mpri)
	return m, err
}
func (pri *oqsPriKey) Unmarshal(m []byte) error {
	mpri := &pb.OqsPriKey{}
	if err := proto.Unmarshal(m, mpri); err != nil {
		return err
	}
	pri.priKey = mpri.GetData()
	pri.mode = mpri.GetMode()
	pri.cipherSize = int(mpri.GetSize())
	return nil
}

type oqsPubKey struct {
	pubKey []byte
	mode   string
}

func (pub *oqsPubKey) Encrypt(data []byte) ([]byte, error) {
	mode := pub.mode
	oqsMan := oqs.KeyEncapsulation{}
	if err := oqsMan.Init(mode, nil); err != nil {
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
func (pub *oqsPubKey) Raw() ([]byte, error) {
	mpub := &pb.OqsKey{
		Data: pub.pubKey,
		Mode: pub.mode,
	}
	m, err := proto.Marshal(mpub)
	return m, err
}
func (pub *oqsPubKey) Unmarshal(m []byte) error {
	mpub := &pb.OqsKey{}
	if err := proto.Unmarshal(m, mpub); err != nil {
		return err
	}
	pub.pubKey = mpub.GetData()
	pub.mode = mpub.GetMode()
	return nil
}
