module go.pact.im/x/crypt

go 1.18

require (
	go.pact.im/x/option v0.0.6
	go.pact.im/x/phcformat v0.0.0-00010101000000-000000000000
	golang.org/x/crypto v0.4.0
)

require golang.org/x/sys v0.3.0 // indirect

replace go.pact.im/x/phcformat => ../phcformat
