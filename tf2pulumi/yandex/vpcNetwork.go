// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Manages a network within the Yandex.Cloud. For more information, see
// [the official documentation](https://cloud.yandex.com/docs/vpc/concepts/network#network).
//
// * How-to Guides
//     * [Cloud Networking](https://cloud.yandex.com/docs/vpc/)
//     * [VPC Addressing](https://cloud.yandex.com/docs/vpc/concepts/address)
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
// 		_, err := yandex.NewVpcNetwork(ctx, "default", nil)
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
// A network can be imported using the `id` of the resource, e.g.
//
// ```sh
//  $ pulumi import yandex:index/vpcNetwork:VpcNetwork default network_id
// ```
type VpcNetwork struct {
	pulumi.CustomResourceState

	// Creation timestamp of the key.
	CreatedAt pulumi.StringOutput `pulumi:"createdAt"`
	// ID of default Security Group of this network.
	DefaultSecurityGroupId pulumi.StringOutput `pulumi:"defaultSecurityGroupId"`
	// An optional description of this resource. Provide this property when
	// you create the resource.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// ID of the folder that the resource belongs to. If it
	// is not provided, the default provider folder is used.
	FolderId pulumi.StringOutput `pulumi:"folderId"`
	// Labels to apply to this network. A list of key/value pairs.
	Labels pulumi.StringMapOutput `pulumi:"labels"`
	// Name of the network. Provided by the client when the network is created.
	Name      pulumi.StringOutput      `pulumi:"name"`
	SubnetIds pulumi.StringArrayOutput `pulumi:"subnetIds"`
	Timeouts  VpcNetworkTimeoutsOutput `pulumi:"timeouts"`
}

// NewVpcNetwork registers a new resource with the given unique name, arguments, and options.
func NewVpcNetwork(ctx *pulumi.Context,
	name string, args *VpcNetworkArgs, opts ...pulumi.ResourceOption) (*VpcNetwork, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	var resource VpcNetwork
	err := ctx.RegisterResource("yandex:index/vpcNetwork:VpcNetwork", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetVpcNetwork gets an existing VpcNetwork resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetVpcNetwork(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *VpcNetworkState, opts ...pulumi.ResourceOption) (*VpcNetwork, error) {
	var resource VpcNetwork
	err := ctx.ReadResource("yandex:index/vpcNetwork:VpcNetwork", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering VpcNetwork resources.
type vpcNetworkState struct {
	// Creation timestamp of the key.
	CreatedAt *string `pulumi:"createdAt"`
	// ID of default Security Group of this network.
	DefaultSecurityGroupId *string `pulumi:"defaultSecurityGroupId"`
	// An optional description of this resource. Provide this property when
	// you create the resource.
	Description *string `pulumi:"description"`
	// ID of the folder that the resource belongs to. If it
	// is not provided, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// Labels to apply to this network. A list of key/value pairs.
	Labels map[string]string `pulumi:"labels"`
	// Name of the network. Provided by the client when the network is created.
	Name      *string             `pulumi:"name"`
	SubnetIds []string            `pulumi:"subnetIds"`
	Timeouts  *VpcNetworkTimeouts `pulumi:"timeouts"`
}

type VpcNetworkState struct {
	// Creation timestamp of the key.
	CreatedAt pulumi.StringPtrInput
	// ID of default Security Group of this network.
	DefaultSecurityGroupId pulumi.StringPtrInput
	// An optional description of this resource. Provide this property when
	// you create the resource.
	Description pulumi.StringPtrInput
	// ID of the folder that the resource belongs to. If it
	// is not provided, the default provider folder is used.
	FolderId pulumi.StringPtrInput
	// Labels to apply to this network. A list of key/value pairs.
	Labels pulumi.StringMapInput
	// Name of the network. Provided by the client when the network is created.
	Name      pulumi.StringPtrInput
	SubnetIds pulumi.StringArrayInput
	Timeouts  VpcNetworkTimeoutsPtrInput
}

func (VpcNetworkState) ElementType() reflect.Type {
	return reflect.TypeOf((*vpcNetworkState)(nil)).Elem()
}

type vpcNetworkArgs struct {
	// An optional description of this resource. Provide this property when
	// you create the resource.
	Description *string `pulumi:"description"`
	// ID of the folder that the resource belongs to. If it
	// is not provided, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// Labels to apply to this network. A list of key/value pairs.
	Labels map[string]string `pulumi:"labels"`
	// Name of the network. Provided by the client when the network is created.
	Name     *string            `pulumi:"name"`
	Timeouts VpcNetworkTimeouts `pulumi:"timeouts"`
}

// The set of arguments for constructing a VpcNetwork resource.
type VpcNetworkArgs struct {
	// An optional description of this resource. Provide this property when
	// you create the resource.
	Description pulumi.StringPtrInput
	// ID of the folder that the resource belongs to. If it
	// is not provided, the default provider folder is used.
	FolderId pulumi.StringPtrInput
	// Labels to apply to this network. A list of key/value pairs.
	Labels pulumi.StringMapInput
	// Name of the network. Provided by the client when the network is created.
	Name     pulumi.StringPtrInput
	Timeouts VpcNetworkTimeoutsInput
}

func (VpcNetworkArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*vpcNetworkArgs)(nil)).Elem()
}

type VpcNetworkInput interface {
	pulumi.Input

	ToVpcNetworkOutput() VpcNetworkOutput
	ToVpcNetworkOutputWithContext(ctx context.Context) VpcNetworkOutput
}

func (*VpcNetwork) ElementType() reflect.Type {
	return reflect.TypeOf((**VpcNetwork)(nil)).Elem()
}

func (i *VpcNetwork) ToVpcNetworkOutput() VpcNetworkOutput {
	return i.ToVpcNetworkOutputWithContext(context.Background())
}

func (i *VpcNetwork) ToVpcNetworkOutputWithContext(ctx context.Context) VpcNetworkOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VpcNetworkOutput)
}

// VpcNetworkArrayInput is an input type that accepts VpcNetworkArray and VpcNetworkArrayOutput values.
// You can construct a concrete instance of `VpcNetworkArrayInput` via:
//
//          VpcNetworkArray{ VpcNetworkArgs{...} }
type VpcNetworkArrayInput interface {
	pulumi.Input

	ToVpcNetworkArrayOutput() VpcNetworkArrayOutput
	ToVpcNetworkArrayOutputWithContext(context.Context) VpcNetworkArrayOutput
}

type VpcNetworkArray []VpcNetworkInput

func (VpcNetworkArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*VpcNetwork)(nil)).Elem()
}

func (i VpcNetworkArray) ToVpcNetworkArrayOutput() VpcNetworkArrayOutput {
	return i.ToVpcNetworkArrayOutputWithContext(context.Background())
}

func (i VpcNetworkArray) ToVpcNetworkArrayOutputWithContext(ctx context.Context) VpcNetworkArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VpcNetworkArrayOutput)
}

// VpcNetworkMapInput is an input type that accepts VpcNetworkMap and VpcNetworkMapOutput values.
// You can construct a concrete instance of `VpcNetworkMapInput` via:
//
//          VpcNetworkMap{ "key": VpcNetworkArgs{...} }
type VpcNetworkMapInput interface {
	pulumi.Input

	ToVpcNetworkMapOutput() VpcNetworkMapOutput
	ToVpcNetworkMapOutputWithContext(context.Context) VpcNetworkMapOutput
}

type VpcNetworkMap map[string]VpcNetworkInput

func (VpcNetworkMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*VpcNetwork)(nil)).Elem()
}

func (i VpcNetworkMap) ToVpcNetworkMapOutput() VpcNetworkMapOutput {
	return i.ToVpcNetworkMapOutputWithContext(context.Background())
}

func (i VpcNetworkMap) ToVpcNetworkMapOutputWithContext(ctx context.Context) VpcNetworkMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VpcNetworkMapOutput)
}

type VpcNetworkOutput struct{ *pulumi.OutputState }

func (VpcNetworkOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**VpcNetwork)(nil)).Elem()
}

func (o VpcNetworkOutput) ToVpcNetworkOutput() VpcNetworkOutput {
	return o
}

func (o VpcNetworkOutput) ToVpcNetworkOutputWithContext(ctx context.Context) VpcNetworkOutput {
	return o
}

// Creation timestamp of the key.
func (o VpcNetworkOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v *VpcNetwork) pulumi.StringOutput { return v.CreatedAt }).(pulumi.StringOutput)
}

// ID of default Security Group of this network.
func (o VpcNetworkOutput) DefaultSecurityGroupId() pulumi.StringOutput {
	return o.ApplyT(func(v *VpcNetwork) pulumi.StringOutput { return v.DefaultSecurityGroupId }).(pulumi.StringOutput)
}

// An optional description of this resource. Provide this property when
// you create the resource.
func (o VpcNetworkOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VpcNetwork) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// ID of the folder that the resource belongs to. If it
// is not provided, the default provider folder is used.
func (o VpcNetworkOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v *VpcNetwork) pulumi.StringOutput { return v.FolderId }).(pulumi.StringOutput)
}

// Labels to apply to this network. A list of key/value pairs.
func (o VpcNetworkOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v *VpcNetwork) pulumi.StringMapOutput { return v.Labels }).(pulumi.StringMapOutput)
}

// Name of the network. Provided by the client when the network is created.
func (o VpcNetworkOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *VpcNetwork) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

func (o VpcNetworkOutput) SubnetIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *VpcNetwork) pulumi.StringArrayOutput { return v.SubnetIds }).(pulumi.StringArrayOutput)
}

func (o VpcNetworkOutput) Timeouts() VpcNetworkTimeoutsOutput {
	return o.ApplyT(func(v *VpcNetwork) VpcNetworkTimeoutsOutput { return v.Timeouts }).(VpcNetworkTimeoutsOutput)
}

type VpcNetworkArrayOutput struct{ *pulumi.OutputState }

func (VpcNetworkArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*VpcNetwork)(nil)).Elem()
}

func (o VpcNetworkArrayOutput) ToVpcNetworkArrayOutput() VpcNetworkArrayOutput {
	return o
}

func (o VpcNetworkArrayOutput) ToVpcNetworkArrayOutputWithContext(ctx context.Context) VpcNetworkArrayOutput {
	return o
}

func (o VpcNetworkArrayOutput) Index(i pulumi.IntInput) VpcNetworkOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *VpcNetwork {
		return vs[0].([]*VpcNetwork)[vs[1].(int)]
	}).(VpcNetworkOutput)
}

type VpcNetworkMapOutput struct{ *pulumi.OutputState }

func (VpcNetworkMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*VpcNetwork)(nil)).Elem()
}

func (o VpcNetworkMapOutput) ToVpcNetworkMapOutput() VpcNetworkMapOutput {
	return o
}

func (o VpcNetworkMapOutput) ToVpcNetworkMapOutputWithContext(ctx context.Context) VpcNetworkMapOutput {
	return o
}

func (o VpcNetworkMapOutput) MapIndex(k pulumi.StringInput) VpcNetworkOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *VpcNetwork {
		return vs[0].(map[string]*VpcNetwork)[vs[1].(string)]
	}).(VpcNetworkOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*VpcNetworkInput)(nil)).Elem(), &VpcNetwork{})
	pulumi.RegisterInputType(reflect.TypeOf((*VpcNetworkArrayInput)(nil)).Elem(), VpcNetworkArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*VpcNetworkMapInput)(nil)).Elem(), VpcNetworkMap{})
	pulumi.RegisterOutputType(VpcNetworkOutput{})
	pulumi.RegisterOutputType(VpcNetworkArrayOutput{})
	pulumi.RegisterOutputType(VpcNetworkMapOutput{})
}
