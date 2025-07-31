// Package httpdebug provides HTTP handler for debug endpoints.
package httpdebug

import (
	"expvar"
	"io"
	"net/http"
	"net/http/pprof"
	"runtime/debug"
	"strconv"

	"golang.org/x/net/trace"
)

// New returns a new HTTP handler for debug endpoints.
func New() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	mux.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	mux.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	mux.Handle("/debug/pprof/block", pprof.Handler("block"))

	mux.Handle("/debug/vars", expvar.Handler())

	mux.HandleFunc("/debug/events", trace.Events)
	mux.HandleFunc("/debug/requests", trace.Traces)

	mux.HandleFunc("/debug/buildinfo", func(w http.ResponseWriter, _ *http.Request) {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			w.WriteHeader(http.StatusNotImplemented)
			return
		}

		modinfo := info.String()

		h := w.Header()
		h.Set("Content-Type", "text/plain; charset=utf-8")
		h.Set("Content-Length", strconv.Itoa(len(modinfo)))

		_, _ = io.WriteString(w, modinfo)
	})

	mux.HandleFunc("/debug/", func(w http.ResponseWriter, _ *http.Request) {
		h := w.Header()
		h.Set("Content-Type", "text/html; charset=utf-8")
		h.Set("Content-Length", strconv.Itoa(len(debugPage)))

		_, _ = io.WriteString(w, debugPage)
	})

	return mux
}

const debugPage = `<!doctype html>
<title>Debug</title>
<meta name=viewport content="width=device-width">
<ul>
 <li><a href=/debug/pprof>pprof</a></li>
 <li><a href=/debug/events>events</a></li>
 <li><a href=/debug/requests>traces</a></li>
</ul>
`
