module go.pact.im/x/zapjournal/tests

go 1.18

require (
	github.com/valyala/fastjson v1.6.3
	go.pact.im/x/zapjournal v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.23.0
	golang.org/x/sys v0.1.0
)

require (
	go.uber.org/atomic v1.10.0 // indirect
	go.uber.org/multierr v1.8.0 // indirect
)

replace go.pact.im/x/zapjournal => ../
