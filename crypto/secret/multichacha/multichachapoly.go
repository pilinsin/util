package multichachapoly

import (
	proto "google.golang.org/protobuf/proto"
	pb "github.com/pilinsin/util/crypto/pb"

	isec "github.com/pilinsin/util/crypto/secret"
	chacha "github.com/pilinsin/util/crypto/secret/chacha"
	hash "github.com/pilinsin/util/crypto/hash"
)

func arange(start, stop, step int) []int {
	if start < 0 {
		start = 0
	}
	if stop < 0 {
		stop = 1
	}
	if step < 0 {
		step = 1
	}

	var arr []int
	for i := start; i < stop; i += step {
		arr = append(arr, i)
	}
	return arr
}

func multiSeedPadding(seed []byte) []byte {
	paddedSize := len(seed) + isec.SecretKeySize - (len(seed) % isec.SecretKeySize)
	indices := arange(0, len(seed), 2)
	salt := make([]byte, len(indices))
	for salt_idx, seed_idx := range indices {
		salt[salt_idx] = seed[seed_idx]
	}
	return hash.HashWithSize(seed, salt, paddedSize)
}
func splitMultiSecretKeySeed(seed []byte) [][isec.SecretKeySize]byte {
	if len(seed)%isec.SecretKeySize != 0{
		return nil
	}
	nSeeds := len(seed) / isec.SecretKeySize

	seeds := make([][isec.SecretKeySize]byte, nSeeds)
	for idx, _ := range seeds {
		begin := idx * isec.SecretKeySize
		end := (idx + 1) * isec.SecretKeySize
		copy(seeds[idx][:], seed[begin:end])
	}
	return seeds
}
func reverse(seeds [][isec.SecretKeySize]byte) [][isec.SecretKeySize]byte {
	seeds2 := make([][isec.SecretKeySize]byte, len(seeds))
	for idx, _ := range seeds {
		seeds2[len(seeds)-1-idx] = seeds[idx]
	}
	return seeds2
}

type multiChachaSecretKey struct {
	seeds [][isec.SecretKeySize]byte
}

func NewSecretKey(seed []byte) isec.ISecretKey {
	seed = multiSeedPadding(seed)
	seeds := splitMultiSecretKeySeed(seed)
	return &multiChachaSecretKey{seeds}
}
func (key multiChachaSecretKey) Encrypt(data []byte) ([]byte, error) {
	var err, tmpErr error
	for _, seed := range key.seeds {
		cha := chacha.NewSecretKey(seed)
		if data, tmpErr = cha.Encrypt(data); tmpErr != nil {
			err = tmpErr
		}
	}

	return data, err
}
func (key multiChachaSecretKey) Decrypt(m []byte) ([]byte, error) {
	var err, tmpErr error
	for _, seed := range reverse(key.seeds) {
		cha := chacha.NewSecretKey(seed)
		if m, tmpErr = cha.Decrypt(m); tmpErr != nil {
			err = tmpErr
		}
	}
	return m, err
}

func (key multiChachaSecretKey) Raw() ([]byte, error) {
	seeds := make([][]byte, len(key.seeds))
	for idx, seed := range key.seeds{
		seeds[idx] = seed[:]
	}
	mSeeds := &pb.MultiChachaKey{
		Seeds: seeds,
	}
	m, err := proto.Marshal(mSeeds)
	return m, err
}
func (key *multiChachaSecretKey) Unmarshal(m []byte) error {
	mSeeds := &pb.MultiChachaKey{}
	if err := proto.Unmarshal(m, mSeeds); err != nil {
		return err
	}

	seeds := make([][isec.SecretKeySize]byte, len(mSeeds.GetSeeds()))
	for idx, seed := range mSeeds.GetSeeds(){
		copy(seeds[idx][:], seed)
	}

	key.seeds = seeds
	return nil
}
func UnmarshalSecretKey(m []byte) (isec.ISecretKey, error){
	sk := &multiChachaSecretKey{}
	err := sk.Unmarshal(m)
	return sk, err
}
