package sidh

/*
import (
	"errors"

	"github.com/open-quantum-safe/liboqs-go/oqs"
	iexch "github.com/pilinsin/util/exchange"
	pb "github.com/pilinsin/util/pb"
	proto "google.golang.org/protobuf/proto"
)

type sidhKeyPair struct {
	pubKey *sidhPubKey
	priKey *sidhPriKey
}

func NewKeyPair() iexch.IExchangeKeyPair {
	mode := "SIDH-p751-compressed"
	cipherSize := 335

	oqsMan := oqs.KeyEncapsulation{}
	defer oqsMan.Clean()
	if err := oqsMan.Init(mode, nil); err != nil {
		return nil
	}
	pub, _ := oqsMan.GenerateKeyPair()
	sk := oqsMan.ExportSecretKey()
	pri := make([]byte, len(sk))
	copy(pri, sk)

	pubKey := &sidhPubKey{pub, mode}
	priKey := &sidhPriKey{pri, mode, cipherSize}
	return &sidhKeyPair{pubKey, priKey}
}
func (kp *sidhKeyPair) Private() iexch.IPriKey {
	return kp.priKey
}
func (kp *sidhKeyPair) Public() iexch.IPubKey {
	return kp.pubKey
}

type sidhPriKey struct {
	priKey     []byte
	mode       string
	cipherSize int
}

func (pri *sidhPriKey) Decrypt(m []byte) ([]byte, error) {
	if len(m) != pri.cipherSize {
		return nil, errors.New("decrypt fail: len(m) != cipherSize")
	}

	priKey := make([]byte, len(pri.priKey))
	copy(priKey, pri.priKey)
	mode := pri.mode

	oqsMan := oqs.KeyEncapsulation{}
	if err := oqsMan.Init(mode, priKey); err != nil {
		return nil, err
	}
	defer oqsMan.Clean()

	//share: secret infomation
	return oqsMan.DecapSecret(m)
}
func (pri *sidhPriKey) Raw() ([]byte, error) {
	mpri := &pb.OqsPriKey{
		Data: pri.priKey,
		Mode: pri.mode,
		Size: int32(pri.cipherSize),
	}
	m, err := proto.Marshal(mpri)
	return m, err
}
func (pri *sidhPriKey) Unmarshal(m []byte) error {
	mpri := &pb.OqsPriKey{}
	if err := proto.Unmarshal(m, mpri); err != nil {
		return err
	}
	pri.priKey = mpri.GetData()
	pri.mode = mpri.GetMode()
	pri.cipherSize = int(mpri.GetSize())
	return nil
}
func UnmarshalPriKey(m []byte) (iexch.IPriKey, error) {
	pri := &sidhPriKey{}
	err := pri.Unmarshal(m)
	return pri, err
}

type sidhPubKey struct {
	pubKey []byte
	mode   string
}

func (pub *sidhPubKey) Encrypt() ([]byte, []byte, error) {
	mode := pub.mode
	oqsMan := oqs.KeyEncapsulation{}
	if err := oqsMan.Init(mode, nil); err != nil {
		return nil, nil, err
	}
	defer oqsMan.Clean()
	cipher, share, err := oqsMan.EncapSecret(pub.pubKey)
	if err != nil {
		return nil, nil, err
	}

	//cipher: public, share: private
	return cipher, share, nil
}
func (pub *sidhPubKey) Raw() ([]byte, error) {
	mpub := &pb.OqsKey{
		Data: pub.pubKey,
		Mode: pub.mode,
	}
	m, err := proto.Marshal(mpub)
	return m, err
}
func (pub *sidhPubKey) Unmarshal(m []byte) error {
	mpub := &pb.OqsKey{}
	if err := proto.Unmarshal(m, mpub); err != nil {
		return err
	}
	pub.pubKey = mpub.GetData()
	pub.mode = mpub.GetMode()
	return nil
}
func UnmarshalPubKey(m []byte) (iexch.IPubKey, error) {
	pub := &sidhPubKey{}
	err := pub.Unmarshal(m)
	return pub, err
}

*/
