// Goupdate is a tool for updating Go module dependencies. It is supports Go
// workspaces and is an alternative to running go get -u that attempts to
// produce cleaner results along with the apidiff reports. In particular, it
// checks for newer module versions in the current build graph using go list,
// and so it does not update modules that are present in the Go workspace.
//
// Note that it loops until there are no updates available in the build graph.
// This is inefficient but does the job good enough even compared to go get.
//
// This tool updates dependencies in-place and may leave workspace in broken
// state on failure. Use with extreme caution if there are unsaved changes. It
// is recommended to run goupdate as part of the scheduled CI pipeline instead.
//
// In fact, goupdate is arguably superior to the GitHub’s dependabot since the
// latter creates pull request per dependency and does not support modules with
// pseudo-versions like the golang.org/x/crypto that does not publish semver
// releases (see https://github.com/dependabot/dependabot-core/issues/3017). It
// also does not support Go workspaces without manual configuration management.
// Integrating goupdate with GitHub Actions or other CI solutions that support
// scheduled pipelines should be straightforward and requires Go and gorelease
// tool to be installed (see golang.org/x/exp/cmd/gorelease).
//
// Goupdate uses gorelease, if installed, to generate Markdown reports with
// mixed HTML for expandable details section that contains module API changes.
package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/alessio/shellescape"
)

var errStop = errors.New("stop")

type update struct {
	ModulePath string
	OldVersion string
	NewVersion string
}

type state struct {
	// gowork is either a Go workspace directory or an empty string.
	gowork string
	// workspace is a list of module disk path in the current workspace.
	workspace []string
	// updates is a set of updated modules. The key is ModulePath for
	// lookups.
	updates map[string]update
}

func main() {
	var (
		chdir     string
		limit     uint
		outPath   string
		goVersion string
	)
	flag.StringVar(&chdir, "chdir", "", "change working directory")
	flag.UintVar(&limit, "limit", 30, "loop iterations limit; set to 0 for unlimited")
	flag.StringVar(&outPath, "o", "", "report output path")
	flag.StringVar(&goVersion, "go", "", "update Go version")
	flag.Parse()

	if chdir != "" {
		must(os.Chdir(chdir))
	}

	if _, err := exec.LookPath("go"); err != nil {
		must(err)
	}

	var apidiff bool
	if _, err := exec.LookPath("gorelease"); err != nil && !errors.Is(err, exec.ErrNotFound) {
		must(err)
	}

	if !apidiff && outPath != "" {
		must(fmt.Errorf("fatal: report output was requested (-o=%q flag) but gorelease tool was not found not in PATH", outPath))
	}

	s := state{
		updates: make(map[string]update),
	}
	for i := uint(0); ; i++ {
		log("iteration " + strconv.FormatUint(uint64(i)+1, 10))
		err := run(goVersion, &s)
		if err == errStop {
			break
		}
		if limit != 0 && i > limit {
			must(errors.New("loop iteration limit exceeded"))
		}
		must(err)
	}

	if len(s.updates) == 0 {
		log("all up to date")
		return
	}

	maxPath := 0
	updates := make([]update, 0, len(s.updates))
	for _, u := range s.updates {
		updates = append(updates, u)
		if n := len(u.ModulePath); maxPath < n {
			maxPath = n
		}
	}
	sort.Slice(updates, func(i, j int) bool {
		return updates[i].ModulePath < updates[j].ModulePath
	})

	if apidiff {
		must(gorelease(outPath, updates))
	}

	for _, u := range updates {
		sep := strings.Repeat(" ", maxPath-len(u.ModulePath)+1)
		log(fmt.Sprintf("%s%s[%s => %s]", u.ModulePath, sep, u.OldVersion, u.NewVersion))
	}
}

func run(goVersion string, s *state) error {
	// Load initial workspace state. Note that gowork returns at least one
	// path on success.
	if len(s.workspace) == 0 {
		var err error
		s.gowork, s.workspace, err = gowork()
		if err != nil {
			return err
		}

		if goVersion != "" {
			for _, dir := range s.workspace {
				c := verboseCommand("go", "mod", "edit", "-go", goVersion)
				c.Dir = dir
				if err := c.Run(); err != nil {
					return err
				}
				if err := gomodtidy(dir); err != nil {
					return err
				}
			}
		}
	}

	deps, err := golist()
	if err != nil {
		return err
	}

	seen := make(map[string]bool, len(deps))
	var updates []goModule
	for _, m := range deps {
		seen[m.Path] = true
		if m.Update == nil {
			continue
		}

		// A dirty hack to skip unused modules.
		buf, err := system("go", "mod", "why", "-m", m.Path)
		if err != nil {
			return err
		}
		if bytes.Contains(buf.Bytes(), []byte("main module does not need module")) {
			continue
		}

		oldVersion := m.Version
		if u, ok := s.updates[m.Path]; ok {
			oldVersion = u.OldVersion
		}
		s.updates[m.Path] = update{
			ModulePath: m.Path,
			OldVersion: oldVersion,
			NewVersion: m.Update.Version,
		}

		updates = append(updates, m)
	}

	// Remove modules from updated set if they are no longer in the build
	// graph.
	for path := range s.updates {
		if seen[path] {
			continue
		}
		delete(s.updates, path)
	}

	// Stop the loop if there are no updates to apply. As a finishing touch,
	// run go work sync if we are in a Go workspace.
	if len(updates) == 0 {
		if s.gowork != "" {
			err := os.Remove(filepath.Join(s.gowork, "go.work.sum"))
			if err != nil && !os.IsNotExist(err) {
				return err
			}
			c := verboseCommand("go", "work", "sync")
			c.Dir = s.gowork
			if err := c.Run(); err != nil {
				return err
			}
		}
		return errStop
	}

	for _, dir := range s.workspace {
		log(dir)
		for _, m := range updates {
			c := verboseCommand("go", "mod", "edit", "-require", m.Path+"@"+m.Update.Version)
			c.Dir = dir
			if err := c.Run(); err != nil {
				return err
			}
		}
		if err := gomodtidy(dir); err != nil {
			return err
		}
	}
	return nil
}

func gomodtidy(dir string) error {
	c := verboseCommand("go", "mod", "tidy")
	c.Dir = dir
	return c.Run()
}

func system(name string, args ...string) (*bytes.Buffer, error) {
	var buf bytes.Buffer
	c := verboseCommand(name, args...)
	c.Stdout = &buf
	if err := c.Run(); err != nil {
		return nil, err
	}
	return &buf, nil
}

func verboseCommand(name string, args ...string) *exec.Cmd {
	fmt.Fprintln(os.Stderr, "$ "+name+" "+shellescape.QuoteCommand(args))
	c := exec.Command(name, args...)
	c.Stderr = os.Stderr
	return c
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
