package basicauth

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStatic(t *testing.T) {
	sentinel := "cheeki breeki"
	user := "user"
	pass := "pass"

	handler := Static(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n, err := io.WriteString(w, sentinel)
		assert.NoError(t, err)
		assert.Equal(t, len(sentinel), n)
	}), user, pass)

	testCases := []struct {
		Name string
		Code int
		Body string
		Auth bool
		User string
		Pass string
	}{
		{
			Name: "AcceptsValidCreds",
			Code: http.StatusOK,
			Body: sentinel,
			Auth: true,
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
			Auth: true,
		},
		{
			Name: "RejectsInvalidUser",
			Code: http.StatusUnauthorized,
			Auth: true,
			Pass: pass,
		},
		{
			Name: "RejectsInvalidPass",
			Code: http.StatusUnauthorized,
			Auth: true,
			User: user,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.Name, func(t *testing.T) {
			r, err := http.NewRequest("GET", "/", nil)
			require.NoError(t, err)

			if tc.Auth {
				r.SetBasicAuth(tc.User, tc.Pass)
			}

			w := httptest.NewRecorder()
			handler.ServeHTTP(w, r)

			resp := w.Result()
			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err)

			assert.Equal(t, tc.Code, resp.StatusCode)
			assert.Equal(t, tc.Body, string(body))
		})
	}
}
