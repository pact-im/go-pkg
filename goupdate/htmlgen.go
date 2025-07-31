package main

import (
	"io"
	"strings"

	"golang.org/x/tools/go/packages"
)

var htmlEscaper = strings.NewReplacer(
	"\n", "&#10;", // for Markdown compatibility
	// From Go’s html package.
	`&`, "&amp;",
	`'`, "&#39;", // "&#39;" is shorter than "&apos;" and apos was not in HTML until HTML5.
	`<`, "&lt;",
	`>`, "&gt;",
	`"`, "&#34;", // "&#34;" is shorter than "&quot;".
)

type htmlGenerator struct {
	w io.Writer

	docsURL string
}

func generateHTML(w io.Writer, r *report) error {
	gen := &htmlGenerator{w, "https://pkg.go.dev/"}
	if err := gen.visitReport(r); err != nil {
		return err
	}
	return nil
}

func (g *htmlGenerator) visitReport(r *report) error {
	if err := g.visitStateDiff(r.State); err != nil {
		return err
	}
	if err := g.visitTestReport(r.Tests); err != nil {
		return err
	}
	return nil
}

func (g *htmlGenerator) visitStateDiff(s stateDiff) error {
	if len(s.Updated) == 0 && len(s.Removed) == 0 && len(s.Added) == 0 {
		return nil
	}
	if err := g.visitStateDiffUpdated(s.Updated); err != nil {
		return err
	}
	if err := g.visitStateDiffAdded(s.Added); err != nil {
		return err
	}
	if err := g.visitStateDiffRemoved(s.Removed); err != nil {
		return err
	}
	return nil
}

func (g *htmlGenerator) visitStateDiffUpdated(updated []moduleDiff) error {
	var unchanged, changed []moduleDiff
	for _, d := range updated {
		if d.Unchanged() {
			unchanged = append(unchanged, d)
		} else {
			changed = append(changed, d)
		}
	}
	for _, d := range changed {
		if err := g.visitModuleDiff(d); err != nil {
			return err
		}
	}
	for _, d := range unchanged {
		if err := g.visitModuleDiff(d); err != nil {
			return err
		}
	}
	return nil
}

func (g *htmlGenerator) visitModuleDiff(d moduleDiff) error {
	return g.htmlDetails(func() error {
		if err := g.moduleUpgradeSummary(d.Path, d.Old.Version, d.New.Version, d.Unchanged()); err != nil {
			return err
		}
		if err := g.visitModuleDiffChanged(d.Changed); err != nil {
			return err
		}
		if err := g.visitModuleDiffRemoved(d.Removed); err != nil {
			return err
		}
		if err := g.visitModuleDiffAdded(d.Removed); err != nil {
			return err
		}
		if err := g.visitModuleDiffOldBroken(d.OldBroken); err != nil {
			return err
		}
		if err := g.visitModuleDiffNewBroken(d.NewBroken); err != nil {
			return err
		}
		return nil
	})
}

func (g *htmlGenerator) moduleUpgradeSummary(modulePath, oldVersion, newVersion string, unchanged bool) error {
	return g.htmlSummary(func() error {
		var prefix string
		if unchanged {
			prefix = "\U0001F7F0 " // heavy equals sign
		} else {
			prefix = "\u2797 " // heavy divide sign
		}
		if err := g.writeString(prefix); err != nil {
			return err
		}
		if err := g.htmlAnchorText(g.docsURL+modulePath, modulePath); err != nil {
			return err
		}
		if err := g.writeString("@"); err != nil {
			return err
		}
		if err := g.htmlAnchorText(g.docsURL+modulePath+"@"+oldVersion, oldVersion); err != nil {
			return err
		}
		if err := g.writeString(" \u27f9 "); err != nil { // long rightwards double arrow
			return err
		}
		if err := g.htmlAnchorText(g.docsURL+modulePath+"@"+newVersion, newVersion); err != nil {
			return err
		}
		return nil
	})
}

func (g *htmlGenerator) visitModuleDiffChanged(changes []packageDiff) error {
	if len(changes) == 0 {
		return nil
	}
	if err := g.htmlHeaderText("Changed packages"); err != nil {
		return err
	}
	if err := g.listPackageDiffs(changes); err != nil {
		return err
	}
	return nil
}

func (g *htmlGenerator) listPackageDiffs(changes []packageDiff) error {
	return g.htmlUnorderedList(func() error {
		for _, d := range changes {
			if err := g.visitPackageDiff(d); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *htmlGenerator) visitPackageDiff(d packageDiff) error {
	return g.htmlListItem(func() error {
		if err := g.htmlCodeText(d.Path); err != nil {
			return err
		}
		if err := g.htmlPreSampText(d.Diff.String()); err != nil {
			return err
		}
		return nil
	})
}

func (g *htmlGenerator) visitModuleDiffAdded(pkgs []*packages.Package) error {
	return g.listPackages("Added packages", pkgs)
}

func (g *htmlGenerator) visitModuleDiffRemoved(pkgs []*packages.Package) error {
	return g.listPackages("Removed packages", pkgs)
}

func (g *htmlGenerator) visitModuleDiffOldBroken(pkgs []*packages.Package) error {
	return g.listPackages("Old package errors", pkgs)
}

func (g *htmlGenerator) visitModuleDiffNewBroken(pkgs []*packages.Package) error {
	return g.listPackages("New package errors", pkgs)
}

func (g *htmlGenerator) listPackages(header string, pkgs []*packages.Package) error {
	if len(pkgs) == 0 {
		return nil
	}
	if err := g.htmlHeaderText(header); err != nil {
		return err
	}
	return g.htmlUnorderedList(func() error {
		for _, pkg := range pkgs {
			if err := g.listItemPackageWithErrors(pkg); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *htmlGenerator) listItemPackageWithErrors(pkg *packages.Package) error {
	return g.htmlListItem(func() error {
		if err := g.htmlCodeText(pkg.PkgPath); err != nil {
			return err
		}
		for _, e := range pkg.Errors {
			if err := g.htmlPreSampText(e.Error()); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *htmlGenerator) visitStateDiffAdded(modules []packages.Module) error {
	return g.listModules("\u2795", modules) // heavy plus sign
}

func (g *htmlGenerator) visitStateDiffRemoved(modules []packages.Module) error {
	return g.listModules("\u2796", modules) // heavy minus sign
}

func (g *htmlGenerator) listModules(prefix string, modules []packages.Module) error {
	for _, m := range modules {
		if err := g.emptyModuleDetails(m, prefix); err != nil {
			return err
		}
	}
	return nil
}

func (g *htmlGenerator) emptyModuleDetails(m packages.Module, prefix string) error {
	return g.htmlDetails(func() error {
		return g.htmlSummary(func() error {
			if err := g.writeString(prefix + " "); err != nil {
				return err
			}
			if err := g.htmlAnchorText(g.docsURL+m.Path, m.Path); err != nil {
				return err
			}
			if err := g.writeString("@"); err != nil {
				return err
			}
			if err := g.htmlAnchorText(g.docsURL+m.Path+"@"+m.Version, m.Version); err != nil {
				return err
			}
			return nil
		})
	})
}

func (g *htmlGenerator) visitTestReport(t *testReport) error {
	if t == nil {
		return nil
	}
	if err := g.htmlHorizontalRule(); err != nil {
		return err
	}
	return g.htmlDetails(func() error {
		if err := g.visitTestReportFailed(t.Failed, t.Decode != nil); err != nil {
			return err
		}
		if err := g.visitTestReportEvents(t.Events); err != nil {
			return err
		}
		if err := g.visitTestReportDecode(t.Decode, t.Buffer); err != nil {
			return err
		}
		if err := g.visitTestReportStderr(t.Stderr); err != nil {
			return err
		}
		return nil
	})
}

func (g *htmlGenerator) visitTestReportFailed(failed, decodeError bool) error {
	var text string
	switch {
	case failed:
		text = "\u274C Tests: failed" // cross mark
	case decodeError:
		text = "\u26A0\uFE0F Tests: passed (malformed output)" // warning sign
	default:
		text = "\u2714\uFE0F Tests: passed" // heavy check mark
	}
	return g.htmlSummaryText(text)
}

func (g *htmlGenerator) visitTestReportEvents(events []testEvent) error {
	return g.htmlPreSamp(func() error {
		for _, ev := range events {
			if ev.Output == nil {
				continue
			}
			if err := g.writeString(*ev.Output); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *htmlGenerator) visitTestReportDecode(decodeError error, buffer []byte) error {
	if decodeError == nil {
		return nil
	}
	if err := g.htmlPreSampText(decodeError.Error()); err != nil {
		return err
	}
	return g.htmlDetails(func() error {
		if err := g.htmlSummaryText("Malformed output"); err != nil {
			return err
		}
		if err := g.writeBytes(buffer); err != nil {
			return err
		}
		return nil
	})
}

func (g *htmlGenerator) visitTestReportStderr(stderr []byte) error {
	if len(stderr) == 0 {
		return nil
	}
	return g.htmlPreSamp(func() error {
		return g.writeBytes(stderr)
	})
}

func (g *htmlGenerator) htmlHeaderText(s string) error {
	return g.html("h4", nil, nil, func() error {
		return g.writeString(s)
	}, true)
}

func (g *htmlGenerator) htmlUnorderedList(body func() error) error {
	return g.html("ul", nil, nil, body, true)
}

func (g *htmlGenerator) htmlListItem(body func() error) error {
	return g.html("li", nil, nil, body, true)
}

func (g *htmlGenerator) htmlDetails(body func() error) error {
	return g.html("details", nil, nil, body, true)
}

func (g *htmlGenerator) htmlSummaryText(s string) error {
	return g.htmlSummary(func() error {
		return g.writeString(s)
	})
}

func (g *htmlGenerator) htmlSummary(body func() error) error {
	return g.html("summary", nil, nil, body, true)
}

func (g *htmlGenerator) htmlPreSampText(s string) error {
	return g.htmlPreSamp(func() error {
		return g.writeString(strings.TrimSuffix(s, "\n"))
	})
}

func (g *htmlGenerator) htmlPreSamp(body func() error) error {
	return g.html("pre", nil, nil, func() error {
		return g.html("samp", nil, nil, body, true)
	}, true)
}

func (g *htmlGenerator) htmlCodeText(s string) error {
	return g.html("code", nil, nil, func() error {
		return g.writeString(s)
	}, true)
}

func (g *htmlGenerator) htmlAnchorText(href, s string) error {
	return g.htmlAnchor(href, func() error {
		return g.writeString(s)
	})
}

func (g *htmlGenerator) htmlAnchor(href string, body func() error) error {
	return g.html("a", []string{"href"}, map[string]string{"href": href}, body, true)
}

func (g *htmlGenerator) htmlHorizontalRule() error {
	return g.html("hr", nil, nil, nil, false)
}

func (g *htmlGenerator) html(tag string, attrNames []string, attrValues map[string]string, body func() error, closeTag bool) error {
	if err := g.writeRawString(`<`); err != nil {
		return err
	}
	if err := g.writeRawString(tag); err != nil {
		return err
	}
	for _, attr := range attrNames {
		if err := g.writeRawString(" "); err != nil {
			return err
		}
		if err := g.writeRawString(attr); err != nil {
			return err
		}
		val, ok := attrValues[attr]
		if !ok {
			continue
		}
		if err := g.writeRawString("="); err != nil {
			return err
		}
		if err := g.writeString(val); err != nil {
			return err
		}
	}
	if err := g.writeRawString(`>`); err != nil {
		return err
	}
	if body != nil {
		if err := body(); err != nil {
			return err
		}
	}
	if !closeTag {
		return nil
	}
	if err := g.writeRawString(`</`); err != nil {
		return err
	}
	if err := g.writeRawString(tag); err != nil {
		return err
	}
	if err := g.writeRawString(`>`); err != nil {
		return err
	}
	return nil
}

func (g *htmlGenerator) writeBytes(s []byte) error {
	return g.writeString(string(s))
}

func (g *htmlGenerator) writeString(s string) error {
	_, err := htmlEscaper.WriteString(g.w, s)
	return err
}

func (g *htmlGenerator) writeRawString(s string) error {
	_, err := io.WriteString(g.w, s)
	return err
}
