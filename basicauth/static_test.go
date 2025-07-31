package basicauth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"gotest.tools/v3/assert"
)

func TestStatic(t *testing.T) {
	sentinel := "cheeki breeki"
	user := "user"
	pass := "pass"

	handler := Static(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		n, err := io.WriteString(w, sentinel)
		assert.NilError(t, err)
		assert.Equal(t, len(sentinel), n)
	}), user, pass)

	testCases := []struct {
		Name   string
		Code   int
		Body   string
		User   string
		Pass   string
		NoAuth bool
	}{
		{
			Name: "AcceptsValidCreds",
			Code: http.StatusOK,
			Body: sentinel,
			User: user,
			Pass: pass,
		},
		{
			Name: "RejectsWithoutCreds",
			Code: http.StatusUnauthorized,
		},
		{
			Name: "RejectsInvalidCreds",
			Code: http.StatusUnauthorized,
		},
		{
			Name: "RejectsInvalidUser",
			Code: http.StatusUnauthorized,
			Pass: pass,
		},
		{
			Name: "RejectsInvalidPass",
			Code: http.StatusUnauthorized,
			User: user,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if !tc.NoAuth {
				r.SetBasicAuth(tc.User, tc.Pass)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			assert.NilError(t, err)

			assert.Equal(t, tc.Code, resp.StatusCode)
			assert.Equal(t, tc.Body, string(body))
		})
	}
}
