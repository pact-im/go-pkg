// Command plumb generates dependency-wiring functions from //plumb:<name>
// providers.
//
// Tag the declarations that produce values with a //plumb:<name> directive
// (functions, methods, variables, constants, conversions, struct fields, and
// struct types), and plumb generates a function named <name> that calls them in
// dependency order. Unlike [wire], you do not write the wiring function or declare
// its signature. plumb infers it: the parameters are what the providers need but
// nobody supplies (inputs), and the results are what they produce but nobody
// consumes (outputs).
//
// # Running plumb
//
// Invoke plumb as:
//
//	plumb [-package-name=<name>] [-import-path=<path>] [-output=<file>] [-v] [<packagePattern>...]
//
// The flags are:
//
//   - -package-name: the name in the generated "package <name>" clause. Defaults to
//     the name of the package the file is generated into (the scanned package, when
//     -import-path is inferred or names one); required when that package is not
//     scanned, and rejected when it contradicts a scanned destination’s real name.
//   - -import-path: the import path of the package the generated file belongs to.
//     Its own providers are referenced unqualified; everything else is imported and
//     qualified. Defaults to the sole scanned package, so a single-package go:generate
//     needs neither flag. Scanning more than one is ambiguous and must name the
//     destination. A path matching no scanned package makes a separate package. It
//     must be a valid Go import path. See “Where the code goes”, below.
//   - -output: file to write; defaults to stdout. Its parent directory must already
//     exist, since plumb does not create intermediate directories.
//   - -v: print a discovery and inference report to stderr: packages scanned,
//     providers found and where, and each generated signature.
//   - <packagePattern>...: zero or more go/packages patterns (".", "./service",
//     "./...", "work", or a full import path). Each may match several packages, and
//     you may pass several patterns. With none, plumb scans the current directory (".").
//
// plumb finds every set across all the matched packages and writes one file
// containing one function per set. Commit the generated file, and re-run plumb when
// providers change. The usual way to wire it up is a go:generate directive. Placed in
// the providers package, it generates into that package with nothing but an output file,
// inferring both the package name and the import path from where it sits:
//
//	//go:generate plumb -output=plumb_gen.go
//
// # Directives
//
// A provider is tagged with a line comment:
//
//	//plumb:build
//
// “build” is the set name and becomes the generated function’s name. It must be a
// valid unexported identifier, starting with a lowercase letter. This is not a style
// choice: //plumb:Build is not a real Go directive. gofmt rewrites it to
// "// plumb:Build" (a plain comment), silently disabling it, so keep the name
// lowercase. To get an exported wiring function, call the generated one from a small
// hand-written wrapper. The directive names a set and nothing else; trailing text
// (//plumb:build extra) is rejected.
//
// A provider may carry several directives to join several sets:
//
//	//plumb:build
//	//plumb:test
//	func NewClock() Clock { return realClock{} }
//
// # Function providers
//
// The common case: a constructor. Its parameters are what it consumes; its results
// are what it produces.
//
//	//plumb:build
//	func NewDB(cfg *Config) (*sql.DB, error) { ... }
//
//	//plumb:build
//	func NewStore(db *sql.DB) *Store { ... }
//
// Variadic parameters are supported; a variadic parameter is treated as a slice ([]T).
//
// A function or method with no results is a side-effect provider: it consumes its inputs
// and is emitted as a bare call (Register(mux, h)) that produces nothing.
//
// # Variable providers
//
// A package-level variable is a value provider, equivalent to a function that returns
// it. The generated code references the variable directly. If a consumer needs a
// pointer to its type, plumb takes the address of a fresh copy, so the pointer points
// at that copy rather than the variable’s own storage (see “Value and pointer”, below).
//
//	//plumb:build
//	var stdin = os.Stdin
//
// is equivalent to
//
//	//plumb:build
//	func stdin() *os.File { return os.Stdin }
//
// # Constant providers
//
// A typed package-level constant is a value provider too, just like a variable. The
// generated code references the constant directly.
//
//	//plumb:build
//	const DefaultPort Port = 8080
//
// provides Port. The constant must have an explicit type. An untyped constant
// (const DefaultPort = 8080) has no determinate provided type, so plumb rejects it.
//
// # Conversions
//
// A blank variable with an explicit target type is a conversion provider: it consumes
// the source type and produces the target type, emitting the Go conversion T(src).
// Binding a concrete type to an interface (the analog of wire.Bind) is the common
// case, but it is not special. Any conversion the blank-variable assignment
// type-checks is allowed.
//
//	//plumb:build
//	var _ io.Reader = (*os.File)(nil)
//
// is roughly equivalent to
//
//	//plumb:build
//	func bind(f *os.File) io.Reader { return f }
//
// except the generated code is a direct conversion, "reader := io.Reader(file)",
// rather than a function call. The value on the right-hand side is only a type hint:
// (*os.File)(nil) is the idiom for “the source type is *os.File”. Now anything that
// needs an io.Reader is satisfied by whatever produces a *os.File. The value must name
// a source type, so the untyped nil (var _ io.Reader = nil) is rejected
// because it names none; write (*os.File)(nil).
//
// # Struct field providers
//
// A struct field can be a provider. It consumes the struct and produces the
// field’s type.
//
//	type Config struct {
//		//plumb:build
//		Port int
//	}
//
// is equivalent to
//
//	//plumb:build
//	func port(c *Config) int { return c.Port }
//
// plumb consumes the struct as a value or a pointer, whichever your set already has.
// It prefers *Config if some other provider produces or consumes *Config; otherwise
// it takes Config by value. Field access works for both.
//
// # Struct type providers
//
// The inverse: a //plumb: directive on a struct type makes the type build itself from
// the set. Its exported fields become inputs, as if they were constructor arguments.
// Unexported fields are left as zero values.
//
//	//plumb:build
//	type Server struct {
//		Addr string      // exported → input
//		DB   *DB         // exported → input
//		log  *log.Logger // unexported → left zero
//	}
//
// is equivalent to
//
//	//plumb:build
//	func newServer(addr string, db *DB) Server { return Server{Addr: addr, DB: db} }
//
// It provides the value Server. If another provider in the set needs *Server, plumb
// also makes the pointer form available (it builds the value once and takes its
// address), so you can consume either Server or *Server without writing two providers.
// A generic struct works too, instantiated like any other generic provider (see
// “Generic providers”, below). A directive on a non-struct type (including a type
// alias whose target is not a struct) is rejected.
//
// # Value and pointer
//
// plumb wires by copy, and bridges between a value T and a pointer *T automatically, in
// both directions, where T is a named type, a type parameter, or a predeclared basic type
// (int⇄*int bridges just like MyInt⇄*MyInt; an inline composite such as []byte does not).
// When a provider yields a value T and a
// consumer needs *T, plumb holds the value in a local and takes its address (&v).
// Because the local is a fresh copy, *T never aliases a variable’s or field’s own
// storage, so mutating through
// the pointer never changes the original. When a provider yields *T and a consumer needs
// T, plumb dereferences it (*ptr); that assumes the pointer is non-nil, and a nil panics
// at run time, which plumb cannot rule out statically. The bridge is one level only
// (T⇄*T, never T⇄**T) and is not interface satisfaction.
//
// # Method providers
//
// A method is a provider whose receiver is its first input:
//
//	//plumb:build
//	func (db *DB) Store() *Store { return ... }
//
// is equivalent to
//
//	//plumb:build
//	func store(db *DB) *Store { return db.Store() }
//
// An interface method works the same way, with the interface as the receiver:
//
//	type Reporter interface {
//		//plumb:build
//		Report() *Summary
//	}
//
// requires a Reporter and provides a *Summary. The generated code calls the method on
// the receiver value (db.Store(), reporter.Report()), so the receiver is never
// package-qualified, but a method must be exported when the generated file lives in a
// different package. Methods on generic receivers are supported (see “Generic
// providers”, below).
//
// # Generic providers
//
// A generic provider is a template: plumb instantiates it as the wiring demands. A type
// parameter pinned to a concrete type is instantiated there, and the same provider can
// be instantiated at several types in one function. The generated call writes the type
// arguments explicitly, which is always valid Go:
//
//	//plumb:build
//	func NewCache[T any]() *Cache[T] { return ... }
//
//	//plumb:build
//	func NewUserService(c *Cache[User]) *UserService { return ... }
//
//	//plumb:build
//	func NewSessionService(c *Cache[Session]) *SessionService { return ... }
//
// NewUserService needs *Cache[User] and NewSessionService needs *Cache[Session], so
// NewCache is called at both NewCache[User]() and NewCache[Session](); the two results
// are distinct types and never collide.
//
// A type parameter that nothing pins is carried to the generated function, which
// becomes generic over it. plumb keeps the template’s own parameter name where it can:
//
//	//plumb:build
//	func NewPool[T any]() *Pool[T] { return ... }
//
//	//plumb:build
//	func Summarize[T any](p *Pool[T]) *Report { return ... }
//
//	func build[T any]() (report *Report) {
//		pool := NewPool[T]()
//		report2 := Summarize[T](pool)
//		report = report2
//		return
//	}
//
// Methods on generic receivers work too. The receiver value carries its type
// arguments, so the call needs none: func (c *Cache[T]) Snapshot() *Snapshot[T] emits
// cache.Snapshot(). A field on a generic struct is likewise a template, equivalent to
// func(Foo[T]) Baz[T], with the struct as its input, taken by value or pointer just like
// a field on a non-generic struct:
//
//	type Foo[T any] struct {
//		//plumb:build
//		Bar Baz[T]
//	}
//
// A method on a generic basic interface (its type set is exactly its methods) is also
// a template, equivalent to func(Source[T]) Result[T], with the interface as the input:
//
//	type Source[T any] interface {
//		//plumb:build
//		Fetch() Result[T]
//	}
//
// A method on a generic constraint interface is rejected: a general, non-basic
// interface, say one embedding comparable. A non-basic interface can only be used as a
// type constraint, never as the type of a value, so it cannot be a receiver.
//
// A template with several results over different type parameters is wired by joint
// pinning: each result is pinned independently, by its own consumer. Given
//
//	//plumb:build
//	func NewKVPair[T, U any]() (Key[T], Value[U]) { ... }
//
// plumb instantiates it once one consumer pins Key[T] and another pins Value[U]:
// NewKVPair[int, string]() when the set needs Key[int] and Value[string].
//
// A demand can match a template’s result by shape yet pin a type its constraint
// rejects, so the template cannot build it. plumb treats that the way it treats any
// type nobody produces: the demand becomes an input to the generated function (see
// “Rules and gotchas”), while the template still serves the demands it can. It only
// becomes an error when that leaves the template with no consumer at all. Then the
// diagnostic names the failed constraint.
//
// plumb rejects, with a clear message, a result-generic template that nothing
// instantiates (unused), one with a bare type-parameter result (func Default[T]() T),
// two templates that match a demand ambiguously, and an instantiation that does not
// terminate (a provider manufacturing unboundedly larger types).
//
// # The generated function
//
// plumb treats the wiring as a bag of typed values, one value per type. From the
// providers in a set it computes:
//
//   - inputs: types required by some provider but produced by none. They become
//     the generated function’s parameters.
//   - outputs: types produced by some provider but consumed by none. They become
//     its results.
//
// So adding a provider that consumes an existing output, or removing the last consumer
// of a type, changes the inferred signature. Run with -v to see it.
//
// # Errors
//
// If any provider returns a predeclared error, the generated function returns a
// trailing error. Each fallible call is checked and the error returned immediately.
// Every output is held in a local and copied to its named result only on the success
// path, so an error never returns a partially built value; the other results stay
// zero:
//
//	func build(cfg *Config) (store *Store, err error) {
//		db, e := NewDB(cfg)
//		if e != nil {
//			err = e
//			return
//		}
//		store2 := NewStore(db)
//		store = store2
//		return
//	}
//
// # Cleanup functions
//
// A provider can return a func() to release whatever it just acquired, the same
// (T, func(), error) convention as wire. plumb recognizes the bare func() (no
// parameters, no results) as a teardown hook rather than a value, so it never enters
// the bag. Instead the generated function returns a single aggregated func(), placed
// after the outputs and before any error.
//
// A resource worth tearing down can usually also fail to open, so cleanup providers
// typically return (T, func(), error). Each error check unwinds whatever was acquired
// before it:
//
//	//plumb:build
//	func OpenConn() (*Conn, func(), error) { ... }
//
//	//plumb:build
//	func OpenPool(c *Conn) (*Pool, func(), error) { ... }
//
//	//plumb:build
//	func OpenServer(p *Pool) (*Server, func(), error) { ... }
//
// generates:
//
//	func build() (server *Server, cleanup func(), err error) {
//		conn, openConnCleanup, e := OpenConn()
//		if e != nil {
//			err = e
//			return
//		}
//		pool, openPoolCleanup, e := OpenPool(conn)
//		if e != nil {
//			openConnCleanup()
//			err = e
//			return
//		}
//		server2, openServerCleanup, e := OpenServer(pool)
//		if e != nil {
//			openPoolCleanup()
//			openConnCleanup()
//			err = e
//			return
//		}
//		server = server2
//		cleanup = func() {
//			openServerCleanup()
//			openPoolCleanup()
//			openConnCleanup()
//		}
//		return
//	}
//
// Call it, check the error, then defer the cleanup:
//
//	server, cleanup, err := build()
//	if err != nil {
//		return err
//	}
//	defer cleanup()
//
// Two things to note. First, cleanups run last-acquired-first, both in the aggregated
// cleanup and on a mid-build error: above, OpenPool failing runs openConnCleanup
// before returning, and OpenServer failing runs openPoolCleanup then openConnCleanup.
// A provider that itself fails is assumed to have cleaned up after itself, so its own
// returned cleanup is not run. Second, the reservation is only of a bare func() result:
// a function or method returning func() yields a teardown hook, not a value, so to return
// a func() as an ordinary value give the type a name (type Handler func()). A func()-typed
// variable, field, or conversion is already an ordinary value.
//
// A cleanup may also report its own failure by returning an error: func() error, like
// (*sql.DB).Close. plumb recognizes that too. When any cleanup in a set is failable, the
// aggregated cleanup becomes func() error. A single failable cleanup with nothing to
// combine returns its error directly:
//
//	//plumb:setup
//	func OpenDB(dsn string) (*sql.DB, func() error, error) {
//		db, err := sql.Open("pgx", dsn)
//		if err != nil {
//			return nil, nil, err
//		}
//		return db, db.Close, nil
//	}
//
// generates:
//
//	func setup(dsn string) (db *sql.DB, cleanup func() error, err error) {
//		db2, openDBCleanup, e := OpenDB(dsn)
//		if e != nil {
//			err = e
//			return
//		}
//		db = db2
//		cleanup = func() error {
//			cleanupErr := openDBCleanup()
//			return cleanupErr
//		}
//		return
//	}
//
// When two or more errors can combine (two failable cleanups, or a failable cleanup
// plus a mid-build error), they are joined newest-first with errors.Join, which drops
// nil errors, so a clean teardown still returns nil.
//
// # Where the code goes
//
// The generated file lives in some package, and that determines what gets imported:
//
//   - A provider in the same package as the generated file is called unqualified
//     and not imported.
//   - A provider in any other package is imported and qualified (pkg.NewX), so it must
//     be exported, as must any type that appears in the generated signature.
//
// plumb keys on the generated file’s import path, not its output location: two packages
// can share a name at different paths, and -output may even be stdout. It takes that
// path from -import-path, or infers the sole scanned package when -import-path is
// omitted, so a go:generate in that package generates “in package” (like wire’s
// wire_gen.go), its providers unqualified, with no flags. A -import-path matching no
// scanned package makes a separate package that imports and qualifies every provider.
//
// Re-running over existing output is safe. In same-package mode plumb re-scans the file
// it is about to replace, so a stale generated file that no longer compiles (a provider
// whose signature changed, or one since removed) does not stop plumb from regenerating
// the corrected wiring. It tolerates those type errors and rewrites the file.
//
// In same-package mode the generated imports and function names share the package block
// with the destination package’s own declarations, so plumb reads that package’s
// top-level identifiers and avoids them. An import qualifier that would collide is
// aliased. A generated function name cannot be aliased (the set names it), so a
// set whose name is already declared in the destination package (or used as an
// import qualifier in one of its hand-written files) is rejected, asking you to rename
// it; the file plumb is about to overwrite (-output) is excluded, so regenerating over
// plumb’s own prior function still works. With output to stdout, where plumb cannot tell
// which file it would replace, that name check is skipped.
//
// These safeguards need the destination to be one of the scanned packages. A -import-path
// matching none of them (a separate package, or one the patterns did not load) leaves
// plumb blind to its declarations, so a generated qualifier or function name may collide
// with one; include the destination package in the patterns for collision-safe output.
//
// # Multiple packages
//
// A pattern such as "./..." scans many packages at once, and a set’s providers may
// live in different packages:
//
//	// package store
//	//plumb:build
//	func NewDB(cfg *Config) (*DB, error) { ... }
//
//	// package web (imports store)
//	//plumb:build
//	func NewServer(db *store.DB) *Server { ... }
//
//	plumb -package-name=app -import-path=example.com/app -output=app/plumb_gen.go ./...
//
// generates a separate app package. Its import path matches neither store nor web, so
// both are imported. When scanning more than one package, you must name the destination:
// -import-path is required, and -package-name too, since app is not a scanned package.
// To generate into web instead, name it: "-import-path=example.com/web" makes web’s
// providers unqualified and imports only store, with -package-name inferred as web.
//
// # Rules and gotchas
//
//   - One value per type. Two parameters of the same type get the same value; you
//     cannot wire two distinct values of one type. Use a slice, a fixed-size array, or
//     a distinct defined type per value (see “Not supported”, below).
//   - One producer per type. If two providers produce the same type, plumb reports an
//     ambiguity. This includes a function returning two values of the same type and a
//     multi-name var/const declaring two names of one type. (Repeating the same directive
//     on one declaration is instead a distinct duplicate-directive error.)
//   - Unprovided means input. There is no “missing dependency” error: a typo in a
//     type, or a forgotten provider, silently widens the signature rather than
//     failing. -v prints the inferred inputs so you can catch this.
//   - Cycles are errors. If the providers cannot be ordered (A needs something B
//     produces and vice versa), plumb reports the cycle.
//   - Matching is by exact type, not Go’s implements relation. A consumer of io.Reader
//     is not matched by a producer of *os.File; add a conversion.
//   - The value/pointer bridge is producer-driven. plumb bridges T and *T when a
//     producer of one covers a demand for the other, but if neither is produced and
//     both are demanded, they become two separate injector inputs (a T and a *T);
//     -v shows both. Give one a producer or a conversion to collapse them.
//   - Directives attach to package-level functions and methods, package-level variables
//     and constants, conversions, struct fields, struct types, and interface methods.
//     Locals are never providers, and embedded fields/interfaces cannot carry one.
//   - Test files are never scanned. A //plumb: directive in a _test.go file is
//     ignored with no diagnostic. This is the one place plumb drops a directive
//     silently, a deliberate limitation. Put a provider a test needs in a regular
//     file, or build it by hand in the test.
//   - Build tags follow the build context: GOFLAGS=-tags=integration plumb … scans
//     tag-gated files, and the output then compiles under that same context.
//   - The destination must be able to import every provider’s package. plumb checks
//     that cross-package names are exported, but not importability, so a provider in a
//     package main or in an internal package outside the destination’s subtree
//     generates an import the compiler, not plumb, rejects.
//   - An -output written outside the destination package’s directory must not share a
//     base name with one of that package’s source files: the file plumb skips as its
//     own prior output is identified by base name, so a matching source file would be
//     silently ignored.
//
// # Not supported
//
// Multiple distinct values of the same type in one set will never be supported. This
// is fundamental to plumb’s type-keyed model: each type maps to exactly one value, so
// two parameters of the same type always receive the same value, and two producers of
// one type are an ambiguity error. To wire several values of one underlying type, use a
// slice ([]T) or fixed-size array ([N]T) for a homogeneous collection, or a distinct
// defined type per use case (type ReadTimeout time.Duration and type WriteTimeout
// time.Duration rather than two time.Duration values). The same applies to decorators:
// func(T) T next to another producer of T is ambiguous, so give the decorated value a
// distinct type.
//
// Method providers on generic constraint interfaces are not supported either, since a
// constraint type cannot be a value. Functions, generic-receiver methods, generic
// struct fields, and generic basic interfaces are all supported (see “Generic
// providers”, above).
//
// # Example
//
// Providers in ./providers:
//
//	package providers
//
//	import (
//		"io"
//		"os"
//	)
//
//	type Config struct {
//		//plumb:build
//		Addr string
//	}
//
//	type Server struct{}
//
//	//plumb:build
//	var Stdin = os.Stdin
//
//	//plumb:build
//	var _ io.Reader = (*os.File)(nil)
//
//	//plumb:build
//	func NewServer(in io.Reader, addr string, cfg *Config) (*Server, error) {
//		return &Server{}, nil
//	}
//
// Run:
//
//	plumb -package-name=app -import-path=example.com/app -output=app/plumb_gen.go ./providers
//
// Generated app/plumb_gen.go:
//
//	// Code generated by plumb. DO NOT EDIT.
//
//	package app
//
//	import (
//		"example.com/providers"
//		"io"
//	)
//
//	func build(config *providers.Config) (server *providers.Server, err error) {
//		string2 := config.Addr
//		file := providers.Stdin
//		reader := io.Reader(file)
//		server2, e := providers.NewServer(reader, string2, config)
//		if e != nil {
//			err = e
//			return
//		}
//		server = server2
//		return
//	}
//
// *providers.Config is the only input, since nothing produces it; *providers.Server and
// error are the outputs. Call build(cfg) from your app package.
//
// [wire]: https://go.dev/blog/wire
package main

import (
	"os"

	"go.pact.im/x/plumb/internal/cli"
)

func main() {
	os.Exit(cli.Run(os.Args[1:], os.Stdout, os.Stderr, os.Stderr))
}
