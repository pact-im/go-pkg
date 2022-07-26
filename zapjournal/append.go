package zapjournal

import (
	"encoding/base64"
	"math"
	"strconv"
)

// strconvAppendComplex is an append-style function for strconv.FormatComplex
// that is missing in the standard library.
func strconvAppendComplex(dst []byte, c complex128, fmt byte, prec, bitSize int) []byte {
	if bitSize != 64 && bitSize != 128 {
		panic("invalid bitSize")
	}
	bitSize >>= 1

	dst = append(dst, '(')
	dst = strconv.AppendFloat(dst, real(c), fmt, prec, bitSize)

	im := imag(c)
	if math.IsNaN(im) || (!math.IsInf(im, 0) && im > 0) {
		dst = append(dst, '+')
	}

	dst = strconv.AppendFloat(dst, im, fmt, prec, bitSize)
	dst = append(dst, 'i', ')')
	return dst
}

// strconvAppendUintptr is an append-style convenience function for uintptr
// values and strconv.AppendUint.
func strconvAppendUintptr(dst []byte, v uintptr) []byte {
	dst = append(dst, []byte("0x")...)
	dst = strconv.AppendUint(dst, uint64(v), 16)
	return dst
}

// base64Append is an append-style function for base64 encoding.
//
// See also https://go.dev/issue/19366
func base64Append(e *base64.Encoding, dst, src []byte) []byte {
	k := len(dst)
	n := e.EncodedLen(len(src))
	if cap(dst)-k < n {
		dst = append(dst, make([]byte, n)...)
	} else {
		dst = dst[:k+n]
	}
	e.Encode(dst[k:], src)
	return dst
}
