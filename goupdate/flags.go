package main

import (
	"errors"
	"flag"
)

const (
	runTestsNone = "none"
	runTestsWork = "work"
	runTestsAll  = "all"
)

const (
	formatHTML = "html"
	formatJSON = "json"
	formatText = "text"
)

// errParse was copied from flag package. It is used for Func flags to indicate
// parse error.
var errParse = errors.New("parse error")

type flags struct {
	goVersion string
	outPath   string
	outFormat string
	runTests  string
	testShort bool
}

const usage = `goupdate [flags]
  -go string
        If non-empty, sets go directive in all workspace’s go.mod files.
  -o string
        Report output path relative to the working directory. By default, report
        is written to stdout in "text" format. If this flag is set, the default
        format is "html".
  -format string
        Report output format ("html", "json", "text"). Newlines in "html" format
        are escaped for Markdown compatibility.
  -test string
        Run tests after updating modules ("none", "work", "all"). By default,
        no tests are run. Setting this flag to "work" runs tests for packages
        in the current workspace; "all" runs tests for all packages in the build
        graph, i.e. go test all. Test results are parsed and included in report
        output.
  -test-short
        Pass -short flag to go test when running tests. If -test is not already
        set, it defaults to "work".
`

func parseFlags(args []string) (*flags, error) {
	var f flags

	fs := flag.NewFlagSet("goupdate", flag.ContinueOnError)
	fs.Usage = func() {}
	fs.StringVar(&f.goVersion, "go", "", "")
	fs.StringVar(&f.outPath, "o", "", "")
	fs.BoolVar(&f.testShort, "test-short", false, "")
	fs.Func("format", "", func(v string) error {
		switch v {
		case formatHTML:
		case formatJSON:
		case formatText:
		default:
			return errParse
		}
		f.outFormat = v
		return nil
	})
	fs.Func("test", "run tests; values are none (default), work, all", func(v string) error {
		switch v {
		case runTestsNone:
		case runTestsWork:
		case runTestsAll:
		default:
			return errParse
		}
		f.runTests = v
		return nil
	})
	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	if f.runTests == "" {
		if f.testShort {
			f.runTests = runTestsWork
		} else {
			f.runTests = runTestsNone
		}
	}
	if f.outFormat == "" {
		if f.outPath != "" {
			f.outFormat = formatHTML
		} else {
			f.outFormat = formatText
		}
	}
	return &f, nil
}
