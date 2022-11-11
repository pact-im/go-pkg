// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Creates a backend group in the specified folder and adds the specified backends to it.
// For more information, see [the official documentation](https://cloud.yandex.com/en/docs/application-load-balancer/concepts/backend-group).
//
// ## Example Usage
//
// ```go
// package main
//
// import (
// 	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
// 	"go.pact.im/x/tf2pulumi/yandex"
// )
//
// func main() {
// 	pulumi.Run(func(ctx *pulumi.Context) error {
// 		_, err := yandex.NewAlbBackendGroup(ctx, "test-backend-group", &yandex.AlbBackendGroupArgs{
// 			HttpBackends: AlbBackendGroupHttpBackendArray{
// 				&AlbBackendGroupHttpBackendArgs{
// 					Healthcheck: &AlbBackendGroupHttpBackendHealthcheckArgs{
// 						HttpHealthcheck: &AlbBackendGroupHttpBackendHealthcheckHttpHealthcheckArgs{
// 							Path: pulumi.String("/"),
// 						},
// 						Interval: pulumi.String("1s"),
// 						Timeout:  pulumi.String("1s"),
// 					},
// 					Http2: pulumi.Bool(true),
// 					LoadBalancingConfig: &AlbBackendGroupHttpBackendLoadBalancingConfigArgs{
// 						PanicThreshold: pulumi.Float64(50),
// 					},
// 					Name: pulumi.String("test-http-backend"),
// 					Port: pulumi.Float64(8080),
// 					TargetGroupIds: pulumi.StringArray{
// 						pulumi.Any(yandex_alb_target_group.TestTargetGroup.Id),
// 					},
// 					Tls: &AlbBackendGroupHttpBackendTlsArgs{
// 						Sni: pulumi.String("backend-domain.internal"),
// 					},
// 					Weight: pulumi.Float64(1),
// 				},
// 			},
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
//
// ## Import
//
// A backend group can be imported using the `id` of the resource, e.g.
//
// ```sh
//  $ pulumi import yandex:index/albBackendGroup:AlbBackendGroup default backend_group_id
// ```
type AlbBackendGroup struct {
	pulumi.CustomResourceState

	// The backend group creation timestamp.
	CreatedAt pulumi.StringOutput `pulumi:"createdAt"`
	// Description of the backend group.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// Folder that the resource belongs to. If value is omitted, the default provider folder is used.
	FolderId pulumi.StringOutput `pulumi:"folderId"`
	// Grpc backend specification that will be used by the ALB Backend Group. Structure is documented below.
	GrpcBackends AlbBackendGroupGrpcBackendArrayOutput `pulumi:"grpcBackends"`
	// Http backend specification that will be used by the ALB Backend Group. Structure is documented below.
	HttpBackends AlbBackendGroupHttpBackendArrayOutput `pulumi:"httpBackends"`
	// Labels to assign to this backend group.
	Labels pulumi.StringMapOutput `pulumi:"labels"`
	// Name of the backend.
	Name            pulumi.StringOutput                     `pulumi:"name"`
	SessionAffinity AlbBackendGroupSessionAffinityPtrOutput `pulumi:"sessionAffinity"`
	// Stream backend specification that will be used by the ALB Backend Group. Structure is documented below.
	StreamBackends AlbBackendGroupStreamBackendArrayOutput `pulumi:"streamBackends"`
	Timeouts       AlbBackendGroupTimeoutsOutput           `pulumi:"timeouts"`
}

// NewAlbBackendGroup registers a new resource with the given unique name, arguments, and options.
func NewAlbBackendGroup(ctx *pulumi.Context,
	name string, args *AlbBackendGroupArgs, opts ...pulumi.ResourceOption) (*AlbBackendGroup, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	var resource AlbBackendGroup
	err := ctx.RegisterResource("yandex:index/albBackendGroup:AlbBackendGroup", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetAlbBackendGroup gets an existing AlbBackendGroup resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetAlbBackendGroup(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *AlbBackendGroupState, opts ...pulumi.ResourceOption) (*AlbBackendGroup, error) {
	var resource AlbBackendGroup
	err := ctx.ReadResource("yandex:index/albBackendGroup:AlbBackendGroup", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering AlbBackendGroup resources.
type albBackendGroupState struct {
	// The backend group creation timestamp.
	CreatedAt *string `pulumi:"createdAt"`
	// Description of the backend group.
	Description *string `pulumi:"description"`
	// Folder that the resource belongs to. If value is omitted, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// Grpc backend specification that will be used by the ALB Backend Group. Structure is documented below.
	GrpcBackends []AlbBackendGroupGrpcBackend `pulumi:"grpcBackends"`
	// Http backend specification that will be used by the ALB Backend Group. Structure is documented below.
	HttpBackends []AlbBackendGroupHttpBackend `pulumi:"httpBackends"`
	// Labels to assign to this backend group.
	Labels map[string]string `pulumi:"labels"`
	// Name of the backend.
	Name            *string                         `pulumi:"name"`
	SessionAffinity *AlbBackendGroupSessionAffinity `pulumi:"sessionAffinity"`
	// Stream backend specification that will be used by the ALB Backend Group. Structure is documented below.
	StreamBackends []AlbBackendGroupStreamBackend `pulumi:"streamBackends"`
	Timeouts       *AlbBackendGroupTimeouts       `pulumi:"timeouts"`
}

type AlbBackendGroupState struct {
	// The backend group creation timestamp.
	CreatedAt pulumi.StringPtrInput
	// Description of the backend group.
	Description pulumi.StringPtrInput
	// Folder that the resource belongs to. If value is omitted, the default provider folder is used.
	FolderId pulumi.StringPtrInput
	// Grpc backend specification that will be used by the ALB Backend Group. Structure is documented below.
	GrpcBackends AlbBackendGroupGrpcBackendArrayInput
	// Http backend specification that will be used by the ALB Backend Group. Structure is documented below.
	HttpBackends AlbBackendGroupHttpBackendArrayInput
	// Labels to assign to this backend group.
	Labels pulumi.StringMapInput
	// Name of the backend.
	Name            pulumi.StringPtrInput
	SessionAffinity AlbBackendGroupSessionAffinityPtrInput
	// Stream backend specification that will be used by the ALB Backend Group. Structure is documented below.
	StreamBackends AlbBackendGroupStreamBackendArrayInput
	Timeouts       AlbBackendGroupTimeoutsPtrInput
}

func (AlbBackendGroupState) ElementType() reflect.Type {
	return reflect.TypeOf((*albBackendGroupState)(nil)).Elem()
}

type albBackendGroupArgs struct {
	// Description of the backend group.
	Description *string `pulumi:"description"`
	// Folder that the resource belongs to. If value is omitted, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// Grpc backend specification that will be used by the ALB Backend Group. Structure is documented below.
	GrpcBackends []AlbBackendGroupGrpcBackend `pulumi:"grpcBackends"`
	// Http backend specification that will be used by the ALB Backend Group. Structure is documented below.
	HttpBackends []AlbBackendGroupHttpBackend `pulumi:"httpBackends"`
	// Labels to assign to this backend group.
	Labels map[string]string `pulumi:"labels"`
	// Name of the backend.
	Name            *string                         `pulumi:"name"`
	SessionAffinity *AlbBackendGroupSessionAffinity `pulumi:"sessionAffinity"`
	// Stream backend specification that will be used by the ALB Backend Group. Structure is documented below.
	StreamBackends []AlbBackendGroupStreamBackend `pulumi:"streamBackends"`
	Timeouts       AlbBackendGroupTimeouts        `pulumi:"timeouts"`
}

// The set of arguments for constructing a AlbBackendGroup resource.
type AlbBackendGroupArgs struct {
	// Description of the backend group.
	Description pulumi.StringPtrInput
	// Folder that the resource belongs to. If value is omitted, the default provider folder is used.
	FolderId pulumi.StringPtrInput
	// Grpc backend specification that will be used by the ALB Backend Group. Structure is documented below.
	GrpcBackends AlbBackendGroupGrpcBackendArrayInput
	// Http backend specification that will be used by the ALB Backend Group. Structure is documented below.
	HttpBackends AlbBackendGroupHttpBackendArrayInput
	// Labels to assign to this backend group.
	Labels pulumi.StringMapInput
	// Name of the backend.
	Name            pulumi.StringPtrInput
	SessionAffinity AlbBackendGroupSessionAffinityPtrInput
	// Stream backend specification that will be used by the ALB Backend Group. Structure is documented below.
	StreamBackends AlbBackendGroupStreamBackendArrayInput
	Timeouts       AlbBackendGroupTimeoutsInput
}

func (AlbBackendGroupArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*albBackendGroupArgs)(nil)).Elem()
}

type AlbBackendGroupInput interface {
	pulumi.Input

	ToAlbBackendGroupOutput() AlbBackendGroupOutput
	ToAlbBackendGroupOutputWithContext(ctx context.Context) AlbBackendGroupOutput
}

func (*AlbBackendGroup) ElementType() reflect.Type {
	return reflect.TypeOf((**AlbBackendGroup)(nil)).Elem()
}

func (i *AlbBackendGroup) ToAlbBackendGroupOutput() AlbBackendGroupOutput {
	return i.ToAlbBackendGroupOutputWithContext(context.Background())
}

func (i *AlbBackendGroup) ToAlbBackendGroupOutputWithContext(ctx context.Context) AlbBackendGroupOutput {
	return pulumi.ToOutputWithContext(ctx, i).(AlbBackendGroupOutput)
}

// AlbBackendGroupArrayInput is an input type that accepts AlbBackendGroupArray and AlbBackendGroupArrayOutput values.
// You can construct a concrete instance of `AlbBackendGroupArrayInput` via:
//
//          AlbBackendGroupArray{ AlbBackendGroupArgs{...} }
type AlbBackendGroupArrayInput interface {
	pulumi.Input

	ToAlbBackendGroupArrayOutput() AlbBackendGroupArrayOutput
	ToAlbBackendGroupArrayOutputWithContext(context.Context) AlbBackendGroupArrayOutput
}

type AlbBackendGroupArray []AlbBackendGroupInput

func (AlbBackendGroupArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*AlbBackendGroup)(nil)).Elem()
}

func (i AlbBackendGroupArray) ToAlbBackendGroupArrayOutput() AlbBackendGroupArrayOutput {
	return i.ToAlbBackendGroupArrayOutputWithContext(context.Background())
}

func (i AlbBackendGroupArray) ToAlbBackendGroupArrayOutputWithContext(ctx context.Context) AlbBackendGroupArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(AlbBackendGroupArrayOutput)
}

// AlbBackendGroupMapInput is an input type that accepts AlbBackendGroupMap and AlbBackendGroupMapOutput values.
// You can construct a concrete instance of `AlbBackendGroupMapInput` via:
//
//          AlbBackendGroupMap{ "key": AlbBackendGroupArgs{...} }
type AlbBackendGroupMapInput interface {
	pulumi.Input

	ToAlbBackendGroupMapOutput() AlbBackendGroupMapOutput
	ToAlbBackendGroupMapOutputWithContext(context.Context) AlbBackendGroupMapOutput
}

type AlbBackendGroupMap map[string]AlbBackendGroupInput

func (AlbBackendGroupMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*AlbBackendGroup)(nil)).Elem()
}

func (i AlbBackendGroupMap) ToAlbBackendGroupMapOutput() AlbBackendGroupMapOutput {
	return i.ToAlbBackendGroupMapOutputWithContext(context.Background())
}

func (i AlbBackendGroupMap) ToAlbBackendGroupMapOutputWithContext(ctx context.Context) AlbBackendGroupMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(AlbBackendGroupMapOutput)
}

type AlbBackendGroupOutput struct{ *pulumi.OutputState }

func (AlbBackendGroupOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**AlbBackendGroup)(nil)).Elem()
}

func (o AlbBackendGroupOutput) ToAlbBackendGroupOutput() AlbBackendGroupOutput {
	return o
}

func (o AlbBackendGroupOutput) ToAlbBackendGroupOutputWithContext(ctx context.Context) AlbBackendGroupOutput {
	return o
}

// The backend group creation timestamp.
func (o AlbBackendGroupOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v *AlbBackendGroup) pulumi.StringOutput { return v.CreatedAt }).(pulumi.StringOutput)
}

// Description of the backend group.
func (o AlbBackendGroupOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *AlbBackendGroup) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// Folder that the resource belongs to. If value is omitted, the default provider folder is used.
func (o AlbBackendGroupOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v *AlbBackendGroup) pulumi.StringOutput { return v.FolderId }).(pulumi.StringOutput)
}

// Grpc backend specification that will be used by the ALB Backend Group. Structure is documented below.
func (o AlbBackendGroupOutput) GrpcBackends() AlbBackendGroupGrpcBackendArrayOutput {
	return o.ApplyT(func(v *AlbBackendGroup) AlbBackendGroupGrpcBackendArrayOutput { return v.GrpcBackends }).(AlbBackendGroupGrpcBackendArrayOutput)
}

// Http backend specification that will be used by the ALB Backend Group. Structure is documented below.
func (o AlbBackendGroupOutput) HttpBackends() AlbBackendGroupHttpBackendArrayOutput {
	return o.ApplyT(func(v *AlbBackendGroup) AlbBackendGroupHttpBackendArrayOutput { return v.HttpBackends }).(AlbBackendGroupHttpBackendArrayOutput)
}

// Labels to assign to this backend group.
func (o AlbBackendGroupOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v *AlbBackendGroup) pulumi.StringMapOutput { return v.Labels }).(pulumi.StringMapOutput)
}

// Name of the backend.
func (o AlbBackendGroupOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *AlbBackendGroup) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

func (o AlbBackendGroupOutput) SessionAffinity() AlbBackendGroupSessionAffinityPtrOutput {
	return o.ApplyT(func(v *AlbBackendGroup) AlbBackendGroupSessionAffinityPtrOutput { return v.SessionAffinity }).(AlbBackendGroupSessionAffinityPtrOutput)
}

// Stream backend specification that will be used by the ALB Backend Group. Structure is documented below.
func (o AlbBackendGroupOutput) StreamBackends() AlbBackendGroupStreamBackendArrayOutput {
	return o.ApplyT(func(v *AlbBackendGroup) AlbBackendGroupStreamBackendArrayOutput { return v.StreamBackends }).(AlbBackendGroupStreamBackendArrayOutput)
}

func (o AlbBackendGroupOutput) Timeouts() AlbBackendGroupTimeoutsOutput {
	return o.ApplyT(func(v *AlbBackendGroup) AlbBackendGroupTimeoutsOutput { return v.Timeouts }).(AlbBackendGroupTimeoutsOutput)
}

type AlbBackendGroupArrayOutput struct{ *pulumi.OutputState }

func (AlbBackendGroupArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*AlbBackendGroup)(nil)).Elem()
}

func (o AlbBackendGroupArrayOutput) ToAlbBackendGroupArrayOutput() AlbBackendGroupArrayOutput {
	return o
}

func (o AlbBackendGroupArrayOutput) ToAlbBackendGroupArrayOutputWithContext(ctx context.Context) AlbBackendGroupArrayOutput {
	return o
}

func (o AlbBackendGroupArrayOutput) Index(i pulumi.IntInput) AlbBackendGroupOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *AlbBackendGroup {
		return vs[0].([]*AlbBackendGroup)[vs[1].(int)]
	}).(AlbBackendGroupOutput)
}

type AlbBackendGroupMapOutput struct{ *pulumi.OutputState }

func (AlbBackendGroupMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*AlbBackendGroup)(nil)).Elem()
}

func (o AlbBackendGroupMapOutput) ToAlbBackendGroupMapOutput() AlbBackendGroupMapOutput {
	return o
}

func (o AlbBackendGroupMapOutput) ToAlbBackendGroupMapOutputWithContext(ctx context.Context) AlbBackendGroupMapOutput {
	return o
}

func (o AlbBackendGroupMapOutput) MapIndex(k pulumi.StringInput) AlbBackendGroupOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *AlbBackendGroup {
		return vs[0].(map[string]*AlbBackendGroup)[vs[1].(string)]
	}).(AlbBackendGroupOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*AlbBackendGroupInput)(nil)).Elem(), &AlbBackendGroup{})
	pulumi.RegisterInputType(reflect.TypeOf((*AlbBackendGroupArrayInput)(nil)).Elem(), AlbBackendGroupArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*AlbBackendGroupMapInput)(nil)).Elem(), AlbBackendGroupMap{})
	pulumi.RegisterOutputType(AlbBackendGroupOutput{})
	pulumi.RegisterOutputType(AlbBackendGroupArrayOutput{})
	pulumi.RegisterOutputType(AlbBackendGroupMapOutput{})
}