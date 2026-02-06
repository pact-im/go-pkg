module go.pact.im/x/zapjournal/tests

go 1.25.0

require (
	github.com/valyala/fastjson v1.6.7
	go.pact.im/x/zapjournal v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.27.1
	golang.org/x/sys v0.40.0
)

require go.uber.org/multierr v1.11.0 // indirect

replace go.pact.im/x/zapjournal => ../
