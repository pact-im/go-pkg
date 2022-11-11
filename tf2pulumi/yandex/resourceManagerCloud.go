// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Allows creation and management of Cloud resources for an existing Yandex.Cloud Organization. See [the official documentation](https://cloud.yandex.com/docs/resource-manager/concepts/resources-hierarchy) for additional info.
// Note: deletion of clouds may take up to 30 minutes as it requires a lot of communication between cloud services.
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
// 		_, err := yandex.NewResourceManagerCloud(ctx, "cloud1", &yandex.ResourceManagerCloudArgs{
// 			OrganizationId: pulumi.String("my_organization_id"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
type ResourceManagerCloud struct {
	pulumi.CustomResourceState

	CreatedAt pulumi.StringOutput `pulumi:"createdAt"`
	// A description of the Cloud.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// A set of key/value label pairs to assign to the Cloud.
	Labels pulumi.StringMapOutput `pulumi:"labels"`
	// The name of the Cloud.
	Name pulumi.StringOutput `pulumi:"name"`
	// Yandex.Cloud Organization that the cloud belongs to. If value is omitted, the default provider Organization ID is used.
	OrganizationId pulumi.StringOutput                `pulumi:"organizationId"`
	Timeouts       ResourceManagerCloudTimeoutsOutput `pulumi:"timeouts"`
}

// NewResourceManagerCloud registers a new resource with the given unique name, arguments, and options.
func NewResourceManagerCloud(ctx *pulumi.Context,
	name string, args *ResourceManagerCloudArgs, opts ...pulumi.ResourceOption) (*ResourceManagerCloud, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	var resource ResourceManagerCloud
	err := ctx.RegisterResource("yandex:index/resourceManagerCloud:ResourceManagerCloud", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetResourceManagerCloud gets an existing ResourceManagerCloud resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetResourceManagerCloud(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *ResourceManagerCloudState, opts ...pulumi.ResourceOption) (*ResourceManagerCloud, error) {
	var resource ResourceManagerCloud
	err := ctx.ReadResource("yandex:index/resourceManagerCloud:ResourceManagerCloud", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering ResourceManagerCloud resources.
type resourceManagerCloudState struct {
	CreatedAt *string `pulumi:"createdAt"`
	// A description of the Cloud.
	Description *string `pulumi:"description"`
	// A set of key/value label pairs to assign to the Cloud.
	Labels map[string]string `pulumi:"labels"`
	// The name of the Cloud.
	Name *string `pulumi:"name"`
	// Yandex.Cloud Organization that the cloud belongs to. If value is omitted, the default provider Organization ID is used.
	OrganizationId *string                       `pulumi:"organizationId"`
	Timeouts       *ResourceManagerCloudTimeouts `pulumi:"timeouts"`
}

type ResourceManagerCloudState struct {
	CreatedAt pulumi.StringPtrInput
	// A description of the Cloud.
	Description pulumi.StringPtrInput
	// A set of key/value label pairs to assign to the Cloud.
	Labels pulumi.StringMapInput
	// The name of the Cloud.
	Name pulumi.StringPtrInput
	// Yandex.Cloud Organization that the cloud belongs to. If value is omitted, the default provider Organization ID is used.
	OrganizationId pulumi.StringPtrInput
	Timeouts       ResourceManagerCloudTimeoutsPtrInput
}

func (ResourceManagerCloudState) ElementType() reflect.Type {
	return reflect.TypeOf((*resourceManagerCloudState)(nil)).Elem()
}

type resourceManagerCloudArgs struct {
	// A description of the Cloud.
	Description *string `pulumi:"description"`
	// A set of key/value label pairs to assign to the Cloud.
	Labels map[string]string `pulumi:"labels"`
	// The name of the Cloud.
	Name *string `pulumi:"name"`
	// Yandex.Cloud Organization that the cloud belongs to. If value is omitted, the default provider Organization ID is used.
	OrganizationId *string                      `pulumi:"organizationId"`
	Timeouts       ResourceManagerCloudTimeouts `pulumi:"timeouts"`
}

// The set of arguments for constructing a ResourceManagerCloud resource.
type ResourceManagerCloudArgs struct {
	// A description of the Cloud.
	Description pulumi.StringPtrInput
	// A set of key/value label pairs to assign to the Cloud.
	Labels pulumi.StringMapInput
	// The name of the Cloud.
	Name pulumi.StringPtrInput
	// Yandex.Cloud Organization that the cloud belongs to. If value is omitted, the default provider Organization ID is used.
	OrganizationId pulumi.StringPtrInput
	Timeouts       ResourceManagerCloudTimeoutsInput
}

func (ResourceManagerCloudArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*resourceManagerCloudArgs)(nil)).Elem()
}

type ResourceManagerCloudInput interface {
	pulumi.Input

	ToResourceManagerCloudOutput() ResourceManagerCloudOutput
	ToResourceManagerCloudOutputWithContext(ctx context.Context) ResourceManagerCloudOutput
}

func (*ResourceManagerCloud) ElementType() reflect.Type {
	return reflect.TypeOf((**ResourceManagerCloud)(nil)).Elem()
}

func (i *ResourceManagerCloud) ToResourceManagerCloudOutput() ResourceManagerCloudOutput {
	return i.ToResourceManagerCloudOutputWithContext(context.Background())
}

func (i *ResourceManagerCloud) ToResourceManagerCloudOutputWithContext(ctx context.Context) ResourceManagerCloudOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourceManagerCloudOutput)
}

// ResourceManagerCloudArrayInput is an input type that accepts ResourceManagerCloudArray and ResourceManagerCloudArrayOutput values.
// You can construct a concrete instance of `ResourceManagerCloudArrayInput` via:
//
//          ResourceManagerCloudArray{ ResourceManagerCloudArgs{...} }
type ResourceManagerCloudArrayInput interface {
	pulumi.Input

	ToResourceManagerCloudArrayOutput() ResourceManagerCloudArrayOutput
	ToResourceManagerCloudArrayOutputWithContext(context.Context) ResourceManagerCloudArrayOutput
}

type ResourceManagerCloudArray []ResourceManagerCloudInput

func (ResourceManagerCloudArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ResourceManagerCloud)(nil)).Elem()
}

func (i ResourceManagerCloudArray) ToResourceManagerCloudArrayOutput() ResourceManagerCloudArrayOutput {
	return i.ToResourceManagerCloudArrayOutputWithContext(context.Background())
}

func (i ResourceManagerCloudArray) ToResourceManagerCloudArrayOutputWithContext(ctx context.Context) ResourceManagerCloudArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourceManagerCloudArrayOutput)
}

// ResourceManagerCloudMapInput is an input type that accepts ResourceManagerCloudMap and ResourceManagerCloudMapOutput values.
// You can construct a concrete instance of `ResourceManagerCloudMapInput` via:
//
//          ResourceManagerCloudMap{ "key": ResourceManagerCloudArgs{...} }
type ResourceManagerCloudMapInput interface {
	pulumi.Input

	ToResourceManagerCloudMapOutput() ResourceManagerCloudMapOutput
	ToResourceManagerCloudMapOutputWithContext(context.Context) ResourceManagerCloudMapOutput
}

type ResourceManagerCloudMap map[string]ResourceManagerCloudInput

func (ResourceManagerCloudMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ResourceManagerCloud)(nil)).Elem()
}

func (i ResourceManagerCloudMap) ToResourceManagerCloudMapOutput() ResourceManagerCloudMapOutput {
	return i.ToResourceManagerCloudMapOutputWithContext(context.Background())
}

func (i ResourceManagerCloudMap) ToResourceManagerCloudMapOutputWithContext(ctx context.Context) ResourceManagerCloudMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourceManagerCloudMapOutput)
}

type ResourceManagerCloudOutput struct{ *pulumi.OutputState }

func (ResourceManagerCloudOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**ResourceManagerCloud)(nil)).Elem()
}

func (o ResourceManagerCloudOutput) ToResourceManagerCloudOutput() ResourceManagerCloudOutput {
	return o
}

func (o ResourceManagerCloudOutput) ToResourceManagerCloudOutputWithContext(ctx context.Context) ResourceManagerCloudOutput {
	return o
}

func (o ResourceManagerCloudOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v *ResourceManagerCloud) pulumi.StringOutput { return v.CreatedAt }).(pulumi.StringOutput)
}

// A description of the Cloud.
func (o ResourceManagerCloudOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *ResourceManagerCloud) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// A set of key/value label pairs to assign to the Cloud.
func (o ResourceManagerCloudOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v *ResourceManagerCloud) pulumi.StringMapOutput { return v.Labels }).(pulumi.StringMapOutput)
}

// The name of the Cloud.
func (o ResourceManagerCloudOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *ResourceManagerCloud) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// Yandex.Cloud Organization that the cloud belongs to. If value is omitted, the default provider Organization ID is used.
func (o ResourceManagerCloudOutput) OrganizationId() pulumi.StringOutput {
	return o.ApplyT(func(v *ResourceManagerCloud) pulumi.StringOutput { return v.OrganizationId }).(pulumi.StringOutput)
}

func (o ResourceManagerCloudOutput) Timeouts() ResourceManagerCloudTimeoutsOutput {
	return o.ApplyT(func(v *ResourceManagerCloud) ResourceManagerCloudTimeoutsOutput { return v.Timeouts }).(ResourceManagerCloudTimeoutsOutput)
}

type ResourceManagerCloudArrayOutput struct{ *pulumi.OutputState }

func (ResourceManagerCloudArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ResourceManagerCloud)(nil)).Elem()
}

func (o ResourceManagerCloudArrayOutput) ToResourceManagerCloudArrayOutput() ResourceManagerCloudArrayOutput {
	return o
}

func (o ResourceManagerCloudArrayOutput) ToResourceManagerCloudArrayOutputWithContext(ctx context.Context) ResourceManagerCloudArrayOutput {
	return o
}

func (o ResourceManagerCloudArrayOutput) Index(i pulumi.IntInput) ResourceManagerCloudOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *ResourceManagerCloud {
		return vs[0].([]*ResourceManagerCloud)[vs[1].(int)]
	}).(ResourceManagerCloudOutput)
}

type ResourceManagerCloudMapOutput struct{ *pulumi.OutputState }

func (ResourceManagerCloudMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ResourceManagerCloud)(nil)).Elem()
}

func (o ResourceManagerCloudMapOutput) ToResourceManagerCloudMapOutput() ResourceManagerCloudMapOutput {
	return o
}

func (o ResourceManagerCloudMapOutput) ToResourceManagerCloudMapOutputWithContext(ctx context.Context) ResourceManagerCloudMapOutput {
	return o
}

func (o ResourceManagerCloudMapOutput) MapIndex(k pulumi.StringInput) ResourceManagerCloudOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *ResourceManagerCloud {
		return vs[0].(map[string]*ResourceManagerCloud)[vs[1].(string)]
	}).(ResourceManagerCloudOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ResourceManagerCloudInput)(nil)).Elem(), &ResourceManagerCloud{})
	pulumi.RegisterInputType(reflect.TypeOf((*ResourceManagerCloudArrayInput)(nil)).Elem(), ResourceManagerCloudArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*ResourceManagerCloudMapInput)(nil)).Elem(), ResourceManagerCloudMap{})
	pulumi.RegisterOutputType(ResourceManagerCloudOutput{})
	pulumi.RegisterOutputType(ResourceManagerCloudArrayOutput{})
	pulumi.RegisterOutputType(ResourceManagerCloudMapOutput{})
}
