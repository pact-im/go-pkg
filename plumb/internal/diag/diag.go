// Package diag holds plumb's diagnostic vocabulary: the located, user-facing
// error type, the sentinel errors every phase wraps, and the position helpers
// that order diagnostics deterministically. It is the leaf every other package
// reaches for when it must reject input or compare source positions.
package diag

import (
	"cmp"
	"errors"
	"fmt"
	"go/token"
)

// ErrorKind is a sentinel diagnostic kind: a constant string that implements
// error. Making it a named string type lets the sentinels below be true compile-
// time constants while still being wrappable and matchable with errors.Is (a
// string is comparable, so errors.Is compares them by value).
type ErrorKind string

// Error returns the sentinel's message.
func (e ErrorKind) Error() string { return string(e) }

// Sentinel errors for every condition plumb can reject. Each diagnostic wraps
// one of these, so callers and tests can classify a failure with errors.Is
// rather than matching message text. The human-readable specifics travel in the
// wrapped detail; the source position travels on *Error.
const (
	ErrNoDirectives              ErrorKind = "no //plumb: directives found"
	ErrInvalidSetName            ErrorKind = "invalid set name"
	ErrMisplacedDirective        ErrorKind = "directive in an unsupported position"
	ErrDuplicateDirective        ErrorKind = "duplicate provider directive"
	ErrInvalidConversion         ErrorKind = "invalid conversion"
	ErrUntypedConstant           ErrorKind = "untyped constant provider"
	ErrBlankProvider             ErrorKind = "blank provider name"
	ErrInitProvider              ErrorKind = "provider directive on func init"
	ErrEmbeddedField             ErrorKind = "provider directive on embedded field"
	ErrEmbeddedInterface         ErrorKind = "provider directive on embedded interface"
	ErrConstraintInterfaceMethod ErrorKind = "method provider on constraint interface"
	ErrStructProvider            ErrorKind = "invalid struct-type provider"
	ErrAmbiguousProducer         ErrorKind = "ambiguous producer"
	ErrMultipleErrors            ErrorKind = "multiple error results"
	ErrDependencyCycle           ErrorKind = "dependency cycle"
	ErrUnusedTemplate            ErrorKind = "unused result-generic template"
	ErrBareTypeParamResult       ErrorKind = "bare type-parameter result"
	ErrAmbiguousTemplates        ErrorKind = "ambiguous templates"
	ErrNonTerminating            ErrorKind = "non-terminating generic instantiation"
	ErrUnexportedProvider        ErrorKind = "unexported provider referenced across a package boundary"
	ErrUnreachableType           ErrorKind = "type not reachable from the destination package"
	ErrInvalidType               ErrorKind = "type does not type-check"
	ErrReservedName              ErrorKind = "reserved function name with a non-empty signature"
	ErrShadowsPredeclared        ErrorKind = "set name shadows a predeclared identifier"
	ErrDestShadowsPredeclared    ErrorKind = "destination declaration shadows a predeclared identifier"
	ErrSetNameCollision          ErrorKind = "set name collides with an existing declaration"
)

// Error is a located, user-facing diagnostic. It wraps one of the sentinel
// errors above (reachable with errors.Is and errors.Unwrap) and carries the
// source position where the problem occurred. Violations of plumb's own
// internal invariants are panics instead, never *Error.
type Error struct {
	pos token.Position // the offending source position; zero if none applies
	err error          // sentinel wrapped with the specific detail
}

func (e *Error) Error() string {
	if e.pos.IsValid() {
		return fmt.Sprintf("%s: %s", e.pos, e.err.Error())
	}
	return e.err.Error()
}

// Unwrap exposes the wrapped chain so errors.Is/As reach the sentinel.
func (e *Error) Unwrap() error { return e.err }

// Errorf builds a located *Error wrapping the given sentinel kind. detail is a
// Printf-style format describing the specifics; it should not repeat the
// sentinel's own words.
func Errorf(pos token.Position, kind ErrorKind, detail string, args ...any) *Error {
	wrapped := fmt.Errorf("%w: "+detail, append([]any{kind}, args...)...)
	return &Error{pos: pos, err: wrapped}
}

// Sentinel returns a positionless *Error wrapping msg. It exists for internal
// control-flow signals that ride the *Error return channel but are compared by
// identity, never by content: unlike a zero-value *Error (whose nil wrapped err
// makes Error() panic), a Sentinel is safe to format if one is ever accidentally
// printed or logged.
func Sentinel(msg string) *Error {
	return &Error{err: errors.New(msg)}
}

// PosIn resolves p to a token.Position within fset, returning the zero Position
// when fset is nil or p is invalid.
func PosIn(fset *token.FileSet, p token.Pos) token.Position {
	if fset == nil || !p.IsValid() {
		return token.Position{}
	}
	return fset.Position(p)
}

// CmpPos orders positions deterministically: by file name, then byte offset.
// This is the single ordering authority for source-anchored decisions. Read
// order (of packages, files, or map entries) never participates. It is the
// int-returning comparator backing slices.SortFunc/SortStableFunc/MinFunc.
func CmpPos(a, b token.Position) int {
	return cmp.Or(cmp.Compare(a.Filename, b.Filename), cmp.Compare(a.Offset, b.Offset))
}

// Earlier returns the source-earlier of two diagnostics by CmpPos, treating a nil
// *Error as "no fault" (later than any real one); it returns nil only when both
// are nil. Use it to report the fault a reader would fix first when a phase finds
// faults out of source order.
func Earlier(a, b *Error) *Error {
	switch {
	case a == nil:
		return b
	case b == nil:
		return a
	case CmpPos(b.pos, a.pos) < 0:
		return b
	default:
		return a
	}
}
