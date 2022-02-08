package crypto

import (
	"github.com/pilinsin/util"
)

func multiSeedPadding(seed []byte) []byte {
	paddedSize := len(seed) + SharedKeySize - (len(seed) % SharedKeySize)
	indices := util.Arange(0, len(seed), 2)
	salt := make([]byte, len(indices))
	for salt_idx, seed_idx := range indices {
		salt[salt_idx] = seed[seed_idx]
	}
	return HashWithSize(seed, salt, paddedSize)
}
func splitMultiSharedKeySeed(seed []byte) [][SharedKeySize]byte {
	if len(seed)%SharedKeySize != 0 {
		return nil
	}
	nSeeds := len(seed) / SharedKeySize

	seeds := make([][SharedKeySize]byte, nSeeds)
	for idx, _ := range seeds {
		begin := idx * SharedKeySize
		end := (idx + 1) * SharedKeySize
		copy(seeds[idx][:], seed[begin:end])
	}
	return seeds
}
func reverse(seeds [][SharedKeySize]byte) [][SharedKeySize]byte {
	seeds2 := make([][SharedKeySize]byte, len(seeds))
	for idx, _ := range seeds {
		seeds2[len(seeds)-1-idx] = seeds[idx]
	}
	return seeds2
}

type multiChachaSharedKey struct {
	seeds [][SharedKeySize]byte
}

func newMultiChaChaSharedKey(seed []byte) ISharedKey {
	seed = multiSeedPadding(seed)
	seeds := splitMultiSharedKeySeed(seed)
	return &multiChachaSharedKey{seeds}
}
func (key multiChachaSharedKey) Encrypt(data []byte) ([]byte, error) {
	var err, tmpErr error
	for _, seed := range key.seeds {
		cha := newChaChaSharedKey(seed)
		if data, tmpErr = cha.Encrypt(data); tmpErr != nil {
			err = tmpErr
		}
	}
	return data, err
}
func (key multiChachaSharedKey) Decrypt(m []byte) ([]byte, error) {
	var err, tmpErr error
	for _, seed := range reverse(key.seeds) {
		cha := newChaChaSharedKey(seed)
		if m, tmpErr = cha.Decrypt(m); tmpErr != nil {
			err = tmpErr
		}
	}
	return m, err
}

func (key multiChachaSharedKey) Equals(key2 ISharedKey) bool {
	return util.ConstTimeBytesEqual(key.Marshal(), key2.Marshal())
}

func (key multiChachaSharedKey) Marshal() []byte {
	m, _ := util.Marshal(key.seeds)
	return m
}
func (key multiChachaSharedKey) Unmarshal(m []byte) error {
	var seeds [][SharedKeySize]byte
	if err := util.Unmarshal(m, seeds); err != nil {
		return err
	}

	key.seeds = seeds
	return nil
}
