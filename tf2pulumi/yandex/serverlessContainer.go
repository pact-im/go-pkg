// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Allows management of Yandex Cloud Serverless Containers
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
// 		_, err := yandex.NewServerlessContainer(ctx, "test-container", &yandex.ServerlessContainerArgs{
// 			CoreFraction:     pulumi.Float64(100),
// 			Cores:            pulumi.Float64(1),
// 			Description:      pulumi.String("any description"),
// 			ExecutionTimeout: pulumi.String("15s"),
// 			Image: &ServerlessContainerImageArgs{
// 				Url: pulumi.String("cr.yandex/yc/test-image:v1"),
// 			},
// 			Memory:           pulumi.Float64(256),
// 			ServiceAccountId: pulumi.String("are1service2account3id"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
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
// 		_, err := yandex.NewServerlessContainer(ctx, "test-container-with-digest", &yandex.ServerlessContainerArgs{
// 			Image: &ServerlessContainerImageArgs{
// 				Digest: pulumi.String("sha256:e1d772fa8795adac847a2420c87d0d2e3d38fb02f168cab8c0b5fe2fb95c47f4"),
// 				Url:    pulumi.String("cr.yandex/yc/test-image:v1"),
// 			},
// 			Memory: pulumi.Float64(128),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
type ServerlessContainer struct {
	pulumi.CustomResourceState

	// Concurrency of Yandex Cloud Serverless Container
	Concurrency pulumi.Float64PtrOutput `pulumi:"concurrency"`
	// Core fraction (**0...100**) of the Yandex Cloud Serverless Container
	CoreFraction pulumi.Float64Output    `pulumi:"coreFraction"`
	Cores        pulumi.Float64PtrOutput `pulumi:"cores"`
	// Creation timestamp of the Yandex Cloud Serverless Container
	CreatedAt pulumi.StringOutput `pulumi:"createdAt"`
	// Description of the Yandex Cloud Serverless Container
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// Execution timeout in seconds (**duration format**) for Yandex Cloud Serverless Container
	ExecutionTimeout pulumi.StringOutput `pulumi:"executionTimeout"`
	// Folder ID for the Yandex Cloud Serverless Container
	FolderId pulumi.StringOutput `pulumi:"folderId"`
	// Revision deployment image for Yandex Cloud Serverless Container
	// * `image.0.url` (Required) - URL of image that will be deployed as Yandex Cloud Serverless Container
	// * `image.0.work_dir` - Working directory for Yandex Cloud Serverless Container
	// * `image.0.digest` - Digest of image that will be deployed as Yandex Cloud Serverless Container.
	//   If presented, should be equal to digest that will be resolved at server side by URL.
	//   Container will be updated on digest change even if `image.0.url` stays the same.
	//   If field not specified then its value will be computed.
	// * `image.0.command` - List of commands for Yandex Cloud Serverless Container
	// * `image.0.args` - List of arguments for Yandex Cloud Serverless Container
	// * `image.0.environment` -  A set of key/value environment variable pairs for Yandex Cloud Serverless Container
	Image ServerlessContainerImageOutput `pulumi:"image"`
	// A set of key/value label pairs to assign to the Yandex Cloud Serverless Container
	Labels pulumi.StringMapOutput `pulumi:"labels"`
	// Memory in megabytes (**aligned to 128MB**) for Yandex Cloud Serverless Container
	Memory pulumi.Float64Output `pulumi:"memory"`
	// Yandex Cloud Serverless Container name
	Name pulumi.StringOutput `pulumi:"name"`
	// Last revision ID of the Yandex Cloud Serverless Container
	RevisionId pulumi.StringOutput `pulumi:"revisionId"`
	// Secrets for Yandex Cloud Serverless Container
	Secrets ServerlessContainerSecretArrayOutput `pulumi:"secrets"`
	// Service account ID for Yandex Cloud Serverless Container
	ServiceAccountId pulumi.StringPtrOutput            `pulumi:"serviceAccountId"`
	Timeouts         ServerlessContainerTimeoutsOutput `pulumi:"timeouts"`
	// Invoke URL for the Yandex Cloud Serverless Container
	Url pulumi.StringOutput `pulumi:"url"`
}

// NewServerlessContainer registers a new resource with the given unique name, arguments, and options.
func NewServerlessContainer(ctx *pulumi.Context,
	name string, args *ServerlessContainerArgs, opts ...pulumi.ResourceOption) (*ServerlessContainer, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Image == nil {
		return nil, errors.New("invalid value for required argument 'Image'")
	}
	if args.Memory == nil {
		return nil, errors.New("invalid value for required argument 'Memory'")
	}
	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	var resource ServerlessContainer
	err := ctx.RegisterResource("yandex:index/serverlessContainer:ServerlessContainer", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetServerlessContainer gets an existing ServerlessContainer resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetServerlessContainer(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *ServerlessContainerState, opts ...pulumi.ResourceOption) (*ServerlessContainer, error) {
	var resource ServerlessContainer
	err := ctx.ReadResource("yandex:index/serverlessContainer:ServerlessContainer", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering ServerlessContainer resources.
type serverlessContainerState struct {
	// Concurrency of Yandex Cloud Serverless Container
	Concurrency *float64 `pulumi:"concurrency"`
	// Core fraction (**0...100**) of the Yandex Cloud Serverless Container
	CoreFraction *float64 `pulumi:"coreFraction"`
	Cores        *float64 `pulumi:"cores"`
	// Creation timestamp of the Yandex Cloud Serverless Container
	CreatedAt *string `pulumi:"createdAt"`
	// Description of the Yandex Cloud Serverless Container
	Description *string `pulumi:"description"`
	// Execution timeout in seconds (**duration format**) for Yandex Cloud Serverless Container
	ExecutionTimeout *string `pulumi:"executionTimeout"`
	// Folder ID for the Yandex Cloud Serverless Container
	FolderId *string `pulumi:"folderId"`
	// Revision deployment image for Yandex Cloud Serverless Container
	// * `image.0.url` (Required) - URL of image that will be deployed as Yandex Cloud Serverless Container
	// * `image.0.work_dir` - Working directory for Yandex Cloud Serverless Container
	// * `image.0.digest` - Digest of image that will be deployed as Yandex Cloud Serverless Container.
	//   If presented, should be equal to digest that will be resolved at server side by URL.
	//   Container will be updated on digest change even if `image.0.url` stays the same.
	//   If field not specified then its value will be computed.
	// * `image.0.command` - List of commands for Yandex Cloud Serverless Container
	// * `image.0.args` - List of arguments for Yandex Cloud Serverless Container
	// * `image.0.environment` -  A set of key/value environment variable pairs for Yandex Cloud Serverless Container
	Image *ServerlessContainerImage `pulumi:"image"`
	// A set of key/value label pairs to assign to the Yandex Cloud Serverless Container
	Labels map[string]string `pulumi:"labels"`
	// Memory in megabytes (**aligned to 128MB**) for Yandex Cloud Serverless Container
	Memory *float64 `pulumi:"memory"`
	// Yandex Cloud Serverless Container name
	Name *string `pulumi:"name"`
	// Last revision ID of the Yandex Cloud Serverless Container
	RevisionId *string `pulumi:"revisionId"`
	// Secrets for Yandex Cloud Serverless Container
	Secrets []ServerlessContainerSecret `pulumi:"secrets"`
	// Service account ID for Yandex Cloud Serverless Container
	ServiceAccountId *string                      `pulumi:"serviceAccountId"`
	Timeouts         *ServerlessContainerTimeouts `pulumi:"timeouts"`
	// Invoke URL for the Yandex Cloud Serverless Container
	Url *string `pulumi:"url"`
}

type ServerlessContainerState struct {
	// Concurrency of Yandex Cloud Serverless Container
	Concurrency pulumi.Float64PtrInput
	// Core fraction (**0...100**) of the Yandex Cloud Serverless Container
	CoreFraction pulumi.Float64PtrInput
	Cores        pulumi.Float64PtrInput
	// Creation timestamp of the Yandex Cloud Serverless Container
	CreatedAt pulumi.StringPtrInput
	// Description of the Yandex Cloud Serverless Container
	Description pulumi.StringPtrInput
	// Execution timeout in seconds (**duration format**) for Yandex Cloud Serverless Container
	ExecutionTimeout pulumi.StringPtrInput
	// Folder ID for the Yandex Cloud Serverless Container
	FolderId pulumi.StringPtrInput
	// Revision deployment image for Yandex Cloud Serverless Container
	// * `image.0.url` (Required) - URL of image that will be deployed as Yandex Cloud Serverless Container
	// * `image.0.work_dir` - Working directory for Yandex Cloud Serverless Container
	// * `image.0.digest` - Digest of image that will be deployed as Yandex Cloud Serverless Container.
	//   If presented, should be equal to digest that will be resolved at server side by URL.
	//   Container will be updated on digest change even if `image.0.url` stays the same.
	//   If field not specified then its value will be computed.
	// * `image.0.command` - List of commands for Yandex Cloud Serverless Container
	// * `image.0.args` - List of arguments for Yandex Cloud Serverless Container
	// * `image.0.environment` -  A set of key/value environment variable pairs for Yandex Cloud Serverless Container
	Image ServerlessContainerImagePtrInput
	// A set of key/value label pairs to assign to the Yandex Cloud Serverless Container
	Labels pulumi.StringMapInput
	// Memory in megabytes (**aligned to 128MB**) for Yandex Cloud Serverless Container
	Memory pulumi.Float64PtrInput
	// Yandex Cloud Serverless Container name
	Name pulumi.StringPtrInput
	// Last revision ID of the Yandex Cloud Serverless Container
	RevisionId pulumi.StringPtrInput
	// Secrets for Yandex Cloud Serverless Container
	Secrets ServerlessContainerSecretArrayInput
	// Service account ID for Yandex Cloud Serverless Container
	ServiceAccountId pulumi.StringPtrInput
	Timeouts         ServerlessContainerTimeoutsPtrInput
	// Invoke URL for the Yandex Cloud Serverless Container
	Url pulumi.StringPtrInput
}

func (ServerlessContainerState) ElementType() reflect.Type {
	return reflect.TypeOf((*serverlessContainerState)(nil)).Elem()
}

type serverlessContainerArgs struct {
	// Concurrency of Yandex Cloud Serverless Container
	Concurrency *float64 `pulumi:"concurrency"`
	// Core fraction (**0...100**) of the Yandex Cloud Serverless Container
	CoreFraction *float64 `pulumi:"coreFraction"`
	Cores        *float64 `pulumi:"cores"`
	// Description of the Yandex Cloud Serverless Container
	Description *string `pulumi:"description"`
	// Execution timeout in seconds (**duration format**) for Yandex Cloud Serverless Container
	ExecutionTimeout *string `pulumi:"executionTimeout"`
	// Folder ID for the Yandex Cloud Serverless Container
	FolderId *string `pulumi:"folderId"`
	// Revision deployment image for Yandex Cloud Serverless Container
	// * `image.0.url` (Required) - URL of image that will be deployed as Yandex Cloud Serverless Container
	// * `image.0.work_dir` - Working directory for Yandex Cloud Serverless Container
	// * `image.0.digest` - Digest of image that will be deployed as Yandex Cloud Serverless Container.
	//   If presented, should be equal to digest that will be resolved at server side by URL.
	//   Container will be updated on digest change even if `image.0.url` stays the same.
	//   If field not specified then its value will be computed.
	// * `image.0.command` - List of commands for Yandex Cloud Serverless Container
	// * `image.0.args` - List of arguments for Yandex Cloud Serverless Container
	// * `image.0.environment` -  A set of key/value environment variable pairs for Yandex Cloud Serverless Container
	Image ServerlessContainerImage `pulumi:"image"`
	// A set of key/value label pairs to assign to the Yandex Cloud Serverless Container
	Labels map[string]string `pulumi:"labels"`
	// Memory in megabytes (**aligned to 128MB**) for Yandex Cloud Serverless Container
	Memory float64 `pulumi:"memory"`
	// Yandex Cloud Serverless Container name
	Name *string `pulumi:"name"`
	// Secrets for Yandex Cloud Serverless Container
	Secrets []ServerlessContainerSecret `pulumi:"secrets"`
	// Service account ID for Yandex Cloud Serverless Container
	ServiceAccountId *string                     `pulumi:"serviceAccountId"`
	Timeouts         ServerlessContainerTimeouts `pulumi:"timeouts"`
}

// The set of arguments for constructing a ServerlessContainer resource.
type ServerlessContainerArgs struct {
	// Concurrency of Yandex Cloud Serverless Container
	Concurrency pulumi.Float64PtrInput
	// Core fraction (**0...100**) of the Yandex Cloud Serverless Container
	CoreFraction pulumi.Float64PtrInput
	Cores        pulumi.Float64PtrInput
	// Description of the Yandex Cloud Serverless Container
	Description pulumi.StringPtrInput
	// Execution timeout in seconds (**duration format**) for Yandex Cloud Serverless Container
	ExecutionTimeout pulumi.StringPtrInput
	// Folder ID for the Yandex Cloud Serverless Container
	FolderId pulumi.StringPtrInput
	// Revision deployment image for Yandex Cloud Serverless Container
	// * `image.0.url` (Required) - URL of image that will be deployed as Yandex Cloud Serverless Container
	// * `image.0.work_dir` - Working directory for Yandex Cloud Serverless Container
	// * `image.0.digest` - Digest of image that will be deployed as Yandex Cloud Serverless Container.
	//   If presented, should be equal to digest that will be resolved at server side by URL.
	//   Container will be updated on digest change even if `image.0.url` stays the same.
	//   If field not specified then its value will be computed.
	// * `image.0.command` - List of commands for Yandex Cloud Serverless Container
	// * `image.0.args` - List of arguments for Yandex Cloud Serverless Container
	// * `image.0.environment` -  A set of key/value environment variable pairs for Yandex Cloud Serverless Container
	Image ServerlessContainerImageInput
	// A set of key/value label pairs to assign to the Yandex Cloud Serverless Container
	Labels pulumi.StringMapInput
	// Memory in megabytes (**aligned to 128MB**) for Yandex Cloud Serverless Container
	Memory pulumi.Float64Input
	// Yandex Cloud Serverless Container name
	Name pulumi.StringPtrInput
	// Secrets for Yandex Cloud Serverless Container
	Secrets ServerlessContainerSecretArrayInput
	// Service account ID for Yandex Cloud Serverless Container
	ServiceAccountId pulumi.StringPtrInput
	Timeouts         ServerlessContainerTimeoutsInput
}

func (ServerlessContainerArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*serverlessContainerArgs)(nil)).Elem()
}

type ServerlessContainerInput interface {
	pulumi.Input

	ToServerlessContainerOutput() ServerlessContainerOutput
	ToServerlessContainerOutputWithContext(ctx context.Context) ServerlessContainerOutput
}

func (*ServerlessContainer) ElementType() reflect.Type {
	return reflect.TypeOf((**ServerlessContainer)(nil)).Elem()
}

func (i *ServerlessContainer) ToServerlessContainerOutput() ServerlessContainerOutput {
	return i.ToServerlessContainerOutputWithContext(context.Background())
}

func (i *ServerlessContainer) ToServerlessContainerOutputWithContext(ctx context.Context) ServerlessContainerOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ServerlessContainerOutput)
}

// ServerlessContainerArrayInput is an input type that accepts ServerlessContainerArray and ServerlessContainerArrayOutput values.
// You can construct a concrete instance of `ServerlessContainerArrayInput` via:
//
//          ServerlessContainerArray{ ServerlessContainerArgs{...} }
type ServerlessContainerArrayInput interface {
	pulumi.Input

	ToServerlessContainerArrayOutput() ServerlessContainerArrayOutput
	ToServerlessContainerArrayOutputWithContext(context.Context) ServerlessContainerArrayOutput
}

type ServerlessContainerArray []ServerlessContainerInput

func (ServerlessContainerArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ServerlessContainer)(nil)).Elem()
}

func (i ServerlessContainerArray) ToServerlessContainerArrayOutput() ServerlessContainerArrayOutput {
	return i.ToServerlessContainerArrayOutputWithContext(context.Background())
}

func (i ServerlessContainerArray) ToServerlessContainerArrayOutputWithContext(ctx context.Context) ServerlessContainerArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ServerlessContainerArrayOutput)
}

// ServerlessContainerMapInput is an input type that accepts ServerlessContainerMap and ServerlessContainerMapOutput values.
// You can construct a concrete instance of `ServerlessContainerMapInput` via:
//
//          ServerlessContainerMap{ "key": ServerlessContainerArgs{...} }
type ServerlessContainerMapInput interface {
	pulumi.Input

	ToServerlessContainerMapOutput() ServerlessContainerMapOutput
	ToServerlessContainerMapOutputWithContext(context.Context) ServerlessContainerMapOutput
}

type ServerlessContainerMap map[string]ServerlessContainerInput

func (ServerlessContainerMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ServerlessContainer)(nil)).Elem()
}

func (i ServerlessContainerMap) ToServerlessContainerMapOutput() ServerlessContainerMapOutput {
	return i.ToServerlessContainerMapOutputWithContext(context.Background())
}

func (i ServerlessContainerMap) ToServerlessContainerMapOutputWithContext(ctx context.Context) ServerlessContainerMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ServerlessContainerMapOutput)
}

type ServerlessContainerOutput struct{ *pulumi.OutputState }

func (ServerlessContainerOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**ServerlessContainer)(nil)).Elem()
}

func (o ServerlessContainerOutput) ToServerlessContainerOutput() ServerlessContainerOutput {
	return o
}

func (o ServerlessContainerOutput) ToServerlessContainerOutputWithContext(ctx context.Context) ServerlessContainerOutput {
	return o
}

// Concurrency of Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) Concurrency() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.Float64PtrOutput { return v.Concurrency }).(pulumi.Float64PtrOutput)
}

// Core fraction (**0...100**) of the Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) CoreFraction() pulumi.Float64Output {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.Float64Output { return v.CoreFraction }).(pulumi.Float64Output)
}

func (o ServerlessContainerOutput) Cores() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.Float64PtrOutput { return v.Cores }).(pulumi.Float64PtrOutput)
}

// Creation timestamp of the Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.StringOutput { return v.CreatedAt }).(pulumi.StringOutput)
}

// Description of the Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// Execution timeout in seconds (**duration format**) for Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) ExecutionTimeout() pulumi.StringOutput {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.StringOutput { return v.ExecutionTimeout }).(pulumi.StringOutput)
}

// Folder ID for the Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.StringOutput { return v.FolderId }).(pulumi.StringOutput)
}

// Revision deployment image for Yandex Cloud Serverless Container
// * `image.0.url` (Required) - URL of image that will be deployed as Yandex Cloud Serverless Container
// * `image.0.work_dir` - Working directory for Yandex Cloud Serverless Container
// * `image.0.digest` - Digest of image that will be deployed as Yandex Cloud Serverless Container.
//   If presented, should be equal to digest that will be resolved at server side by URL.
//   Container will be updated on digest change even if `image.0.url` stays the same.
//   If field not specified then its value will be computed.
// * `image.0.command` - List of commands for Yandex Cloud Serverless Container
// * `image.0.args` - List of arguments for Yandex Cloud Serverless Container
// * `image.0.environment` -  A set of key/value environment variable pairs for Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) Image() ServerlessContainerImageOutput {
	return o.ApplyT(func(v *ServerlessContainer) ServerlessContainerImageOutput { return v.Image }).(ServerlessContainerImageOutput)
}

// A set of key/value label pairs to assign to the Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.StringMapOutput { return v.Labels }).(pulumi.StringMapOutput)
}

// Memory in megabytes (**aligned to 128MB**) for Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) Memory() pulumi.Float64Output {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.Float64Output { return v.Memory }).(pulumi.Float64Output)
}

// Yandex Cloud Serverless Container name
func (o ServerlessContainerOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// Last revision ID of the Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) RevisionId() pulumi.StringOutput {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.StringOutput { return v.RevisionId }).(pulumi.StringOutput)
}

// Secrets for Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) Secrets() ServerlessContainerSecretArrayOutput {
	return o.ApplyT(func(v *ServerlessContainer) ServerlessContainerSecretArrayOutput { return v.Secrets }).(ServerlessContainerSecretArrayOutput)
}

// Service account ID for Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) ServiceAccountId() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.StringPtrOutput { return v.ServiceAccountId }).(pulumi.StringPtrOutput)
}

func (o ServerlessContainerOutput) Timeouts() ServerlessContainerTimeoutsOutput {
	return o.ApplyT(func(v *ServerlessContainer) ServerlessContainerTimeoutsOutput { return v.Timeouts }).(ServerlessContainerTimeoutsOutput)
}

// Invoke URL for the Yandex Cloud Serverless Container
func (o ServerlessContainerOutput) Url() pulumi.StringOutput {
	return o.ApplyT(func(v *ServerlessContainer) pulumi.StringOutput { return v.Url }).(pulumi.StringOutput)
}

type ServerlessContainerArrayOutput struct{ *pulumi.OutputState }

func (ServerlessContainerArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ServerlessContainer)(nil)).Elem()
}

func (o ServerlessContainerArrayOutput) ToServerlessContainerArrayOutput() ServerlessContainerArrayOutput {
	return o
}

func (o ServerlessContainerArrayOutput) ToServerlessContainerArrayOutputWithContext(ctx context.Context) ServerlessContainerArrayOutput {
	return o
}

func (o ServerlessContainerArrayOutput) Index(i pulumi.IntInput) ServerlessContainerOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *ServerlessContainer {
		return vs[0].([]*ServerlessContainer)[vs[1].(int)]
	}).(ServerlessContainerOutput)
}

type ServerlessContainerMapOutput struct{ *pulumi.OutputState }

func (ServerlessContainerMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ServerlessContainer)(nil)).Elem()
}

func (o ServerlessContainerMapOutput) ToServerlessContainerMapOutput() ServerlessContainerMapOutput {
	return o
}

func (o ServerlessContainerMapOutput) ToServerlessContainerMapOutputWithContext(ctx context.Context) ServerlessContainerMapOutput {
	return o
}

func (o ServerlessContainerMapOutput) MapIndex(k pulumi.StringInput) ServerlessContainerOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *ServerlessContainer {
		return vs[0].(map[string]*ServerlessContainer)[vs[1].(string)]
	}).(ServerlessContainerOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ServerlessContainerInput)(nil)).Elem(), &ServerlessContainer{})
	pulumi.RegisterInputType(reflect.TypeOf((*ServerlessContainerArrayInput)(nil)).Elem(), ServerlessContainerArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*ServerlessContainerMapInput)(nil)).Elem(), ServerlessContainerMap{})
	pulumi.RegisterOutputType(ServerlessContainerOutput{})
	pulumi.RegisterOutputType(ServerlessContainerArrayOutput{})
	pulumi.RegisterOutputType(ServerlessContainerMapOutput{})
}
