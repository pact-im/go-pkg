package crypt

import (
	"fmt"
	"io"
	"math"
	"strconv"

	"golang.org/x/crypto/argon2"

	"go.pact.im/x/option"
	"go.pact.im/x/phcformat"
	"go.pact.im/x/phcformat/encode"

	"go.pact.im/x/crypt/crypterrors"
)

// crypterArgon2 is the parsedCrypter implementation for Argon2 functions.
type crypterArgon2 struct {
	rand io.Reader
}

// Crypt implements the parsedCrypter interface.
func (c *crypterArgon2) parsedCrypt(k string, h phcformat.Hash) (string, error) {
	if v, ok := h.Version.Unwrap(); ok {
		version, err := strconv.ParseUint(v, 10, 32)
		if err != nil {
			return "", fmt.Errorf("parse version: %w", err)
		}
		if version != argon2.Version {
			return "", &crypterrors.UnsupportedVersionError{
				Parsed:    uint(version),
				Suggested: argon2.Version,
			}
		}
	}

	var time, memory uint32
	var threads uint8

	it := phcformat.IterParams(option.UnwrapOrZero(h.Params))
	for ; it.Valid; it = it.Next() {
		desc := "[1;math.MaxUint32]"

		var n uint64
		var unimplemented bool
		var err error
		switch it.Name {
		case "m":
			n, err = strconv.ParseUint(it.Value, 10, 32)
			memory = uint32(n)
		case "t":
			n, err = strconv.ParseUint(it.Value, 10, 32)
			time = uint32(n)
		case "p":
			n, err = strconv.ParseUint(it.Value, 10, 8)
			threads = uint8(n)
			desc = "[1;math.MaxUint8]"
		case "keyid", "data":
			unimplemented = true
			fallthrough
		default:
			return "", &crypterrors.UnsupportedParameterError{
				Name:          it.Name,
				Unimplemented: unimplemented,
			}
		}
		if err != nil {
			return "", fmt.Errorf("parse parameter %q: %w", it.Name, err)
		}
		if n < 1 {
			return "", &crypterrors.InvalidParameterValueError{
				Name:     it.Name,
				Value:    it.Value,
				Expected: desc,
			}
		}
	}
	if it.After != "" {
		return "", &crypterrors.MalformedParametersError{
			Unparsed: it.After,
		}
	}
	if time == 0 || memory == 0 || threads == 0 {
		return "", &crypterrors.MissingRequiredParametersError{
			Required: "m, t, p",
		}
	}

	var salt []byte
	if v, ok := h.Salt.Unwrap(); ok {
		var err error
		salt, err = b64.DecodeString(v)
		if err != nil {
			return "", fmt.Errorf("decode salt: %w", err)
		}
	} else {
		salt = make([]byte, 32)
		_, err := io.ReadFull(c.rand, salt)
		if err != nil {
			return "", fmt.Errorf("generate salt: %w", err)
		}
	}

	keyLen := uint32(32)
	if v, ok := h.Output.Unwrap(); ok {
		n := b64.DecodedLen(len(v))
		if n <= 0 || n > math.MaxUint32 {
			return "", &crypterrors.InvalidOutputLengthError{
				Length:   n,
				Expected: "non-zero unsigned 32-bit integer",
			}
		}
		keyLen = uint32(n)
	}

	var output []byte
	if h.ID == schemeArgon2id {
		output = argon2.IDKey([]byte(k), salt, time, memory, threads, keyLen)
	} else {
		output = argon2.Key([]byte(k), salt, time, memory, threads, keyLen)
	}

	rawLen := len(h.Raw)
	if option.IsNil(h.Salt) {
		rawLen += 1 + b64.EncodedLen(len(salt))
	}
	if v, ok := h.Output.Unwrap(); ok {
		rawLen -= 1 + len(v)
	}
	rawLen += 1 + b64.EncodedLen(len(output))

	return string(phcformat.Append(make([]byte, 0, rawLen),
		encode.NewString(h.ID),
		option.Map(h.Version, encode.NewString),
		option.Map(h.Params, encode.NewString),
		option.Value(encode.NewBase64(salt)),
		option.Value(encode.NewBase64(output)),
	)), nil
}
