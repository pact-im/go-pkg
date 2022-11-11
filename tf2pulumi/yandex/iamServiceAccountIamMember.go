// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// When managing IAM roles, you can treat a service account either as a resource or as an identity.
// This resource is used to add IAM policy bindings to a service account resource to configure permissions
// that define who can edit the service account.
//
// There are three different resources that help you manage your IAM policy for a service account.
// Each of these resources is used for a different use case:
//
// * yandex_iam_service_account_iam_policy: Authoritative. Sets the IAM policy for the service account and replaces any existing policy already attached.
// * yandex_iam_service_account_iam_binding: Authoritative for a given role. Updates the IAM policy to grant a role to a list of members. Other roles within the IAM policy for the service account are preserved.
// * yandex_iam_service_account_iam_member: Non-authoritative. Updates the IAM policy to grant a role to a new member. Other members for the role of the service account are preserved.
//
// > **Note:** `IamServiceAccountIamPolicy` **cannot** be used in conjunction with `IamServiceAccountIamBinding` and `IamServiceAccountIamMember` or they will conflict over what your policy should be.
//
// > **Note:** `IamServiceAccountIamBinding` resources **can be** used in conjunction with `IamServiceAccountIamMember` resources **only if** they do not grant privileges to the same role.
//
// ## yandex\_service\_account\_iam\_member
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
// 		_, err := yandex.NewIamServiceAccountIamMember(ctx, "admin-account-iam", &yandex.IamServiceAccountIamMemberArgs{
// 			Member:           pulumi.String("userAccount:bar_user_id"),
// 			Role:             pulumi.String("admin"),
// 			ServiceAccountId: pulumi.String("your-service-account-id"),
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
// Service account IAM member resources can be imported using the service account ID, role and member.
//
// ```sh
//  $ pulumi import yandex:index/iamServiceAccountIamMember:IamServiceAccountIamMember admin-account-iam "service_account_id roles/editor foo@example.com"
// ```
type IamServiceAccountIamMember struct {
	pulumi.CustomResourceState

	// Identity that will be granted the privilege in `role`.
	// Entry can have one of the following values:
	// * **userAccount:{user_id}**: A unique user ID that represents a specific Yandex account.
	// * **serviceAccount:{service_account_id}**: A unique service account ID.
	Member pulumi.StringOutput `pulumi:"member"`
	// The role that should be applied. Only one
	// `IamServiceAccountIamBinding` can be used per role.
	Role pulumi.StringOutput `pulumi:"role"`
	// The service account ID to apply a policy to.
	ServiceAccountId pulumi.StringOutput     `pulumi:"serviceAccountId"`
	SleepAfter       pulumi.Float64PtrOutput `pulumi:"sleepAfter"`
}

// NewIamServiceAccountIamMember registers a new resource with the given unique name, arguments, and options.
func NewIamServiceAccountIamMember(ctx *pulumi.Context,
	name string, args *IamServiceAccountIamMemberArgs, opts ...pulumi.ResourceOption) (*IamServiceAccountIamMember, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Member == nil {
		return nil, errors.New("invalid value for required argument 'Member'")
	}
	if args.Role == nil {
		return nil, errors.New("invalid value for required argument 'Role'")
	}
	if args.ServiceAccountId == nil {
		return nil, errors.New("invalid value for required argument 'ServiceAccountId'")
	}
	var resource IamServiceAccountIamMember
	err := ctx.RegisterResource("yandex:index/iamServiceAccountIamMember:IamServiceAccountIamMember", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetIamServiceAccountIamMember gets an existing IamServiceAccountIamMember resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetIamServiceAccountIamMember(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *IamServiceAccountIamMemberState, opts ...pulumi.ResourceOption) (*IamServiceAccountIamMember, error) {
	var resource IamServiceAccountIamMember
	err := ctx.ReadResource("yandex:index/iamServiceAccountIamMember:IamServiceAccountIamMember", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering IamServiceAccountIamMember resources.
type iamServiceAccountIamMemberState struct {
	// Identity that will be granted the privilege in `role`.
	// Entry can have one of the following values:
	// * **userAccount:{user_id}**: A unique user ID that represents a specific Yandex account.
	// * **serviceAccount:{service_account_id}**: A unique service account ID.
	Member *string `pulumi:"member"`
	// The role that should be applied. Only one
	// `IamServiceAccountIamBinding` can be used per role.
	Role *string `pulumi:"role"`
	// The service account ID to apply a policy to.
	ServiceAccountId *string  `pulumi:"serviceAccountId"`
	SleepAfter       *float64 `pulumi:"sleepAfter"`
}

type IamServiceAccountIamMemberState struct {
	// Identity that will be granted the privilege in `role`.
	// Entry can have one of the following values:
	// * **userAccount:{user_id}**: A unique user ID that represents a specific Yandex account.
	// * **serviceAccount:{service_account_id}**: A unique service account ID.
	Member pulumi.StringPtrInput
	// The role that should be applied. Only one
	// `IamServiceAccountIamBinding` can be used per role.
	Role pulumi.StringPtrInput
	// The service account ID to apply a policy to.
	ServiceAccountId pulumi.StringPtrInput
	SleepAfter       pulumi.Float64PtrInput
}

func (IamServiceAccountIamMemberState) ElementType() reflect.Type {
	return reflect.TypeOf((*iamServiceAccountIamMemberState)(nil)).Elem()
}

type iamServiceAccountIamMemberArgs struct {
	// Identity that will be granted the privilege in `role`.
	// Entry can have one of the following values:
	// * **userAccount:{user_id}**: A unique user ID that represents a specific Yandex account.
	// * **serviceAccount:{service_account_id}**: A unique service account ID.
	Member string `pulumi:"member"`
	// The role that should be applied. Only one
	// `IamServiceAccountIamBinding` can be used per role.
	Role string `pulumi:"role"`
	// The service account ID to apply a policy to.
	ServiceAccountId string   `pulumi:"serviceAccountId"`
	SleepAfter       *float64 `pulumi:"sleepAfter"`
}

// The set of arguments for constructing a IamServiceAccountIamMember resource.
type IamServiceAccountIamMemberArgs struct {
	// Identity that will be granted the privilege in `role`.
	// Entry can have one of the following values:
	// * **userAccount:{user_id}**: A unique user ID that represents a specific Yandex account.
	// * **serviceAccount:{service_account_id}**: A unique service account ID.
	Member pulumi.StringInput
	// The role that should be applied. Only one
	// `IamServiceAccountIamBinding` can be used per role.
	Role pulumi.StringInput
	// The service account ID to apply a policy to.
	ServiceAccountId pulumi.StringInput
	SleepAfter       pulumi.Float64PtrInput
}

func (IamServiceAccountIamMemberArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*iamServiceAccountIamMemberArgs)(nil)).Elem()
}

type IamServiceAccountIamMemberInput interface {
	pulumi.Input

	ToIamServiceAccountIamMemberOutput() IamServiceAccountIamMemberOutput
	ToIamServiceAccountIamMemberOutputWithContext(ctx context.Context) IamServiceAccountIamMemberOutput
}

func (*IamServiceAccountIamMember) ElementType() reflect.Type {
	return reflect.TypeOf((**IamServiceAccountIamMember)(nil)).Elem()
}

func (i *IamServiceAccountIamMember) ToIamServiceAccountIamMemberOutput() IamServiceAccountIamMemberOutput {
	return i.ToIamServiceAccountIamMemberOutputWithContext(context.Background())
}

func (i *IamServiceAccountIamMember) ToIamServiceAccountIamMemberOutputWithContext(ctx context.Context) IamServiceAccountIamMemberOutput {
	return pulumi.ToOutputWithContext(ctx, i).(IamServiceAccountIamMemberOutput)
}

// IamServiceAccountIamMemberArrayInput is an input type that accepts IamServiceAccountIamMemberArray and IamServiceAccountIamMemberArrayOutput values.
// You can construct a concrete instance of `IamServiceAccountIamMemberArrayInput` via:
//
//          IamServiceAccountIamMemberArray{ IamServiceAccountIamMemberArgs{...} }
type IamServiceAccountIamMemberArrayInput interface {
	pulumi.Input

	ToIamServiceAccountIamMemberArrayOutput() IamServiceAccountIamMemberArrayOutput
	ToIamServiceAccountIamMemberArrayOutputWithContext(context.Context) IamServiceAccountIamMemberArrayOutput
}

type IamServiceAccountIamMemberArray []IamServiceAccountIamMemberInput

func (IamServiceAccountIamMemberArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*IamServiceAccountIamMember)(nil)).Elem()
}

func (i IamServiceAccountIamMemberArray) ToIamServiceAccountIamMemberArrayOutput() IamServiceAccountIamMemberArrayOutput {
	return i.ToIamServiceAccountIamMemberArrayOutputWithContext(context.Background())
}

func (i IamServiceAccountIamMemberArray) ToIamServiceAccountIamMemberArrayOutputWithContext(ctx context.Context) IamServiceAccountIamMemberArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(IamServiceAccountIamMemberArrayOutput)
}

// IamServiceAccountIamMemberMapInput is an input type that accepts IamServiceAccountIamMemberMap and IamServiceAccountIamMemberMapOutput values.
// You can construct a concrete instance of `IamServiceAccountIamMemberMapInput` via:
//
//          IamServiceAccountIamMemberMap{ "key": IamServiceAccountIamMemberArgs{...} }
type IamServiceAccountIamMemberMapInput interface {
	pulumi.Input

	ToIamServiceAccountIamMemberMapOutput() IamServiceAccountIamMemberMapOutput
	ToIamServiceAccountIamMemberMapOutputWithContext(context.Context) IamServiceAccountIamMemberMapOutput
}

type IamServiceAccountIamMemberMap map[string]IamServiceAccountIamMemberInput

func (IamServiceAccountIamMemberMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*IamServiceAccountIamMember)(nil)).Elem()
}

func (i IamServiceAccountIamMemberMap) ToIamServiceAccountIamMemberMapOutput() IamServiceAccountIamMemberMapOutput {
	return i.ToIamServiceAccountIamMemberMapOutputWithContext(context.Background())
}

func (i IamServiceAccountIamMemberMap) ToIamServiceAccountIamMemberMapOutputWithContext(ctx context.Context) IamServiceAccountIamMemberMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(IamServiceAccountIamMemberMapOutput)
}

type IamServiceAccountIamMemberOutput struct{ *pulumi.OutputState }

func (IamServiceAccountIamMemberOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**IamServiceAccountIamMember)(nil)).Elem()
}

func (o IamServiceAccountIamMemberOutput) ToIamServiceAccountIamMemberOutput() IamServiceAccountIamMemberOutput {
	return o
}

func (o IamServiceAccountIamMemberOutput) ToIamServiceAccountIamMemberOutputWithContext(ctx context.Context) IamServiceAccountIamMemberOutput {
	return o
}

// Identity that will be granted the privilege in `role`.
// Entry can have one of the following values:
// * **userAccount:{user_id}**: A unique user ID that represents a specific Yandex account.
// * **serviceAccount:{service_account_id}**: A unique service account ID.
func (o IamServiceAccountIamMemberOutput) Member() pulumi.StringOutput {
	return o.ApplyT(func(v *IamServiceAccountIamMember) pulumi.StringOutput { return v.Member }).(pulumi.StringOutput)
}

// The role that should be applied. Only one
// `IamServiceAccountIamBinding` can be used per role.
func (o IamServiceAccountIamMemberOutput) Role() pulumi.StringOutput {
	return o.ApplyT(func(v *IamServiceAccountIamMember) pulumi.StringOutput { return v.Role }).(pulumi.StringOutput)
}

// The service account ID to apply a policy to.
func (o IamServiceAccountIamMemberOutput) ServiceAccountId() pulumi.StringOutput {
	return o.ApplyT(func(v *IamServiceAccountIamMember) pulumi.StringOutput { return v.ServiceAccountId }).(pulumi.StringOutput)
}

func (o IamServiceAccountIamMemberOutput) SleepAfter() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v *IamServiceAccountIamMember) pulumi.Float64PtrOutput { return v.SleepAfter }).(pulumi.Float64PtrOutput)
}

type IamServiceAccountIamMemberArrayOutput struct{ *pulumi.OutputState }

func (IamServiceAccountIamMemberArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*IamServiceAccountIamMember)(nil)).Elem()
}

func (o IamServiceAccountIamMemberArrayOutput) ToIamServiceAccountIamMemberArrayOutput() IamServiceAccountIamMemberArrayOutput {
	return o
}

func (o IamServiceAccountIamMemberArrayOutput) ToIamServiceAccountIamMemberArrayOutputWithContext(ctx context.Context) IamServiceAccountIamMemberArrayOutput {
	return o
}

func (o IamServiceAccountIamMemberArrayOutput) Index(i pulumi.IntInput) IamServiceAccountIamMemberOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *IamServiceAccountIamMember {
		return vs[0].([]*IamServiceAccountIamMember)[vs[1].(int)]
	}).(IamServiceAccountIamMemberOutput)
}

type IamServiceAccountIamMemberMapOutput struct{ *pulumi.OutputState }

func (IamServiceAccountIamMemberMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*IamServiceAccountIamMember)(nil)).Elem()
}

func (o IamServiceAccountIamMemberMapOutput) ToIamServiceAccountIamMemberMapOutput() IamServiceAccountIamMemberMapOutput {
	return o
}

func (o IamServiceAccountIamMemberMapOutput) ToIamServiceAccountIamMemberMapOutputWithContext(ctx context.Context) IamServiceAccountIamMemberMapOutput {
	return o
}

func (o IamServiceAccountIamMemberMapOutput) MapIndex(k pulumi.StringInput) IamServiceAccountIamMemberOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *IamServiceAccountIamMember {
		return vs[0].(map[string]*IamServiceAccountIamMember)[vs[1].(string)]
	}).(IamServiceAccountIamMemberOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*IamServiceAccountIamMemberInput)(nil)).Elem(), &IamServiceAccountIamMember{})
	pulumi.RegisterInputType(reflect.TypeOf((*IamServiceAccountIamMemberArrayInput)(nil)).Elem(), IamServiceAccountIamMemberArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*IamServiceAccountIamMemberMapInput)(nil)).Elem(), IamServiceAccountIamMemberMap{})
	pulumi.RegisterOutputType(IamServiceAccountIamMemberOutput{})
	pulumi.RegisterOutputType(IamServiceAccountIamMemberArrayOutput{})
	pulumi.RegisterOutputType(IamServiceAccountIamMemberMapOutput{})
}
