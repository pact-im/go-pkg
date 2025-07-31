package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type report struct {
	Tests *testReport
	State stateDiff
}

// generateReport generates the report for the requested flags.
func generateReport(f *flags, r *report) error {
	var buf bytes.Buffer
	switch f.outFormat {
	case formatJSON:
		enc := json.NewEncoder(&buf)
		if f.outPath == "" {
			enc.SetIndent("", "\t")
		}
		if err := enc.Encode(r); err != nil {
			return err
		}
	case formatHTML:
		if err := generateHTML(&buf, r); err != nil {
			return err
		}
	case formatText:
		if err := generateText(&buf, r); err != nil {
			return err
		}
	}
	out := bytes.TrimSpace(buf.Bytes())
	if f.outPath == "" {
		if len(out) == 0 {
			return nil
		}
		_, err := fmt.Fprintln(os.Stdout, string(out))
		return err
	}
	return os.WriteFile(f.outPath, out, 0o640)
}
