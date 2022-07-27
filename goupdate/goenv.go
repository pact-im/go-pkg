package main

import (
	"encoding/json"
)

func goenv(vars ...string) (map[string]string, error) {
	buf, err := system("go", append([]string{"env", "-json"}, vars...)...)
	if err != nil {
		return nil, err
	}
	out := make(map[string]string)
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		return nil, err
	}
	return out, nil
}
