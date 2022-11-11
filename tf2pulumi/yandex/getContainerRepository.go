// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Get information about a Yandex Container Repository. For more information, see
// [the official documentation](https://cloud.yandex.com/docs/container-registry/concepts/repository)
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
// 		_, err = yandex.LookupContainerRepository(ctx, &GetContainerRepositoryArgs{
// 			Name: pulumi.StringRef("some_repository_name"),
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		_, err = yandex.LookupContainerRepository(ctx, &GetContainerRepositoryArgs{
// 			RepositoryId: pulumi.StringRef("some_repository_id"),
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
func LookupContainerRepository(ctx *pulumi.Context, args *LookupContainerRepositoryArgs, opts ...pulumi.InvokeOption) (*LookupContainerRepositoryResult, error) {
	var rv LookupContainerRepositoryResult
	err := ctx.Invoke("yandex:index/getContainerRepository:getContainerRepository", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getContainerRepository.
type LookupContainerRepositoryArgs struct {
	// Name of the repository. The name of the repository should start with id of a container registry and match the name of the images in the repository.
	Name *string `pulumi:"name"`
	// The ID of a specific repository.
	RepositoryId *string `pulumi:"repositoryId"`
}

// A collection of values returned by getContainerRepository.
type LookupContainerRepositoryResult struct {
	Id           string `pulumi:"id"`
	Name         string `pulumi:"name"`
	RepositoryId string `pulumi:"repositoryId"`
}

func LookupContainerRepositoryOutput(ctx *pulumi.Context, args LookupContainerRepositoryOutputArgs, opts ...pulumi.InvokeOption) LookupContainerRepositoryResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (LookupContainerRepositoryResult, error) {
			args := v.(LookupContainerRepositoryArgs)
			r, err := LookupContainerRepository(ctx, &args, opts...)
			var s LookupContainerRepositoryResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(LookupContainerRepositoryResultOutput)
}

// A collection of arguments for invoking getContainerRepository.
type LookupContainerRepositoryOutputArgs struct {
	// Name of the repository. The name of the repository should start with id of a container registry and match the name of the images in the repository.
	Name pulumi.StringPtrInput `pulumi:"name"`
	// The ID of a specific repository.
	RepositoryId pulumi.StringPtrInput `pulumi:"repositoryId"`
}

func (LookupContainerRepositoryOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupContainerRepositoryArgs)(nil)).Elem()
}

// A collection of values returned by getContainerRepository.
type LookupContainerRepositoryResultOutput struct{ *pulumi.OutputState }

func (LookupContainerRepositoryResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupContainerRepositoryResult)(nil)).Elem()
}

func (o LookupContainerRepositoryResultOutput) ToLookupContainerRepositoryResultOutput() LookupContainerRepositoryResultOutput {
	return o
}

func (o LookupContainerRepositoryResultOutput) ToLookupContainerRepositoryResultOutputWithContext(ctx context.Context) LookupContainerRepositoryResultOutput {
	return o
}

func (o LookupContainerRepositoryResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v LookupContainerRepositoryResult) string { return v.Id }).(pulumi.StringOutput)
}

func (o LookupContainerRepositoryResultOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v LookupContainerRepositoryResult) string { return v.Name }).(pulumi.StringOutput)
}

func (o LookupContainerRepositoryResultOutput) RepositoryId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupContainerRepositoryResult) string { return v.RepositoryId }).(pulumi.StringOutput)
}

func init() {
	pulumi.RegisterOutputType(LookupContainerRepositoryResultOutput{})
}