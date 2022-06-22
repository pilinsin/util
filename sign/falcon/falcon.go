package falcon

/*
import (
	"github.com/open-quantum-safe/liboqs-go/oqs"
	pb "github.com/pilinsin/util/pb"
	isign "github.com/pilinsin/util/sign"
	proto "google.golang.org/protobuf/proto"
)

type falconKeyPair struct {
	signKey *falconSignKey
	verfKey *falconVerfKey
}

func NewKeyPair() isign.ISignKeyPair {
	mode := "Falcon-1024"

	oqsMan := oqs.Signature{}
	defer oqsMan.Clean()
	if err := oqsMan.Init(mode, nil); err != nil {
		return nil
	}
	verf, _ := oqsMan.GenerateKeyPair()
	sk := oqsMan.ExportSecretKey()
	sign := make([]byte, len(sk))
	copy(sign, sk)

	verfKey := &falconVerfKey{verf, mode}
	signKey := &falconSignKey{sign, mode}
	return &falconKeyPair{signKey, verfKey}
}
func (kp *falconKeyPair) Sign() isign.ISignKey {
	return kp.signKey
}
func (kp *falconKeyPair) Verify() isign.IVerfKey {
	return kp.verfKey
}

type falconSignKey struct {
	signKey []byte
	mode    string
}

func (sk *falconSignKey) Sign(data []byte) ([]byte, error) {
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
func (sk *falconSignKey) Raw() ([]byte, error) {
	marshalSignKey := &pb.OqsKey{
		Data: sk.signKey,
		Mode: sk.mode,
	}
	m, err := proto.Marshal(marshalSignKey)
	return m, err
}
func (sk *falconSignKey) Unmarshal(m []byte) error {
	marshalSignKey := &pb.OqsKey{}
	if err := proto.Unmarshal(m, marshalSignKey); err != nil {
		return err
	}
	sk.signKey = marshalSignKey.GetData()
	sk.mode = marshalSignKey.GetMode()
	return nil
}
func UnmarshalSignKey(m []byte) (isign.ISignKey, error) {
	sk := &falconSignKey{}
	err := sk.Unmarshal(m)
	return sk, err
}

type falconVerfKey struct {
	verfKey []byte
	mode    string
}

func (vk *falconVerfKey) Verify(data, sign []byte) (bool, error) {
	mode := vk.mode

	oqsMan := oqs.Signature{}
	defer oqsMan.Clean()
	if err := oqsMan.Init(mode, nil); err != nil {
		return false, err
	}
	return oqsMan.Verify(data, sign, vk.verfKey)
}
func (vk *falconVerfKey) Raw() ([]byte, error) {
	marshalVerfKey := &pb.OqsKey{
		Data: vk.verfKey,
		Mode: vk.mode,
	}
	m, err := proto.Marshal(marshalVerfKey)
	return m, err
}
func (vk *falconVerfKey) Unmarshal(m []byte) error {
	marshalVerfKey := &pb.OqsKey{}
	if err := proto.Unmarshal(m, marshalVerfKey); err != nil {
		return err
	}

	vk.verfKey = marshalVerfKey.GetData()
	vk.mode = marshalVerfKey.GetMode()
	return nil
}
func UnmarshalVerfKey(m []byte) (isign.IVerfKey, error) {
	sk := &falconVerfKey{}
	err := sk.Unmarshal(m)
	return sk, err
}
*/
