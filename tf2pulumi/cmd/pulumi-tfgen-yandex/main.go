// The pulumi-tfgen-yandex command implements design-time Pulumi code generator
// for Yandex Cloud.
//
// Note that currently both Pulumi and Terraform protocol definitions are linked
// with the resulting binary, so it must be built with the following flags:
//
//  -ldflags=-X=google.golang.org/protobuf/reflect/protoregistry.conflictPolicy=ignore
//
// See https://developers.google.com/protocol-buffers/docs/reference/go/faq#namespace-conflict
package main

import (
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfgen"

	"go.pact.im/x/tf2pulumi/internal/yandexprovider"
)

func main() {
	tfgen.Main(yandexprovider.Package(), yandexprovider.Version(), yandexprovider.Provider())
}
