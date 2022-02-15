package flaky_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"

	"go.pact.im/x/flaky"
)

func ExampleRetry() {
	ctx := context.Background()

	executor := flaky.Retry(flaky.Limit(5, flaky.Exp2(time.Second)))

	req, _ := http.NewRequest(http.MethodGet, "https://example.com", nil)
	client := http.DefaultClient

	var dump []byte
	var err error
	err = executor.Execute(ctx, func(ctx context.Context) error {
		ctx, cancel := context.WithTimeout(ctx, time.Second)
		defer cancel()

		var resp *http.Response
		resp, err = client.Do(req.WithContext(ctx))
		if err != nil {
			return err
		}

		dump, err = httputil.DumpResponse(resp, false)
		return err
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(dump))
}
