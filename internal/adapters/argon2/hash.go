package argon2

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"
)

const (
	SaltLength = 16
	format     = "$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s"
)

func (a *Argon2Hasher) Hash(password string) (string, error) {
	salt := make([]byte, SaltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	hashedBytes := argon2.IDKey([]byte(password), salt, a.iterations, a.memory, a.parallelism, a.keyLength)

	saltBase64 := base64.RawStdEncoding.EncodeToString(salt)
	hashBase64 := base64.RawStdEncoding.EncodeToString(hashedBytes)

	fullHash := fmt.Sprintf(format, a.memory, a.iterations, a.parallelism, saltBase64, hashBase64)
	return fullHash, nil
}

func (a *Argon2Hasher) Check(password, hashedPassword string) (bool, error) {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 6 {
		return false, errors.New("invalid hash format")
	}
	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	originalHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		a.iterations,
		a.memory,
		a.parallelism,
		a.keyLength,
	)
	if subtle.ConstantTimeCompare(computedHash, originalHash) == 1 {
		return true, nil
	}
	return false, nil
}
