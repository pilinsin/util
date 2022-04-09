package crypto

import (
	pb "github.com/pilinsin/util/crypto/pb"
	proto "google.golang.org/protobuf/proto"

	"github.com/LoCCS/bliss"
	"github.com/LoCCS/bliss/sampler"
	"github.com/pilinsin/util"
)

type blissKeyPair struct {
	signKey *blissSignKey
	verfKey *blissVerfKey
}

func newBlissKeyPair() ISignKeyPair {
	seedSize := int(sampler.SHA_512_DIGEST_LENGTH)
	seed := util.BytesToUint8s(util.GenRandomBytes(seedSize))
	e, _ := sampler.NewEntropy(seed)

	priKey, _ := bliss.GeneratePrivateKey(4, e)
	signKey := &blissSignKey{priKey, seed}
	pubKey := priKey.PublicKey()
	verfKey := &blissVerfKey{pubKey}

	return &blissKeyPair{signKey, verfKey}
}
func (kp *blissKeyPair) Sign() ISignKey {
	return kp.signKey
}
func (kp *blissKeyPair) Verify() IVerfKey {
	return kp.verfKey
}

type blissSignKey struct {
	signKey *bliss.PrivateKey
	seed    []uint8
}

func (sk *blissSignKey) Close() {
	sk.signKey.Destroy()
}
func (sk *blissSignKey) Sign(data []byte) ([]byte, error) {
	e, _ := sampler.NewEntropy(sk.seed)
	if signature, err := sk.signKey.SignAgainstSideChannel(data, e); err != nil {
		return nil, err
	} else {
		return signature.Encode(), nil
	}
}
func (sk *blissSignKey) Verify() IVerfKey {
	verfKey := sk.signKey.PublicKey()
	return &blissVerfKey{verfKey}
}
func (sk *blissSignKey) Raw() ([]byte, error) {
	marshalSignKey := &pb.BlissSignKey{
		Data: sk.signKey.Encode(),
		Seed: sk.seed,
	}
	m, err := proto.Marshal(marshalSignKey)
	return m, err
}
func (sk *blissSignKey) Unmarshal(m []byte) error {
	marshalSignKey := &pb.BlissSignKey{}
	if err := proto.Unmarshal(m, marshalSignKey); err != nil {
		return err
	}
	signKey, err := bliss.DecodePrivateKey(marshalSignKey.GetData())
	if err != nil {
		return err
	}

	sk.signKey = signKey
	sk.seed = marshalSignKey.GetSeed()
	return nil
}

type blissVerfKey struct {
	verfKey *bliss.PublicKey
}

func (vk *blissVerfKey) Verify(data, sign []byte) (bool, error) {
	if signature, err := bliss.DecodeSignature(sign); err != nil {
		return false, err
	} else {
		ok, err := vk.verfKey.Verify(data, signature)
		return ok, err
	}
}
func (vk *blissVerfKey) Raw() ([]byte, error) {
	return vk.verfKey.Encode(), nil
}
func (vk *blissVerfKey) Unmarshal(m []byte) error {
	if verfKey, err := bliss.DecodePublicKey(m); err != nil {
		return nil
	} else {
		vk.verfKey = verfKey
		return nil
	}
}
