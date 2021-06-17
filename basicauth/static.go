package basicauth

import (
	"crypto/subtle"
	"net/http"

	"golang.org/x/crypto/sha3"
)

const hashSize = 64

var _ http.Handler = (*StaticHandler)(nil)

type StaticHandler struct {
	// Handler is the handler protected by HTTP basic authentication.
	Handler http.Handler
	// UserHash is username hash to compare against.
	UserHash []byte
	// PassHash is password hash to compare against.
	PassHash []byte
}

// Static returns a handler that serves requests from underlying handler after
// successful HTTP basic authentication with static credentials.
func Static(handler http.Handler, user, pass string) *StaticHandler {
	h := &StaticHandler{
		Handler:  handler,
		UserHash: make([]byte, hashSize),
		PassHash: make([]byte, hashSize),
	}
	sha3.ShakeSum256(h.UserHash, []byte(user))
	sha3.ShakeSum256(h.PassHash, []byte(pass))
	return h
}

func (h *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	u, p, ok := r.BasicAuth()

	// Use constant time comparison to guard against brute-force timing
	// attacks. Also short-circuit on ok in case there is no auth header.
	//
	var cmp int
	if ok {
		var userHash, passHash [hashSize]byte

		sha3.ShakeSum256(userHash[:], []byte(u))
		sha3.ShakeSum256(passHash[:], []byte(p))

		cmp += subtle.ConstantTimeCompare(h.UserHash, userHash[:])
		cmp += subtle.ConstantTimeCompare(h.PassHash, passHash[:])
	}
	if !ok || cmp != 2 {
		w.Header().Add("WWW-Authenticate", "Basic")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	h.Handler.ServeHTTP(w, r)
}
