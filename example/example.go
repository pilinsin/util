package main

import (
	"fmt"

	"github.com/pilinsin/util"
	"github.com/pilinsin/util/crypto"
)

func main() {
	data := []byte("meow meow ^.^")

	crypto.SelectedPubEncMode = crypto.Sidh
	exampleEncrypt(data)
	crypto.SelectedPubEncMode = crypto.Ntru
	exampleEncrypt(data)
	crypto.SelectedSignMode = crypto.Sphincs
	exampleSign(data)
	crypto.SelectedSignMode = crypto.Falcon
	exampleSign(data)
}

func exampleEncrypt(data []byte) {
	kp := crypto.NewPubEncryptKeyPair()
	fmt.Println(kp.Public() != nil, kp.Private() != nil)
	enc, err := kp.Public().Encrypt(data)
	util.CheckError(err)
	data1, err := kp.Private().Decrypt(enc)
	util.CheckError(err)
	fmt.Println(util.ConstTimeBytesEqual(data, data1))
	mpu := kp.Public().Marshal()
	pub1, err := crypto.UnmarshalPubKey(mpu)
	util.CheckError(err)
	fmt.Println(kp.Public().Equals(pub1))
	mpr := kp.Private().Marshal()
	pri1, err := crypto.UnmarshalPriKey(mpr)
	util.CheckError(err)
	fmt.Println(kp.Private().Equals(pri1))
}

func exampleSign(data []byte) {
	kp := crypto.NewSignKeyPair()
	sg, err := kp.Sign().Sign(data)
	util.CheckError(err)
	ok, err := kp.Verify().Verify(data, sg)
	util.CheckError(err)
	fmt.Println(ok)
	mv := kp.Verify().Marshal()
	v1, err := crypto.UnmarshalVerfKey(mv)
	util.CheckError(err)
	fmt.Println(kp.Verify().Equals(v1))
	ms := kp.Sign().Marshal()
	s1, err := crypto.UnmarshalSignKey(ms)
	util.CheckError(err)
	fmt.Println(kp.Sign().Equals(s1))
}
