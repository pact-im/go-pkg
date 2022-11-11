// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Get information about a Yandex Cloud Function. For more information about Yandex Cloud Functions, see
// [Yandex Cloud Functions](https://cloud.yandex.com/docs/functions/).
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
// 		_, err = yandex.LookupFunction(ctx, &GetFunctionArgs{
// 			FunctionId: pulumi.StringRef("are1samplefunction11"),
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
//
// This data source is used to define [Yandex Cloud Function](https://cloud.yandex.com/docs/functions/concepts/function) that can be used by other resources.
func LookupFunction(ctx *pulumi.Context, args *LookupFunctionArgs, opts ...pulumi.InvokeOption) (*LookupFunctionResult, error) {
	var rv LookupFunctionResult
	err := ctx.Invoke("yandex:index/getFunction:getFunction", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getFunction.
type LookupFunctionArgs struct {
	// Folder ID for the Yandex Cloud Function
	FolderId *string `pulumi:"folderId"`
	// Yandex Cloud Function id used to define function
	FunctionId *string `pulumi:"functionId"`
	// Yandex Cloud Function name used to define function
	Name *string `pulumi:"name"`
	// Secrets for Yandex Cloud Function.
	Secrets []GetFunctionSecret `pulumi:"secrets"`
}

// A collection of values returned by getFunction.
type LookupFunctionResult struct {
	// Creation timestamp of the Yandex Cloud Function
	CreatedAt string `pulumi:"createdAt"`
	// Description of the Yandex Cloud Function
	Description string `pulumi:"description"`
	// Entrypoint for Yandex Cloud Function
	Entrypoint string `pulumi:"entrypoint"`
	// A set of key/value environment variables for Yandex Cloud Function
	Environment map[string]string `pulumi:"environment"`
	// Execution timeout in seconds for Yandex Cloud Function
	ExecutionTimeout string  `pulumi:"executionTimeout"`
	FolderId         *string `pulumi:"folderId"`
	FunctionId       *string `pulumi:"functionId"`
	Id               string  `pulumi:"id"`
	// Image size for Yandex Cloud Function.
	ImageSize float64 `pulumi:"imageSize"`
	// A set of key/value label pairs to assign to the Yandex Cloud Function
	Labels map[string]string `pulumi:"labels"`
	// Log group ID size for Yandex Cloud Function.
	LoggroupId string `pulumi:"loggroupId"`
	// Memory in megabytes (**aligned to 128MB**) for Yandex Cloud Function
	Memory float64 `pulumi:"memory"`
	Name   *string `pulumi:"name"`
	// Runtime for Yandex Cloud Function
	Runtime string `pulumi:"runtime"`
	// Secrets for Yandex Cloud Function.
	Secrets []GetFunctionSecret `pulumi:"secrets"`
	// Service account ID for Yandex Cloud Function
	ServiceAccountId string `pulumi:"serviceAccountId"`
	// Tags for Yandex Cloud Function. Tag "$latest" isn't returned.
	Tags []string `pulumi:"tags"`
	// Version for Yandex Cloud Function.
	Version string `pulumi:"version"`
}

func LookupFunctionOutput(ctx *pulumi.Context, args LookupFunctionOutputArgs, opts ...pulumi.InvokeOption) LookupFunctionResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (LookupFunctionResult, error) {
			args := v.(LookupFunctionArgs)
			r, err := LookupFunction(ctx, &args, opts...)
			var s LookupFunctionResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(LookupFunctionResultOutput)
}

// A collection of arguments for invoking getFunction.
type LookupFunctionOutputArgs struct {
	// Folder ID for the Yandex Cloud Function
	FolderId pulumi.StringPtrInput `pulumi:"folderId"`
	// Yandex Cloud Function id used to define function
	FunctionId pulumi.StringPtrInput `pulumi:"functionId"`
	// Yandex Cloud Function name used to define function
	Name pulumi.StringPtrInput `pulumi:"name"`
	// Secrets for Yandex Cloud Function.
	Secrets GetFunctionSecretArrayInput `pulumi:"secrets"`
}

func (LookupFunctionOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupFunctionArgs)(nil)).Elem()
}

// A collection of values returned by getFunction.
type LookupFunctionResultOutput struct{ *pulumi.OutputState }

func (LookupFunctionResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupFunctionResult)(nil)).Elem()
}

func (o LookupFunctionResultOutput) ToLookupFunctionResultOutput() LookupFunctionResultOutput {
	return o
}

func (o LookupFunctionResultOutput) ToLookupFunctionResultOutputWithContext(ctx context.Context) LookupFunctionResultOutput {
	return o
}

// Creation timestamp of the Yandex Cloud Function
func (o LookupFunctionResultOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v LookupFunctionResult) string { return v.CreatedAt }).(pulumi.StringOutput)
}

// Description of the Yandex Cloud Function
func (o LookupFunctionResultOutput) Description() pulumi.StringOutput {
	return o.ApplyT(func(v LookupFunctionResult) string { return v.Description }).(pulumi.StringOutput)
}

// Entrypoint for Yandex Cloud Function
func (o LookupFunctionResultOutput) Entrypoint() pulumi.StringOutput {
	return o.ApplyT(func(v LookupFunctionResult) string { return v.Entrypoint }).(pulumi.StringOutput)
}

// A set of key/value environment variables for Yandex Cloud Function
func (o LookupFunctionResultOutput) Environment() pulumi.StringMapOutput {
	return o.ApplyT(func(v LookupFunctionResult) map[string]string { return v.Environment }).(pulumi.StringMapOutput)
}

// Execution timeout in seconds for Yandex Cloud Function
func (o LookupFunctionResultOutput) ExecutionTimeout() pulumi.StringOutput {
	return o.ApplyT(func(v LookupFunctionResult) string { return v.ExecutionTimeout }).(pulumi.StringOutput)
}

func (o LookupFunctionResultOutput) FolderId() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LookupFunctionResult) *string { return v.FolderId }).(pulumi.StringPtrOutput)
}

func (o LookupFunctionResultOutput) FunctionId() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LookupFunctionResult) *string { return v.FunctionId }).(pulumi.StringPtrOutput)
}

func (o LookupFunctionResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v LookupFunctionResult) string { return v.Id }).(pulumi.StringOutput)
}

// Image size for Yandex Cloud Function.
func (o LookupFunctionResultOutput) ImageSize() pulumi.Float64Output {
	return o.ApplyT(func(v LookupFunctionResult) float64 { return v.ImageSize }).(pulumi.Float64Output)
}

// A set of key/value label pairs to assign to the Yandex Cloud Function
func (o LookupFunctionResultOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v LookupFunctionResult) map[string]string { return v.Labels }).(pulumi.StringMapOutput)
}

// Log group ID size for Yandex Cloud Function.
func (o LookupFunctionResultOutput) LoggroupId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupFunctionResult) string { return v.LoggroupId }).(pulumi.StringOutput)
}

// Memory in megabytes (**aligned to 128MB**) for Yandex Cloud Function
func (o LookupFunctionResultOutput) Memory() pulumi.Float64Output {
	return o.ApplyT(func(v LookupFunctionResult) float64 { return v.Memory }).(pulumi.Float64Output)
}

func (o LookupFunctionResultOutput) Name() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LookupFunctionResult) *string { return v.Name }).(pulumi.StringPtrOutput)
}

// Runtime for Yandex Cloud Function
func (o LookupFunctionResultOutput) Runtime() pulumi.StringOutput {
	return o.ApplyT(func(v LookupFunctionResult) string { return v.Runtime }).(pulumi.StringOutput)
}

// Secrets for Yandex Cloud Function.
func (o LookupFunctionResultOutput) Secrets() GetFunctionSecretArrayOutput {
	return o.ApplyT(func(v LookupFunctionResult) []GetFunctionSecret { return v.Secrets }).(GetFunctionSecretArrayOutput)
}

// Service account ID for Yandex Cloud Function
func (o LookupFunctionResultOutput) ServiceAccountId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupFunctionResult) string { return v.ServiceAccountId }).(pulumi.StringOutput)
}

// Tags for Yandex Cloud Function. Tag "$latest" isn't returned.
func (o LookupFunctionResultOutput) Tags() pulumi.StringArrayOutput {
	return o.ApplyT(func(v LookupFunctionResult) []string { return v.Tags }).(pulumi.StringArrayOutput)
}

// Version for Yandex Cloud Function.
func (o LookupFunctionResultOutput) Version() pulumi.StringOutput {
	return o.ApplyT(func(v LookupFunctionResult) string { return v.Version }).(pulumi.StringOutput)
}

func init() {
	pulumi.RegisterOutputType(LookupFunctionResultOutput{})
}