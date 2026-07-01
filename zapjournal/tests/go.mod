module go.pact.im/x/zapjournal/tests

go 1.26.0

require (
	github.com/valyala/fastjson v1.6.10
	go.pact.im/x/zapjournal v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.28.0
	golang.org/x/sys v0.46.0
)

require go.uber.org/multierr v1.11.0 // indirect

replace go.pact.im/x/zapjournal => ../
