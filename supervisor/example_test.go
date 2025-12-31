package supervisor

import (
	"context"

	"go.pact.im/x/process"
)

func ExampleGroupBuilder() {
	ctx := context.Background()

	// Run background processes in Group under the same context as Run.
	builder := NewBuilder(GroupBuilder)
	_ = builder.Run(ctx, func(_ context.Context) error {
		g, _ := builder.Runner()
		g.Go(process.Nop(), nil)
		return nil
	})

	// Output:
}
