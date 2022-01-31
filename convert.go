package util

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
)

func Uint8sToBytes(u8s []uint8) []byte {
	bs := make([]byte, len(u8s))
	for idx, u8 := range u8s {
		bs[idx] = byte(u8)
	}
	return bs
}
func BytesToUint8s(bs []byte) []uint8 {
	u8s := make([]uint8, len(bs))
	for idx, b := range bs {
		u8s[idx] = uint8(b)
	}
	return u8s
}

func BytesToReader(b []byte) io.Reader {
	return bytes.NewBuffer(b)
}

func ReaderToBytes(reader io.Reader) []byte {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(reader)
	if err != nil {
		return nil
	}

	return buf.Bytes()
}

func AnyBytes64ToStr(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}
func StrToAnyBytes64(str string) []byte {
	b, err := base64.StdEncoding.DecodeString(str)
	if err != nil {
		return nil
	} else {
		return b
	}
}
func Bytes64ToAnyStr(b []byte) string {
	b, err := base64.StdEncoding.DecodeString(string(b))
	if err != nil {
		return ""
	} else {
		return string(b)
	}
}
func AnyStrToBytes64(str string) []byte {
	b64Str := base64.StdEncoding.EncodeToString([]byte(str))
	return []byte(b64Str)
}

func Marshal(objWithPublicMembers interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	err := enc.Encode(objWithPublicMembers)
	return buf.Bytes(), err
}
func Unmarshal(b []byte, objWithPublicMembers interface{}) error {
	dec := json.NewDecoder(bytes.NewBuffer(b))
	dec.DisallowUnknownFields()
	err := dec.Decode(objWithPublicMembers)
	return err
}
