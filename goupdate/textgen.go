package main

import (
	"io"
	"strings"

	"github.com/kr/text"
	"golang.org/x/exp/apidiff"
	"golang.org/x/tools/go/packages"
)

type textGenerator struct {
	w *sectionIndentWriter
}

func generateText(w io.Writer, r *report) error {
	gen := &textGenerator{newSectionIndentWriter(w, "    ")} // four space indent
	if err := gen.visitReport(r); err != nil {
		return err
	}
	return nil
}

func (g *textGenerator) visitReport(r *report) error {
	if err := g.visitTestReport(r.Tests); err != nil {
		return err
	}
	if err := g.visitStateDiff(r.State); err != nil {
		return err
	}
	return nil
}

func (g *textGenerator) visitTestReport(t *testReport) error {
	if t == nil {
		return nil
	}
	return g.section("TESTS", func() error {
		if err := g.visitTestReportEvents(t.Events, t.Failed); err != nil {
			return err
		}
		if err := g.visitTestReportDecodeError(t.Decode); err != nil {
			return err
		}
		if err := g.visitTestReportBuffer(t.Buffer); err != nil {
			return err
		}
		if err := g.visitTestReportStderr(t.Stderr); err != nil {
			return err
		}
		return nil
	})
}

func (g *textGenerator) visitTestReportEvents(events []testEvent, failed bool) error {
	if len(events) == 0 {
		return nil
	}
	return g.section("OUTPUT", func() error {
		for _, ev := range events {
			if ev.Output == nil {
				continue
			}
			out := *ev.Output
			if !failed && !strings.HasPrefix(out, "ok  \t") && !strings.HasPrefix(out, "?   \t") {
				continue
			}
			if err := g.writeRawString(out); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *textGenerator) visitTestReportDecodeError(decodeError error) error {
	if decodeError == nil {
		return nil
	}
	return g.section("INTERNAL ERROR", func() error {
		return g.writeString(decodeError.Error())
	})
}

func (g *textGenerator) visitTestReportBuffer(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	return g.section("MALFORMED OUTPUT", func() error {
		return g.writeBytes(buf)
	})
}

func (g *textGenerator) visitTestReportStderr(buf []byte) error {
	if len(buf) == 0 {
		return nil
	}
	return g.section("STDERR", func() error {
		return g.writeBytes(buf)
	})
}

func (g *textGenerator) visitStateDiff(s stateDiff) error {
	if len(s.Updated) == 0 && len(s.Removed) == 0 && len(s.Added) == 0 {
		return nil
	}
	return g.section("DEPENDENCIES", func() error {
		if err := g.visitStateDiffRemoved(s.Removed); err != nil {
			return nil
		}
		if err := g.visitStateDiffAdded(s.Added); err != nil {
			return nil
		}
		if err := g.visitStateDiffUpdated(s.Updated); err != nil {
			return nil
		}
		return nil
	})
}

func (g *textGenerator) visitStateDiffRemoved(modules []packages.Module) error {
	return g.sectionStateDiffModules("REMOVED", modules)
}

func (g *textGenerator) visitStateDiffAdded(modules []packages.Module) error {
	return g.sectionStateDiffModules("ADDED", modules)
}

func (g *textGenerator) sectionStateDiffModules(title string, modules []packages.Module) error {
	if len(modules) == 0 {
		return nil
	}
	return g.section(title, func() error {
		for _, m := range modules {
			if err := g.writeString(m.Path + "@" + m.Version); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *textGenerator) visitStateDiffUpdated(updates []moduleDiff) error {
	if len(updates) == 0 {
		return nil
	}
	return g.section("UPDATED", func() error {
		for _, u := range updates {
			if err := g.visitModuleDiff(u); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *textGenerator) visitModuleDiff(d moduleDiff) error {
	return g.section(d.Path+"@"+d.Old.Version+" => "+d.New.Version, func() error {
		if len(d.Removed) == 0 && len(d.Added) == 0 && len(d.Changed) == 0 && len(d.NewBroken) == 0 && len(d.OldBroken) == 0 {
			return g.writeString("no changes")
		}
		if err := g.visitModuleDiffRemoved(d.Removed); err != nil {
			return err
		}
		if err := g.visitModuleDiffAdded(d.Added); err != nil {
			return err
		}
		if err := g.visitModuleDiffChanged(d.Changed); err != nil {
			return err
		}
		if err := g.visitModuleDiffOldBroken(d.OldBroken); err != nil {
			return err
		}
		if err := g.visitModuleDiffNewBroken(d.OldBroken); err != nil {
			return err
		}
		return nil
	})
}

func (g *textGenerator) visitModuleDiffRemoved(pkgs []*packages.Package) error {
	return g.sectionModuleDiffPackages("REMOVED PACKAGES", pkgs)
}

func (g *textGenerator) visitModuleDiffAdded(pkgs []*packages.Package) error {
	return g.sectionModuleDiffPackages("ADDED PACKAGES", pkgs)
}

func (g *textGenerator) sectionModuleDiffPackages(title string, pkgs []*packages.Package) error {
	if len(pkgs) == 0 {
		return nil
	}
	return g.section(title, func() error {
		for _, p := range pkgs {
			if err := g.writeString(p.PkgPath); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *textGenerator) visitModuleDiffChanged(changes []packageDiff) error {
	if len(changes) == 0 {
		return nil
	}
	return g.section("CHANGED", func() error {
		for _, c := range changes {
			if err := g.visitPackageDiff(c); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *textGenerator) visitModuleDiffOldBroken(pkgs []*packages.Package) error {
	return g.sectionBrokenPackages("OLD BROKEN", pkgs)
}

func (g *textGenerator) visitModuleDiffNewBroken(pkgs []*packages.Package) error {
	return g.sectionBrokenPackages("NEW BROKEN", pkgs)
}

func (g *textGenerator) sectionBrokenPackages(title string, pkgs []*packages.Package) error {
	if len(pkgs) == 0 {
		return nil
	}
	return g.section(title, func() error {
		for _, p := range pkgs {
			if err := g.sectionPackageErrors(p.PkgPath, p.Errors); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *textGenerator) sectionPackageErrors(title string, errs []packages.Error) error {
	if len(errs) == 0 {
		return nil
	}
	return g.section(title, func() error {
		for _, e := range errs {
			if err := g.writeString(e.Error()); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *textGenerator) visitPackageDiff(d packageDiff) error {
	if len(d.Diff.Changes) == 0 {
		return nil
	}
	return g.section(d.Path, func() error {
		if err := g.visitPacakgeDiffIncompatible(d.Incompatible()); err != nil {
			return err
		}
		if err := g.visitPacakgeDiffCompatible(d.Compatible()); err != nil {
			return err
		}
		return nil
	})
}

func (g *textGenerator) visitPacakgeDiffIncompatible(changes []apidiff.Change) error {
	return g.sectionPackageDiffChanges("INCOMPATIBLE", changes)
}

func (g *textGenerator) visitPacakgeDiffCompatible(changes []apidiff.Change) error {
	return g.sectionPackageDiffChanges("COMPATIBLE", changes)
}

func (g *textGenerator) sectionPackageDiffChanges(title string, changes []apidiff.Change) error {
	if len(changes) == 0 {
		return nil
	}
	return g.section(title, func() error {
		for _, c := range changes {
			if err := g.writeString(c.Message); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *textGenerator) section(title string, f func() error) error {
	if err := g.writeString(title); err != nil {
		return err
	}
	if err := g.w.Section(f); err != nil {
		return err
	}
	return nil
}

func (g *textGenerator) writeString(s string) error {
	return g.writeBytes([]byte(s))
}

func (g *textGenerator) writeBytes(s []byte) error {
	if _, err := g.w.Write(s); err != nil {
		return err
	}
	if len(s) == 0 || s[len(s)-1] != '\n' {
		_, err := g.w.Write([]byte{'\n'})
		return err
	}
	return nil
}

func (g *textGenerator) writeRawString(s string) error {
	return g.writeRawBytes([]byte(s))
}

func (g *textGenerator) writeRawBytes(s []byte) error {
	if _, err := g.w.Write(s); err != nil {
		return err
	}
	return nil
}

type sectionIndentWriter struct {
	wr     io.Writer
	indent io.Writer
	pre    []byte
}

func newSectionIndentWriter(wr io.Writer, pre string) *sectionIndentWriter {
	return &sectionIndentWriter{
		wr:  wr,
		pre: []byte(pre),
	}
}

func (w *sectionIndentWriter) Section(f func() error) error {
	prev := w.indent
	indent := text.NewIndentWriter(w.writer(), w.pre)
	w.indent = indent
	err := f()
	w.indent = prev
	return err
}

func (w *sectionIndentWriter) Write(p []byte) (int, error) {
	return w.writer().Write(p)
}

func (w *sectionIndentWriter) writer() io.Writer {
	wr := w.indent
	if wr == nil {
		wr = w.wr
	}
	return wr
}
