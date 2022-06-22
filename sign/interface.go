package sign

type ISignKeyPair interface {
	Sign() ISignKey
	Verify() IVerfKey
}
type ISignKey interface {
	Raw() ([]byte, error)
	Sign(m []byte) ([]byte, error)
}
type IVerfKey interface {
	Raw() ([]byte, error)
	Verify(data, sign []byte) (bool, error)
}
