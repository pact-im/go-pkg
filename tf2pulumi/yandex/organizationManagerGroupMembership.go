// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Allows members management of a single Yandex.Cloud Organization Manager Group. For more information, see [the official documentation](https://cloud.yandex.com/en-ru/docs/organization/manage-groups#add-member).
//
// > **Note:** Multiple `yandexOrganizationmanagerGroupIamBinding` resources with the same group id will produce inconsistent behavior!
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
// 		_, err := yandex.NewOrganizationManagerGroupMembership(ctx, "group", &yandex.OrganizationManagerGroupMembershipArgs{
// 			GroupId: pulumi.String("sdf4*********3fr"),
// 			Members: pulumi.StringArray{
// 				pulumi.String("xdf********123"),
// 			},
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
type OrganizationManagerGroupMembership struct {
	pulumi.CustomResourceState

	// The Group to add/remove members to/from.
	GroupId pulumi.StringOutput `pulumi:"groupId"`
	// A set of members of the Group. Each member is represented by an id.
	Members  pulumi.StringArrayOutput                         `pulumi:"members"`
	Timeouts OrganizationManagerGroupMembershipTimeoutsOutput `pulumi:"timeouts"`
}

// NewOrganizationManagerGroupMembership registers a new resource with the given unique name, arguments, and options.
func NewOrganizationManagerGroupMembership(ctx *pulumi.Context,
	name string, args *OrganizationManagerGroupMembershipArgs, opts ...pulumi.ResourceOption) (*OrganizationManagerGroupMembership, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.GroupId == nil {
		return nil, errors.New("invalid value for required argument 'GroupId'")
	}
	if args.Members == nil {
		return nil, errors.New("invalid value for required argument 'Members'")
	}
	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	var resource OrganizationManagerGroupMembership
	err := ctx.RegisterResource("yandex:index/organizationManagerGroupMembership:OrganizationManagerGroupMembership", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetOrganizationManagerGroupMembership gets an existing OrganizationManagerGroupMembership resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetOrganizationManagerGroupMembership(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *OrganizationManagerGroupMembershipState, opts ...pulumi.ResourceOption) (*OrganizationManagerGroupMembership, error) {
	var resource OrganizationManagerGroupMembership
	err := ctx.ReadResource("yandex:index/organizationManagerGroupMembership:OrganizationManagerGroupMembership", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering OrganizationManagerGroupMembership resources.
type organizationManagerGroupMembershipState struct {
	// The Group to add/remove members to/from.
	GroupId *string `pulumi:"groupId"`
	// A set of members of the Group. Each member is represented by an id.
	Members  []string                                    `pulumi:"members"`
	Timeouts *OrganizationManagerGroupMembershipTimeouts `pulumi:"timeouts"`
}

type OrganizationManagerGroupMembershipState struct {
	// The Group to add/remove members to/from.
	GroupId pulumi.StringPtrInput
	// A set of members of the Group. Each member is represented by an id.
	Members  pulumi.StringArrayInput
	Timeouts OrganizationManagerGroupMembershipTimeoutsPtrInput
}

func (OrganizationManagerGroupMembershipState) ElementType() reflect.Type {
	return reflect.TypeOf((*organizationManagerGroupMembershipState)(nil)).Elem()
}

type organizationManagerGroupMembershipArgs struct {
	// The Group to add/remove members to/from.
	GroupId string `pulumi:"groupId"`
	// A set of members of the Group. Each member is represented by an id.
	Members  []string                                   `pulumi:"members"`
	Timeouts OrganizationManagerGroupMembershipTimeouts `pulumi:"timeouts"`
}

// The set of arguments for constructing a OrganizationManagerGroupMembership resource.
type OrganizationManagerGroupMembershipArgs struct {
	// The Group to add/remove members to/from.
	GroupId pulumi.StringInput
	// A set of members of the Group. Each member is represented by an id.
	Members  pulumi.StringArrayInput
	Timeouts OrganizationManagerGroupMembershipTimeoutsInput
}

func (OrganizationManagerGroupMembershipArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*organizationManagerGroupMembershipArgs)(nil)).Elem()
}

type OrganizationManagerGroupMembershipInput interface {
	pulumi.Input

	ToOrganizationManagerGroupMembershipOutput() OrganizationManagerGroupMembershipOutput
	ToOrganizationManagerGroupMembershipOutputWithContext(ctx context.Context) OrganizationManagerGroupMembershipOutput
}

func (*OrganizationManagerGroupMembership) ElementType() reflect.Type {
	return reflect.TypeOf((**OrganizationManagerGroupMembership)(nil)).Elem()
}

func (i *OrganizationManagerGroupMembership) ToOrganizationManagerGroupMembershipOutput() OrganizationManagerGroupMembershipOutput {
	return i.ToOrganizationManagerGroupMembershipOutputWithContext(context.Background())
}

func (i *OrganizationManagerGroupMembership) ToOrganizationManagerGroupMembershipOutputWithContext(ctx context.Context) OrganizationManagerGroupMembershipOutput {
	return pulumi.ToOutputWithContext(ctx, i).(OrganizationManagerGroupMembershipOutput)
}

// OrganizationManagerGroupMembershipArrayInput is an input type that accepts OrganizationManagerGroupMembershipArray and OrganizationManagerGroupMembershipArrayOutput values.
// You can construct a concrete instance of `OrganizationManagerGroupMembershipArrayInput` via:
//
//          OrganizationManagerGroupMembershipArray{ OrganizationManagerGroupMembershipArgs{...} }
type OrganizationManagerGroupMembershipArrayInput interface {
	pulumi.Input

	ToOrganizationManagerGroupMembershipArrayOutput() OrganizationManagerGroupMembershipArrayOutput
	ToOrganizationManagerGroupMembershipArrayOutputWithContext(context.Context) OrganizationManagerGroupMembershipArrayOutput
}

type OrganizationManagerGroupMembershipArray []OrganizationManagerGroupMembershipInput

func (OrganizationManagerGroupMembershipArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*OrganizationManagerGroupMembership)(nil)).Elem()
}

func (i OrganizationManagerGroupMembershipArray) ToOrganizationManagerGroupMembershipArrayOutput() OrganizationManagerGroupMembershipArrayOutput {
	return i.ToOrganizationManagerGroupMembershipArrayOutputWithContext(context.Background())
}

func (i OrganizationManagerGroupMembershipArray) ToOrganizationManagerGroupMembershipArrayOutputWithContext(ctx context.Context) OrganizationManagerGroupMembershipArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(OrganizationManagerGroupMembershipArrayOutput)
}

// OrganizationManagerGroupMembershipMapInput is an input type that accepts OrganizationManagerGroupMembershipMap and OrganizationManagerGroupMembershipMapOutput values.
// You can construct a concrete instance of `OrganizationManagerGroupMembershipMapInput` via:
//
//          OrganizationManagerGroupMembershipMap{ "key": OrganizationManagerGroupMembershipArgs{...} }
type OrganizationManagerGroupMembershipMapInput interface {
	pulumi.Input

	ToOrganizationManagerGroupMembershipMapOutput() OrganizationManagerGroupMembershipMapOutput
	ToOrganizationManagerGroupMembershipMapOutputWithContext(context.Context) OrganizationManagerGroupMembershipMapOutput
}

type OrganizationManagerGroupMembershipMap map[string]OrganizationManagerGroupMembershipInput

func (OrganizationManagerGroupMembershipMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*OrganizationManagerGroupMembership)(nil)).Elem()
}

func (i OrganizationManagerGroupMembershipMap) ToOrganizationManagerGroupMembershipMapOutput() OrganizationManagerGroupMembershipMapOutput {
	return i.ToOrganizationManagerGroupMembershipMapOutputWithContext(context.Background())
}

func (i OrganizationManagerGroupMembershipMap) ToOrganizationManagerGroupMembershipMapOutputWithContext(ctx context.Context) OrganizationManagerGroupMembershipMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(OrganizationManagerGroupMembershipMapOutput)
}

type OrganizationManagerGroupMembershipOutput struct{ *pulumi.OutputState }

func (OrganizationManagerGroupMembershipOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**OrganizationManagerGroupMembership)(nil)).Elem()
}

func (o OrganizationManagerGroupMembershipOutput) ToOrganizationManagerGroupMembershipOutput() OrganizationManagerGroupMembershipOutput {
	return o
}

func (o OrganizationManagerGroupMembershipOutput) ToOrganizationManagerGroupMembershipOutputWithContext(ctx context.Context) OrganizationManagerGroupMembershipOutput {
	return o
}

// The Group to add/remove members to/from.
func (o OrganizationManagerGroupMembershipOutput) GroupId() pulumi.StringOutput {
	return o.ApplyT(func(v *OrganizationManagerGroupMembership) pulumi.StringOutput { return v.GroupId }).(pulumi.StringOutput)
}

// A set of members of the Group. Each member is represented by an id.
func (o OrganizationManagerGroupMembershipOutput) Members() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *OrganizationManagerGroupMembership) pulumi.StringArrayOutput { return v.Members }).(pulumi.StringArrayOutput)
}

func (o OrganizationManagerGroupMembershipOutput) Timeouts() OrganizationManagerGroupMembershipTimeoutsOutput {
	return o.ApplyT(func(v *OrganizationManagerGroupMembership) OrganizationManagerGroupMembershipTimeoutsOutput {
		return v.Timeouts
	}).(OrganizationManagerGroupMembershipTimeoutsOutput)
}

type OrganizationManagerGroupMembershipArrayOutput struct{ *pulumi.OutputState }

func (OrganizationManagerGroupMembershipArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*OrganizationManagerGroupMembership)(nil)).Elem()
}

func (o OrganizationManagerGroupMembershipArrayOutput) ToOrganizationManagerGroupMembershipArrayOutput() OrganizationManagerGroupMembershipArrayOutput {
	return o
}

func (o OrganizationManagerGroupMembershipArrayOutput) ToOrganizationManagerGroupMembershipArrayOutputWithContext(ctx context.Context) OrganizationManagerGroupMembershipArrayOutput {
	return o
}

func (o OrganizationManagerGroupMembershipArrayOutput) Index(i pulumi.IntInput) OrganizationManagerGroupMembershipOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *OrganizationManagerGroupMembership {
		return vs[0].([]*OrganizationManagerGroupMembership)[vs[1].(int)]
	}).(OrganizationManagerGroupMembershipOutput)
}

type OrganizationManagerGroupMembershipMapOutput struct{ *pulumi.OutputState }

func (OrganizationManagerGroupMembershipMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*OrganizationManagerGroupMembership)(nil)).Elem()
}

func (o OrganizationManagerGroupMembershipMapOutput) ToOrganizationManagerGroupMembershipMapOutput() OrganizationManagerGroupMembershipMapOutput {
	return o
}

func (o OrganizationManagerGroupMembershipMapOutput) ToOrganizationManagerGroupMembershipMapOutputWithContext(ctx context.Context) OrganizationManagerGroupMembershipMapOutput {
	return o
}

func (o OrganizationManagerGroupMembershipMapOutput) MapIndex(k pulumi.StringInput) OrganizationManagerGroupMembershipOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *OrganizationManagerGroupMembership {
		return vs[0].(map[string]*OrganizationManagerGroupMembership)[vs[1].(string)]
	}).(OrganizationManagerGroupMembershipOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*OrganizationManagerGroupMembershipInput)(nil)).Elem(), &OrganizationManagerGroupMembership{})
	pulumi.RegisterInputType(reflect.TypeOf((*OrganizationManagerGroupMembershipArrayInput)(nil)).Elem(), OrganizationManagerGroupMembershipArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*OrganizationManagerGroupMembershipMapInput)(nil)).Elem(), OrganizationManagerGroupMembershipMap{})
	pulumi.RegisterOutputType(OrganizationManagerGroupMembershipOutput{})
	pulumi.RegisterOutputType(OrganizationManagerGroupMembershipArrayOutput{})
	pulumi.RegisterOutputType(OrganizationManagerGroupMembershipMapOutput{})
}
