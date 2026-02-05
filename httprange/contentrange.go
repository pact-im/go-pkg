package httprange

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/http/httpguts"
)

// contentRangeHeader represents a Content-Range header value.
// It consists of a unit and a range response (e.g. "bytes 0-99/100").
type contentRangeHeader struct {
	Unit string
	Resp string
}

func parseContentRange(s string) (contentRangeHeader, bool) {
	// Note: RFC 9110 has a regression in Content-Range grammar: it only
	// allows bytes-like range-resp. We use the definition from RFC 7233
	// that allows other units.
	//
	// https://datatracker.ietf.org/doc/html/rfc9110
	// https://datatracker.ietf.org/doc/html/rfc7233

	before, after, ok := strings.Cut(s, " ")
	if !ok || before == "" {
		return contentRangeHeader{}, false
	}

	for _, r := range before {
		if httpguts.IsTokenRune(r) {
			continue
		}
		return contentRangeHeader{}, false
	}

	for _, r := range after {
		if r > 0 && r < 0x80 {
			continue
		}
		return contentRangeHeader{}, false
	}

	return contentRangeHeader{
		Unit: before,
		Resp: after,
	}, true
}

func parseContentRangeFromHeader(h http.Header, unit string) (string, error) {
	value := h.Get(httpHeaderContentRange)
	r, ok := parseContentRange(value)
	if !ok {
		return "", fmt.Errorf(
			"httprange: invalid Content-Range header value %q",
			value,
		)
	}
	if !equalFoldASCII(r.Unit, unit) {
		return "", fmt.Errorf(
			"httprange: unexpected range unit %q (expected %q)",
			r.Unit, unit,
		)
	}
	return r.Resp, nil
}

type unsatisfiedRangeResp struct {
	CompleteLength int64
}

func parseUnsatisfiedRangeResp(str string) (unsatisfiedRangeResp, bool) {
	completeLengthStr, ok := strings.CutPrefix(str, "*/")
	if !ok || !isDigits(completeLengthStr) {
		return unsatisfiedRangeResp{}, false
	}
	completeLength, err := strconv.ParseInt(completeLengthStr, 10, 64)
	if err != nil {
		return unsatisfiedRangeResp{}, false
	}
	return unsatisfiedRangeResp{
		CompleteLength: completeLength,
	}, true
}

type bytesRangeResp struct {
	First, Last    int64
	CompleteLength int64 // zero if unknown
}

func parseBytesRangeResp(str string) (bytesRangeResp, bool) {
	inclRange, completeLengthStr, ok := strings.Cut(str, "/")
	if !ok {
		return bytesRangeResp{}, false
	}

	hasLength := completeLengthStr != "*"
	if hasLength && !isDigits(completeLengthStr) {
		return bytesRangeResp{}, false
	}

	firstStr, lastStr, ok := strings.Cut(inclRange, "-")
	if !ok || !isDigits(firstStr) || !isDigits(lastStr) {
		return bytesRangeResp{}, false
	}

	first, err := strconv.ParseInt(firstStr, 10, 64)
	if err != nil {
		return bytesRangeResp{}, false
	}
	last, err := strconv.ParseInt(lastStr, 10, 64)
	if err != nil {
		return bytesRangeResp{}, false
	}

	if first > last {
		return bytesRangeResp{}, false
	}

	var completeLength int64
	if hasLength {
		completeLength, err = strconv.ParseInt(completeLengthStr, 10, 64)
		if err != nil {
			return bytesRangeResp{}, false
		}
		if last >= completeLength {
			return bytesRangeResp{}, false
		}
	}

	return bytesRangeResp{
		First:          first,
		Last:           last,
		CompleteLength: completeLength,
	}, true
}

func isDigits(s string) bool {
	if s == "" {
		return false
	}
	for i := range len(s) {
		if !isDigitASCII(s[i]) {
			return false
		}
	}
	return true
}

func isDigitASCII(b byte) bool {
	return '0' <= b && b <= '9'
}
