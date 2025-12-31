// Package httpdebug provides HTTP handler for debug endpoints.
package httpdebug

import (
	"expvar"
	"io"
	"net/http"
	"net/http/pprof"
	"runtime"
	"runtime/debug"
	"strconv"

	"golang.org/x/net/trace"
)

// New returns a new HTTP handler for debug endpoints.
func New() http.Handler {
	mux := http.NewServeMux()
	handlePprof(mux)
	handleExpvar(mux)
	handleNetTrace(mux, true)
	handleBuildInfo(mux)
	handleIndex(mux)
	return mux
}

func handlePprof(mux *http.ServeMux) {
	mux.HandleFunc("GET /debug/pprof/", pprof.Index)
	mux.HandleFunc("GET /debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("GET /debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("GET /debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("GET /debug/pprof/trace", pprof.Trace)
}

func handleExpvar(mux *http.ServeMux) {
	mux.Handle("GET /debug/vars", expvar.Handler())
}

func handleNetTrace(mux *http.ServeMux, sensitive bool) {
	mux.HandleFunc("GET /debug/requests", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		trace.Render(w, r, sensitive)
	})
	mux.HandleFunc("GET /debug/events", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		trace.RenderEvents(w, r, sensitive)
	})
}

func handleBuildInfo(mux *http.ServeMux) {
	mux.HandleFunc("GET /debug/buildinfo", serveBuildInfo)
}

func serveBuildInfo(w http.ResponseWriter, _ *http.Request) {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		info = &debug.BuildInfo{GoVersion: runtime.Version()}
	}

	// TODO: support JSON responses; see also https://go.dev/issue/19307

	modinfo := info.String()

	h := w.Header()
	h.Set("Content-Type", "text/plain; charset=utf-8")
	h.Set("Content-Length", strconv.Itoa(len(modinfo)))
	_, _ = io.WriteString(w, modinfo)
}

func handleIndex(mux *http.ServeMux) {
	mux.HandleFunc("/debug/", serveIndex)
}

func serveIndex(w http.ResponseWriter, _ *http.Request) {
	h := w.Header()
	h.Set("Content-Type", "text/html; charset=utf-8")
	h.Set("Content-Length", strconv.Itoa(len(debugPage)))
	_, _ = io.WriteString(w, debugPage)
}

const debugPage = `<!doctype html>
<html lang=en>
<title>Debug</title>
<meta name=viewport content="width=device-width">
<ul>
 <li><a href=/debug/pprof>pprof</a></li>
 <li><a href=/debug/vars>expvar</a></li>
 <li><a href=/debug/events>events</a></li>
 <li><a href=/debug/requests>traces</a></li>
 <li><a href=/debug/buildinfo>build info</a></li>
</ul>
`
