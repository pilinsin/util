package crypto

type iBaseKey interface{
	Raw() ([]byte, error)
}

type ISignKeyPair interface {
	Sign() ISignKey
	Verify() IVerfKey
}
type ISignKey interface {
	iBaseKey
	Sign(m []byte) ([]byte, error)
}
type IVerfKey interface {
	iBaseKey
	Verify(data, sign []byte) (bool, error)
}

type IPubEncryptKeyPair interface {
	Public() IPubKey
	Private() IPriKey
}
type IPriKey interface {
	iBaseKey
	Decrypt(m []byte) ([]byte, error)
}
type IPubKey interface {
	iBaseKey
	Encrypt(data []byte) ([]byte, error)
}

type ISharedKey interface {
	iBaseKey
	Encrypt(data []byte) ([]byte, error)
	Decrypt(m []byte) ([]byte, error)
}
