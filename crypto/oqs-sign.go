package crypto

import (
	pb "github.com/pilinsin/util/crypto/pb"
	proto "google.golang.org/protobuf/proto"

	"github.com/open-quantum-safe/liboqs-go/oqs"
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
	signKey := make([]byte, len(sk.signKey))
	copy(signKey, sk.signKey)
	mode := sk.mode

	oqsMan := oqs.Signature{}
	defer oqsMan.Clean()
	if err := oqsMan.Init(mode, signKey); err != nil {
		return nil, err
	}
	return oqsMan.Sign(data)
}
func (sk *oqsSignKey) Raw() ([]byte, error) {
	marshalSignKey := &pb.OqsKey{
		Data: sk.signKey,
		Mode: sk.mode,
	}
	m, err := proto.Marshal(marshalSignKey)
	return m, err
}
func (sk *oqsSignKey) Unmarshal(m []byte) error {
	marshalSignKey := &pb.OqsKey{}
	if err := proto.Unmarshal(m, marshalSignKey); err != nil {
		return err
	}
	sk.signKey = marshalSignKey.GetData()
	sk.mode = marshalSignKey.GetMode()
	return nil
}

type oqsVerfKey struct {
	verfKey []byte
	mode    string
}

func (vk *oqsVerfKey) Verify(data, sign []byte) (bool, error) {
	mode := vk.mode

	oqsMan := oqs.Signature{}
	defer oqsMan.Clean()
	if err := oqsMan.Init(mode, nil); err != nil {
		return false, err
	}
	return oqsMan.Verify(data, sign, vk.verfKey)
}
func (vk *oqsVerfKey) Raw() ([]byte, error) {
	marshalVerfKey := &pb.OqsKey{
		Data: vk.verfKey,
		Mode: vk.mode,
	}
	m, err := proto.Marshal(marshalVerfKey)
	return m, err
}
func (vk *oqsVerfKey) Unmarshal(m []byte) error {
	marshalVerfKey := &pb.OqsKey{}
	if err := proto.Unmarshal(m, marshalVerfKey); err != nil {
		return err
	}

	vk.verfKey = marshalVerfKey.GetData()
	vk.mode = marshalVerfKey.GetMode()
	return nil
}
