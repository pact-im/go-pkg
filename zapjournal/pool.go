package zapjournal

import (
	"sync"
)

// pooledBufferCapacity is the capacity required for the encoder to be pooled.
// This is necessary since the operation of sync.Pool assumes that the memory
// cost of each element is approximately the same in order to be efficient.
// See https://go.dev/issue/23199
const pooledBufferCapacity = 8192

// varsEncoderPool is a pool of varsEncoder objects with fixed buffer capacity.
var varsEncoderPool = sync.Pool{
	New: func() any {
		e := &varsEncoder{
			buf: make([]byte, 0, pooledBufferCapacity),
		}
		e.json.bufp = &e.buf
		return e
	},
}

// getVarsEncoder returns a new varsEncoder instance from the pool with the
// given variable prefix set.
func getVarsEncoder(prefix string) *varsEncoder {
	e := varsEncoderPool.Get().(*varsEncoder)
	e.prefix = prefix
	return e
}

// putVarsEncoder puts the varsEncoder instance back to the pool if possible.
func putVarsEncoder(e *varsEncoder) {
	if cap(e.buf) != pooledBufferCapacity {
		return
	}
	e.prefix = ""
	e.buf = e.buf[:0]
	varsEncoderPool.Put(e)
}

// cloneVarsEncoder clones the given varsEncoder, using the pooled objects when
// possible. It must not be called on varsEncoders between beginVar and endVar
// calls.
func cloneVarsEncoder(e *varsEncoder) *varsEncoder {
	if cap(e.buf) != pooledBufferCapacity {
		buf := make([]byte, len(e.buf), cap(e.buf))
		_ = copy(buf, e.buf)
		// Note that we do not have to copy hdr and json fields since
		// they are reset to zero values after use in each method.
		return &varsEncoder{
			prefix: e.prefix,
			buf:    buf,
		}
	}

	clone := getVarsEncoder(e.prefix)
	clone.buf = clone.buf[:len(e.buf)]
	_ = copy(clone.buf, e.buf)
	return clone
}
