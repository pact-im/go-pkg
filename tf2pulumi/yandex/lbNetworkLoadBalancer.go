// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Creates a network load balancer in the specified folder using the data specified in the config.
// For more information, see [the official documentation](https://cloud.yandex.com/docs/load-balancer/concepts).
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
// 		_, err := yandex.NewLbNetworkLoadBalancer(ctx, "foo", &yandex.LbNetworkLoadBalancerArgs{
// 			AttachedTargetGroups: LbNetworkLoadBalancerAttachedTargetGroupArray{
// 				&LbNetworkLoadBalancerAttachedTargetGroupArgs{
// 					Healthchecks: LbNetworkLoadBalancerAttachedTargetGroupHealthcheckArray{
// 						&LbNetworkLoadBalancerAttachedTargetGroupHealthcheckArgs{
// 							HttpOptions: &LbNetworkLoadBalancerAttachedTargetGroupHealthcheckHttpOptionsArgs{
// 								Path: pulumi.String("/ping"),
// 								Port: pulumi.Float64(8080),
// 							},
// 							Name: pulumi.String("http"),
// 						},
// 					},
// 					TargetGroupId: pulumi.Any(yandex_lb_target_group.MyTargetGroup.Id),
// 				},
// 			},
// 			Listeners: LbNetworkLoadBalancerListenerArray{
// 				&LbNetworkLoadBalancerListenerArgs{
// 					ExternalAddressSpec: &LbNetworkLoadBalancerListenerExternalAddressSpecArgs{
// 						IpVersion: pulumi.String("ipv4"),
// 					},
// 					Name: pulumi.String("my-listener"),
// 					Port: pulumi.Float64(8080),
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
// A network load balancer can be imported using the `id` of the resource, e.g.
//
// ```sh
//  $ pulumi import yandex:index/lbNetworkLoadBalancer:LbNetworkLoadBalancer default network_load_balancer_id
// ```
type LbNetworkLoadBalancer struct {
	pulumi.CustomResourceState

	// An AttachedTargetGroup resource. The structure is documented below.
	AttachedTargetGroups LbNetworkLoadBalancerAttachedTargetGroupArrayOutput `pulumi:"attachedTargetGroups"`
	// The network load balancer creation timestamp.
	CreatedAt pulumi.StringOutput `pulumi:"createdAt"`
	// An optional description of the network load balancer. Provide this property when
	// you create the resource.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// The ID of the folder to which the resource belongs.
	// If omitted, the provider folder is used.
	FolderId pulumi.StringOutput `pulumi:"folderId"`
	// Labels to assign to this network load balancer. A list of key/value pairs.
	Labels pulumi.StringMapOutput `pulumi:"labels"`
	// Listener specification that will be used by a network load balancer. The structure is documented below.
	Listeners LbNetworkLoadBalancerListenerArrayOutput `pulumi:"listeners"`
	// Name of the listener. The name must be unique for each listener on a single load balancer.
	Name pulumi.StringOutput `pulumi:"name"`
	// ID of the availability zone where the network load balancer resides.
	// The default is 'ru-central1'.
	RegionId pulumi.StringPtrOutput              `pulumi:"regionId"`
	Timeouts LbNetworkLoadBalancerTimeoutsOutput `pulumi:"timeouts"`
	// Type of the network load balancer. Must be one of 'external' or 'internal'. The default is 'external'.
	Type pulumi.StringPtrOutput `pulumi:"type"`
}

// NewLbNetworkLoadBalancer registers a new resource with the given unique name, arguments, and options.
func NewLbNetworkLoadBalancer(ctx *pulumi.Context,
	name string, args *LbNetworkLoadBalancerArgs, opts ...pulumi.ResourceOption) (*LbNetworkLoadBalancer, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	var resource LbNetworkLoadBalancer
	err := ctx.RegisterResource("yandex:index/lbNetworkLoadBalancer:LbNetworkLoadBalancer", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetLbNetworkLoadBalancer gets an existing LbNetworkLoadBalancer resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetLbNetworkLoadBalancer(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *LbNetworkLoadBalancerState, opts ...pulumi.ResourceOption) (*LbNetworkLoadBalancer, error) {
	var resource LbNetworkLoadBalancer
	err := ctx.ReadResource("yandex:index/lbNetworkLoadBalancer:LbNetworkLoadBalancer", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering LbNetworkLoadBalancer resources.
type lbNetworkLoadBalancerState struct {
	// An AttachedTargetGroup resource. The structure is documented below.
	AttachedTargetGroups []LbNetworkLoadBalancerAttachedTargetGroup `pulumi:"attachedTargetGroups"`
	// The network load balancer creation timestamp.
	CreatedAt *string `pulumi:"createdAt"`
	// An optional description of the network load balancer. Provide this property when
	// you create the resource.
	Description *string `pulumi:"description"`
	// The ID of the folder to which the resource belongs.
	// If omitted, the provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// Labels to assign to this network load balancer. A list of key/value pairs.
	Labels map[string]string `pulumi:"labels"`
	// Listener specification that will be used by a network load balancer. The structure is documented below.
	Listeners []LbNetworkLoadBalancerListener `pulumi:"listeners"`
	// Name of the listener. The name must be unique for each listener on a single load balancer.
	Name *string `pulumi:"name"`
	// ID of the availability zone where the network load balancer resides.
	// The default is 'ru-central1'.
	RegionId *string                        `pulumi:"regionId"`
	Timeouts *LbNetworkLoadBalancerTimeouts `pulumi:"timeouts"`
	// Type of the network load balancer. Must be one of 'external' or 'internal'. The default is 'external'.
	Type *string `pulumi:"type"`
}

type LbNetworkLoadBalancerState struct {
	// An AttachedTargetGroup resource. The structure is documented below.
	AttachedTargetGroups LbNetworkLoadBalancerAttachedTargetGroupArrayInput
	// The network load balancer creation timestamp.
	CreatedAt pulumi.StringPtrInput
	// An optional description of the network load balancer. Provide this property when
	// you create the resource.
	Description pulumi.StringPtrInput
	// The ID of the folder to which the resource belongs.
	// If omitted, the provider folder is used.
	FolderId pulumi.StringPtrInput
	// Labels to assign to this network load balancer. A list of key/value pairs.
	Labels pulumi.StringMapInput
	// Listener specification that will be used by a network load balancer. The structure is documented below.
	Listeners LbNetworkLoadBalancerListenerArrayInput
	// Name of the listener. The name must be unique for each listener on a single load balancer.
	Name pulumi.StringPtrInput
	// ID of the availability zone where the network load balancer resides.
	// The default is 'ru-central1'.
	RegionId pulumi.StringPtrInput
	Timeouts LbNetworkLoadBalancerTimeoutsPtrInput
	// Type of the network load balancer. Must be one of 'external' or 'internal'. The default is 'external'.
	Type pulumi.StringPtrInput
}

func (LbNetworkLoadBalancerState) ElementType() reflect.Type {
	return reflect.TypeOf((*lbNetworkLoadBalancerState)(nil)).Elem()
}

type lbNetworkLoadBalancerArgs struct {
	// An AttachedTargetGroup resource. The structure is documented below.
	AttachedTargetGroups []LbNetworkLoadBalancerAttachedTargetGroup `pulumi:"attachedTargetGroups"`
	// An optional description of the network load balancer. Provide this property when
	// you create the resource.
	Description *string `pulumi:"description"`
	// The ID of the folder to which the resource belongs.
	// If omitted, the provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// Labels to assign to this network load balancer. A list of key/value pairs.
	Labels map[string]string `pulumi:"labels"`
	// Listener specification that will be used by a network load balancer. The structure is documented below.
	Listeners []LbNetworkLoadBalancerListener `pulumi:"listeners"`
	// Name of the listener. The name must be unique for each listener on a single load balancer.
	Name *string `pulumi:"name"`
	// ID of the availability zone where the network load balancer resides.
	// The default is 'ru-central1'.
	RegionId *string                       `pulumi:"regionId"`
	Timeouts LbNetworkLoadBalancerTimeouts `pulumi:"timeouts"`
	// Type of the network load balancer. Must be one of 'external' or 'internal'. The default is 'external'.
	Type *string `pulumi:"type"`
}

// The set of arguments for constructing a LbNetworkLoadBalancer resource.
type LbNetworkLoadBalancerArgs struct {
	// An AttachedTargetGroup resource. The structure is documented below.
	AttachedTargetGroups LbNetworkLoadBalancerAttachedTargetGroupArrayInput
	// An optional description of the network load balancer. Provide this property when
	// you create the resource.
	Description pulumi.StringPtrInput
	// The ID of the folder to which the resource belongs.
	// If omitted, the provider folder is used.
	FolderId pulumi.StringPtrInput
	// Labels to assign to this network load balancer. A list of key/value pairs.
	Labels pulumi.StringMapInput
	// Listener specification that will be used by a network load balancer. The structure is documented below.
	Listeners LbNetworkLoadBalancerListenerArrayInput
	// Name of the listener. The name must be unique for each listener on a single load balancer.
	Name pulumi.StringPtrInput
	// ID of the availability zone where the network load balancer resides.
	// The default is 'ru-central1'.
	RegionId pulumi.StringPtrInput
	Timeouts LbNetworkLoadBalancerTimeoutsInput
	// Type of the network load balancer. Must be one of 'external' or 'internal'. The default is 'external'.
	Type pulumi.StringPtrInput
}

func (LbNetworkLoadBalancerArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*lbNetworkLoadBalancerArgs)(nil)).Elem()
}

type LbNetworkLoadBalancerInput interface {
	pulumi.Input

	ToLbNetworkLoadBalancerOutput() LbNetworkLoadBalancerOutput
	ToLbNetworkLoadBalancerOutputWithContext(ctx context.Context) LbNetworkLoadBalancerOutput
}

func (*LbNetworkLoadBalancer) ElementType() reflect.Type {
	return reflect.TypeOf((**LbNetworkLoadBalancer)(nil)).Elem()
}

func (i *LbNetworkLoadBalancer) ToLbNetworkLoadBalancerOutput() LbNetworkLoadBalancerOutput {
	return i.ToLbNetworkLoadBalancerOutputWithContext(context.Background())
}

func (i *LbNetworkLoadBalancer) ToLbNetworkLoadBalancerOutputWithContext(ctx context.Context) LbNetworkLoadBalancerOutput {
	return pulumi.ToOutputWithContext(ctx, i).(LbNetworkLoadBalancerOutput)
}

// LbNetworkLoadBalancerArrayInput is an input type that accepts LbNetworkLoadBalancerArray and LbNetworkLoadBalancerArrayOutput values.
// You can construct a concrete instance of `LbNetworkLoadBalancerArrayInput` via:
//
//          LbNetworkLoadBalancerArray{ LbNetworkLoadBalancerArgs{...} }
type LbNetworkLoadBalancerArrayInput interface {
	pulumi.Input

	ToLbNetworkLoadBalancerArrayOutput() LbNetworkLoadBalancerArrayOutput
	ToLbNetworkLoadBalancerArrayOutputWithContext(context.Context) LbNetworkLoadBalancerArrayOutput
}

type LbNetworkLoadBalancerArray []LbNetworkLoadBalancerInput

func (LbNetworkLoadBalancerArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*LbNetworkLoadBalancer)(nil)).Elem()
}

func (i LbNetworkLoadBalancerArray) ToLbNetworkLoadBalancerArrayOutput() LbNetworkLoadBalancerArrayOutput {
	return i.ToLbNetworkLoadBalancerArrayOutputWithContext(context.Background())
}

func (i LbNetworkLoadBalancerArray) ToLbNetworkLoadBalancerArrayOutputWithContext(ctx context.Context) LbNetworkLoadBalancerArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(LbNetworkLoadBalancerArrayOutput)
}

// LbNetworkLoadBalancerMapInput is an input type that accepts LbNetworkLoadBalancerMap and LbNetworkLoadBalancerMapOutput values.
// You can construct a concrete instance of `LbNetworkLoadBalancerMapInput` via:
//
//          LbNetworkLoadBalancerMap{ "key": LbNetworkLoadBalancerArgs{...} }
type LbNetworkLoadBalancerMapInput interface {
	pulumi.Input

	ToLbNetworkLoadBalancerMapOutput() LbNetworkLoadBalancerMapOutput
	ToLbNetworkLoadBalancerMapOutputWithContext(context.Context) LbNetworkLoadBalancerMapOutput
}

type LbNetworkLoadBalancerMap map[string]LbNetworkLoadBalancerInput

func (LbNetworkLoadBalancerMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*LbNetworkLoadBalancer)(nil)).Elem()
}

func (i LbNetworkLoadBalancerMap) ToLbNetworkLoadBalancerMapOutput() LbNetworkLoadBalancerMapOutput {
	return i.ToLbNetworkLoadBalancerMapOutputWithContext(context.Background())
}

func (i LbNetworkLoadBalancerMap) ToLbNetworkLoadBalancerMapOutputWithContext(ctx context.Context) LbNetworkLoadBalancerMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(LbNetworkLoadBalancerMapOutput)
}

type LbNetworkLoadBalancerOutput struct{ *pulumi.OutputState }

func (LbNetworkLoadBalancerOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**LbNetworkLoadBalancer)(nil)).Elem()
}

func (o LbNetworkLoadBalancerOutput) ToLbNetworkLoadBalancerOutput() LbNetworkLoadBalancerOutput {
	return o
}

func (o LbNetworkLoadBalancerOutput) ToLbNetworkLoadBalancerOutputWithContext(ctx context.Context) LbNetworkLoadBalancerOutput {
	return o
}

// An AttachedTargetGroup resource. The structure is documented below.
func (o LbNetworkLoadBalancerOutput) AttachedTargetGroups() LbNetworkLoadBalancerAttachedTargetGroupArrayOutput {
	return o.ApplyT(func(v *LbNetworkLoadBalancer) LbNetworkLoadBalancerAttachedTargetGroupArrayOutput {
		return v.AttachedTargetGroups
	}).(LbNetworkLoadBalancerAttachedTargetGroupArrayOutput)
}

// The network load balancer creation timestamp.
func (o LbNetworkLoadBalancerOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v *LbNetworkLoadBalancer) pulumi.StringOutput { return v.CreatedAt }).(pulumi.StringOutput)
}

// An optional description of the network load balancer. Provide this property when
// you create the resource.
func (o LbNetworkLoadBalancerOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *LbNetworkLoadBalancer) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// The ID of the folder to which the resource belongs.
// If omitted, the provider folder is used.
func (o LbNetworkLoadBalancerOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v *LbNetworkLoadBalancer) pulumi.StringOutput { return v.FolderId }).(pulumi.StringOutput)
}

// Labels to assign to this network load balancer. A list of key/value pairs.
func (o LbNetworkLoadBalancerOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v *LbNetworkLoadBalancer) pulumi.StringMapOutput { return v.Labels }).(pulumi.StringMapOutput)
}

// Listener specification that will be used by a network load balancer. The structure is documented below.
func (o LbNetworkLoadBalancerOutput) Listeners() LbNetworkLoadBalancerListenerArrayOutput {
	return o.ApplyT(func(v *LbNetworkLoadBalancer) LbNetworkLoadBalancerListenerArrayOutput { return v.Listeners }).(LbNetworkLoadBalancerListenerArrayOutput)
}

// Name of the listener. The name must be unique for each listener on a single load balancer.
func (o LbNetworkLoadBalancerOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *LbNetworkLoadBalancer) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// ID of the availability zone where the network load balancer resides.
// The default is 'ru-central1'.
func (o LbNetworkLoadBalancerOutput) RegionId() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *LbNetworkLoadBalancer) pulumi.StringPtrOutput { return v.RegionId }).(pulumi.StringPtrOutput)
}

func (o LbNetworkLoadBalancerOutput) Timeouts() LbNetworkLoadBalancerTimeoutsOutput {
	return o.ApplyT(func(v *LbNetworkLoadBalancer) LbNetworkLoadBalancerTimeoutsOutput { return v.Timeouts }).(LbNetworkLoadBalancerTimeoutsOutput)
}

// Type of the network load balancer. Must be one of 'external' or 'internal'. The default is 'external'.
func (o LbNetworkLoadBalancerOutput) Type() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *LbNetworkLoadBalancer) pulumi.StringPtrOutput { return v.Type }).(pulumi.StringPtrOutput)
}

type LbNetworkLoadBalancerArrayOutput struct{ *pulumi.OutputState }

func (LbNetworkLoadBalancerArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*LbNetworkLoadBalancer)(nil)).Elem()
}

func (o LbNetworkLoadBalancerArrayOutput) ToLbNetworkLoadBalancerArrayOutput() LbNetworkLoadBalancerArrayOutput {
	return o
}

func (o LbNetworkLoadBalancerArrayOutput) ToLbNetworkLoadBalancerArrayOutputWithContext(ctx context.Context) LbNetworkLoadBalancerArrayOutput {
	return o
}

func (o LbNetworkLoadBalancerArrayOutput) Index(i pulumi.IntInput) LbNetworkLoadBalancerOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *LbNetworkLoadBalancer {
		return vs[0].([]*LbNetworkLoadBalancer)[vs[1].(int)]
	}).(LbNetworkLoadBalancerOutput)
}

type LbNetworkLoadBalancerMapOutput struct{ *pulumi.OutputState }

func (LbNetworkLoadBalancerMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*LbNetworkLoadBalancer)(nil)).Elem()
}

func (o LbNetworkLoadBalancerMapOutput) ToLbNetworkLoadBalancerMapOutput() LbNetworkLoadBalancerMapOutput {
	return o
}

func (o LbNetworkLoadBalancerMapOutput) ToLbNetworkLoadBalancerMapOutputWithContext(ctx context.Context) LbNetworkLoadBalancerMapOutput {
	return o
}

func (o LbNetworkLoadBalancerMapOutput) MapIndex(k pulumi.StringInput) LbNetworkLoadBalancerOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *LbNetworkLoadBalancer {
		return vs[0].(map[string]*LbNetworkLoadBalancer)[vs[1].(string)]
	}).(LbNetworkLoadBalancerOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*LbNetworkLoadBalancerInput)(nil)).Elem(), &LbNetworkLoadBalancer{})
	pulumi.RegisterInputType(reflect.TypeOf((*LbNetworkLoadBalancerArrayInput)(nil)).Elem(), LbNetworkLoadBalancerArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*LbNetworkLoadBalancerMapInput)(nil)).Elem(), LbNetworkLoadBalancerMap{})
	pulumi.RegisterOutputType(LbNetworkLoadBalancerOutput{})
	pulumi.RegisterOutputType(LbNetworkLoadBalancerArrayOutput{})
	pulumi.RegisterOutputType(LbNetworkLoadBalancerMapOutput{})
}
