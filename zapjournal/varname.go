package zapjournal

import (
	"strings"
)

// appendVarName appends sanitized variable names to the dst and returns the
// resulting byte slice.
func appendVarName(dst []byte, parts ...string) []byte {
	var appended bool
	for _, s := range parts {
		s = sanitizeVarName(s)
		if s == "" {
			continue
		}
		if appended {
			dst = append(dst, '_')
		}
		dst = appendUpperSnakeVarName(dst, s)
		appended = true
	}
	if !appended {
		dst = append(dst, defaultPrefix...)
	}
	return dst
}

// sanitizeVarName sanitizes the variable name to be a valid input for
// appendUpperSnakeVarName. It also trims characters that would be converted to
// underscore.
func sanitizeVarName(s string) string {
	// Trim invalid bytes and underscore. Check whether the trimmed string
	// has invalid characters. If so, we need to allocate a copy with these
	// characters removed.
	ltrim := 0
	rtrim := len(s)
	ngrow := 0
	for i := 0; i < len(s); i++ {
		c := s[i]
		if isAllowedVarNameByte(c) {
			break
		}
		ltrim++
	}
	for i := 0; i < len(s); i++ {
		c := s[len(s)-1-i]
		if isAllowedVarNameByte(c) {
			break
		}
		rtrim--
	}
	for i := ltrim; i < rtrim; i++ {
		c := s[i]
		if shouldConvertToUnderscore(c) || isAllowedVarNameByte(c) {
			ngrow++
		}
	}
	s = s[ltrim:rtrim]
	if ngrow == len(s) || ngrow == 0 {
		return s
	}

	var b strings.Builder
	b.Grow(ngrow)
	for i := 0; i < len(s); i++ {
		c := s[i]
		if !shouldConvertToUnderscore(c) && !isAllowedVarNameByte(c) {
			continue
		}
		_ = b.WriteByte(c)
	}
	return b.String()
}

// isAllowedVarNameByte checks whether c is allowed variable name character
// before converting to upper case snake.
func isAllowedVarNameByte(c byte) bool {
	switch {
	case isDigitASCII(c):
	case isUpperASCII(c):
	case isLowerASCII(c):
	default:
		return false
	}
	return true
}

// shouldConvertToUnderscore checks whether c should be converted to underscore.
func shouldConvertToUnderscore(c byte) bool {
	switch c {
	case '_': // underscore
	case '@': // at sign
	case '$': // dollar
	case '%': // percent
	case '#': // hash
	case ':': // colon
	case '+': // plus
	case '-': // dash
	case '.': // dot
	case ' ': // space
	default:
		return false
	}
	return true
}

// appendUpperSnakeVarName appends s in upper snake case to dst.
func appendUpperSnakeVarName(dst []byte, s string) []byte {
	for i := 0; i < len(s); i++ {
		c := s[i]

		var underscore bool
		if i > 0 && !isDigitASCII(c) {
			underscore = isDigitASCII(s[i-1])
		}

		switch {
		case isDigitASCII(c):
			// No-op.
		case isLowerASCII(c):
			c = lowerToUpperASCII(c)
		case shouldConvertToUnderscore(c):
			c = '_'
			for i+1 < len(s) && shouldConvertToUnderscore(s[i+1]) {
				i++
			}
		case isUpperASCII(c):
			if i == 0 || underscore {
				break
			}
			prev := s[i-1]

			if isLowerASCII(prev) {
				underscore = true
				break
			}

			if !isUpperASCII(prev) {
				break
			}

			if i+1 == len(s) {
				break
			}
			next := s[i+1]

			if !isLowerASCII(next) {
				break
			}
			underscore = true
		default:
			continue // unreachable
		}
		if underscore {
			dst = append(dst, '_')
		}
		dst = append(dst, c)
	}
	return dst
}

// isLowerASCII checks whether c is a lower case ASCII letter.
func isLowerASCII(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// isUpperASCII checks whether c is an upper case ASCII letter.
func isUpperASCII(c byte) bool {
	return 'A' <= c && c <= 'Z'
}

// isDigitASCII checks whether c is an ASCII digit.
func isDigitASCII(c byte) bool {
	return '0' <= c && c <= '9'
}

// lowerToUpperASCII converts c to an upper case ASCII character, assuming that c is
// a lower case letter.
func lowerToUpperASCII(c byte) byte {
	const d = 'a' - 'A'
	return c - d
}
