// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Get information about a Yandex Managed PostgreSQL user. For more information, see
// [the official documentation](https://cloud.yandex.com/docs/managed-postgresql/).
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
// 		foo, err := yandex.LookupMdbPostgresqlUser(ctx, &GetMdbPostgresqlUserArgs{
// 			ClusterId: "some_cluster_id",
// 			Name:      "test",
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		ctx.Export("permission", foo.Permissions)
// 		return nil
// 	})
// }
// ```
func LookupMdbPostgresqlUser(ctx *pulumi.Context, args *LookupMdbPostgresqlUserArgs, opts ...pulumi.InvokeOption) (*LookupMdbPostgresqlUserResult, error) {
	var rv LookupMdbPostgresqlUserResult
	err := ctx.Invoke("yandex:index/getMdbPostgresqlUser:getMdbPostgresqlUser", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getMdbPostgresqlUser.
type LookupMdbPostgresqlUserArgs struct {
	// The ID of the PostgreSQL cluster.
	ClusterId string `pulumi:"clusterId"`
	// The maximum number of connections per user.
	ConnLimit *float64 `pulumi:"connLimit"`
	// User's ability to login.
	Login *bool `pulumi:"login"`
	// The name of the PostgreSQL user.
	Name string `pulumi:"name"`
	// Map of user settings.
	Settings map[string]string `pulumi:"settings"`
}

// A collection of values returned by getMdbPostgresqlUser.
type LookupMdbPostgresqlUserResult struct {
	ClusterId string `pulumi:"clusterId"`
	// The maximum number of connections per user.
	ConnLimit *float64 `pulumi:"connLimit"`
	// List of the user's grants.
	Grants []string `pulumi:"grants"`
	Id     string   `pulumi:"id"`
	// User's ability to login.
	Login *bool  `pulumi:"login"`
	Name  string `pulumi:"name"`
	// The password of the user.
	Password string `pulumi:"password"`
	// Set of permissions granted to the user. The structure is documented below.
	Permissions []GetMdbPostgresqlUserPermission `pulumi:"permissions"`
	// Map of user settings.
	Settings map[string]string `pulumi:"settings"`
}

func LookupMdbPostgresqlUserOutput(ctx *pulumi.Context, args LookupMdbPostgresqlUserOutputArgs, opts ...pulumi.InvokeOption) LookupMdbPostgresqlUserResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (LookupMdbPostgresqlUserResult, error) {
			args := v.(LookupMdbPostgresqlUserArgs)
			r, err := LookupMdbPostgresqlUser(ctx, &args, opts...)
			var s LookupMdbPostgresqlUserResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(LookupMdbPostgresqlUserResultOutput)
}

// A collection of arguments for invoking getMdbPostgresqlUser.
type LookupMdbPostgresqlUserOutputArgs struct {
	// The ID of the PostgreSQL cluster.
	ClusterId pulumi.StringInput `pulumi:"clusterId"`
	// The maximum number of connections per user.
	ConnLimit pulumi.Float64PtrInput `pulumi:"connLimit"`
	// User's ability to login.
	Login pulumi.BoolPtrInput `pulumi:"login"`
	// The name of the PostgreSQL user.
	Name pulumi.StringInput `pulumi:"name"`
	// Map of user settings.
	Settings pulumi.StringMapInput `pulumi:"settings"`
}

func (LookupMdbPostgresqlUserOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupMdbPostgresqlUserArgs)(nil)).Elem()
}

// A collection of values returned by getMdbPostgresqlUser.
type LookupMdbPostgresqlUserResultOutput struct{ *pulumi.OutputState }

func (LookupMdbPostgresqlUserResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupMdbPostgresqlUserResult)(nil)).Elem()
}

func (o LookupMdbPostgresqlUserResultOutput) ToLookupMdbPostgresqlUserResultOutput() LookupMdbPostgresqlUserResultOutput {
	return o
}

func (o LookupMdbPostgresqlUserResultOutput) ToLookupMdbPostgresqlUserResultOutputWithContext(ctx context.Context) LookupMdbPostgresqlUserResultOutput {
	return o
}

func (o LookupMdbPostgresqlUserResultOutput) ClusterId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbPostgresqlUserResult) string { return v.ClusterId }).(pulumi.StringOutput)
}

// The maximum number of connections per user.
func (o LookupMdbPostgresqlUserResultOutput) ConnLimit() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v LookupMdbPostgresqlUserResult) *float64 { return v.ConnLimit }).(pulumi.Float64PtrOutput)
}

// List of the user's grants.
func (o LookupMdbPostgresqlUserResultOutput) Grants() pulumi.StringArrayOutput {
	return o.ApplyT(func(v LookupMdbPostgresqlUserResult) []string { return v.Grants }).(pulumi.StringArrayOutput)
}

func (o LookupMdbPostgresqlUserResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbPostgresqlUserResult) string { return v.Id }).(pulumi.StringOutput)
}

// User's ability to login.
func (o LookupMdbPostgresqlUserResultOutput) Login() pulumi.BoolPtrOutput {
	return o.ApplyT(func(v LookupMdbPostgresqlUserResult) *bool { return v.Login }).(pulumi.BoolPtrOutput)
}

func (o LookupMdbPostgresqlUserResultOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbPostgresqlUserResult) string { return v.Name }).(pulumi.StringOutput)
}

// The password of the user.
func (o LookupMdbPostgresqlUserResultOutput) Password() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbPostgresqlUserResult) string { return v.Password }).(pulumi.StringOutput)
}

// Set of permissions granted to the user. The structure is documented below.
func (o LookupMdbPostgresqlUserResultOutput) Permissions() GetMdbPostgresqlUserPermissionArrayOutput {
	return o.ApplyT(func(v LookupMdbPostgresqlUserResult) []GetMdbPostgresqlUserPermission { return v.Permissions }).(GetMdbPostgresqlUserPermissionArrayOutput)
}

// Map of user settings.
func (o LookupMdbPostgresqlUserResultOutput) Settings() pulumi.StringMapOutput {
	return o.ApplyT(func(v LookupMdbPostgresqlUserResult) map[string]string { return v.Settings }).(pulumi.StringMapOutput)
}

func init() {
	pulumi.RegisterOutputType(LookupMdbPostgresqlUserResultOutput{})
}
