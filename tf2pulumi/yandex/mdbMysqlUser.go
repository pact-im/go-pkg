// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Manages a MySQL user within the Yandex.Cloud. For more information, see
// [the official documentation](https://cloud.yandex.com/docs/managed-mysql/).
//
// ## Import
//
// A MySQL user can be imported using the following format
//
// ```sh
//  $ pulumi import yandex:index/mdbMysqlUser:MdbMysqlUser foo {{cluster_id}}:{{username}}
// ```
type MdbMysqlUser struct {
	pulumi.CustomResourceState

	// Authentication plugin. Allowed values: `MYSQL_NATIVE_PASSWORD`, `CACHING_SHA2_PASSWORD`, `SHA256_PASSWORD` (for version 5.7 `MYSQL_NATIVE_PASSWORD`, `SHA256_PASSWORD`)
	AuthenticationPlugin pulumi.StringOutput `pulumi:"authenticationPlugin"`
	ClusterId            pulumi.StringOutput `pulumi:"clusterId"`
	// User's connection limits. The structure is documented below.
	// If the attribute is not specified there will be no changes.
	ConnectionLimits MdbMysqlUserConnectionLimitsPtrOutput `pulumi:"connectionLimits"`
	// List user's global permissions\
	// Allowed permissions:  `REPLICATION_CLIENT`, `REPLICATION_SLAVE`, `PROCESS` for clear list use empty list.
	// If the attribute is not specified there will be no changes.
	GlobalPermissions pulumi.StringArrayOutput `pulumi:"globalPermissions"`
	// The name of the user.
	Name pulumi.StringOutput `pulumi:"name"`
	// The password of the user.
	Password pulumi.StringOutput `pulumi:"password"`
	// Set of permissions granted to the user. The structure is documented below.
	Permissions MdbMysqlUserPermissionArrayOutput `pulumi:"permissions"`
	Timeouts    MdbMysqlUserTimeoutsOutput        `pulumi:"timeouts"`
}

// NewMdbMysqlUser registers a new resource with the given unique name, arguments, and options.
func NewMdbMysqlUser(ctx *pulumi.Context,
	name string, args *MdbMysqlUserArgs, opts ...pulumi.ResourceOption) (*MdbMysqlUser, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.ClusterId == nil {
		return nil, errors.New("invalid value for required argument 'ClusterId'")
	}
	if args.Password == nil {
		return nil, errors.New("invalid value for required argument 'Password'")
	}
	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	if args.Password != nil {
		args.Password = pulumi.ToSecret(args.Password).(pulumi.StringOutput)
	}
	secrets := pulumi.AdditionalSecretOutputs([]string{
		"password",
	})
	opts = append(opts, secrets)
	var resource MdbMysqlUser
	err := ctx.RegisterResource("yandex:index/mdbMysqlUser:MdbMysqlUser", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetMdbMysqlUser gets an existing MdbMysqlUser resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetMdbMysqlUser(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *MdbMysqlUserState, opts ...pulumi.ResourceOption) (*MdbMysqlUser, error) {
	var resource MdbMysqlUser
	err := ctx.ReadResource("yandex:index/mdbMysqlUser:MdbMysqlUser", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering MdbMysqlUser resources.
type mdbMysqlUserState struct {
	// Authentication plugin. Allowed values: `MYSQL_NATIVE_PASSWORD`, `CACHING_SHA2_PASSWORD`, `SHA256_PASSWORD` (for version 5.7 `MYSQL_NATIVE_PASSWORD`, `SHA256_PASSWORD`)
	AuthenticationPlugin *string `pulumi:"authenticationPlugin"`
	ClusterId            *string `pulumi:"clusterId"`
	// User's connection limits. The structure is documented below.
	// If the attribute is not specified there will be no changes.
	ConnectionLimits *MdbMysqlUserConnectionLimits `pulumi:"connectionLimits"`
	// List user's global permissions\
	// Allowed permissions:  `REPLICATION_CLIENT`, `REPLICATION_SLAVE`, `PROCESS` for clear list use empty list.
	// If the attribute is not specified there will be no changes.
	GlobalPermissions []string `pulumi:"globalPermissions"`
	// The name of the user.
	Name *string `pulumi:"name"`
	// The password of the user.
	Password *string `pulumi:"password"`
	// Set of permissions granted to the user. The structure is documented below.
	Permissions []MdbMysqlUserPermission `pulumi:"permissions"`
	Timeouts    *MdbMysqlUserTimeouts    `pulumi:"timeouts"`
}

type MdbMysqlUserState struct {
	// Authentication plugin. Allowed values: `MYSQL_NATIVE_PASSWORD`, `CACHING_SHA2_PASSWORD`, `SHA256_PASSWORD` (for version 5.7 `MYSQL_NATIVE_PASSWORD`, `SHA256_PASSWORD`)
	AuthenticationPlugin pulumi.StringPtrInput
	ClusterId            pulumi.StringPtrInput
	// User's connection limits. The structure is documented below.
	// If the attribute is not specified there will be no changes.
	ConnectionLimits MdbMysqlUserConnectionLimitsPtrInput
	// List user's global permissions\
	// Allowed permissions:  `REPLICATION_CLIENT`, `REPLICATION_SLAVE`, `PROCESS` for clear list use empty list.
	// If the attribute is not specified there will be no changes.
	GlobalPermissions pulumi.StringArrayInput
	// The name of the user.
	Name pulumi.StringPtrInput
	// The password of the user.
	Password pulumi.StringPtrInput
	// Set of permissions granted to the user. The structure is documented below.
	Permissions MdbMysqlUserPermissionArrayInput
	Timeouts    MdbMysqlUserTimeoutsPtrInput
}

func (MdbMysqlUserState) ElementType() reflect.Type {
	return reflect.TypeOf((*mdbMysqlUserState)(nil)).Elem()
}

type mdbMysqlUserArgs struct {
	// Authentication plugin. Allowed values: `MYSQL_NATIVE_PASSWORD`, `CACHING_SHA2_PASSWORD`, `SHA256_PASSWORD` (for version 5.7 `MYSQL_NATIVE_PASSWORD`, `SHA256_PASSWORD`)
	AuthenticationPlugin *string `pulumi:"authenticationPlugin"`
	ClusterId            string  `pulumi:"clusterId"`
	// User's connection limits. The structure is documented below.
	// If the attribute is not specified there will be no changes.
	ConnectionLimits *MdbMysqlUserConnectionLimits `pulumi:"connectionLimits"`
	// List user's global permissions\
	// Allowed permissions:  `REPLICATION_CLIENT`, `REPLICATION_SLAVE`, `PROCESS` for clear list use empty list.
	// If the attribute is not specified there will be no changes.
	GlobalPermissions []string `pulumi:"globalPermissions"`
	// The name of the user.
	Name *string `pulumi:"name"`
	// The password of the user.
	Password string `pulumi:"password"`
	// Set of permissions granted to the user. The structure is documented below.
	Permissions []MdbMysqlUserPermission `pulumi:"permissions"`
	Timeouts    MdbMysqlUserTimeouts     `pulumi:"timeouts"`
}

// The set of arguments for constructing a MdbMysqlUser resource.
type MdbMysqlUserArgs struct {
	// Authentication plugin. Allowed values: `MYSQL_NATIVE_PASSWORD`, `CACHING_SHA2_PASSWORD`, `SHA256_PASSWORD` (for version 5.7 `MYSQL_NATIVE_PASSWORD`, `SHA256_PASSWORD`)
	AuthenticationPlugin pulumi.StringPtrInput
	ClusterId            pulumi.StringInput
	// User's connection limits. The structure is documented below.
	// If the attribute is not specified there will be no changes.
	ConnectionLimits MdbMysqlUserConnectionLimitsPtrInput
	// List user's global permissions\
	// Allowed permissions:  `REPLICATION_CLIENT`, `REPLICATION_SLAVE`, `PROCESS` for clear list use empty list.
	// If the attribute is not specified there will be no changes.
	GlobalPermissions pulumi.StringArrayInput
	// The name of the user.
	Name pulumi.StringPtrInput
	// The password of the user.
	Password pulumi.StringInput
	// Set of permissions granted to the user. The structure is documented below.
	Permissions MdbMysqlUserPermissionArrayInput
	Timeouts    MdbMysqlUserTimeoutsInput
}

func (MdbMysqlUserArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*mdbMysqlUserArgs)(nil)).Elem()
}

type MdbMysqlUserInput interface {
	pulumi.Input

	ToMdbMysqlUserOutput() MdbMysqlUserOutput
	ToMdbMysqlUserOutputWithContext(ctx context.Context) MdbMysqlUserOutput
}

func (*MdbMysqlUser) ElementType() reflect.Type {
	return reflect.TypeOf((**MdbMysqlUser)(nil)).Elem()
}

func (i *MdbMysqlUser) ToMdbMysqlUserOutput() MdbMysqlUserOutput {
	return i.ToMdbMysqlUserOutputWithContext(context.Background())
}

func (i *MdbMysqlUser) ToMdbMysqlUserOutputWithContext(ctx context.Context) MdbMysqlUserOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MdbMysqlUserOutput)
}

// MdbMysqlUserArrayInput is an input type that accepts MdbMysqlUserArray and MdbMysqlUserArrayOutput values.
// You can construct a concrete instance of `MdbMysqlUserArrayInput` via:
//
//          MdbMysqlUserArray{ MdbMysqlUserArgs{...} }
type MdbMysqlUserArrayInput interface {
	pulumi.Input

	ToMdbMysqlUserArrayOutput() MdbMysqlUserArrayOutput
	ToMdbMysqlUserArrayOutputWithContext(context.Context) MdbMysqlUserArrayOutput
}

type MdbMysqlUserArray []MdbMysqlUserInput

func (MdbMysqlUserArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*MdbMysqlUser)(nil)).Elem()
}

func (i MdbMysqlUserArray) ToMdbMysqlUserArrayOutput() MdbMysqlUserArrayOutput {
	return i.ToMdbMysqlUserArrayOutputWithContext(context.Background())
}

func (i MdbMysqlUserArray) ToMdbMysqlUserArrayOutputWithContext(ctx context.Context) MdbMysqlUserArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MdbMysqlUserArrayOutput)
}

// MdbMysqlUserMapInput is an input type that accepts MdbMysqlUserMap and MdbMysqlUserMapOutput values.
// You can construct a concrete instance of `MdbMysqlUserMapInput` via:
//
//          MdbMysqlUserMap{ "key": MdbMysqlUserArgs{...} }
type MdbMysqlUserMapInput interface {
	pulumi.Input

	ToMdbMysqlUserMapOutput() MdbMysqlUserMapOutput
	ToMdbMysqlUserMapOutputWithContext(context.Context) MdbMysqlUserMapOutput
}

type MdbMysqlUserMap map[string]MdbMysqlUserInput

func (MdbMysqlUserMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*MdbMysqlUser)(nil)).Elem()
}

func (i MdbMysqlUserMap) ToMdbMysqlUserMapOutput() MdbMysqlUserMapOutput {
	return i.ToMdbMysqlUserMapOutputWithContext(context.Background())
}

func (i MdbMysqlUserMap) ToMdbMysqlUserMapOutputWithContext(ctx context.Context) MdbMysqlUserMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MdbMysqlUserMapOutput)
}

type MdbMysqlUserOutput struct{ *pulumi.OutputState }

func (MdbMysqlUserOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**MdbMysqlUser)(nil)).Elem()
}

func (o MdbMysqlUserOutput) ToMdbMysqlUserOutput() MdbMysqlUserOutput {
	return o
}

func (o MdbMysqlUserOutput) ToMdbMysqlUserOutputWithContext(ctx context.Context) MdbMysqlUserOutput {
	return o
}

// Authentication plugin. Allowed values: `MYSQL_NATIVE_PASSWORD`, `CACHING_SHA2_PASSWORD`, `SHA256_PASSWORD` (for version 5.7 `MYSQL_NATIVE_PASSWORD`, `SHA256_PASSWORD`)
func (o MdbMysqlUserOutput) AuthenticationPlugin() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbMysqlUser) pulumi.StringOutput { return v.AuthenticationPlugin }).(pulumi.StringOutput)
}

func (o MdbMysqlUserOutput) ClusterId() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbMysqlUser) pulumi.StringOutput { return v.ClusterId }).(pulumi.StringOutput)
}

// User's connection limits. The structure is documented below.
// If the attribute is not specified there will be no changes.
func (o MdbMysqlUserOutput) ConnectionLimits() MdbMysqlUserConnectionLimitsPtrOutput {
	return o.ApplyT(func(v *MdbMysqlUser) MdbMysqlUserConnectionLimitsPtrOutput { return v.ConnectionLimits }).(MdbMysqlUserConnectionLimitsPtrOutput)
}

// List user's global permissions\
// Allowed permissions:  `REPLICATION_CLIENT`, `REPLICATION_SLAVE`, `PROCESS` for clear list use empty list.
// If the attribute is not specified there will be no changes.
func (o MdbMysqlUserOutput) GlobalPermissions() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *MdbMysqlUser) pulumi.StringArrayOutput { return v.GlobalPermissions }).(pulumi.StringArrayOutput)
}

// The name of the user.
func (o MdbMysqlUserOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbMysqlUser) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// The password of the user.
func (o MdbMysqlUserOutput) Password() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbMysqlUser) pulumi.StringOutput { return v.Password }).(pulumi.StringOutput)
}

// Set of permissions granted to the user. The structure is documented below.
func (o MdbMysqlUserOutput) Permissions() MdbMysqlUserPermissionArrayOutput {
	return o.ApplyT(func(v *MdbMysqlUser) MdbMysqlUserPermissionArrayOutput { return v.Permissions }).(MdbMysqlUserPermissionArrayOutput)
}

func (o MdbMysqlUserOutput) Timeouts() MdbMysqlUserTimeoutsOutput {
	return o.ApplyT(func(v *MdbMysqlUser) MdbMysqlUserTimeoutsOutput { return v.Timeouts }).(MdbMysqlUserTimeoutsOutput)
}

type MdbMysqlUserArrayOutput struct{ *pulumi.OutputState }

func (MdbMysqlUserArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*MdbMysqlUser)(nil)).Elem()
}

func (o MdbMysqlUserArrayOutput) ToMdbMysqlUserArrayOutput() MdbMysqlUserArrayOutput {
	return o
}

func (o MdbMysqlUserArrayOutput) ToMdbMysqlUserArrayOutputWithContext(ctx context.Context) MdbMysqlUserArrayOutput {
	return o
}

func (o MdbMysqlUserArrayOutput) Index(i pulumi.IntInput) MdbMysqlUserOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *MdbMysqlUser {
		return vs[0].([]*MdbMysqlUser)[vs[1].(int)]
	}).(MdbMysqlUserOutput)
}

type MdbMysqlUserMapOutput struct{ *pulumi.OutputState }

func (MdbMysqlUserMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*MdbMysqlUser)(nil)).Elem()
}

func (o MdbMysqlUserMapOutput) ToMdbMysqlUserMapOutput() MdbMysqlUserMapOutput {
	return o
}

func (o MdbMysqlUserMapOutput) ToMdbMysqlUserMapOutputWithContext(ctx context.Context) MdbMysqlUserMapOutput {
	return o
}

func (o MdbMysqlUserMapOutput) MapIndex(k pulumi.StringInput) MdbMysqlUserOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *MdbMysqlUser {
		return vs[0].(map[string]*MdbMysqlUser)[vs[1].(string)]
	}).(MdbMysqlUserOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*MdbMysqlUserInput)(nil)).Elem(), &MdbMysqlUser{})
	pulumi.RegisterInputType(reflect.TypeOf((*MdbMysqlUserArrayInput)(nil)).Elem(), MdbMysqlUserArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*MdbMysqlUserMapInput)(nil)).Elem(), MdbMysqlUserMap{})
	pulumi.RegisterOutputType(MdbMysqlUserOutput{})
	pulumi.RegisterOutputType(MdbMysqlUserArrayOutput{})
	pulumi.RegisterOutputType(MdbMysqlUserMapOutput{})
}
