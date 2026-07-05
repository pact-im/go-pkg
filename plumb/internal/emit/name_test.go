package emit

import "testing"

func TestLowerCamel(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"empty", "", ""},
		{"already lower", "user", "user"},
		{"single word", "User", "user"},
		{"single capital", "X", "x"},
		{"all-caps acronym lowered whole", "DB", "db"},
		{"two-letter initialism", "ID", "id"},
		{"unknown acronym then word", "DBConfig", "dbConfig"},
		{"acronym with a trailing digit", "DB2", "db2"},
		{"leading initialism then word", "HTTPServer", "httpServer"},
		{"two initialisms keep the second", "XMLHTTPRequest", "xmlHTTPRequest"},
		{"longer initialism wins", "HTTPSConfig", "httpsConfig"},
		{"initialism is the whole name", "IP", "ip"},
		{"initialism then word", "DNSResolver", "dnsResolver"},
		{"initialism with a trailing digit", "UTF8Decoder", "utf8Decoder"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := lowerCamel(c.in); got != c.want {
				t.Errorf("lowerCamel(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
