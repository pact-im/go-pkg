// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Allows creation and management of the IAM policy for an existing Yandex Resource
// Manager folder.
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
// 		_, err = yandex.LookupResourcemanagerFolder(ctx, &GetResourcemanagerFolderArgs{
// 			FolderId: pulumi.StringRef("my_folder_id"),
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		admin, err := yandex.GetIamPolicy(ctx, &GetIamPolicyArgs{
// 			Bindings: []GetIamPolicyBinding{
// 				GetIamPolicyBinding{
// 					Members: []string{
// 						"userAccount:some_user_id",
// 					},
// 					Role: "editor",
// 				},
// 			},
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		_, err = yandex.NewResourcemanagerFolderIamPolicy(ctx, "folderAdminPolicy", &yandex.ResourcemanagerFolderIamPolicyArgs{
// 			FolderId:   pulumi.Any(data.Yandex_folder.Project1.Id),
// 			PolicyData: pulumi.String(admin.PolicyData),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
type ResourcemanagerFolderIamPolicy struct {
	pulumi.CustomResourceState

	// ID of the folder that the policy is attached to.
	FolderId pulumi.StringOutput `pulumi:"folderId"`
	// The `getIamPolicy` data source that represents
	// the IAM policy that will be applied to the folder. This policy overrides any existing policy applied to the folder.
	PolicyData pulumi.StringOutput `pulumi:"policyData"`
}

// NewResourcemanagerFolderIamPolicy registers a new resource with the given unique name, arguments, and options.
func NewResourcemanagerFolderIamPolicy(ctx *pulumi.Context,
	name string, args *ResourcemanagerFolderIamPolicyArgs, opts ...pulumi.ResourceOption) (*ResourcemanagerFolderIamPolicy, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.FolderId == nil {
		return nil, errors.New("invalid value for required argument 'FolderId'")
	}
	if args.PolicyData == nil {
		return nil, errors.New("invalid value for required argument 'PolicyData'")
	}
	var resource ResourcemanagerFolderIamPolicy
	err := ctx.RegisterResource("yandex:index/resourcemanagerFolderIamPolicy:ResourcemanagerFolderIamPolicy", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetResourcemanagerFolderIamPolicy gets an existing ResourcemanagerFolderIamPolicy resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetResourcemanagerFolderIamPolicy(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *ResourcemanagerFolderIamPolicyState, opts ...pulumi.ResourceOption) (*ResourcemanagerFolderIamPolicy, error) {
	var resource ResourcemanagerFolderIamPolicy
	err := ctx.ReadResource("yandex:index/resourcemanagerFolderIamPolicy:ResourcemanagerFolderIamPolicy", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering ResourcemanagerFolderIamPolicy resources.
type resourcemanagerFolderIamPolicyState struct {
	// ID of the folder that the policy is attached to.
	FolderId *string `pulumi:"folderId"`
	// The `getIamPolicy` data source that represents
	// the IAM policy that will be applied to the folder. This policy overrides any existing policy applied to the folder.
	PolicyData *string `pulumi:"policyData"`
}

type ResourcemanagerFolderIamPolicyState struct {
	// ID of the folder that the policy is attached to.
	FolderId pulumi.StringPtrInput
	// The `getIamPolicy` data source that represents
	// the IAM policy that will be applied to the folder. This policy overrides any existing policy applied to the folder.
	PolicyData pulumi.StringPtrInput
}

func (ResourcemanagerFolderIamPolicyState) ElementType() reflect.Type {
	return reflect.TypeOf((*resourcemanagerFolderIamPolicyState)(nil)).Elem()
}

type resourcemanagerFolderIamPolicyArgs struct {
	// ID of the folder that the policy is attached to.
	FolderId string `pulumi:"folderId"`
	// The `getIamPolicy` data source that represents
	// the IAM policy that will be applied to the folder. This policy overrides any existing policy applied to the folder.
	PolicyData string `pulumi:"policyData"`
}

// The set of arguments for constructing a ResourcemanagerFolderIamPolicy resource.
type ResourcemanagerFolderIamPolicyArgs struct {
	// ID of the folder that the policy is attached to.
	FolderId pulumi.StringInput
	// The `getIamPolicy` data source that represents
	// the IAM policy that will be applied to the folder. This policy overrides any existing policy applied to the folder.
	PolicyData pulumi.StringInput
}

func (ResourcemanagerFolderIamPolicyArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*resourcemanagerFolderIamPolicyArgs)(nil)).Elem()
}

type ResourcemanagerFolderIamPolicyInput interface {
	pulumi.Input

	ToResourcemanagerFolderIamPolicyOutput() ResourcemanagerFolderIamPolicyOutput
	ToResourcemanagerFolderIamPolicyOutputWithContext(ctx context.Context) ResourcemanagerFolderIamPolicyOutput
}

func (*ResourcemanagerFolderIamPolicy) ElementType() reflect.Type {
	return reflect.TypeOf((**ResourcemanagerFolderIamPolicy)(nil)).Elem()
}

func (i *ResourcemanagerFolderIamPolicy) ToResourcemanagerFolderIamPolicyOutput() ResourcemanagerFolderIamPolicyOutput {
	return i.ToResourcemanagerFolderIamPolicyOutputWithContext(context.Background())
}

func (i *ResourcemanagerFolderIamPolicy) ToResourcemanagerFolderIamPolicyOutputWithContext(ctx context.Context) ResourcemanagerFolderIamPolicyOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourcemanagerFolderIamPolicyOutput)
}

// ResourcemanagerFolderIamPolicyArrayInput is an input type that accepts ResourcemanagerFolderIamPolicyArray and ResourcemanagerFolderIamPolicyArrayOutput values.
// You can construct a concrete instance of `ResourcemanagerFolderIamPolicyArrayInput` via:
//
//          ResourcemanagerFolderIamPolicyArray{ ResourcemanagerFolderIamPolicyArgs{...} }
type ResourcemanagerFolderIamPolicyArrayInput interface {
	pulumi.Input

	ToResourcemanagerFolderIamPolicyArrayOutput() ResourcemanagerFolderIamPolicyArrayOutput
	ToResourcemanagerFolderIamPolicyArrayOutputWithContext(context.Context) ResourcemanagerFolderIamPolicyArrayOutput
}

type ResourcemanagerFolderIamPolicyArray []ResourcemanagerFolderIamPolicyInput

func (ResourcemanagerFolderIamPolicyArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ResourcemanagerFolderIamPolicy)(nil)).Elem()
}

func (i ResourcemanagerFolderIamPolicyArray) ToResourcemanagerFolderIamPolicyArrayOutput() ResourcemanagerFolderIamPolicyArrayOutput {
	return i.ToResourcemanagerFolderIamPolicyArrayOutputWithContext(context.Background())
}

func (i ResourcemanagerFolderIamPolicyArray) ToResourcemanagerFolderIamPolicyArrayOutputWithContext(ctx context.Context) ResourcemanagerFolderIamPolicyArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourcemanagerFolderIamPolicyArrayOutput)
}

// ResourcemanagerFolderIamPolicyMapInput is an input type that accepts ResourcemanagerFolderIamPolicyMap and ResourcemanagerFolderIamPolicyMapOutput values.
// You can construct a concrete instance of `ResourcemanagerFolderIamPolicyMapInput` via:
//
//          ResourcemanagerFolderIamPolicyMap{ "key": ResourcemanagerFolderIamPolicyArgs{...} }
type ResourcemanagerFolderIamPolicyMapInput interface {
	pulumi.Input

	ToResourcemanagerFolderIamPolicyMapOutput() ResourcemanagerFolderIamPolicyMapOutput
	ToResourcemanagerFolderIamPolicyMapOutputWithContext(context.Context) ResourcemanagerFolderIamPolicyMapOutput
}

type ResourcemanagerFolderIamPolicyMap map[string]ResourcemanagerFolderIamPolicyInput

func (ResourcemanagerFolderIamPolicyMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ResourcemanagerFolderIamPolicy)(nil)).Elem()
}

func (i ResourcemanagerFolderIamPolicyMap) ToResourcemanagerFolderIamPolicyMapOutput() ResourcemanagerFolderIamPolicyMapOutput {
	return i.ToResourcemanagerFolderIamPolicyMapOutputWithContext(context.Background())
}

func (i ResourcemanagerFolderIamPolicyMap) ToResourcemanagerFolderIamPolicyMapOutputWithContext(ctx context.Context) ResourcemanagerFolderIamPolicyMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ResourcemanagerFolderIamPolicyMapOutput)
}

type ResourcemanagerFolderIamPolicyOutput struct{ *pulumi.OutputState }

func (ResourcemanagerFolderIamPolicyOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**ResourcemanagerFolderIamPolicy)(nil)).Elem()
}

func (o ResourcemanagerFolderIamPolicyOutput) ToResourcemanagerFolderIamPolicyOutput() ResourcemanagerFolderIamPolicyOutput {
	return o
}

func (o ResourcemanagerFolderIamPolicyOutput) ToResourcemanagerFolderIamPolicyOutputWithContext(ctx context.Context) ResourcemanagerFolderIamPolicyOutput {
	return o
}

// ID of the folder that the policy is attached to.
func (o ResourcemanagerFolderIamPolicyOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v *ResourcemanagerFolderIamPolicy) pulumi.StringOutput { return v.FolderId }).(pulumi.StringOutput)
}

// The `getIamPolicy` data source that represents
// the IAM policy that will be applied to the folder. This policy overrides any existing policy applied to the folder.
func (o ResourcemanagerFolderIamPolicyOutput) PolicyData() pulumi.StringOutput {
	return o.ApplyT(func(v *ResourcemanagerFolderIamPolicy) pulumi.StringOutput { return v.PolicyData }).(pulumi.StringOutput)
}

type ResourcemanagerFolderIamPolicyArrayOutput struct{ *pulumi.OutputState }

func (ResourcemanagerFolderIamPolicyArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ResourcemanagerFolderIamPolicy)(nil)).Elem()
}

func (o ResourcemanagerFolderIamPolicyArrayOutput) ToResourcemanagerFolderIamPolicyArrayOutput() ResourcemanagerFolderIamPolicyArrayOutput {
	return o
}

func (o ResourcemanagerFolderIamPolicyArrayOutput) ToResourcemanagerFolderIamPolicyArrayOutputWithContext(ctx context.Context) ResourcemanagerFolderIamPolicyArrayOutput {
	return o
}

func (o ResourcemanagerFolderIamPolicyArrayOutput) Index(i pulumi.IntInput) ResourcemanagerFolderIamPolicyOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *ResourcemanagerFolderIamPolicy {
		return vs[0].([]*ResourcemanagerFolderIamPolicy)[vs[1].(int)]
	}).(ResourcemanagerFolderIamPolicyOutput)
}

type ResourcemanagerFolderIamPolicyMapOutput struct{ *pulumi.OutputState }

func (ResourcemanagerFolderIamPolicyMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ResourcemanagerFolderIamPolicy)(nil)).Elem()
}

func (o ResourcemanagerFolderIamPolicyMapOutput) ToResourcemanagerFolderIamPolicyMapOutput() ResourcemanagerFolderIamPolicyMapOutput {
	return o
}

func (o ResourcemanagerFolderIamPolicyMapOutput) ToResourcemanagerFolderIamPolicyMapOutputWithContext(ctx context.Context) ResourcemanagerFolderIamPolicyMapOutput {
	return o
}

func (o ResourcemanagerFolderIamPolicyMapOutput) MapIndex(k pulumi.StringInput) ResourcemanagerFolderIamPolicyOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *ResourcemanagerFolderIamPolicy {
		return vs[0].(map[string]*ResourcemanagerFolderIamPolicy)[vs[1].(string)]
	}).(ResourcemanagerFolderIamPolicyOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ResourcemanagerFolderIamPolicyInput)(nil)).Elem(), &ResourcemanagerFolderIamPolicy{})
	pulumi.RegisterInputType(reflect.TypeOf((*ResourcemanagerFolderIamPolicyArrayInput)(nil)).Elem(), ResourcemanagerFolderIamPolicyArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*ResourcemanagerFolderIamPolicyMapInput)(nil)).Elem(), ResourcemanagerFolderIamPolicyMap{})
	pulumi.RegisterOutputType(ResourcemanagerFolderIamPolicyOutput{})
	pulumi.RegisterOutputType(ResourcemanagerFolderIamPolicyArrayOutput{})
	pulumi.RegisterOutputType(ResourcemanagerFolderIamPolicyMapOutput{})
}