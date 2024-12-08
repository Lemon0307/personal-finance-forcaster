package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"io"
)

func GenerateSalt(size int) ([]byte, error) {
	salt := make([]byte, size)
	_, err := io.ReadFull(rand.Reader, salt)
	if err != nil {
		return nil, err
	}
	return salt, nil
}

func (user *User) HashPassword(salt []byte) {
	hash := sha256.New()
	hash.Write(salt)
	hash.Write([]byte(user.Password))
	user.Password = base64.RawStdEncoding.EncodeToString(hash.Sum(nil))
}
