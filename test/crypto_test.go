package test

import (
	"github.com/pilinsin/util/crypto"
	"testing"
)

func TestCrypto(t *testing.T) {
	crypto.SelectedPubEncMode = crypto.Sike
	crypto.NewPubEncryptKeyPair()
	crypto.SelectedPubEncMode = crypto.Ntru
	crypto.NewPubEncryptKeyPair()

	crypto.SelectedSignMode = crypto.Sphincs
	crypto.NewSignKeyPair()
	crypto.SelectedSignMode = crypto.Falcon
	crypto.NewSignKeyPair()
}
