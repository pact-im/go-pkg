package phcformat

import (
	"strings"
)

// Parse parses PHC formatted string s. It returns the parsed hash and a boolean
// value indicating either success or parse error.
func Parse(s string) (Hash, bool) {
	if s == "" {
		return Hash{}, false
	}
	if s[0] != '$' {
		return Hash{}, false
	}
	raw, s := s, s[1:]

	var sep bool

	var hashID string
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c == '$' {
			hashID, s = s[:i], s[i+1:]
			sep = true
			break
		}
		if i >= 32 {
			return Hash{}, false
		}
		if validID(c) {
			continue
		}
		return Hash{}, false
	}
	if !sep {
		hashID = s
		return Hash{
			ID:  hashID,
			Raw: raw,
		}, true
	}
	sep = false

	var version, params, salt OptionalString
	var maybeSalt bool

	var cur int
	if prefix := "v="; strings.HasPrefix(s, prefix) {
		n := len(prefix)
		for cur = n; cur < len(s); cur++ {
			c := s[cur]
			if c == '$' {
				version, s = String(s[n:cur]), s[cur+1:]
				sep = true
				cur = 0
				break
			}
			if validVersion(c) {
				continue
			}
			if c == ',' {
				cur++
				goto paramName
			}
			if validParamValue(c) {
				cur++
				goto paramValue
			}
			return Hash{}, false
		}
		if !sep {
			version = String(s[n:])
			return Hash{
				ID:      hashID,
				Version: version,
				Raw:     raw,
			}, true
		}
		sep = false
	}

	maybeSalt = true

paramName:
	for ; cur < len(s); cur++ {
		c := s[cur]
		if validParamName(c) {
			continue
		}
		if c == '=' {
			cur++
			maybeSalt = false
			goto paramValue
		}
		if maybeSalt {
			if c == '$' {
				salt, s = String(s[:cur]), s[cur+1:]
				goto hash
			}
			if validSalt(c) {
				cur++
				goto salt
			}
		}
		return Hash{}, false
	}
	// If we did not find the equals sign in the string and it does not
	// contain commas, it is a salt.
	//
	// Note that we jump to paramName only on comma character after the
	// parameter value (either from version or paramValue states), that is,
	// when maybeSalt is false.
	if maybeSalt {
		salt = String(s)
		return Hash{
			ID:      hashID,
			Version: version,
			Salt:    salt,
			Raw:     raw,
		}, true
	}
	return Hash{}, false
paramValue:
	for ; cur < len(s); cur++ {
		c := s[cur]
		if validParamValue(c) {
			continue
		}
		if c == ',' {
			cur++
			goto paramName
		}
		if c == '$' {
			params, s = String(s[:cur]), s[cur+1:]
			sep = true
			cur = 0
			break
		}
		return Hash{}, false
	}
	if !sep {
		params = String(s)
		return Hash{
			ID:      hashID,
			Version: version,
			Params:  params,
			Raw:     raw,
		}, true
	}
	sep = false

salt:
	for ; cur < len(s); cur++ {
		c := s[cur]
		if validSalt(c) {
			continue
		}
		if c == '$' {
			salt, s = String(s[:cur]), s[cur+1:]
			sep = true
			break
		}
		return Hash{}, false
	}
	if !sep {
		salt = String(s)
		return Hash{
			ID:      hashID,
			Version: version,
			Params:  params,
			Salt:    salt,
			Raw:     raw,
		}, true
	}
	// sep is unused in hash

hash:
	for i := 0; i < len(s); i++ {
		c := s[i]
		if validOutput(c) {
			continue
		}
		return Hash{}, false
	}
	return Hash{
		ID:      hashID,
		Version: version,
		Params:  params,
		Salt:    salt,
		Output:  String(s),
		Raw:     raw,
	}, true
}
