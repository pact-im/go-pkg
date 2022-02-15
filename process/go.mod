module go.pact.im/x/process

go 1.18

require (
	go.pact.im/x/clock v0.0.0-00010101000000-000000000000
	go.pact.im/x/syncx v0.0.0-00010101000000-000000000000
	go.uber.org/atomic v1.9.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	gotest.tools/v3 v3.1.0
)

require (
	github.com/google/go-cmp v0.5.5 // indirect
	github.com/pkg/errors v0.9.1 // indirect
)

replace (
	go.pact.im/x/clock => ../clock
	go.pact.im/x/syncx => ../syncx
)
