package util

import (
	"crypto/subtle"
	"bytes"
	"os"
	"path/filepath"
	"strconv"
)

func ExeDirPath() string {
	exe, _ := os.Executable()
	return filepath.Dir(exe)
}
func PathJoin(base string, adders ...string) string {
	for _, adder := range adders {
		base = filepath.Join(base, adder)
	}
	return base
}

func BoolPtr(b bool) *bool {
	return &b
}

func StrPtr(s string) *string {
	return &s
}
func IntPtr(n int) *int{
	return &n
} 

func ConstTimeBytesEqual(b1, b2 []byte) bool {
	return subtle.ConstantTimeCompare(b1, b2) == 1
}

//m1 > m2
func MapContainMap(m1, m2 map[string][]byte) bool {
	for k2, v2 := range m2 {
		if v1, ok := m1[k2]; !ok {
			return false
		} else if v2 != nil && !bytes.Equal(v1, v2) {
			return false
		}
	}
	return true
}

func StrSliceContain(sl []string, elem string) bool {
	for _, s := range sl {
		if s == elem {
			return true
		}
	}
	return false
}

func Arange(start, stop, step int) []int {
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

func ArangeStr(start, stop, step int) []string {
	if start < 0 {
		start = 0
	}
	if stop < 0 {
		stop = 1
	}
	if step < 0 {
		step = 1
	}

	var arr []string
	for i := start; i < stop; i += step {
		arr = append(arr, strconv.Itoa(i))
	}
	return arr
}
