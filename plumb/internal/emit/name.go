package emit

import (
	"go/token"
	"go/types"
	"maps"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	"go.pact.im/x/plumb/internal/solve"
)

type allocator struct {
	used map[string]bool
}

// reservedIdents returns the identifiers a generated name must not collide with,
// common to import-alias assignment and local-name allocation: the predeclared
// universe, the destination package’s own top-level names, and the lifted
// type-parameter names. Keywords are not seeded here: go/token exposes no
// enumerable keyword set, so they are rejected by the token.IsKeyword arm of each
// collision predicate instead. Each caller augments the result with its own extra
// category (set names for aliases, import aliases for locals).
func reservedIdents(dest *solve.DestInfo, lifted map[string]bool) map[string]bool {
	r := dest.ReservedNames()
	for n := range lifted {
		r[n] = true
	}
	return r
}

func newAllocator(dest *solve.DestInfo, q *qualifier, lifted map[string]bool) *allocator {
	used := reservedIdents(dest, lifted)
	for _, a := range q.aliasNames() {
		used[a] = true
	}
	return &allocator{used: used}
}

// branch returns a sub-allocator that treats every name taken so far as
// reserved but whose own allocations do not leak back to the parent. Names
// allocated in one cleanup branch (an if-block or the cleanup closure) are
// invisible to other branches, so each branch starts numbering afresh.
func (a *allocator) branch() *allocator {
	return &allocator{used: maps.Clone(a.used)}
}

func (a *allocator) alloc(base string) string {
	base = sanitizeIdent(base)
	name := base
	for i := 2; a.used[name] || token.IsKeyword(name); i++ {
		name = base + strconv.Itoa(i)
	}
	a.used[name] = true
	return name
}

func sanitizeIdent(s string) string {
	// A keyword base (a type named Type lowers to "type") is a valid identifier
	// shape; let it through so the allocator, whose collision test rejects every
	// keyword, suffixes it to a readable "type2" rather than collapsing it to "v".
	if s == "" || s == "_" || (!token.IsIdentifier(s) && !token.IsKeyword(s)) {
		return "v"
	}
	return s
}

// baseName derives a readable local-variable base name from a type.
func baseName(t types.Type) string {
	switch u := t.(type) {
	case *types.Pointer:
		return baseName(u.Elem())
	case *types.Slice:
		return baseName(u.Elem()) + "s"
	case *types.Array:
		return baseName(u.Elem()) + "s"
	case *types.Named:
		return lowerCamel(u.Obj().Name())
	case *types.Alias:
		return lowerCamel(u.Obj().Name())
	case *types.Basic:
		return lowerCamel(u.Name())
	case *types.Map:
		return "m"
	case *types.Chan:
		return "ch"
	case *types.Interface:
		return "iface"
	case *types.Signature:
		return "fn"
	case *types.TypeParam:
		return lowerCamel(u.Obj().Name())
	}
	return "v"
}

// cleanupBaseName derives a cleanup local’s base name from the provider function
// that returns it, e.g. OpenConn’s teardown becomes openConnCleanup. A cleanup
// result only ever comes from a function or method provider, so its Fn is always
// set. Naming after the provider (not the value the provider yields) keeps the
// cleanup unambiguous even when a provider returns several values.
func cleanupBaseName(fn *types.Func) string {
	return lowerCamel(fn.Name()) + "Cleanup"
}

// commonInitialisms are acronyms treated as one unit when lowering an
// identifier’s first word. This keeps the camelCase boundary after the whole
// acronym (HTTPServer becomes httpServer) and stops a leading acronym from
// merging with the next one (XMLHTTPRequest becomes xmlHTTPRequest).
var commonInitialisms = map[string]bool{
	"ACL": true, "API": true, "ASCII": true, "CPU": true, "CSS": true,
	"DNS": true, "EOF": true, "GUID": true, "HTML": true, "HTTP": true,
	"HTTPS": true, "ID": true, "IP": true, "JSON": true, "LHS": true,
	"QPS": true, "RAM": true, "RHS": true, "RPC": true, "SLA": true,
	"SMTP": true, "SQL": true, "SSH": true, "TCP": true, "TLS": true,
	"TTL": true, "UDP": true, "UI": true, "UID": true, "UUID": true,
	"URI": true, "URL": true, "UTF8": true, "VM": true, "XML": true,
	"XMPP": true, "XSRF": true, "XSS": true,
}

// lowerCamel lowercases the first word of an identifier while preserving
// camelCase: User → user, DB → db, DB2 → db2, HTTPServer → httpServer,
// XMLHTTPRequest → xmlHTTPRequest. A leading common initialism is lowered as a
// unit. Otherwise the leading run of capitals is lowered, except for its final
// capital when a lowercase letter follows it, since that capital begins the next
// word.
func lowerCamel(s string) string {
	n := firstWordLen(s)
	return strings.ToLower(s[:n]) + s[n:]
}

// firstWordLen returns the byte length of the prefix lowerCamel lowercases,
// which is the identifier’s first word. A leading common initialism is the whole
// word. Otherwise the word is the leading run of capitals, dropping its final
// capital when a lowercase letter follows, since that capital starts the next word.
func firstWordLen(s string) int {
	if a := leadingInitialism(s); a != "" {
		return len(a)
	}
	count, lastUpper := 0, 0
	for i, c := range s {
		if !unicode.IsUpper(c) {
			if count > 1 && unicode.IsLower(c) {
				return lastUpper // the run’s final capital begins the next word
			}
			return i
		}
		count++
		lastUpper = i
	}
	return len(s) // all of s is a single run of capitals
}

// leadingInitialism returns the longest common initialism that prefixes s at a
// word boundary, or "". A word boundary means the initialism is all of s, or the
// rune right after it is uppercase.
func leadingInitialism(s string) string {
	best := ""
	for a := range commonInitialisms {
		if len(a) <= len(best) || !strings.HasPrefix(s, a) {
			continue
		}
		if rest := s[len(a):]; rest == "" || startsUpper(rest) {
			best = a
		}
	}
	return best
}

func startsUpper(s string) bool {
	r, _ := utf8.DecodeRuneInString(s)
	return unicode.IsUpper(r)
}
