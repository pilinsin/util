package public

type IPubEncryptKeyPair interface {
	Public() IPubKey
	Private() IPriKey
}
type IPriKey interface {
	Raw() ([]byte, error)
	Decrypt(m []byte) ([]byte, error)
}
type IPubKey interface {
	Raw() ([]byte, error)
	Encrypt(data []byte) ([]byte, error)
}