package phcformat_test

import (
	"fmt"

	"go.pact.im/x/option"
	"go.pact.im/x/phcformat"
	"go.pact.im/x/phcformat/encode"
)

func ExampleIterParams() {
	params := option.Value("a=b,c=d")
	it := phcformat.IterParams(option.UnwrapOrZero(params))
	for ; it.Valid; it = it.Next() {
		fmt.Println(it.Name, it.Value)
	}
	if it.After != "" {
		panic("parse error")
	}
	// Output:
	// a b
	// c d
}

func ExampleMustParse() {
	h := phcformat.MustParse("$name$v=42$k=v$salt$hash")
	fmt.Println(h)
	fmt.Println(h.ID)
	fmt.Println(option.UnwrapOrZero(h.Version))
	fmt.Println(option.UnwrapOrZero(h.Params))
	fmt.Println(option.UnwrapOrZero(h.Salt))
	fmt.Println(option.UnwrapOrZero(h.Output))
	// Output:
	// $name$v=42$k=v$salt$hash
	// name
	// 42
	// k=v
	// salt
	// hash
}

func ExampleAppend_argon2id() {
	var (
		// version is the version of the Argon2 algorithm.
		//
		// Example:
		//  version := argon2.Version
		version = 19
		// salt is the random salt bytes.
		//
		// Example:
		//  import "crypto/rand"
		//  salt := make([]byte, 16)
		//  if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		//      panic(err)
		//  }
		salt = []byte{
			0x81, 0x98, 0x95, 0xFC, 0xCD, 0x60, 0x3D, 0xCD,
			0xB6, 0x12, 0x50, 0x07, 0xFC, 0x98, 0x75, 0x1F,
		}
		// time is the time constraint parameter.
		time = uint32(2)
		// memory is the time constraint parameter.
		memory = uint32(64 * 1024)
		// threads is the parallelism parameter.
		threads = uint8(1)
		// output is the hash function output.
		//
		// Example:
		//  output := argon2.IDKey([]byte("pass"), salt, time, memory, threads, 32)
		output = []byte{
			0xDA, 0x9F, 0xFE, 0x1E, 0x01, 0x48, 0xD4, 0xE1,
			0x07, 0x32, 0xFD, 0xAA, 0x88, 0xC3, 0x7A, 0x5C,
			0xC8, 0xC0, 0xC4, 0x23, 0xC7, 0xED, 0xA5, 0xC9,
			0x09, 0x78, 0x21, 0xE7, 0xD9, 0x7C, 0x0D, 0xBD,
		}
	)
	newParam := func(k string, v uint) encode.KV[encode.String, encode.Byte, encode.Uint] {
		return encode.NewKV(encode.NewByte('='), encode.NewString(k), encode.NewUint(v))
	}
	fmt.Println(string(phcformat.Append(nil,
		encode.String("argon2id"),
		option.Value(encode.NewUint(uint(version))),
		option.Value(encode.NewList(
			encode.NewByte(','),
			newParam("m", uint(memory)),
			newParam("t", uint(time)),
			newParam("p", uint(threads)),
		)),
		option.Value(encode.NewBase64(salt)),
		option.Value(encode.NewBase64(output)),
	)))
	// Output:
	// $argon2id$v=19$m=65536,t=2,p=1$gZiV/M1gPc22ElAH/Jh1Hw$2p/+HgFI1OEHMv2qiMN6XMjAxCPH7aXJCXgh59l8Db0
}
