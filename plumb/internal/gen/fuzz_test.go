package gen_test

import (
	"errors"
	"maps"
	"strings"
	"testing"

	"go.pact.im/x/plumb/internal/gen"
	"go.pact.im/x/plumb/internal/packagestest"
)

// FuzzGenerate exercises the central invariant on arbitrary input: generation
// must never panic, and whenever it succeeds on type-correct input the emitted
// code must type-check. Tolerated type-error input (which the loader keeps going
// on) is fed in too and must not crash, but is exempt from the compile assertion
// since its output may reproduce the user's own type error.
//
// The three bodies form a multi-package providers graph: package a, package b
// (which may import a), and a pre-existing package dest with its own top-level
// declarations, so generation runs against more than one scanned package and
// across an existing destination. The modes cover the three placements where
// collision and qualification bugs hide: same-package into a scanned providers
// package, a fresh separate package, and a pre-existing scanned destination whose
// declarations and import qualifiers the generated file must avoid.
func FuzzGenerate(f *testing.F) {
	for _, s := range fuzzSeeds {
		f.Add(s.a, s.b, s.dest)
	}

	f.Fuzz(func(t *testing.T, bodyA, bodyB, bodyDest string) {
		files := map[string]string{
			"example.com/a/a.go":           "package a\n" + bodyA,
			"example.com/b/b.go":           "package b\n" + bodyB,
			"example.com/dest/existing.go": "package dest\n" + bodyDest,
		}

		loaded, err := packagestest.Load(files)
		if err != nil {
			t.Skip() // not parseable Go
		}
		if len(loaded.Packages) == 0 {
			t.Skip()
		}
		// The no-crash guarantee covers tolerated type-error input too: the loader
		// keeps going on type errors, so such input still reaches Generate and must
		// not panic. But its output may faithfully reproduce the user's own type
		// error, so the compile assertion is skipped for it.
		typeErrs := len(loaded.TypeErrors) != 0

		modes := []gen.Options{
			{ImportPath: "example.com/a", PackageName: "a", OutputBase: "plumb_gen.go"},       // same-package, one of several scanned
			{ImportPath: "example.com/gen", PackageName: "gen", OutputBase: "plumb_gen.go"},   // separate fresh package
			{ImportPath: "example.com/dest", PackageName: "dest", OutputBase: "plumb_gen.go"}, // pre-existing scanned destination
		}
		for _, opts := range modes {
			res, gerr := gen.Generate(opts, loaded.Packages) // must not panic
			if gerr != nil || typeErrs {
				continue
			}
			fuzzAssertCompiles(t, files, opts, res.Source)
		}
	})
}

// fuzzAssertCompiles overlays the generated source and asserts the result
// type-checks, failing loudly on a miscompile. The overlay sits beside any
// pre-existing files of the destination package, so a generated name colliding
// with one of them surfaces here as a type error.
func fuzzAssertCompiles(t *testing.T, files map[string]string, opts gen.Options, src string) {
	t.Helper()
	aug := map[string]string{}
	maps.Copy(aug, files)
	aug[opts.ImportPath+"/plumb_generated_fuzz.go"] = src
	loaded, err := packagestest.Load(aug)
	if err != nil {
		// A cycle-forming destination is valid input that plumb does not reject:
		// the generated file legitimately imports a package that imports the
		// destination back. plumb leaves that for the Go compiler to report, so a
		// cycle here is expected and the input is skipped, not a miscompile.
		if errors.Is(err, packagestest.ErrImportCycle) {
			t.Skip()
		}
		t.Fatalf("generated code failed to parse:\n%s\nerr: %v", src, err)
	}
	if len(loaded.TypeErrors) != 0 {
		var b strings.Builder
		for _, e := range loaded.TypeErrors {
			b.WriteString(e.Error())
			b.WriteByte('\n')
		}
		t.Fatalf("MISCOMPILE: generated code does not type-check:\n%s\n--- errors ---\n%s", src, b.String())
	}
}

// fuzzSeed is one corpus entry: provider bodies for packages a and b and the
// pre-existing declarations of the destination package dest.
type fuzzSeed struct{ a, b, dest string }

var fuzzSeeds = []fuzzSeed{
	// --- single-package shapes (b/dest empty): provider kinds and generics ----
	{a: `type Config struct{ Addr string }
type Server struct{}
//plumb:build
func NewConfig() *Config { return &Config{} }
//plumb:build
func NewServer(c *Config) (*Server, error) { return &Server{}, nil }`},
	{a: `type Conn struct{}
//plumb:build
func OpenConn() (*Conn, func(), error) { return &Conn{}, func() {}, nil }`},
	{a: `type Cache[T any] struct{}
type User struct{}
//plumb:build
func NewCache[T any]() *Cache[T] { return &Cache[T]{} }
//plumb:build
func Use(c *Cache[User]) int { return 0 }`},
	{a: `type Config struct{ Port int }
//plumb:build
var Default = Config{Port: 8080}
//plumb:build
func NewServer(c *Config) int { return c.Port }`},
	{a: `type Port int
//plumb:build
const DefaultPort Port = 8080`},
	{a: `type Pool[T any] struct{}
//plumb:build
func NewPool[T any]() *Pool[T] { return &Pool[T]{} }
//plumb:build
func PoolSize[T any](p *Pool[T]) int { return 0 }`},
	{a: `type Server struct {
	//plumb:build
	Addr string
}
//plumb:build
func New() *Server { return &Server{} }`},
	{a: `type Key[T any] struct{}
type Value[U any] struct{}
type A struct{}
type B struct{}
//plumb:build
func NewKVPair[T, U any]() (Key[T], Value[U]) { return Key[T]{}, Value[U]{} }
//plumb:build
func NewA(k Key[int], v Value[string]) *A { return &A{} }
//plumb:build
func NewB(k Key[bool], v Value[byte]) *B { return &B{} }`},
	{a: `type Pool[T any] struct{}
//plumb:build
func MakePool[A any, B interface{ ~[]A }](p Pool[A], q Pool[B]) int { return 0 }`},
	{a: `type Inner[T any] struct{}
type Alias[T any] = Inner[T]
//plumb:build
func New[T any]() Alias[T] { return Alias[T]{} }
//plumb:build
func Use(x Alias[int]) int { return 0 }`},
	{a: `type Thing struct{}
//plumb:build
func NewMapper[T any]() func(T) T { return func(x T) T { return x } }
//plumb:build
func Use(m func(int) int) *Thing { return &Thing{} }`},
	{a: `type Sink struct{}
//plumb:build
func Getter[T any]() interface{ Get() T } { return nil }
//plumb:build
func Use(g interface{ Get() int }) *Sink { return &Sink{} }`},
	// constraint near-miss with a free input-only parameter (the lifted-leak case)
	{a: `type Stringer interface{ String() string }
type Name string
func (n Name) String() string { return string(n) }
type Box[T any] struct{ V T }
//plumb:build
func Conv[T Stringer, U any](u U) Box[T] { return Box[T]{} }
//plumb:build
func UseName(b Box[Name]) int { return 0 }
//plumb:build
func UseInt(b Box[int]) string { return "" }`},

	// --- multi-package: providers in a consumed across the boundary by b ------
	{
		a: `type DB struct{}
//plumb:build
func NewDB() *DB { return &DB{} }`,
		b: `import "example.com/a"
type Server struct{}
//plumb:build
func NewServer(db *a.DB) *Server { return &Server{} }`,
	},

	// --- pre-existing scanned destination with colliding declarations ---------
	// dest declares a type named "a", colliding with the import qualifier the
	// generated file needs for package a; plumb must alias the import.
	{
		a: `type DB struct{}
//plumb:build
func NewDB() *DB { return &DB{} }
//plumb:build
func Use(db *DB) int { return 0 }`,
		dest: `type a struct{}`,
	},
	// dest already declares the set name "build", which a generated function
	// cannot alias, so generation should report a collision (handled, not a panic
	// or a miscompile).
	{
		a: `type DB struct{}
//plumb:build
func NewDB() *DB { return &DB{} }`,
		dest: `func build() {}`,
	},

	// --- tolerated type-error input: must be reported, never crash -------------
	// A method with an undefined receiver and a lifted parameter with an undefined
	// constraint are type errors the loader tolerates; each must surface as a
	// located diagnostic, never a panic. Seeds so the fuzzer explores this space.
	{a: `//plumb:build
func (r *Missing) New() int { return 0 }`},
	{a: `type Key[T any] struct{}
type Value[U any] struct{}
//plumb:build
func NewKV[T any, U Undefined]() (Key[T], Value[U]) { return Key[T]{}, Value[U]{} }
//plumb:build
func UseKey(k Key[int]) int { return 0 }`},

	// --- func init provider: rejected, never an undefined init() call -----------
	{a: `type Config struct{}
//plumb:build
func init() {}
//plumb:build
func NewConfig() *Config { return &Config{} }`},

	// --- alias/target reachability: one type, two spellings ---------------------
	// An exported alias over an unexported target, both spelled in one rendered
	// type: reachability judges each spelling on its own name (same-package mode
	// generates; the separate-package modes must reject, never emit lib-internal
	// names). Seed so the fuzzer explores the alias spelling space.
	{a: `type hidden struct{ X int }
type Exported = hidden
type Thing struct{}
//plumb:build
func New(m map[Exported]hidden) *Thing { return &Thing{} }`},

	// --- pinning revision: pins arriving in different fixpoint rounds -----------
	// A joint template whose second pin surfaces only through another template's
	// instantiation, and the full-pin-over-partial shape that once aborted the
	// solve. Seeds so the fuzzer explores the revision/restart paths.
	{a: `type Key[T any] struct{}
type Val[U any] struct{}
type Box[W any] struct{}
type Cog[X any] struct{}
type SA struct{}
type SB struct{}
//plumb:build
func Make[T, U any]() (Key[T], Val[U]) { return Key[T]{}, Val[U]{} }
//plumb:build
func UseKey(k Key[int]) *SA { return &SA{} }
//plumb:build
func Q2[W, X any](v Val[W]) (Box[W], Cog[X]) { return Box[W]{}, Cog[X]{} }
//plumb:build
func UsePair(b Box[string]) *SB { return &SB{} }`},
	{a: `type Key[T any] struct{}
type Pair[T, U any] struct{}
type Thing[W any] struct{}
type SA struct{}
//plumb:build
func Make2[T, U any]() (Key[T], Pair[T, U]) { return Key[T]{}, Pair[T, U]{} }
//plumb:build
func UseKey(k Key[int]) *SA { return &SA{} }
//plumb:build
func MkThing[W any](p Pair[int, W]) Thing[W] { return Thing[W]{} }
//plumb:build
func UseThing(t Thing[string]) int { return 0 }`},

	// --- template arbitration: overlapping result shapes ------------------------
	// Two same-shape templates with disjoint constraints, each viable for exactly
	// one demand, and two joint templates competing for one demand. Seeds so the
	// fuzzer explores the constraint-filtered ambiguity paths.
	{a: `type Ints interface{ ~int }
type Strs interface{ ~string }
type Box[T any] struct{}
//plumb:build
func MakeInt[T Ints]() Box[T] { return Box[T]{} }
//plumb:build
func MakeStr[T Strs]() Box[T] { return Box[T]{} }
//plumb:build
func UseInt(b Box[int]) byte { return 0 }
//plumb:build
func UseStr(b Box[string]) rune { return 0 }`},
	{a: `type Key[T any] struct{}
type Val[U any] struct{}
type Other[U any] struct{}
type Sink struct{}
//plumb:build
func MakeA[T, U any]() (Key[T], Val[U]) { return Key[T]{}, Val[U]{} }
//plumb:build
func MakeB[T, U any]() (Key[T], Other[U]) { return Key[T]{}, Other[U]{} }
//plumb:build
func NewSink(k Key[int], v Val[string], o Other[bool]) *Sink { return nil }`},
}
