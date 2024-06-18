package hash

import (
	"crypto/sha1"
	"fmt"
)

// IStringHasher provides hashing logic for any string
type IStringHasher interface {
	Hash(msg string) (string, error)
}

// SHA1Hasher uses SHA1 to hash strings with provided salt.
type SHA1Hasher struct {
	salt string
}

func NewSHA1Hasher(salt string) *SHA1Hasher {
	return &SHA1Hasher{salt: salt}
}

func (h *SHA1Hasher) Hash(str string) (string, error) {
	hash := sha1.New()
	if _, err := hash.Write([]byte(str)); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum([]byte(h.salt))), nil
}
