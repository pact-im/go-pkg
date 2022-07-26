package tests

import (
	"encoding/json"
	"os"
	"os/exec"
	"strconv"
	"time"

	"github.com/valyala/fastjson"
)

const journalctlTimeout = time.Second

func lastJournalMessage(pid int, message string) (*fastjson.Object, error) {
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}
	defer func() { _ = r.Close() }()
	defer func() { _ = w.Close() }()

	deadline := time.Now().Add(journalctlTimeout)
	if err := r.SetReadDeadline(deadline); err != nil {
		return nil, err
	}

	// Journal adds entries asynchronously from the sending process. We
	// assume that there is at most one entry with the given message per
	// process and follow the log until the first entry.
	c := exec.Command("journalctl",
		"-a", "-f",
		"-o", "json",
		"-n", "1",
		"_PID="+strconv.Itoa(pid),
		"MESSAGE="+message,
	)
	c.Stdout = w
	if err := c.Start(); err != nil {
		return nil, err
	}
	defer func() {
		_ = c.Process.Kill()
		_ = c.Wait()
	}()

	var m json.RawMessage
	if err := json.NewDecoder(r).Decode(&m); err != nil {
		return nil, err
	}

	v, err := fastjson.ParseBytes(m)
	if err != nil {
		return nil, err
	}
	return v.Object()
}
