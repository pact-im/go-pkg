// Goupdate is a tool for updating Go module dependencies. It is supports Go
// workspaces and is an alternative to running go get -t -u that attempts to
// be less conservative and produces apidiff reports.
//
// # Warning
//
// Note that goupdate loops until there are no updates available in the build
// graph. This may be inefficient but should be good enough in practice.
//
// Module dependencies are updated in-place and may leave workspace in broken
// state on failure. Use with extreme caution if there are unsaved changes. It
// is recommended to run goupdate as part of the scheduled CI pipeline instead.
//
// # CI
//
// Integrating goupdate with GitHub Actions or other CI solutions that support
// scheduled pipelines should be straightforward and only requires Go toolchain
// to be installed. Goupdate can generate HTML reports (that can be embedded in
// Markdown) with overview of updated modules and API changes. It can also run
// tests and includes results in the generated reports. Test failures are
// intended to be reviewed by humans and so do not cause goupdate to exit with
// non-zero status code.
//
// # Motivation
//
// Goupdate was written after fighting with dependabot storming repositories
// with pull requests per dependency update per module in workspace. In addition
// to that, dependabot does not support updating modules with pseudo-versions
// like golang.org/x/crypto which does not publish semver-tagged releases, see
// https://github.com/dependabot/dependabot-core/issues/3017. It also does not
// support Go workspaces which we use extensively since their introduction in
// Go 1.18.
package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	f, err := parseFlags(os.Args[1:])
	if err != nil {
		fmt.Fprint(os.Stderr, usage)
		must(err)
	}

	if _, err := exec.LookPath("go"); err != nil {
		must(err)
	}

	w, err := loadWorkspace()
	if err != nil {
		must(err)
	}

	// Update Go version in each module if needed.
	if f.goVersion != "" {
		for _, dir := range w.Paths {
			oldVersion, err := updateGoModVersion(dir, f.goVersion)
			if err != nil {
				must(err)
			}
			if oldVersion != f.goVersion {
				log("updated go directive from " + oldVersion + ": " + dir)
			}
		}
	}

	prev, err := loadState(w.Root(), "all")
	if err != nil {
		must(err)
	}
	final, err := upgradeModules(w, prev)
	if err != nil {
		must(err)
	}

	log("comparing state")
	diff := diffState(prev, final)

	tests, err := runTests(f, w)
	if err != nil {
		must(err)
	}

	must(generateReport(f, &report{
		Tests: tests,
		State: diff,
	}))
}

func must(err error) {
	if err == nil {
		return
	}
	log(err.Error())
	os.Exit(1)
}

func log(s string) {
	fmt.Fprintln(os.Stderr, "goupdate: "+s)
}
