// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Allows management of [Yandex.Cloud IoT Registry](https://cloud.yandex.com/docs/iot-core/quickstart).
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
// 		_, err := yandex.NewIotCoreRegistry(ctx, "myRegistry", &yandex.IotCoreRegistryArgs{
// 			Certificates: pulumi.StringArray{
// 				pulumi.String("public part of certificate1"),
// 				pulumi.String("public part of certificate2"),
// 			},
// 			Description: pulumi.String("any description"),
// 			Labels: pulumi.StringMap{
// 				"my-label": pulumi.String("my-label-value"),
// 			},
// 			Passwords: pulumi.StringArray{
// 				pulumi.String("my-password1"),
// 				pulumi.String("my-password2"),
// 			},
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
type IotCoreRegistry struct {
	pulumi.CustomResourceState

	// A set of certificate's fingerprints for the IoT Core Registry
	Certificates pulumi.StringArrayOutput `pulumi:"certificates"`
	// Creation timestamp of the IoT Core Registry
	CreatedAt pulumi.StringOutput `pulumi:"createdAt"`
	// Description of the IoT Core Registry
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// Folder ID for the IoT Core Registry
	FolderId pulumi.StringOutput `pulumi:"folderId"`
	// A set of key/value label pairs to assign to the IoT Core Registry.
	Labels pulumi.StringMapOutput `pulumi:"labels"`
	// IoT Core Device name used to define registry
	Name pulumi.StringOutput `pulumi:"name"`
	// A set of passwords's id for the IoT Core Registry
	Passwords pulumi.StringArrayOutput      `pulumi:"passwords"`
	Timeouts  IotCoreRegistryTimeoutsOutput `pulumi:"timeouts"`
}

// NewIotCoreRegistry registers a new resource with the given unique name, arguments, and options.
func NewIotCoreRegistry(ctx *pulumi.Context,
	name string, args *IotCoreRegistryArgs, opts ...pulumi.ResourceOption) (*IotCoreRegistry, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	if args.Passwords != nil {
		args.Passwords = pulumi.ToSecret(args.Passwords).(pulumi.StringArrayOutput)
	}
	secrets := pulumi.AdditionalSecretOutputs([]string{
		"passwords",
	})
	opts = append(opts, secrets)
	var resource IotCoreRegistry
	err := ctx.RegisterResource("yandex:index/iotCoreRegistry:IotCoreRegistry", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetIotCoreRegistry gets an existing IotCoreRegistry resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetIotCoreRegistry(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *IotCoreRegistryState, opts ...pulumi.ResourceOption) (*IotCoreRegistry, error) {
	var resource IotCoreRegistry
	err := ctx.ReadResource("yandex:index/iotCoreRegistry:IotCoreRegistry", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering IotCoreRegistry resources.
type iotCoreRegistryState struct {
	// A set of certificate's fingerprints for the IoT Core Registry
	Certificates []string `pulumi:"certificates"`
	// Creation timestamp of the IoT Core Registry
	CreatedAt *string `pulumi:"createdAt"`
	// Description of the IoT Core Registry
	Description *string `pulumi:"description"`
	// Folder ID for the IoT Core Registry
	FolderId *string `pulumi:"folderId"`
	// A set of key/value label pairs to assign to the IoT Core Registry.
	Labels map[string]string `pulumi:"labels"`
	// IoT Core Device name used to define registry
	Name *string `pulumi:"name"`
	// A set of passwords's id for the IoT Core Registry
	Passwords []string                 `pulumi:"passwords"`
	Timeouts  *IotCoreRegistryTimeouts `pulumi:"timeouts"`
}

type IotCoreRegistryState struct {
	// A set of certificate's fingerprints for the IoT Core Registry
	Certificates pulumi.StringArrayInput
	// Creation timestamp of the IoT Core Registry
	CreatedAt pulumi.StringPtrInput
	// Description of the IoT Core Registry
	Description pulumi.StringPtrInput
	// Folder ID for the IoT Core Registry
	FolderId pulumi.StringPtrInput
	// A set of key/value label pairs to assign to the IoT Core Registry.
	Labels pulumi.StringMapInput
	// IoT Core Device name used to define registry
	Name pulumi.StringPtrInput
	// A set of passwords's id for the IoT Core Registry
	Passwords pulumi.StringArrayInput
	Timeouts  IotCoreRegistryTimeoutsPtrInput
}

func (IotCoreRegistryState) ElementType() reflect.Type {
	return reflect.TypeOf((*iotCoreRegistryState)(nil)).Elem()
}

type iotCoreRegistryArgs struct {
	// A set of certificate's fingerprints for the IoT Core Registry
	Certificates []string `pulumi:"certificates"`
	// Description of the IoT Core Registry
	Description *string `pulumi:"description"`
	// Folder ID for the IoT Core Registry
	FolderId *string `pulumi:"folderId"`
	// A set of key/value label pairs to assign to the IoT Core Registry.
	Labels map[string]string `pulumi:"labels"`
	// IoT Core Device name used to define registry
	Name *string `pulumi:"name"`
	// A set of passwords's id for the IoT Core Registry
	Passwords []string                `pulumi:"passwords"`
	Timeouts  IotCoreRegistryTimeouts `pulumi:"timeouts"`
}

// The set of arguments for constructing a IotCoreRegistry resource.
type IotCoreRegistryArgs struct {
	// A set of certificate's fingerprints for the IoT Core Registry
	Certificates pulumi.StringArrayInput
	// Description of the IoT Core Registry
	Description pulumi.StringPtrInput
	// Folder ID for the IoT Core Registry
	FolderId pulumi.StringPtrInput
	// A set of key/value label pairs to assign to the IoT Core Registry.
	Labels pulumi.StringMapInput
	// IoT Core Device name used to define registry
	Name pulumi.StringPtrInput
	// A set of passwords's id for the IoT Core Registry
	Passwords pulumi.StringArrayInput
	Timeouts  IotCoreRegistryTimeoutsInput
}

func (IotCoreRegistryArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*iotCoreRegistryArgs)(nil)).Elem()
}

type IotCoreRegistryInput interface {
	pulumi.Input

	ToIotCoreRegistryOutput() IotCoreRegistryOutput
	ToIotCoreRegistryOutputWithContext(ctx context.Context) IotCoreRegistryOutput
}

func (*IotCoreRegistry) ElementType() reflect.Type {
	return reflect.TypeOf((**IotCoreRegistry)(nil)).Elem()
}

func (i *IotCoreRegistry) ToIotCoreRegistryOutput() IotCoreRegistryOutput {
	return i.ToIotCoreRegistryOutputWithContext(context.Background())
}

func (i *IotCoreRegistry) ToIotCoreRegistryOutputWithContext(ctx context.Context) IotCoreRegistryOutput {
	return pulumi.ToOutputWithContext(ctx, i).(IotCoreRegistryOutput)
}

// IotCoreRegistryArrayInput is an input type that accepts IotCoreRegistryArray and IotCoreRegistryArrayOutput values.
// You can construct a concrete instance of `IotCoreRegistryArrayInput` via:
//
//          IotCoreRegistryArray{ IotCoreRegistryArgs{...} }
type IotCoreRegistryArrayInput interface {
	pulumi.Input

	ToIotCoreRegistryArrayOutput() IotCoreRegistryArrayOutput
	ToIotCoreRegistryArrayOutputWithContext(context.Context) IotCoreRegistryArrayOutput
}

type IotCoreRegistryArray []IotCoreRegistryInput

func (IotCoreRegistryArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*IotCoreRegistry)(nil)).Elem()
}

func (i IotCoreRegistryArray) ToIotCoreRegistryArrayOutput() IotCoreRegistryArrayOutput {
	return i.ToIotCoreRegistryArrayOutputWithContext(context.Background())
}

func (i IotCoreRegistryArray) ToIotCoreRegistryArrayOutputWithContext(ctx context.Context) IotCoreRegistryArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(IotCoreRegistryArrayOutput)
}

// IotCoreRegistryMapInput is an input type that accepts IotCoreRegistryMap and IotCoreRegistryMapOutput values.
// You can construct a concrete instance of `IotCoreRegistryMapInput` via:
//
//          IotCoreRegistryMap{ "key": IotCoreRegistryArgs{...} }
type IotCoreRegistryMapInput interface {
	pulumi.Input

	ToIotCoreRegistryMapOutput() IotCoreRegistryMapOutput
	ToIotCoreRegistryMapOutputWithContext(context.Context) IotCoreRegistryMapOutput
}

type IotCoreRegistryMap map[string]IotCoreRegistryInput

func (IotCoreRegistryMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*IotCoreRegistry)(nil)).Elem()
}

func (i IotCoreRegistryMap) ToIotCoreRegistryMapOutput() IotCoreRegistryMapOutput {
	return i.ToIotCoreRegistryMapOutputWithContext(context.Background())
}

func (i IotCoreRegistryMap) ToIotCoreRegistryMapOutputWithContext(ctx context.Context) IotCoreRegistryMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(IotCoreRegistryMapOutput)
}

type IotCoreRegistryOutput struct{ *pulumi.OutputState }

func (IotCoreRegistryOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**IotCoreRegistry)(nil)).Elem()
}

func (o IotCoreRegistryOutput) ToIotCoreRegistryOutput() IotCoreRegistryOutput {
	return o
}

func (o IotCoreRegistryOutput) ToIotCoreRegistryOutputWithContext(ctx context.Context) IotCoreRegistryOutput {
	return o
}

// A set of certificate's fingerprints for the IoT Core Registry
func (o IotCoreRegistryOutput) Certificates() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *IotCoreRegistry) pulumi.StringArrayOutput { return v.Certificates }).(pulumi.StringArrayOutput)
}

// Creation timestamp of the IoT Core Registry
func (o IotCoreRegistryOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v *IotCoreRegistry) pulumi.StringOutput { return v.CreatedAt }).(pulumi.StringOutput)
}

// Description of the IoT Core Registry
func (o IotCoreRegistryOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *IotCoreRegistry) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// Folder ID for the IoT Core Registry
func (o IotCoreRegistryOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v *IotCoreRegistry) pulumi.StringOutput { return v.FolderId }).(pulumi.StringOutput)
}

// A set of key/value label pairs to assign to the IoT Core Registry.
func (o IotCoreRegistryOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v *IotCoreRegistry) pulumi.StringMapOutput { return v.Labels }).(pulumi.StringMapOutput)
}

// IoT Core Device name used to define registry
func (o IotCoreRegistryOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *IotCoreRegistry) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// A set of passwords's id for the IoT Core Registry
func (o IotCoreRegistryOutput) Passwords() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *IotCoreRegistry) pulumi.StringArrayOutput { return v.Passwords }).(pulumi.StringArrayOutput)
}

func (o IotCoreRegistryOutput) Timeouts() IotCoreRegistryTimeoutsOutput {
	return o.ApplyT(func(v *IotCoreRegistry) IotCoreRegistryTimeoutsOutput { return v.Timeouts }).(IotCoreRegistryTimeoutsOutput)
}

type IotCoreRegistryArrayOutput struct{ *pulumi.OutputState }

func (IotCoreRegistryArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*IotCoreRegistry)(nil)).Elem()
}

func (o IotCoreRegistryArrayOutput) ToIotCoreRegistryArrayOutput() IotCoreRegistryArrayOutput {
	return o
}

func (o IotCoreRegistryArrayOutput) ToIotCoreRegistryArrayOutputWithContext(ctx context.Context) IotCoreRegistryArrayOutput {
	return o
}

func (o IotCoreRegistryArrayOutput) Index(i pulumi.IntInput) IotCoreRegistryOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *IotCoreRegistry {
		return vs[0].([]*IotCoreRegistry)[vs[1].(int)]
	}).(IotCoreRegistryOutput)
}

type IotCoreRegistryMapOutput struct{ *pulumi.OutputState }

func (IotCoreRegistryMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*IotCoreRegistry)(nil)).Elem()
}

func (o IotCoreRegistryMapOutput) ToIotCoreRegistryMapOutput() IotCoreRegistryMapOutput {
	return o
}

func (o IotCoreRegistryMapOutput) ToIotCoreRegistryMapOutputWithContext(ctx context.Context) IotCoreRegistryMapOutput {
	return o
}

func (o IotCoreRegistryMapOutput) MapIndex(k pulumi.StringInput) IotCoreRegistryOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *IotCoreRegistry {
		return vs[0].(map[string]*IotCoreRegistry)[vs[1].(string)]
	}).(IotCoreRegistryOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*IotCoreRegistryInput)(nil)).Elem(), &IotCoreRegistry{})
	pulumi.RegisterInputType(reflect.TypeOf((*IotCoreRegistryArrayInput)(nil)).Elem(), IotCoreRegistryArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*IotCoreRegistryMapInput)(nil)).Elem(), IotCoreRegistryMap{})
	pulumi.RegisterOutputType(IotCoreRegistryOutput{})
	pulumi.RegisterOutputType(IotCoreRegistryArrayOutput{})
	pulumi.RegisterOutputType(IotCoreRegistryMapOutput{})
}