// Package inprocshim provides in-process Terraform-to-Pulumi provider shim
// implementation.
package inprocshim

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	shim "github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim/tfplugin5"
	"github.com/pulumi/pulumi-terraform-bridge/v3/pkg/tfshim/tfplugin5/proto"
	"google.golang.org/grpc"
)

var _ proto.ProviderClient = (*providerClient)(nil)

// providerClient is a Pulumi’s proto.ProviderClient adapter for Terraform’s
// schema.GRPCProviderServer.
type providerClient struct {
	s *schema.GRPCProviderServer
}

// Provider returns a new Pulumi provider shim for the given Terraform provider
// implementation.
func Provider(provider *schema.Provider) shim.Provider {
	grpcProvider := schema.NewGRPCProviderServer(provider)
	shimProvider, err := tfplugin5.NewProvider(context.Background(), &providerClient{grpcProvider}, "")
	if err != nil {
		panic(err)
	}
	return shimProvider
}

// GetSchema implements the proto.ProviderClient interface.
func (c *providerClient) GetSchema(ctx context.Context, in *proto.GetProviderSchema_Request, _ ...grpc.CallOption) (*proto.GetProviderSchema_Response, error) {
	return invoke(ctx, in, c.s.GetProviderSchema,
		func(in *proto.GetProviderSchema_Request) (*tfprotov5.GetProviderSchemaRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.GetProviderSchemaRequest{}, nil
		},
		func(resp *tfprotov5.GetProviderSchemaResponse) (*proto.GetProviderSchema_Response, error) {
			if resp == nil {
				return nil, nil
			}
			var err error
			out := &proto.GetProviderSchema_Response{
				Provider:          call(&err, fromSchema, resp.Provider),
				ResourceSchemas:   call(&err, fromSchemaMap, resp.ResourceSchemas),
				DataSourceSchemas: call(&err, fromSchemaMap, resp.DataSourceSchemas),
				Diagnostics:       call(&err, fromDiagnostics, resp.Diagnostics),
				ProviderMeta:      call(&err, fromSchema, resp.ProviderMeta),
				// NB resp.ServerCapabilities is not defined in Pulumi’s proto.
			}
			return out, err
		},
	)
}

// PrepareProviderConfig implements the proto.ProviderClient interface.
func (c *providerClient) PrepareProviderConfig(ctx context.Context, in *proto.PrepareProviderConfig_Request, _ ...grpc.CallOption) (*proto.PrepareProviderConfig_Response, error) {
	return invoke(ctx, in, c.s.PrepareProviderConfig,
		func(in *proto.PrepareProviderConfig_Request) (*tfprotov5.PrepareProviderConfigRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.PrepareProviderConfigRequest{
				Config: toDynamicValue(in.GetConfig()),
			}, nil
		},
		func(resp *tfprotov5.PrepareProviderConfigResponse) (*proto.PrepareProviderConfig_Response, error) {
			if resp == nil {
				return nil, nil
			}
			var err error
			out := &proto.PrepareProviderConfig_Response{
				PreparedConfig: fromDynamicValue(resp.PreparedConfig),
				Diagnostics:    call(&err, fromDiagnostics, resp.Diagnostics),
			}
			return out, err
		},
	)
}

// ValidateResourceTypeConfig implements the proto.ProviderClient interface.
func (c *providerClient) ValidateResourceTypeConfig(ctx context.Context, in *proto.ValidateResourceTypeConfig_Request, _ ...grpc.CallOption) (*proto.ValidateResourceTypeConfig_Response, error) {
	return invoke(ctx, in, c.s.ValidateResourceTypeConfig,
		func(in *proto.ValidateResourceTypeConfig_Request) (*tfprotov5.ValidateResourceTypeConfigRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.ValidateResourceTypeConfigRequest{
				TypeName: in.GetTypeName(),
				Config:   toDynamicValue(in.GetConfig()),
			}, nil
		},
		func(resp *tfprotov5.ValidateResourceTypeConfigResponse) (*proto.ValidateResourceTypeConfig_Response, error) {
			if resp == nil {
				return nil, nil
			}
			var err error
			out := &proto.ValidateResourceTypeConfig_Response{
				Diagnostics: call(&err, fromDiagnostics, resp.Diagnostics),
			}
			return out, err
		},
	)
}

// ValidateDataSourceConfig implements the proto.ProviderClient interface.
func (c *providerClient) ValidateDataSourceConfig(ctx context.Context, in *proto.ValidateDataSourceConfig_Request, _ ...grpc.CallOption) (*proto.ValidateDataSourceConfig_Response, error) {
	return invoke(ctx, in, c.s.ValidateDataSourceConfig,
		func(in *proto.ValidateDataSourceConfig_Request) (*tfprotov5.ValidateDataSourceConfigRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.ValidateDataSourceConfigRequest{
				TypeName: in.GetTypeName(),
				Config:   toDynamicValue(in.GetConfig()),
			}, nil
		},
		func(resp *tfprotov5.ValidateDataSourceConfigResponse) (*proto.ValidateDataSourceConfig_Response, error) {
			if resp == nil {
				return nil, nil
			}
			var err error
			out := &proto.ValidateDataSourceConfig_Response{
				Diagnostics: call(&err, fromDiagnostics, resp.Diagnostics),
			}
			return out, err
		},
	)
}

// UpgradeResourceState implements the proto.ProviderClient interface.
func (c *providerClient) UpgradeResourceState(ctx context.Context, in *proto.UpgradeResourceState_Request, _ ...grpc.CallOption) (*proto.UpgradeResourceState_Response, error) {
	return invoke(ctx, in, c.s.UpgradeResourceState,
		func(in *proto.UpgradeResourceState_Request) (*tfprotov5.UpgradeResourceStateRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.UpgradeResourceStateRequest{
				TypeName: in.GetTypeName(),
				Version:  in.GetVersion(),
				RawState: toRawState(in.GetRawState()),
			}, nil
		},
		func(resp *tfprotov5.UpgradeResourceStateResponse) (*proto.UpgradeResourceState_Response, error) {
			if resp == nil {
				return nil, nil
			}
			var err error
			out := &proto.UpgradeResourceState_Response{
				UpgradedState: fromDynamicValue(resp.UpgradedState),
				Diagnostics:   call(&err, fromDiagnostics, resp.Diagnostics),
			}
			return out, err
		},
	)
}

// Configure implements the proto.ProviderClient interface.
func (c *providerClient) Configure(ctx context.Context, in *proto.Configure_Request, _ ...grpc.CallOption) (*proto.Configure_Response, error) {
	return invoke(ctx, in, c.s.ConfigureProvider,
		func(in *proto.Configure_Request) (*tfprotov5.ConfigureProviderRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.ConfigureProviderRequest{
				TerraformVersion: in.GetTerraformVersion(),
				Config:           toDynamicValue(in.GetConfig()),
			}, nil
		},
		func(resp *tfprotov5.ConfigureProviderResponse) (*proto.Configure_Response, error) {
			if resp == nil {
				return nil, nil
			}
			var err error
			out := &proto.Configure_Response{
				Diagnostics: call(&err, fromDiagnostics, resp.Diagnostics),
			}
			return out, err
		},
	)
}

// ReadResource implements the proto.ProviderClient interface.
func (c *providerClient) ReadResource(ctx context.Context, in *proto.ReadResource_Request, _ ...grpc.CallOption) (*proto.ReadResource_Response, error) {
	return invoke(ctx, in, c.s.ReadResource,
		func(in *proto.ReadResource_Request) (*tfprotov5.ReadResourceRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.ReadResourceRequest{
				TypeName:     in.GetTypeName(),
				CurrentState: toDynamicValue(in.GetCurrentState()),
				Private:      bytesClone(in.GetPrivate()),
				ProviderMeta: toDynamicValue(in.GetProviderMeta()),
			}, nil
		},
		func(resp *tfprotov5.ReadResourceResponse) (*proto.ReadResource_Response, error) {
			if resp == nil {
				return nil, nil
			}
			var err error
			out := &proto.ReadResource_Response{
				NewState:    fromDynamicValue(resp.NewState),
				Diagnostics: call(&err, fromDiagnostics, resp.Diagnostics),
				Private:     bytesClone(resp.Private),
			}
			return out, err
		},
	)
}

// PlanResourceChange implements the proto.ProviderClient interface.
func (c *providerClient) PlanResourceChange(ctx context.Context, in *proto.PlanResourceChange_Request, _ ...grpc.CallOption) (*proto.PlanResourceChange_Response, error) {
	return invoke(ctx, in, c.s.PlanResourceChange,
		func(in *proto.PlanResourceChange_Request) (*tfprotov5.PlanResourceChangeRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.PlanResourceChangeRequest{
				TypeName:         in.GetTypeName(),
				PriorState:       toDynamicValue(in.GetPriorState()),
				ProposedNewState: toDynamicValue(in.GetProposedNewState()),
				Config:           toDynamicValue(in.GetConfig()),
				PriorPrivate:     bytesClone(in.GetPriorPrivate()),
				ProviderMeta:     toDynamicValue(in.GetProviderMeta()),
			}, nil
		},
		func(resp *tfprotov5.PlanResourceChangeResponse) (*proto.PlanResourceChange_Response, error) {
			if resp == nil {
				return nil, nil
			}
			var err error
			out := &proto.PlanResourceChange_Response{
				PlannedState:     fromDynamicValue(resp.PlannedState),
				RequiresReplace:  call(&err, fromAttributePaths, resp.RequiresReplace),
				PlannedPrivate:   bytesClone(resp.PlannedPrivate),
				Diagnostics:      call(&err, fromDiagnostics, resp.Diagnostics),
				LegacyTypeSystem: resp.UnsafeToUseLegacyTypeSystem, //nolint:staticcheck // Contains deprecated comment.
			}
			return out, err
		},
	)
}

// ApplyResourceChange implements the proto.ProviderClient interface.
func (c *providerClient) ApplyResourceChange(ctx context.Context, in *proto.ApplyResourceChange_Request, _ ...grpc.CallOption) (*proto.ApplyResourceChange_Response, error) {
	return invoke(ctx, in, c.s.ApplyResourceChange,
		func(in *proto.ApplyResourceChange_Request) (*tfprotov5.ApplyResourceChangeRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.ApplyResourceChangeRequest{
				TypeName:       in.GetTypeName(),
				PriorState:     toDynamicValue(in.GetPriorState()),
				PlannedState:   toDynamicValue(in.GetPlannedState()),
				Config:         toDynamicValue(in.GetConfig()),
				PlannedPrivate: bytesClone(in.GetPlannedPrivate()),
				ProviderMeta:   toDynamicValue(in.GetProviderMeta()),
			}, nil
		},
		func(resp *tfprotov5.ApplyResourceChangeResponse) (*proto.ApplyResourceChange_Response, error) {
			if resp == nil {
				return nil, nil
			}
			var err error
			out := &proto.ApplyResourceChange_Response{
				NewState:         fromDynamicValue(resp.NewState),
				Private:          bytesClone(resp.Private),
				Diagnostics:      call(&err, fromDiagnostics, resp.Diagnostics),
				LegacyTypeSystem: resp.UnsafeToUseLegacyTypeSystem, //nolint:staticcheck // Contains deprecated comment.
			}
			return out, err
		},
	)
}

// ImportResourceState implements the proto.ProviderClient interface.
func (c *providerClient) ImportResourceState(ctx context.Context, in *proto.ImportResourceState_Request, _ ...grpc.CallOption) (*proto.ImportResourceState_Response, error) {
	return invoke(ctx, in, c.s.ImportResourceState,
		func(in *proto.ImportResourceState_Request) (*tfprotov5.ImportResourceStateRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.ImportResourceStateRequest{
				TypeName: in.GetTypeName(),
				ID:       in.GetId(),
			}, nil
		},
		func(resp *tfprotov5.ImportResourceStateResponse) (*proto.ImportResourceState_Response, error) {
			if resp == nil {
				return nil, nil
			}
			var err error
			out := &proto.ImportResourceState_Response{
				ImportedResources: fromImportedResources(resp.ImportedResources),
				Diagnostics:       call(&err, fromDiagnostics, resp.Diagnostics),
			}
			return out, err
		},
	)
}

// ReadDataSource implements the proto.ProviderClient interface.
func (c *providerClient) ReadDataSource(ctx context.Context, in *proto.ReadDataSource_Request, _ ...grpc.CallOption) (*proto.ReadDataSource_Response, error) {
	return invoke(ctx, in, c.s.ReadDataSource,
		func(in *proto.ReadDataSource_Request) (*tfprotov5.ReadDataSourceRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.ReadDataSourceRequest{
				TypeName:     in.GetTypeName(),
				Config:       toDynamicValue(in.GetConfig()),
				ProviderMeta: toDynamicValue(in.GetProviderMeta()),
			}, nil
		},
		func(resp *tfprotov5.ReadDataSourceResponse) (*proto.ReadDataSource_Response, error) {
			if resp == nil {
				return nil, nil
			}
			var err error
			out := &proto.ReadDataSource_Response{
				State:       fromDynamicValue(resp.State),
				Diagnostics: call(&err, fromDiagnostics, resp.Diagnostics),
			}
			return out, err
		},
	)
}

// Stop implements the proto.ProviderClient interface.
func (c *providerClient) Stop(ctx context.Context, in *proto.Stop_Request, _ ...grpc.CallOption) (*proto.Stop_Response, error) {
	return invoke(ctx, in, c.s.StopProvider,
		func(in *proto.Stop_Request) (*tfprotov5.StopProviderRequest, error) {
			if in == nil {
				return nil, nil
			}
			return &tfprotov5.StopProviderRequest{}, nil
		},
		func(resp *tfprotov5.StopProviderResponse) (*proto.Stop_Response, error) {
			if resp == nil {
				return nil, nil
			}
			return &proto.Stop_Response{
				Error: resp.Error,
			}, nil
		},
	)
}

func fromDynamicValue(in *tfprotov5.DynamicValue) *proto.DynamicValue {
	if in == nil {
		return nil
	}
	return &proto.DynamicValue{
		Msgpack: bytesClone(in.MsgPack),
		Json:    bytesClone(in.JSON),
	}
}

func toDynamicValue(in *proto.DynamicValue) *tfprotov5.DynamicValue {
	if in == nil {
		return nil
	}
	return &tfprotov5.DynamicValue{
		MsgPack: bytesClone(in.GetMsgpack()),
		JSON:    bytesClone(in.GetJson()),
	}
}

func toRawState(in *proto.RawState) *tfprotov5.RawState {
	if in == nil {
		return nil
	}
	return &tfprotov5.RawState{
		JSON:    bytesClone(in.GetJson()),
		Flatmap: cloneFlatmap(in.GetFlatmap()),
	}
}

func fromDiagnostics(in []*tfprotov5.Diagnostic) ([]*proto.Diagnostic, error) {
	if in == nil {
		return nil, nil
	}
	out := make([]*proto.Diagnostic, len(in))
	for i, v := range in {
		var err error
		out[i], err = fromDiagnostic(v)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func fromDiagnostic(in *tfprotov5.Diagnostic) (*proto.Diagnostic, error) {
	if in == nil {
		return nil, nil
	}
	var err error
	out := &proto.Diagnostic{
		Severity:  proto.Diagnostic_Severity(in.Severity),
		Summary:   in.Summary,
		Detail:    in.Detail,
		Attribute: call(&err, fromAttributePath, in.Attribute),
	}
	return out, err
}

func fromImportedResources(in []*tfprotov5.ImportedResource) []*proto.ImportResourceState_ImportedResource {
	if in == nil {
		return nil
	}
	out := make([]*proto.ImportResourceState_ImportedResource, len(in))
	for i, v := range in {
		out[i] = fromImportedResource(v)
	}
	return out
}

func fromImportedResource(in *tfprotov5.ImportedResource) *proto.ImportResourceState_ImportedResource {
	if in == nil {
		return nil
	}
	return &proto.ImportResourceState_ImportedResource{
		TypeName: in.TypeName,
		State:    fromDynamicValue(in.State),
		Private:  bytesClone(in.Private),
	}
}

func fromAttributePaths(in []*tftypes.AttributePath) ([]*proto.AttributePath, error) {
	if in == nil {
		return nil, nil
	}
	out := make([]*proto.AttributePath, len(in))
	for i, v := range in {
		var err error
		out[i], err = fromAttributePath(v)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func fromAttributePath(in *tftypes.AttributePath) (*proto.AttributePath, error) {
	if in == nil {
		return nil, nil
	}
	var err error
	out := &proto.AttributePath{
		Steps: call(&err, fromAttributePathSteps, in.Steps()),
	}
	return out, err
}

func fromAttributePathSteps(in []tftypes.AttributePathStep) ([]*proto.AttributePath_Step, error) {
	if in == nil {
		return nil, nil
	}
	out := make([]*proto.AttributePath_Step, len(in))
	for i, v := range in {
		var err error
		out[i], err = fromAttributePathStep(v)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func fromAttributePathStep(in tftypes.AttributePathStep) (*proto.AttributePath_Step, error) {
	switch in := in.(type) {
	case tftypes.AttributeName:
		return &proto.AttributePath_Step{
			Selector: &proto.AttributePath_Step_AttributeName{
				AttributeName: string(in),
			},
		}, nil
	case tftypes.ElementKeyString:
		return &proto.AttributePath_Step{
			Selector: &proto.AttributePath_Step_ElementKeyString{
				ElementKeyString: string(in),
			},
		}, nil
	case tftypes.ElementKeyInt:
		return &proto.AttributePath_Step{
			Selector: &proto.AttributePath_Step_ElementKeyInt{
				ElementKeyInt: int64(in),
			},
		}, nil
	default:
		return nil, fmt.Errorf("unknown type %T for AttributePathStep", in)
	}
}

func fromSchemaMap(in map[string]*tfprotov5.Schema) (map[string]*proto.Schema, error) {
	if in == nil {
		return nil, nil
	}
	out := make(map[string]*proto.Schema, len(in))
	for k, v := range in {
		var err error
		out[k], err = fromSchema(v)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func fromSchema(in *tfprotov5.Schema) (*proto.Schema, error) {
	if in == nil {
		return nil, nil
	}
	var err error
	out := &proto.Schema{
		Version: in.Version,
		Block:   call(&err, fromSchemaBlock, in.Block),
	}
	return out, err
}

func fromSchemaBlock(in *tfprotov5.SchemaBlock) (*proto.Schema_Block, error) {
	if in == nil {
		return nil, nil
	}
	var err error
	out := &proto.Schema_Block{
		Version:         in.Version,
		Attributes:      call(&err, fromSchemaAttributes, in.Attributes),
		BlockTypes:      call(&err, fromSchemaNestedBlocks, in.BlockTypes),
		Description:     in.Description,
		DescriptionKind: proto.StringKind(in.DescriptionKind),
		Deprecated:      in.Deprecated,
	}
	return out, err
}

func fromSchemaAttributes(in []*tfprotov5.SchemaAttribute) ([]*proto.Schema_Attribute, error) {
	if in == nil {
		return nil, nil
	}
	out := make([]*proto.Schema_Attribute, len(in))
	for i, v := range in {
		var err error
		out[i], err = fromSchemaAttribute(v)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func fromSchemaAttribute(in *tfprotov5.SchemaAttribute) (*proto.Schema_Attribute, error) {
	if in == nil {
		return nil, nil
	}
	var err error
	out := &proto.Schema_Attribute{
		Name:            in.Name,
		Type:            call(&err, ctyType, in.Type),
		Description:     in.Description,
		Required:        in.Required,
		Optional:        in.Optional,
		Computed:        in.Computed,
		Sensitive:       in.Sensitive,
		DescriptionKind: proto.StringKind(in.DescriptionKind),
		Deprecated:      in.Deprecated,
	}
	return out, nil
}

func fromSchemaNestedBlocks(in []*tfprotov5.SchemaNestedBlock) ([]*proto.Schema_NestedBlock, error) {
	if in == nil {
		return nil, nil
	}
	out := make([]*proto.Schema_NestedBlock, len(in))
	for i, v := range in {
		var err error
		out[i], err = fromSchemaNestedBlock(v)
		if err != nil {
			return nil, err
		}
	}
	return out, nil
}

func fromSchemaNestedBlock(in *tfprotov5.SchemaNestedBlock) (*proto.Schema_NestedBlock, error) {
	if in == nil {
		return nil, nil
	}
	var err error
	out := &proto.Schema_NestedBlock{
		TypeName: in.TypeName,
		Block:    call(&err, fromSchemaBlock, in.Block),
		Nesting:  proto.Schema_NestedBlock_NestingMode(in.Nesting),
		MinItems: in.MinItems,
		MaxItems: in.MaxItems,
	}
	return out, err
}

func ctyType(in tftypes.Type) ([]byte, error) {
	switch {
	case in.Is(tftypes.String), in.Is(tftypes.Bool), in.Is(tftypes.Number),
		in.Is(tftypes.List{}), in.Is(tftypes.Map{}),
		in.Is(tftypes.Set{}), in.Is(tftypes.Object{}),
		in.Is(tftypes.Tuple{}), in.Is(tftypes.DynamicPseudoType):
		return in.MarshalJSON() //nolint:staticcheck
	}
	return nil, fmt.Errorf("unknown type %s", in)
}

func cloneFlatmap(in map[string]string) map[string]string {
	if in == nil {
		return nil
	}
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

// bytesClone is a bytes.Clone function from Go 1.20.
//
// See https://github.com/golang/go/issues/45038
func bytesClone(b []byte) []byte {
	if b == nil {
		return nil
	}
	return append([]byte{}, b...)
}

func invoke[T, U, V, W any](ctx context.Context, in T, f func(context.Context, U) (V, error), encode func(T) (U, error), decode func(V) (W, error)) (W, error) {
	var zero W

	req, err := encode(in)
	if err != nil {
		return zero, err
	}

	resp, err := f(ctx, req)
	if err != nil {
		return zero, err
	}

	out, err := decode(resp)
	if err != nil {
		return zero, err
	}

	return out, nil
}

func call[T, U any](e *error, f func(T) (U, error), in T) U {
	var zero U

	if *e != nil {
		return zero
	}

	out, err := f(in)
	if err != nil {
		*e = err
		return zero
	}

	return out
}
