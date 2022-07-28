package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"path/filepath"
	"time"
)

type testReport struct {
	Events []testEvent
	Stderr []byte // test stderr output
	Failed bool   // test failed (non-zero exit code)
	Decode error  // stdout decoding error
	Buffer []byte // remaining stdout buffer if Decode != nil
}

type testEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  *string // optional
}

func runTests(f *flags, w *workspace) (*testReport, error) {
	if f.runTests == runTestsNone {
		return nil, nil
	}

	var r testReport

	base := []string{"test", "-json", "-short"}
	if !f.testShort {
		base = base[:len(base)-1]
	}
	var args []string
	if f.runTests == runTestsWork {
		args = make([]string, len(w.Paths)+len(base))
		_ = copy(args, base)
		for i, moduleDir := range w.Paths {
			args[i+len(base)] = filepath.Join(moduleDir, "...")
		}
	} else {
		args = make([]string, len(base)+1)
		_ = copy(args, base)
		args[len(args)-1] = "all"
	}

	log("running tests")

	var stdout, stderr bytes.Buffer
	c := exec.Command("go", args...)
	c.Dir = w.Root()
	c.Stdout = &stdout
	c.Stderr = &stderr
	var exit *exec.ExitError
	switch err := c.Run(); {
	case errors.As(err, &exit):
		r.Failed = true
	case err != nil:
		return nil, err
	}

	r.Stderr = stderr.Bytes()

	dec := json.NewDecoder(&stdout)
	for {
		var ev testEvent
		err := dec.Decode(&ev)
		if err == io.EOF {
			break
		}
		if err != nil {
			r.Decode = err
			r.Buffer, _ = io.ReadAll(io.MultiReader(dec.Buffered(), &stdout))
			break
		}
		r.Events = append(r.Events, ev)
	}
	return &r, nil
}
