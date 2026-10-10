package argon2

import "github.com/alvarolucio2007/Scholarly/internal/ports"

var _ ports.PasswordHasher = (*Argon2Hasher)(nil)

type Argon2Hasher struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	keyLength   uint32
}

func NewArgon2Hasher(memory, iterations uint32, parallelism uint8) *Argon2Hasher {
	return &Argon2Hasher{
		memory:      memory,      //64*1024
		iterations:  iterations,  // 3
		parallelism: parallelism, // 2
		keyLength:   32,
	}
}
