package gotp

import (
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/binary"
	"fmt"
	"hash"
)

// ConvAlgoString converts an algorithm string to the hash function
func ConvAlgoString(algorithm string) (func() hash.Hash, error) {
	var algo func() hash.Hash
	switch algorithm {
	case "sha1", "SHA1":
		algo = sha1.New
	case "sha256", "SHA256":
		algo = sha256.New
	case "sha512", "SHA512":
		algo = sha512.New
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algorithm)
	}
	return algo, nil
}

// Itob fills uint64 a byte array
func Itob(v uint64) []byte {
	b := make([]byte, 8) // uint64 length is 8
	binary.BigEndian.PutUint64(b, uint64(v))
	return b
}
