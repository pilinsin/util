package crypto

import(
	"testing"
	"crypto/rand"
)

func checkError(t *testing.T, err error, args ...interface{}) {
	if err != nil {
		args0 := make([]interface{}, len(args)+1)
		args0[0] = err
		copy(args0[1:], args)

		t.Fatal(args0...)
	}
}
func assertError(t *testing.T, cond bool, args ...interface{}){
	if !cond{
		t.Fatal(args...)
	}
}

func encDecTest(t *testing.T, pub IPubKey, priv IPriKey) bool{
	data := []byte("meow meow ^_^")
	enc, err := pub.Encrypt(data)
	checkError(t, err)
	dec, err := priv.Decrypt(enc)
	checkError(t, err)
	return string(data) == string(dec)
}
func pubTest(t *testing.T, mode pubEncryptMode){
	SelectedPubEncMode = mode
	kp := NewPubEncryptKeyPair()
	pub := kp.Public()
	priv := kp.Private()

	t.Log("encrypt-decrypt test")
	ok := encDecTest(t, pub, priv)
	assertError(t, ok, "encrypt-decrypt error")

	t.Log("marshal-unmarshal pub test")
	m, err := MarshalPubKey(pub)
	checkError(t, err)
	pub2, err := UnmarshalPubKey(m)
	checkError(t, err)
	ok = encDecTest(t, pub2, priv)
	assertError(t, ok, "marshal-unmarshal pub error")

	t.Log("marshal-unmarshal priv test")
	m, err = MarshalPriKey(priv)
	checkError(t, err)
	priv2, err := UnmarshalPriKey(m)
	checkError(t, err)
	ok = encDecTest(t, pub, priv2)
	assertError(t, ok, "marshal-unmarshal priv error")
}
func TestPub(t *testing.T){
	pubTest(t, Sidh)
	pubTest(t, Ntru)
	pubTest(t, Sntrup)
	pubTest(t, Ecies)
}


func signVerifyTest(t *testing.T, sk ISignKey, vk IVerfKey) bool{
	data := []byte("meow meow ^_^")
	sig, err := sk.Sign(data)
	checkError(t, err)
	ok, err := vk.Verify(data, sig)
	checkError(t, err)
	return ok
}
func signTest(t *testing.T, mode signMode){
	SelectedSignMode = mode
	kp := NewSignKeyPair()
	sk := kp.Sign()
	vk := kp.Verify()

	t.Log("sign-verify test")
	ok := signVerifyTest(t, sk, vk)
	assertError(t, ok, "sign-verify error")

	t.Log("marshal-unmarshal sign test")
	m, err := MarshalSignKey(sk)
	checkError(t, err)
	sk2, err := UnmarshalSignKey(m)
	checkError(t, err)
	ok = signVerifyTest(t, sk2, vk)
	assertError(t, ok, "marshal-unmarshal sign error")

	t.Log("marshal-unmarshal verify test")
	m, err = MarshalVerfKey(vk)
	checkError(t, err)
	vk2, err := UnmarshalVerfKey(m)
	checkError(t, err)
	ok = signVerifyTest(t, sk, vk2)
	assertError(t, ok, "marshal-unmarshal verify error")
}
func TestSign(t *testing.T){
	signTest(t, Sphincs)
	signTest(t, Falcon)
	signTest(t, Bliss)
	signTest(t, Ed25519)
}


func sharedEncDecTest(t *testing.T, sh ISharedKey) bool{
	data := []byte("meow meow ^_^")
	enc, err := sh.Encrypt(data)
	checkError(t, err)
	dec, err := sh.Decrypt(enc)
	checkError(t, err)
	return string(data) == string(dec)
}
func sharedTest(t *testing.T, mode sharedEncryptMode){
	SelectedSharedEncMode = mode
	seed := make([]byte, SharedKeySize)
	rand.Read(seed)
	sh := NewSharedEncryptKey(seed)

	t.Log("encrypt-decrypt test")
	ok := sharedEncDecTest(t, sh)
	assertError(t, ok, "encrypt-decrypt error")

	t.Log("marshal-unmarshal shared test")
	m, err := MarshalSharedKey(sh)
	checkError(t, err)
	sh2, err := UnmarshalSharedKey(m)
	checkError(t, err)
	ok = sharedEncDecTest(t, sh2)
	assertError(t, ok, "marshal-unmarshal shared error")
}
func TestShared(t *testing.T){
	sharedTest(t, ChaCha)
	sharedTest(t, AES)
}