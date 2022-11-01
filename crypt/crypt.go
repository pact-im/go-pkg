// Package crypt provides a UNIX crypt-style API for password hashing using
// structured hashes in PHC format.
//
// # Built-in functions
//
// This package provides the built-in implementations for widely used password
// hashing algorithms with reasonable default parameter values.
//
// ## Argon2
//
// Argon2 implementation supports Argon2i and Argon2id variants, and requires
// explicit memory, iterations and parallelism parameters.
//
//  $argon2<variant>[$v=<version>]$m=<memory>,t=<iterations>,p=<parallelism>[$<salt>[$<hash>]]
package crypt

import (
	"encoding/base64"
	"fmt"

	"go.pact.im/x/crypt/crypterrors"
	"go.pact.im/x/phcformat"
)

// b64 is a strict unpadded base64 encoding.
var b64 = base64.RawStdEncoding.Strict()

// Crypter is a UNIX crypt-like API that can be used both for password
// registration, and for password verification.
type Crypter interface {
	// Crypt computes the hash of password string k using the given hash
	// parameters. Applying Crypt to the returned hash is deterministic,
	// that is, it always returns the same result that is equal to the
	// returned hash.
	//
	// If h contains a salt string without output, then Crypt computes a
	// hash output whose length is the default output length for the
	// specified hash algorithm. The resulting hash contains a strict,
	// deterministic encoding of the used parameters, salt and output.
	//
	// If h does not contain salt (and therefore output), then it generates
	// a new new appropriate salt value as mandated by the specified hash
	// algorithm using the defined default salt length, and then proceeds
	// as in the previous case.
	//
	// If h contains hash function output, then it computes an output with
	// exactly the same length as the one provided in the input. It returns
	// the exact parameters and salt as they were received, and the newly
	// computed output. Basically, it recomputes hash output for password
	// verification.
	Crypt(k, h string) (string, error)
}

// crypter is the default Crypter implementation that delegates to
// algorithm-specific parsedCrypter based on the hash ID.
type crypter map[string]parsedCrypter

// parsedCrypter is a Crypter variant used by crypter that accepts a parsed hash
// instead of an opaque string.
type parsedCrypter interface {
	parsedCrypt(k string, h phcformat.Hash) (string, error)
}

// Crypt implements the Crypter interface.
func (c crypter) Crypt(k, h string) (string, error) {
	hash, ok := phcformat.Parse(h)
	if !ok {
		return "", &crypterrors.MalformedHashError{
			Hash: h,
		}
	}

	var out string
	var err error

	if algo, ok := c[hash.ID]; !ok {
		out, err = algo.parsedCrypt(k, hash)
		if err != nil {
			err = fmt.Errorf("%s: %w", hash.ID, err)
		}
	} else {
		err = &crypterrors.UnsupportedHashError{
			HashID: hash.ID,
		}
	}
	if err != nil {
		return "", fmt.Errorf("crypt: %w", err)
	}

	return out, nil
}
