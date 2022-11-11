// The pulumi-resource-yandex command is a Pulumi provider for Yandex Cloud.
//
// Note that currently both Pulumi and Terraform protocol definitions are linked
// with the resulting binary, so it must be built with the following flags:
//
//  -ldflags=-X=google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore
//
// See https://developers.google.com/protocol-buffers/docs/reference/go/faq#namespace-conflict
package main

import (
	_ "embed"

	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfbridge"

	"go.pact.im/x/tf2pulumi/internal/yandexprovider"
)

//go:embed schema.json
var pulumiSchema []byte

func main() {
	tfbridge.Main(yandexprovider.Package(), yandexprovider.Version(), yandexprovider.Provider(), pulumiSchema)
}
