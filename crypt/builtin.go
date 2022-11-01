package crypt

import (
	"crypto/rand"
)

const (
	schemeArgon2i  = "argon2i"
	schemeArgon2id = "argon2id"
)

var (
	defaultArgon2        = crypterArgon2{rand.Reader}
	defaultArgon2Crypter = crypter{
		schemeArgon2i:  &defaultArgon2,
		schemeArgon2id: &defaultArgon2,
	}

	defaultCrypter = crypter{
		schemeArgon2i:  &defaultArgon2,
		schemeArgon2id: &defaultArgon2,
	}
)

// Argon2 returns the default Crypter implementation for Argon2i and Argon2id
// functions.
func Argon2() Crypter {
	return &defaultArgon2Crypter
}

// Default returns the default Crypter implementation.
func Default() Crypter {
	return &defaultCrypter
}
