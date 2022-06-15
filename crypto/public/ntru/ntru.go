package ntru

import (
	"errors"
	
	"github.com/open-quantum-safe/liboqs-go/oqs"
	proto "google.golang.org/protobuf/proto"
	pb "github.com/pilinsin/util/crypto/pb"
	ipub "github.com/pilinsin/util/crypto/public"
	mchacha "github.com/pilinsin/util/crypto/secret/multichacha"
)

type ntruKeyPair struct {
	pubKey *ntruPubKey
	priKey *ntruPriKey
}

func NewKeyPair() ipub.IPubEncryptKeyPair {
	mode := "NTRU-HPS-4096-821"
	cipherSize := 1230

	oqsMan := oqs.KeyEncapsulation{}
	defer oqsMan.Clean()
	if err := oqsMan.Init(mode, nil); err != nil {
		return nil
	}
	pub, _ := oqsMan.GenerateKeyPair()
	sk := oqsMan.ExportSecretKey()
	pri := make([]byte, len(sk))
	copy(pri, sk)

	pubKey := &ntruPubKey{pub, mode}
	priKey := &ntruPriKey{pri, mode, cipherSize}
	return &ntruKeyPair{pubKey, priKey}
}
func (kp *ntruKeyPair) Private() ipub.IPriKey {
	return kp.priKey
}
func (kp *ntruKeyPair) Public() ipub.IPubKey {
	return kp.pubKey
}

type ntruPriKey struct {
	priKey     []byte
	mode       string
	cipherSize int
}

func (pri *ntruPriKey) Decrypt(m []byte) ([]byte, error) {
	if len(m) <= pri.cipherSize {
		return nil, errors.New("decrypt fail: len(m) <= cipherSize")
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
	return mchacha.NewSecretKey(share).Decrypt(enc)
}
func (pri *ntruPriKey) Raw() ([]byte, error) {
	mpri := &pb.OqsPriKey{
		Data: pri.priKey,
		Mode: pri.mode,
		Size: int32(pri.cipherSize),
	}
	m, err := proto.Marshal(mpri)
	return m, err
}
func (pri *ntruPriKey) Unmarshal(m []byte) error {
	mpri := &pb.OqsPriKey{}
	if err := proto.Unmarshal(m, mpri); err != nil {
		return err
	}
	pri.priKey = mpri.GetData()
	pri.mode = mpri.GetMode()
	pri.cipherSize = int(mpri.GetSize())
	return nil
}
func UnmarshalPriKey(m []byte) (ipub.IPriKey, error){
	pri := &ntruPriKey{}
	err := pri.Unmarshal(m)
	return pri, err
}


type ntruPubKey struct {
	pubKey []byte
	mode   string
}

func (pub *ntruPubKey) Encrypt(data []byte) ([]byte, error) {
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
	enc, err := mchacha.NewSecretKey(share).Encrypt(data)
	return append(cipher, enc...), err
}
func (pub *ntruPubKey) Raw() ([]byte, error) {
	mpub := &pb.OqsKey{
		Data: pub.pubKey,
		Mode: pub.mode,
	}
	m, err := proto.Marshal(mpub)
	return m, err
}
func (pub *ntruPubKey) Unmarshal(m []byte) error {
	mpub := &pb.OqsKey{}
	if err := proto.Unmarshal(m, mpub); err != nil {
		return err
	}
	pub.pubKey = mpub.GetData()
	pub.mode = mpub.GetMode()
	return nil
}
func UnmarshalPubKey(m []byte) (ipub.IPubKey, error){
	pub := &ntruPubKey{}
	err := pub.Unmarshal(m)
	return pub, err
}
