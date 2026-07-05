package discover

import (
	"errors"
	"go/ast"
	"go/parser"
	"go/token"
	"slices"
	"testing"

	"go.pact.im/x/plumb/internal/diag"
)

// TestDirectiveNames pins the directive-recognition rule: gofmt-demoted comments
// and other tools' directives are ignored, a recognized name must be a valid Go
// identifier, and a directive carries no trailing arguments.
func TestDirectiveNames(t *testing.T) {
	fset := token.NewFileSet()
	tests := []struct {
		name    string
		texts   []string
		want    []string
		wantErr error // sentinel, or nil
	}{
		{"valid", []string{"//plumb:build"}, []string{"build"}, nil},
		{"valid with underscore", []string{"//plumb:my_set"}, []string{"my_set"}, nil},
		{"gofmt-demoted space", []string{"// plumb:build"}, nil, nil},
		{"another tool", []string{"//go:build linux"}, nil, nil},
		{"uppercase is not a directive", []string{"//plumb:Build"}, nil, nil},
		{"blank is not a directive", []string{"//plumb:_"}, nil, nil},
		{"digit-leading name", []string{"//plumb:123"}, nil, diag.ErrInvalidSetName},
		{"keyword name", []string{"//plumb:range"}, nil, diag.ErrInvalidSetName},
		{"trailing arguments", []string{"//plumb:build extra"}, nil, diag.ErrInvalidSetName},
		{"repeated directive is rejected", []string{"//plumb:build", "//plumb:build"}, nil, diag.ErrDuplicateDirective},
		{"distinct sets are preserved", []string{"//plumb:build", "//plumb:other"}, []string{"build", "other"}, nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			doc := &ast.CommentGroup{}
			for i, text := range tc.texts {
				doc.List = append(doc.List, &ast.Comment{Slash: token.Pos(1 + i), Text: text})
			}
			got, err := directiveNames(doc, fset)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("names = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestAppendDirectiveNames pins the group/spec union that specDirectiveNames is
// built on: names already claimed by an enclosing group directive (prior) are
// unioned with a spec's own directives, but a spec directive restating a claimed
// set is a duplicate, the same rule as a repeat within one comment group, now
// across the group/spec boundary.
func TestAppendDirectiveNames(t *testing.T) {
	fset := token.NewFileSet()
	tests := []struct {
		name    string
		prior   []string
		texts   []string
		want    []string
		wantErr error // sentinel, or nil
	}{
		{"extra set is unioned in", []string{"build"}, []string{"//plumb:test"}, []string{"build", "test"}, nil},
		{"restating a group set is rejected", []string{"build"}, []string{"//plumb:build"}, nil, diag.ErrDuplicateDirective},
		{"nil doc keeps prior unchanged", []string{"build"}, nil, []string{"build"}, nil},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var doc *ast.CommentGroup
			if tc.texts != nil {
				doc = &ast.CommentGroup{}
				for i, text := range tc.texts {
					doc.List = append(doc.List, &ast.Comment{Slash: token.Pos(1 + i), Text: text})
				}
			}
			got, err := appendDirectiveNames(slices.Clone(tc.prior), doc, fset)
			if tc.wantErr != nil {
				if !errors.Is(err, tc.wantErr) {
					t.Fatalf("err = %v, want %v", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if !slices.Equal(got, tc.want) {
				t.Errorf("names = %v, want %v", got, tc.want)
			}
		})
	}
}

// TestCollectProviderDocs pins which comment groups the targeted scan claims as
// “read”: the docs of a func, a var/const/type spec, and a struct field or
// interface method, but not a directive on an import spec. This is the gate that
// reportStrayDirectives trusts, so a directive in any unmarked group is a stray.
func TestCollectProviderDocs(t *testing.T) {
	const src = `package p

import (
	//plumb:onimport
	"fmt"
)

//plumb:onfunc
func F() {}

//plumb:onvar
var V = 0

//plumb:ontype
type S struct {
	//plumb:onfield
	X int
}

type I interface {
	//plumb:onmethod
	M()
}
`
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "p.go", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	read := collectProviderDocs(file)
	marked := func(text string) bool {
		for _, g := range file.Comments {
			for _, c := range g.List {
				if c.Text == text {
					return read[g]
				}
			}
		}
		t.Fatalf("comment %q not found", text)
		return false
	}
	tests := map[string]bool{
		"//plumb:onfunc":   true,
		"//plumb:onvar":    true,
		"//plumb:ontype":   true,
		"//plumb:onfield":  true,
		"//plumb:onmethod": true,
		"//plumb:onimport": false, // imports are never providers
	}
	for text, want := range tests {
		if got := marked(text); got != want {
			t.Errorf("marked(%q) = %v, want %v", text, got, want)
		}
	}
}

// TestReportStrayDirectives pins the other half of the pairing: a //plumb:
// directive in a position the scan never reads (a function body, an import, or
// floating free between declarations) is reported, while a file whose directives
// all sit in scanned positions is clean. (Uses only the parsed AST and the file
// set; reportStrayDirectives reads no type information.)
func TestReportStrayDirectives(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{"legit positions only", "package p\n\n//plumb:build\nfunc F() {}\n", false},
		{"function body", "package p\n\nfunc F() {\n\t//plumb:build\n\tx := 0\n\t_ = x\n}\n", true},
		{"import spec", "package p\n\nimport (\n\t//plumb:build\n\t\"fmt\"\n)\n", true},
		{"free-floating", "package p\n\n//plumb:build\n\nvar V = 0\n", true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fset := token.NewFileSet()
			file, err := parser.ParseFile(fset, "p.go", tc.src, parser.ParseComments)
			if err != nil {
				t.Fatal(err)
			}
			gotErr := reportStrayDirectives(&Package{Fset: fset}, file)
			switch {
			case tc.wantErr && !errors.Is(gotErr, diag.ErrMisplacedDirective):
				t.Fatalf("err = %v, want ErrMisplacedDirective", gotErr)
			case !tc.wantErr && gotErr != nil:
				t.Fatalf("unexpected error: %v", gotErr)
			}
		})
	}
}
