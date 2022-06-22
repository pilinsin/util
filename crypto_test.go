package util

import (
	"crypto/rand"
	"testing"

	isec "github.com/pilinsin/util/secret"
	aes "github.com/pilinsin/util/secret/aes"
	chacha "github.com/pilinsin/util/secret/chacha"
	mchacha "github.com/pilinsin/util/secret/multichacha"

	//ntru "github.com/pilinsin/util/exchange/ntru"
	//sidh "github.com/pilinsin/util/exchange/sidh"

	ipub "github.com/pilinsin/util/public"
	ecies "github.com/pilinsin/util/public/ecies"

	//pntru "github.com/pilinsin/util/public/ntru"
	//psidh "github.com/pilinsin/util/public/sidh"

	isign "github.com/pilinsin/util/sign"
	ed25519 "github.com/pilinsin/util/sign/ed25519"
	//falcon "github.com/pilinsin/util/sign/falcon"
	//sphincs "github.com/pilinsin/util/sign/sphincs"
)

func checkError(t *testing.T, err error, args ...interface{}) {
	if err != nil {
		args0 := make([]interface{}, len(args)+1)
		args0[0] = err
		copy(args0[1:], args)

		t.Fatal(args0...)
	}
}
func assertError(t *testing.T, cond bool, args ...interface{}) {
	if !cond {
		t.Fatal(args...)
	}
}

func secretEncDecTest(t *testing.T, sh isec.ISecretKey) bool {
	data := []byte("meow meow ^_^")
	enc, err := sh.Encrypt(data)
	checkError(t, err)
	dec, err := sh.Decrypt(enc)
	checkError(t, err)
	return string(data) == string(dec)
}

type secretKeyGenerator func([]byte) isec.ISecretKey
type secretUnmarshaller func([]byte) (isec.ISecretKey, error)

func secretTest(t *testing.T, kGen secretKeyGenerator, sUm secretUnmarshaller) {
	seed := make([]byte, isec.SecretKeySize)
	rand.Read(seed)
	sk := kGen(seed)

	t.Log("encrypt-decrypt test")
	ok := secretEncDecTest(t, sk)
	assertError(t, ok, "encrypt-decrypt error")

	t.Log("marshal-unmarshal secret test")
	m, err := sk.Raw()
	checkError(t, err)
	sk2, err := sUm(m)
	checkError(t, err)
	ok = secretEncDecTest(t, sk2)
	assertError(t, ok, "marshal-unmarshal secret error")
}
func TestSecret(t *testing.T) {
	aesNewKey := func(seed []byte) isec.ISecretKey {
		tmp := [isec.SecretKeySize]byte{}
		copy(tmp[:], seed)
		return aes.NewSecretKey(tmp)
	}
	secretTest(t, aesNewKey, aes.UnmarshalSecretKey)

	chaNewKey := func(seed []byte) isec.ISecretKey {
		tmp := [isec.SecretKeySize]byte{}
		copy(tmp[:], seed)
		return chacha.NewSecretKey(tmp)
	}
	secretTest(t, chaNewKey, chacha.UnmarshalSecretKey)
	secretTest(t, mchacha.NewSecretKey, mchacha.UnmarshalSecretKey)
}

/*
func cipherTest(t *testing.T, pub iexch.IPubKey, priv iexch.IPriKey) bool {
	cipher, share, err := pub.Encrypt()
	checkError(t, err)
	share2, err := priv.Decrypt(cipher)
	checkError(t, err)
	return string(share) == string(share2)
}

type exPriUnmarshaller func([]byte) (iexch.IPriKey, error)
type exPubUnmarshaller func([]byte) (iexch.IPubKey, error)

func exchTest(t *testing.T, kp iexch.IExchangeKeyPair, pru exPriUnmarshaller, puu exPubUnmarshaller) {
	pub := kp.Public()
	priv := kp.Private()

	t.Log("encrypt-decrypt test")
	ok := cipherTest(t, pub, priv)
	assertError(t, ok, "encrypt-decrypt error")

	t.Log("marshal-unmarshal pub test")
	m, err := pub.Raw()
	checkError(t, err)
	pub2, err := puu(m)
	checkError(t, err)
	ok = cipherTest(t, pub2, priv)
	assertError(t, ok, "marshal-unmarshal pub error")

	t.Log("marshal-unmarshal priv test")
	m, err = priv.Raw()
	checkError(t, err)
	priv2, err := pru(m)
	checkError(t, err)
	ok = cipherTest(t, pub, priv2)
	assertError(t, ok, "marshal-unmarshal priv error")
}
func TestExch(t *testing.T) {
	exchTest(t, sidh.NewKeyPair(), sidh.UnmarshalPriKey, sidh.UnmarshalPubKey)
	exchTest(t, ntru.NewKeyPair(), ntru.UnmarshalPriKey, ntru.UnmarshalPubKey)
}
*/

func encDecTest(t *testing.T, pub ipub.IPubKey, priv ipub.IPriKey) bool {
	data := []byte("meow meow ^_^")
	enc, err := pub.Encrypt(data)
	checkError(t, err)
	dec, err := priv.Decrypt(enc)
	checkError(t, err)
	return string(data) == string(dec)
}

type priUnmarshaller func([]byte) (ipub.IPriKey, error)
type pubUnmarshaller func([]byte) (ipub.IPubKey, error)

func pubTest(t *testing.T, kp ipub.IPubEncryptKeyPair, pru priUnmarshaller, puu pubUnmarshaller) {
	pub := kp.Public()
	priv := kp.Private()

	t.Log("encrypt-decrypt test")
	ok := encDecTest(t, pub, priv)
	assertError(t, ok, "encrypt-decrypt error")

	t.Log("marshal-unmarshal pub test")
	m, err := pub.Raw()
	checkError(t, err)
	pub2, err := puu(m)
	checkError(t, err)
	ok = encDecTest(t, pub2, priv)
	assertError(t, ok, "marshal-unmarshal pub error")

	t.Log("marshal-unmarshal priv test")
	m, err = priv.Raw()
	checkError(t, err)
	priv2, err := pru(m)
	checkError(t, err)
	ok = encDecTest(t, pub, priv2)
	assertError(t, ok, "marshal-unmarshal priv error")
}
func TestPub(t *testing.T) {
	pubTest(t, ecies.NewKeyPair(), ecies.UnmarshalPriKey, ecies.UnmarshalPubKey)
	//pubTest(t, psidh.NewKeyPair(), psidh.UnmarshalPriKey, psidh.UnmarshalPubKey)
	//pubTest(t, pntru.NewKeyPair(), pntru.UnmarshalPriKey, pntru.UnmarshalPubKey)
}

func signVerifyTest(t *testing.T, sk isign.ISignKey, vk isign.IVerfKey) bool {
	data := []byte("meow meow ^_^")
	sig, err := sk.Sign(data)
	checkError(t, err)
	ok, err := vk.Verify(data, sig)
	checkError(t, err)
	return ok
}

type signUnmarshaller func([]byte) (isign.ISignKey, error)
type verfUnmarshaller func([]byte) (isign.IVerfKey, error)

func signTest(t *testing.T, kp isign.ISignKeyPair, sum signUnmarshaller, vum verfUnmarshaller) {
	sk := kp.Sign()
	vk := kp.Verify()

	t.Log("sign-verify test")
	ok := signVerifyTest(t, sk, vk)
	assertError(t, ok, "sign-verify error")

	t.Log("marshal-unmarshal sign test")
	m, err := sk.Raw()
	checkError(t, err)
	sk2, err := sum(m)
	checkError(t, err)
	ok = signVerifyTest(t, sk2, vk)
	assertError(t, ok, "marshal-unmarshal sign error")

	t.Log("marshal-unmarshal verify test")
	m, err = vk.Raw()
	checkError(t, err)
	vk2, err := vum(m)
	checkError(t, err)
	ok = signVerifyTest(t, sk, vk2)
	assertError(t, ok, "marshal-unmarshal verify error")
}
func TestSign(t *testing.T) {
	signTest(t, ed25519.NewKeyPair(), ed25519.UnmarshalSignKey, ed25519.UnmarshalVerfKey)
	//signTest(t, sphincs.NewKeyPair(), sphincs.UnmarshalSignKey, sphincs.UnmarshalVerfKey)
	//signTest(t, falcon.NewKeyPair(), falcon.UnmarshalSignKey, falcon.UnmarshalVerfKey)
}
