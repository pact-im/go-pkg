// Package extraio implements extra I/O utilities.
package extraio

import "io"

var (
	_ io.Reader = (ByteReader)(0)
	_ io.Reader = (*CountReader)(nil)
	_ io.Reader = (*DiscardReader)(nil)
	_ io.Reader = (*HardLimitedReader)(nil)
	_ io.Reader = (*HashReader)(nil)
	_ io.Reader = (*PadReader)(nil)
	_ io.Reader = (ReaderFunc)(nil)
	_ io.Reader = (*TailReader)(nil)
	_ io.Reader = (*UnpadReader)(nil)
)
