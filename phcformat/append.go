package phcformat

import (
	"go.pact.im/x/option"
	"go.pact.im/x/phcformat/encode"
)

// Append appends the given parameters in a PHC string format to dst and returns
// the resulting slice.
//
// The caller is responsible for ensuring that:
//  • name is a sequence of characters in “a-z0-9-”.
//  • version is a sequence of characters in “0-9”.
//  • params is a sequence of comma-separated name and value pairs (separated by
//    equals sign) where name is a sequence of characters in “a-z0-9-” and value
//    is a sequence of characters in “a-zA-Z0-9/+.-”. If version is not set and
//    only a single parameter named “v” is given, to avoid ambiguity, its value
//    must not be a sequence of characters in “0-9” (as in version).
//  • salt is a sequence of characters in “a-zA-Z0-9/+.-”.
//  • output is a sequence of characters in “A-Za-z0-9+/” (base64 character set)
//    and salt is set. That is, it must be a base64-encoded string and implies
//    that salt is set.
func Append[NameAppender, VersionAppender, ParamsAppender, SaltAppender, OutputAppender encode.Appender](
	dst []byte,
	name NameAppender,
	version option.Of[VersionAppender],
	params option.Of[ParamsAppender],
	salt option.Of[SaltAppender],
	output option.Of[OutputAppender],
) []byte {
	prefix := encode.NewString("$")
	prefixV := encode.NewString("$v=")

	prefixedName := encode.NewConcat(prefix, name)
	prefixedVersion := optionAppenderWithPrefix(prefixV, version)
	prefixedParams := optionAppenderWithPrefix(prefix, params)
	prefixedSalt := optionAppenderWithPrefix(prefix, salt)
	prefixedOutput := optionAppenderWithPrefix(prefix, output)

	dst = prefixedName.Append(dst)
	dst = prefixedVersion.Append(dst)
	dst = prefixedParams.Append(dst)
	dst = prefixedSalt.Append(dst)
	dst = prefixedOutput.Append(dst)

	return dst
}

// optionAppenderWithPrefix returns an encode.Appender that, if set, is prefixed
// with the given prefix.
func optionAppenderWithPrefix[PrefixAppender, OptionAppender encode.Appender](
	prefix PrefixAppender,
	opt option.Of[OptionAppender],
) encode.Option[encode.Concat[PrefixAppender, OptionAppender]] {
	if v, ok := opt.Unwrap(); ok {
		return encode.NewOption(option.Value(encode.NewConcat(prefix, v)))
	}
	return encode.NewOption(option.Nil[encode.Concat[PrefixAppender, OptionAppender]]())
}
