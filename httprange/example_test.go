package httprange_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"time"

	"go.pact.im/x/httprange"
)

type handlerTransport struct {
	handler http.Handler
}

func (t *handlerTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	w := httptest.NewRecorder()
	t.handler.ServeHTTP(w, r)
	return w.Result(), nil
}

func newClient(handler http.HandlerFunc) *http.Client {
	return &http.Client{
		Transport: &handlerTransport{
			handler: handler,
		},
	}
}

func ExampleBytesResource() {
	client := newClient(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Etag", `"512aedef7775096f9e152526a30a0ce7"`)
		http.ServeContent(w, r,
			"shakespeare.txt",
			time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Add(-time.Second),
			bytes.NewReader([]byte("To be, or not to be, that is the question.")),
		)
	})

	ctx := context.Background()

	resource, err := httprange.BuildFromRawURL(ctx, "https://example.com", client)
	if err != nil {
		panic(err)
	}

	reader := resource.Reader(ctx)
	defer func() {
		_ = reader.Close()
	}()

	buf1 := make([]byte, 5)
	n, err := reader.ReadAt(buf1, 0)
	if err != nil {
		panic(err)
	}
	fmt.Println("Bytes 0-4:", string(buf1[:n]))

	buf2 := make([]byte, 12)
	n, err = reader.ReadAt(buf2, 7)
	if err != nil {
		panic(err)
	}
	fmt.Println("Bytes 7-18:", string(buf2[:n]))

	buf3 := make([]byte, 42)
	n, err = reader.ReadAt(buf3, 41)
	if err != io.EOF {
		panic(err)
	}
	fmt.Println("Bytes 41-41:", string(buf3[:n]))
	// Output:
	// Bytes 0-4: To be
	// Bytes 7-18: or not to be
	// Bytes 41-41: .
}
