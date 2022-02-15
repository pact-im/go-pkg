module go.pact.im/x/flaky

go 1.18

require (
	github.com/golang/mock v1.6.0
	github.com/robfig/cron/v3 v3.0.1
	go.pact.im/x/clock v0.0.0-00010101000000-000000000000
	gotest.tools/v3 v3.1.0
)

require (
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)

replace go.pact.im/x/clock => ../clock
