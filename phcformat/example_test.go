package phcformat_test

import (
	"fmt"
	"strconv"

	"go.pact.im/x/phcformat"
)

func ExampleParamsIterator_Collect() {
	params := phcformat.IterParams(phcformat.String("a=b,c=d")).Collect()
	for _, p := range params {
		fmt.Println(p.Name)
		fmt.Println(p.Value)
	}
	// Output:
	// a
	// b
	// c
	// d
}

func ExampleIterParams() {
	params := phcformat.String("a=b,c=d")
	for it := phcformat.IterParams(params); it.Valid; it = it.Next() {
		p := it.Param
		fmt.Println(p.Name)
		fmt.Println(p.Value)
	}
	// Output:
	// a
	// b
	// c
	// d
}

func ExampleEncode_argon2id() {
	var (
		version = "19" // strconv.Itoa(argon2.Version)
		salt    = []byte("\x81\x98\x95\xFC\xCD`=\xCD\xB6\x12P\a\xFC\x98u\x1F")
		time    = uint32(2)
		memory  = uint32(64 * 1024)
		threads = uint8(1)
		output  = []byte{
			0xDA, 0x9F, 0xFE, 0x1E, 0x01, 0x48, 0xD4, 0xE1,
			0x07, 0x32, 0xFD, 0xAA, 0x88, 0xC3, 0x7A, 0x5C,
			0xC8, 0xC0, 0xC4, 0x23, 0xC7, 0xED, 0xA5, 0xC9,
			0x09, 0x78, 0x21, 0xE7, 0xD9, 0x7C, 0x0D, 0xBD,
		} // argon2.IDKey([]byte("pass"), salt, time, memory, threads, 32)
	)
	fmt.Println(phcformat.Encode("argon2id",
		phcformat.String(version),
		phcformat.HashSalt{Format: phcformat.HashSaltFormatBase64, Bytes: salt},
		output,
		phcformat.HashParam{"m", strconv.Itoa(int(memory))},
		phcformat.HashParam{"t", strconv.Itoa(int(time))},
		phcformat.HashParam{"p", strconv.Itoa(int(threads))},
	))
	// Output:
	// $argon2id$v=19$m=65536,t=2,p=1$gZiV/M1gPc22ElAH/Jh1Hw$2p/+HgFI1OEHMv2qiMN6XMjAxCPH7aXJCXgh59l8Db0 true
}
