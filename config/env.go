package config

import (
	"bufio"
	"io"
	"strings"
)

func parseEnv(r io.Reader) ([]string, error) {
	var env []string

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if i := strings.IndexByte(line, '#'); i >= 0 {
			line = line[:i]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		env = append(env, line)
	}
	return env, scanner.Err()
}
