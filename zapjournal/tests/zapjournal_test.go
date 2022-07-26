//go:build linux
// +build linux

package tests

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"sort"
	"testing"
	"time"

	"github.com/valyala/fastjson"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"

	"go.pact.im/x/zapjournal"
)

func TestCore(t *testing.T) {
	ok, err := zapjournal.Available()
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Skip("journal is not available")
	}
	pid := os.Getpid()

	conn, err := zapjournal.Bind()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	fakeTime := time.Date(2009, time.November, 10, 23, 0, 0, 42, time.UTC)

	log := zap.New(zapjournal.NewCore(conn),
		zap.ErrorOutput(failTestWriter{t}),
		zap.WithClock(&fakeClock{fakeTime}),
	)

	t.Run("Huge", func(t *testing.T) {
		f, err := conn.File()
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = f.Close() }()

		wmem, err := unix.GetsockoptInt(int(f.Fd()), unix.SOL_SOCKET, unix.SO_SNDBUF)
		if err != nil {
			t.Fatal(err)
		}

		// NB encoding overhead exceeds wmem so we expect the memfd path.
		huge := make([]byte, wmem)
		log.Info(t.Name(), zap.Binary("huge", huge))

		v, err := lastJournalMessage(pid, t.Name())
		if err != nil {
			t.Fatal(err)
		}

		field := v.Get("X_HUGE")
		if field == nil {
			t.Fatal("missing X_HUGE field")
		}
		xs, err := field.Array()
		if err != nil {
			t.Fatal(err)
		}
		buf, err := fastjsonArrayToBytes(xs)
		if err != nil {
			t.Fatal(err)
		}
		if !bytes.Equal(buf, huge) {
			t.Fatal("bytes are not equal")
		}
	})

	t.Run("Send", func(t *testing.T) {
		allBytes := make([]byte, 256)
		for i := 0; i < len(allBytes); i++ {
			allBytes[i] = byte(i)
		}

		log.Named("zapjournal").Info(t.Name(),
			zap.Reflect("reflect", map[string]string{"foo": "bar"}),
			zap.Strings("array", []string{"foo", "bar", "baz"}),
			zap.Binary("binary", allBytes),
			zap.Bool("bool", true),
			zap.ByteString("byteString", []byte("byte string")),
			zap.String("string", "string"),
			zap.String("multiline", "multi\nline"),
			zap.Complex128("complex128", complex(math.NaN(), math.Inf(1))),
			zap.Complex64("complex64", 5+7i),
			zap.Float64("float64", math.Inf(1)),
			zap.Float32("float32", math.MaxFloat32),
			zap.Duration("duration", time.Second),
			zap.Uint("uint", 42),
			zap.Int("int", -42),
			zap.Time("time", fakeTime),
			zap.Error(errors.New("something bad")),
		)

		v, err := lastJournalMessage(pid, t.Name())
		if err != nil {
			t.Fatal(err)
		}

		binaryFields := map[string][]byte{
			"X_BINARY": allBytes,
		}
		stringFields := map[string]string{
			"PRIORITY":      "6",
			"MESSAGE":       t.Name(),
			"LOG_LEVEL":     "info",
			"LOG_NAME":      "zapjournal",
			"TIMESTAMP":     "2009-11-10T23:00:00.000000042Z",
			"X_REFLECT":     `{"foo":"bar"}`,
			"X_ARRAY":       `["foo","bar","baz"]`,
			"X_BOOL":        "true",
			"X_BYTE_STRING": "byte string",
			"X_STRING":      "string",
			"X_MULTILINE":   "multi\nline",
			"X_COMPLEX128":  "(NaN+Infi)",
			"X_COMPLEX64":   "(5+7i)",
			"X_FLOAT64":     "+Inf",
			"X_FLOAT32":     "3.4028235e+38",
			"X_DURATION":    "1s",
			"X_UINT":        "42",
			"X_INT":         "-42",
			"X_TIME":        "2009-11-10T23:00:00.000000042Z",
			"X_ERROR":       "something bad",
		}
		v.Visit(func(key []byte, v *fastjson.Value) {
			k := string(key)
			if s, ok := binaryFields[k]; ok {
				expectBytes(t, k, v, s)
				delete(binaryFields, k)
			}
			if s, ok := stringFields[k]; ok {
				expectString(t, k, v, s)
				delete(stringFields, k)
			}
		})
		if n := len(binaryFields) + len(stringFields); n != 0 {
			missing := make([]string, 0, n)
			for k := range binaryFields {
				missing = append(missing, k)
			}
			for k := range stringFields {
				missing = append(missing, k)
			}
			sort.Strings(missing)
			t.Fatalf("missing fields: %q", missing)
		}
	})
}

func expectString(t *testing.T, k string, v *fastjson.Value, expected string) {
	s, err := v.StringBytes()
	if err != nil {
		t.Fatalf("%s: %v", k, err)
	}
	got := string(s)
	if got != expected {
		t.Fatalf("%s: expected %q, got %q", k, expected, got)
	}
}

func expectBytes(t *testing.T, k string, v *fastjson.Value, expected []byte) {
	xs, err := v.Array()
	if err != nil {
		t.Fatalf("%s: %v", k, err)
	}
	got, err := fastjsonArrayToBytes(xs)
	if err != nil {
		t.Fatalf("%s: %v", k, err)
	}
	if !bytes.Equal(got, expected) {
		t.Fatalf("%s: expected %q, got %q", k, expected, got)
	}
}

func fastjsonArrayToBytes(xs []*fastjson.Value) ([]byte, error) {
	buf := make([]byte, len(xs))
	for i, x := range xs {
		n, err := x.Int()
		if err == nil && (n < 0 || n > 255) {
			err = fmt.Errorf("expected byte-sized integer, got %d", n)
		}
		if err != nil {
			return nil, fmt.Errorf("value at index %d: %w", i, err)
		}
		buf[i] = byte(n)
	}
	return buf, nil
}
