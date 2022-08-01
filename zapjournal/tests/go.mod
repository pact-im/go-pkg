module go.pact.im/x/zapjournal/tests

go 1.18

require (
	github.com/valyala/fastjson v1.6.3
	go.pact.im/x/zapjournal v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.21.0
	golang.org/x/sys v0.0.0-20220731174439-a90be440212d
)

require (
	go.uber.org/atomic v1.9.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace go.pact.im/x/zapjournal => ../
