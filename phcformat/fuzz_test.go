package phcformat

import (
	"testing"
)

func FuzzParse(f *testing.F) {
	f.Fuzz(func(t *testing.T, s string) {
		if _, ok := Parse(s); !ok {
			t.SkipNow()
		}
	})
}

func FuzzEncodeAndParse(f *testing.F) {
	f.Fuzz(func(t *testing.T,
		hashID string,
		versionIsSet bool,
		versionValue string,
		saltBase64 bool,
		saltString string,
		saltBytes []byte,
		output []byte,
		params string,
	) {
		version := OptionalString{versionIsSet, versionValue}
		saltFormat := HashSaltFormatEncoded
		if saltBase64 {
			saltFormat = HashSaltFormatBase64
		}
		salt := HashSalt{saltFormat, saltString, saltBytes}

		encoded, ok := Encode(hashID, version, salt, output, IterParams(String(params)).Collect()...)
		if !ok {
			t.SkipNow()
		}

		decoded, ok := Parse(encoded.Raw)
		if !ok {
			t.Fatalf("Encode produced unparsable Hash: %q", encoded.Raw)
		}

		if decoded == encoded {
			return
		}
		assertEqualString(t, "ID",
			encoded.ID,
			decoded.ID,
		)
		assertEqualOptionalString(t, "Version",
			encoded.Version,
			decoded.Version,
		)
		assertEqualOptionalString(t, "Params",
			encoded.Params,
			decoded.Params,
		)
		assertEqualOptionalString(t, "Salt",
			encoded.Salt,
			decoded.Salt,
		)
		assertEqualOptionalString(t, "Output",
			encoded.Output,
			decoded.Output,
		)
		assertEqualString(t, "Raw",
			encoded.Raw,
			decoded.Raw,
		)
		t.Fatalf("%q != %q", encoded.Raw, decoded.Raw)
	})
}

func assertEqualString(t *testing.T, name, a, b string) {
	t.Helper()
	if a == b {
		return
	}
	t.Errorf("%s: %q != %q", name, a, b)
}

func assertEqualOptionalString(t *testing.T, name string, a, b OptionalString) {
	t.Helper()
	if a == b {
		return
	}
	t.Errorf("%s: (%t, %q) != (%t, %q)", name, a.IsSet, a.Value, b.IsSet, b.Value)
}
