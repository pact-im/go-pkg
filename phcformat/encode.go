package phcformat

import (
	"encoding/base64"
)

// Encode encodes the given parameters in a PHC string format. It returns false
// if one of the parameters contains invalid characters or output is set without
// salt. If version is not set and there is only one parameter where the name is
// "v" and the value is a sequence of digits, it returns false.
func Encode(hashID string, version OptionalString, salt HashSalt, output []byte, params ...HashParam) (Hash, bool) {
	if len(hashID) > 32 {
		return Hash{}, false
	}
	for i := 0; i < len(hashID); i++ {
		if validID(hashID[i]) {
			continue
		}
		return Hash{}, false
	}
	if version.IsSet {
		for i := 0; i < len(version.Value); i++ {
			if validVersion(version.Value[i]) {
				continue
			}
			return Hash{}, false
		}
	}
	for _, p := range params {
		for i := 0; i < len(p.Name); i++ {
			if validParamName(p.Name[i]) {
				continue
			}
			return Hash{}, false
		}
		for i := 0; i < len(p.Value); i++ {
			if validParamValue(p.Value[i]) {
				continue
			}
			return Hash{}, false
		}
	}
	if !version.IsSet && len(params) == 1 && params[0].Name == "v" {
		value := params[0].Value
		valid := true
		for i := 0; i < len(value); i++ {
			if validVersion(value[i]) {
				continue
			}
			valid = false
			break
		}
		if valid {
			return Hash{}, false
		}
	}
	if salt.IsSet() {
		if !salt.Valid() {
			return Hash{}, false
		}
	} else if output != nil {
		return Hash{}, false
	}

	const (
		prefix   = "$"
		prefixV  = "$v="
		paramSep = ","
		paramDef = "="
	)

	n := len(prefix) + len(hashID)
	if version.IsSet {
		n += len(prefixV) + len(version.Value)
	}
	if count := len(params); count > 0 {
		n += len(prefix)
		if count > 1 {
			n += len(paramSep) * (count - 1)
		}
		for _, p := range params {
			n += len(p.Name) + len(paramDef) + len(p.Value)
		}
	}
	if salt.IsSet() {
		n += len(prefix) + salt.EncodedLen()
	}
	if output != nil {
		n += len(prefix) + base64.RawStdEncoding.EncodedLen(len(output))
	}

	var paramsStart, paramsEnd int
	var saltStart, saltEnd int
	var outStart int

	i, buf := 0, make([]byte, n)
	i += copy(buf, prefix)
	i += copy(buf[i:], hashID)
	if version.IsSet {
		i += copy(buf[i:], prefixV)
		i += copy(buf[i:], version.Value)
	}
	for pos, p := range params {
		if pos == 0 {
			i += copy(buf[i:], prefix)
			paramsStart = i
		} else {
			i += copy(buf[i:], paramSep)
		}
		i += copy(buf[i:], p.Name)
		i += copy(buf[i:], paramDef)
		i += copy(buf[i:], p.Value)
	}
	if len(params) > 0 {
		paramsEnd = i
	}
	if salt.IsSet() {
		i += copy(buf[i:], prefix)
		saltStart = i
		salt.Encode(buf[i:])
		i += salt.EncodedLen()
		saltEnd = i
	}
	if output != nil {
		i += copy(buf[i:], prefix)
		outStart = i
		base64.RawStdEncoding.Encode(buf[i:], output)
	}

	raw := string(buf)
	if !version.IsSet {
		version.Value = ""
	}
	var hashParams, hashSalt, hashOutput OptionalString
	if paramsStart != 0 {
		hashParams = String(raw[paramsStart:paramsEnd])
	}
	if saltStart != 0 {
		hashSalt = String(raw[saltStart:saltEnd])
	}
	if outStart != 0 {
		hashOutput = String(raw[outStart:])
	}
	return Hash{
		ID:      hashID,
		Version: version,
		Params:  hashParams,
		Salt:    hashSalt,
		Output:  hashOutput,
		Raw:     raw,
	}, true
}
