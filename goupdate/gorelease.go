package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html"
	"os"
	"os/exec"
	"strings"
)

func gorelease(out string, updates []update) error {
	dirs := make([]string, len(updates))
	for i, u := range updates {
		buf, err := system("go", "mod", "download", "-json", u.ModulePath+"@"+u.NewVersion)
		if err != nil {
			return err
		}
		var out struct{ Dir string }
		if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
			return err
		}
		dirs[i] = out.Dir
	}

	write := out != ""
	var report string
	for i, u := range updates {
		var buf bytes.Buffer

		c := verboseCommand("gorelease", "-base", u.ModulePath+"@"+u.OldVersion, "-version", u.NewVersion)
		c.Dir = dirs[i]
		if write {
			c.Stdout = &buf
		} else {
			c.Stdout = os.Stdout
		}
		if err := c.Run(); err != nil && !errors.As(err, new(*exec.ExitError)) {
			return err
		}
		if !write {
			continue
		}

		apidiff := buf.Bytes()

		// Drop version suggestion that is in diagnostics section.
		diag := fmt.Sprintf("\nVersion %s is lower than most pseudo-versions. Consider releasing v0.1.0-0 instead.\n", u.NewVersion)
		apidiff = bytes.ReplaceAll(apidiff, []byte(diag), []byte("\n"))
		apidiff = bytes.TrimSpace(apidiff)

		// Keep diagnostic message that at the end of output but after
		// summary section that we remove completely.
		suffix := "Errors were found in the base version. Some API changes may be omitted."
		if suffixBytes := []byte("\n" + suffix); bytes.HasSuffix(apidiff, suffixBytes) {
			apidiff = bytes.TrimSuffix(apidiff, suffixBytes)
			apidiff = bytes.TrimSpace(apidiff)
		} else {
			suffix = ""
		}

		// Drop summary section from gorelease output. Note that it is
		// the last section and contains mostly version suggestions.
		header := "# summary"
		if i := bytes.LastIndex(apidiff, []byte("\n"+header)); i >= 0 {
			apidiff = apidiff[:i]
		} else if bytes.HasPrefix(apidiff, []byte(header)) {
			apidiff = nil
		}
		apidiff = bytes.TrimSpace(apidiff)

		// Trim diagnostics section if it contains only the version
		// suggestion we dropped earlier.
		if header := "# diagnostics"; bytes.Equal(apidiff, []byte(header)) {
			apidiff = nil
		} else {
			apidiff = bytes.TrimSuffix(apidiff, []byte("\n# diagnostics"))
		}
		apidiff = bytes.TrimSpace(apidiff)

		pkggodev := `https://pkg.go.dev/`
		moduleURL := htmlAnchor(pkggodev+u.ModulePath, u.ModulePath)
		oldVersion := htmlAnchor(pkggodev+u.ModulePath+"@"+u.OldVersion, u.OldVersion)
		newVersion := htmlAnchor(pkggodev+u.ModulePath+"@"+u.NewVersion, u.NewVersion)
		var warning string
		if len(apidiff) > 0 || suffix != "" {
			warning = " \u26A0"
		}
		summary := moduleURL + "@" + oldVersion + " \u27F9 " + newVersion + warning
		report += `<details><summary>` + summary + `</summary>`
		if len(apidiff) > 0 {
			indent := "\n\t"
			report += "\n"
			report += indent + strings.ReplaceAll(string(apidiff), "\n", indent)
			report += "\n"
		}
		if suffix != "" {
			report += `<p>` + suffix + `</p>`
		}
		report += `</details>`
	}
	if write {
		return os.WriteFile(out, []byte(report), 0o640)
	}
	return nil
}

func htmlAnchor(href, text string) string {
	href = html.EscapeString(href)
	text = html.EscapeString(text)
	return fmt.Sprintf(`<a href="%s">%s</a>`, href, text)
}
