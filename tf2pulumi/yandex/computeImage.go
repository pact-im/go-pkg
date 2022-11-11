// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Creates a virtual machine image resource for the Yandex Compute Cloud service from an existing
// tarball. For more information, see [the official documentation](https://cloud.yandex.com/docs/compute/concepts/image).
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
// 		_, err := yandex.NewComputeImage(ctx, "foo-image", &yandex.ComputeImageArgs{
// 			SourceUrl: pulumi.String("https://storage.yandexcloud.net/lucky-images/kube-it.img"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		_, err = yandex.NewComputeInstance(ctx, "vm", &yandex.ComputeInstanceArgs{
// 			BootDisk: &ComputeInstanceBootDiskArgs{
// 				InitializeParams: &ComputeInstanceBootDiskInitializeParamsArgs{
// 					ImageId: foo_image.ID(),
// 				},
// 			},
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
// A VM image can be imported using the `id` of the resource, e.g.
//
// ```sh
//  $ pulumi import yandex:index/computeImage:ComputeImage web-image image_id
// ```
type ComputeImage struct {
	pulumi.CustomResourceState

	// Creation timestamp of the image.
	CreatedAt pulumi.StringOutput `pulumi:"createdAt"`
	// An optional description of the image. Provide this property when
	// you create a resource.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// The name of the image family to which this image belongs.
	Family pulumi.StringPtrOutput `pulumi:"family"`
	// The ID of the folder that the resource belongs to. If it
	// is not provided, the default provider folder is used.
	FolderId pulumi.StringOutput `pulumi:"folderId"`
	// A set of key/value label pairs to assign to the image.
	Labels pulumi.StringMapOutput `pulumi:"labels"`
	// Minimum size in GB of the disk that will be created from this image.
	MinDiskSize pulumi.Float64Output `pulumi:"minDiskSize"`
	// Name of the disk.
	Name pulumi.StringOutput `pulumi:"name"`
	// Operating system type that is contained in the image. Possible values: "LINUX", "WINDOWS".
	OsType pulumi.StringOutput `pulumi:"osType"`
	// Optimize the image to create a disk.
	Pooled pulumi.BoolOutput `pulumi:"pooled"`
	// License IDs that indicate which licenses are
	// attached to this image.
	ProductIds pulumi.StringArrayOutput `pulumi:"productIds"`
	// The size of the image, specified in GB.
	Size pulumi.Float64Output `pulumi:"size"`
	// The ID of a disk to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceDisk pulumi.StringOutput `pulumi:"sourceDisk"`
	// The name of the family to use as the source of the new image.
	// The ID of the latest image is taken from the "standard-images" folder. Changing the family forces
	// a new resource to be created.
	SourceFamily pulumi.StringOutput `pulumi:"sourceFamily"`
	// The ID of an existing image to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceImage pulumi.StringOutput `pulumi:"sourceImage"`
	// The ID of a snapshot to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceSnapshot pulumi.StringOutput `pulumi:"sourceSnapshot"`
	// The URL to use as the source of the
	// image. Changing this URL forces a new resource to be created.
	SourceUrl pulumi.StringOutput `pulumi:"sourceUrl"`
	// The status of the image.
	Status   pulumi.StringOutput        `pulumi:"status"`
	Timeouts ComputeImageTimeoutsOutput `pulumi:"timeouts"`
}

// NewComputeImage registers a new resource with the given unique name, arguments, and options.
func NewComputeImage(ctx *pulumi.Context,
	name string, args *ComputeImageArgs, opts ...pulumi.ResourceOption) (*ComputeImage, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	var resource ComputeImage
	err := ctx.RegisterResource("yandex:index/computeImage:ComputeImage", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetComputeImage gets an existing ComputeImage resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetComputeImage(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *ComputeImageState, opts ...pulumi.ResourceOption) (*ComputeImage, error) {
	var resource ComputeImage
	err := ctx.ReadResource("yandex:index/computeImage:ComputeImage", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering ComputeImage resources.
type computeImageState struct {
	// Creation timestamp of the image.
	CreatedAt *string `pulumi:"createdAt"`
	// An optional description of the image. Provide this property when
	// you create a resource.
	Description *string `pulumi:"description"`
	// The name of the image family to which this image belongs.
	Family *string `pulumi:"family"`
	// The ID of the folder that the resource belongs to. If it
	// is not provided, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// A set of key/value label pairs to assign to the image.
	Labels map[string]string `pulumi:"labels"`
	// Minimum size in GB of the disk that will be created from this image.
	MinDiskSize *float64 `pulumi:"minDiskSize"`
	// Name of the disk.
	Name *string `pulumi:"name"`
	// Operating system type that is contained in the image. Possible values: "LINUX", "WINDOWS".
	OsType *string `pulumi:"osType"`
	// Optimize the image to create a disk.
	Pooled *bool `pulumi:"pooled"`
	// License IDs that indicate which licenses are
	// attached to this image.
	ProductIds []string `pulumi:"productIds"`
	// The size of the image, specified in GB.
	Size *float64 `pulumi:"size"`
	// The ID of a disk to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceDisk *string `pulumi:"sourceDisk"`
	// The name of the family to use as the source of the new image.
	// The ID of the latest image is taken from the "standard-images" folder. Changing the family forces
	// a new resource to be created.
	SourceFamily *string `pulumi:"sourceFamily"`
	// The ID of an existing image to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceImage *string `pulumi:"sourceImage"`
	// The ID of a snapshot to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceSnapshot *string `pulumi:"sourceSnapshot"`
	// The URL to use as the source of the
	// image. Changing this URL forces a new resource to be created.
	SourceUrl *string `pulumi:"sourceUrl"`
	// The status of the image.
	Status   *string               `pulumi:"status"`
	Timeouts *ComputeImageTimeouts `pulumi:"timeouts"`
}

type ComputeImageState struct {
	// Creation timestamp of the image.
	CreatedAt pulumi.StringPtrInput
	// An optional description of the image. Provide this property when
	// you create a resource.
	Description pulumi.StringPtrInput
	// The name of the image family to which this image belongs.
	Family pulumi.StringPtrInput
	// The ID of the folder that the resource belongs to. If it
	// is not provided, the default provider folder is used.
	FolderId pulumi.StringPtrInput
	// A set of key/value label pairs to assign to the image.
	Labels pulumi.StringMapInput
	// Minimum size in GB of the disk that will be created from this image.
	MinDiskSize pulumi.Float64PtrInput
	// Name of the disk.
	Name pulumi.StringPtrInput
	// Operating system type that is contained in the image. Possible values: "LINUX", "WINDOWS".
	OsType pulumi.StringPtrInput
	// Optimize the image to create a disk.
	Pooled pulumi.BoolPtrInput
	// License IDs that indicate which licenses are
	// attached to this image.
	ProductIds pulumi.StringArrayInput
	// The size of the image, specified in GB.
	Size pulumi.Float64PtrInput
	// The ID of a disk to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceDisk pulumi.StringPtrInput
	// The name of the family to use as the source of the new image.
	// The ID of the latest image is taken from the "standard-images" folder. Changing the family forces
	// a new resource to be created.
	SourceFamily pulumi.StringPtrInput
	// The ID of an existing image to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceImage pulumi.StringPtrInput
	// The ID of a snapshot to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceSnapshot pulumi.StringPtrInput
	// The URL to use as the source of the
	// image. Changing this URL forces a new resource to be created.
	SourceUrl pulumi.StringPtrInput
	// The status of the image.
	Status   pulumi.StringPtrInput
	Timeouts ComputeImageTimeoutsPtrInput
}

func (ComputeImageState) ElementType() reflect.Type {
	return reflect.TypeOf((*computeImageState)(nil)).Elem()
}

type computeImageArgs struct {
	// An optional description of the image. Provide this property when
	// you create a resource.
	Description *string `pulumi:"description"`
	// The name of the image family to which this image belongs.
	Family *string `pulumi:"family"`
	// The ID of the folder that the resource belongs to. If it
	// is not provided, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// A set of key/value label pairs to assign to the image.
	Labels map[string]string `pulumi:"labels"`
	// Minimum size in GB of the disk that will be created from this image.
	MinDiskSize *float64 `pulumi:"minDiskSize"`
	// Name of the disk.
	Name *string `pulumi:"name"`
	// Operating system type that is contained in the image. Possible values: "LINUX", "WINDOWS".
	OsType *string `pulumi:"osType"`
	// Optimize the image to create a disk.
	Pooled *bool `pulumi:"pooled"`
	// License IDs that indicate which licenses are
	// attached to this image.
	ProductIds []string `pulumi:"productIds"`
	// The ID of a disk to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceDisk *string `pulumi:"sourceDisk"`
	// The name of the family to use as the source of the new image.
	// The ID of the latest image is taken from the "standard-images" folder. Changing the family forces
	// a new resource to be created.
	SourceFamily *string `pulumi:"sourceFamily"`
	// The ID of an existing image to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceImage *string `pulumi:"sourceImage"`
	// The ID of a snapshot to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceSnapshot *string `pulumi:"sourceSnapshot"`
	// The URL to use as the source of the
	// image. Changing this URL forces a new resource to be created.
	SourceUrl *string              `pulumi:"sourceUrl"`
	Timeouts  ComputeImageTimeouts `pulumi:"timeouts"`
}

// The set of arguments for constructing a ComputeImage resource.
type ComputeImageArgs struct {
	// An optional description of the image. Provide this property when
	// you create a resource.
	Description pulumi.StringPtrInput
	// The name of the image family to which this image belongs.
	Family pulumi.StringPtrInput
	// The ID of the folder that the resource belongs to. If it
	// is not provided, the default provider folder is used.
	FolderId pulumi.StringPtrInput
	// A set of key/value label pairs to assign to the image.
	Labels pulumi.StringMapInput
	// Minimum size in GB of the disk that will be created from this image.
	MinDiskSize pulumi.Float64PtrInput
	// Name of the disk.
	Name pulumi.StringPtrInput
	// Operating system type that is contained in the image. Possible values: "LINUX", "WINDOWS".
	OsType pulumi.StringPtrInput
	// Optimize the image to create a disk.
	Pooled pulumi.BoolPtrInput
	// License IDs that indicate which licenses are
	// attached to this image.
	ProductIds pulumi.StringArrayInput
	// The ID of a disk to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceDisk pulumi.StringPtrInput
	// The name of the family to use as the source of the new image.
	// The ID of the latest image is taken from the "standard-images" folder. Changing the family forces
	// a new resource to be created.
	SourceFamily pulumi.StringPtrInput
	// The ID of an existing image to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceImage pulumi.StringPtrInput
	// The ID of a snapshot to use as the source of the
	// image. Changing this ID forces a new resource to be created.
	SourceSnapshot pulumi.StringPtrInput
	// The URL to use as the source of the
	// image. Changing this URL forces a new resource to be created.
	SourceUrl pulumi.StringPtrInput
	Timeouts  ComputeImageTimeoutsInput
}

func (ComputeImageArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*computeImageArgs)(nil)).Elem()
}

type ComputeImageInput interface {
	pulumi.Input

	ToComputeImageOutput() ComputeImageOutput
	ToComputeImageOutputWithContext(ctx context.Context) ComputeImageOutput
}

func (*ComputeImage) ElementType() reflect.Type {
	return reflect.TypeOf((**ComputeImage)(nil)).Elem()
}

func (i *ComputeImage) ToComputeImageOutput() ComputeImageOutput {
	return i.ToComputeImageOutputWithContext(context.Background())
}

func (i *ComputeImage) ToComputeImageOutputWithContext(ctx context.Context) ComputeImageOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ComputeImageOutput)
}

// ComputeImageArrayInput is an input type that accepts ComputeImageArray and ComputeImageArrayOutput values.
// You can construct a concrete instance of `ComputeImageArrayInput` via:
//
//          ComputeImageArray{ ComputeImageArgs{...} }
type ComputeImageArrayInput interface {
	pulumi.Input

	ToComputeImageArrayOutput() ComputeImageArrayOutput
	ToComputeImageArrayOutputWithContext(context.Context) ComputeImageArrayOutput
}

type ComputeImageArray []ComputeImageInput

func (ComputeImageArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ComputeImage)(nil)).Elem()
}

func (i ComputeImageArray) ToComputeImageArrayOutput() ComputeImageArrayOutput {
	return i.ToComputeImageArrayOutputWithContext(context.Background())
}

func (i ComputeImageArray) ToComputeImageArrayOutputWithContext(ctx context.Context) ComputeImageArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ComputeImageArrayOutput)
}

// ComputeImageMapInput is an input type that accepts ComputeImageMap and ComputeImageMapOutput values.
// You can construct a concrete instance of `ComputeImageMapInput` via:
//
//          ComputeImageMap{ "key": ComputeImageArgs{...} }
type ComputeImageMapInput interface {
	pulumi.Input

	ToComputeImageMapOutput() ComputeImageMapOutput
	ToComputeImageMapOutputWithContext(context.Context) ComputeImageMapOutput
}

type ComputeImageMap map[string]ComputeImageInput

func (ComputeImageMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ComputeImage)(nil)).Elem()
}

func (i ComputeImageMap) ToComputeImageMapOutput() ComputeImageMapOutput {
	return i.ToComputeImageMapOutputWithContext(context.Background())
}

func (i ComputeImageMap) ToComputeImageMapOutputWithContext(ctx context.Context) ComputeImageMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(ComputeImageMapOutput)
}

type ComputeImageOutput struct{ *pulumi.OutputState }

func (ComputeImageOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**ComputeImage)(nil)).Elem()
}

func (o ComputeImageOutput) ToComputeImageOutput() ComputeImageOutput {
	return o
}

func (o ComputeImageOutput) ToComputeImageOutputWithContext(ctx context.Context) ComputeImageOutput {
	return o
}

// Creation timestamp of the image.
func (o ComputeImageOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringOutput { return v.CreatedAt }).(pulumi.StringOutput)
}

// An optional description of the image. Provide this property when
// you create a resource.
func (o ComputeImageOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// The name of the image family to which this image belongs.
func (o ComputeImageOutput) Family() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringPtrOutput { return v.Family }).(pulumi.StringPtrOutput)
}

// The ID of the folder that the resource belongs to. If it
// is not provided, the default provider folder is used.
func (o ComputeImageOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringOutput { return v.FolderId }).(pulumi.StringOutput)
}

// A set of key/value label pairs to assign to the image.
func (o ComputeImageOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringMapOutput { return v.Labels }).(pulumi.StringMapOutput)
}

// Minimum size in GB of the disk that will be created from this image.
func (o ComputeImageOutput) MinDiskSize() pulumi.Float64Output {
	return o.ApplyT(func(v *ComputeImage) pulumi.Float64Output { return v.MinDiskSize }).(pulumi.Float64Output)
}

// Name of the disk.
func (o ComputeImageOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// Operating system type that is contained in the image. Possible values: "LINUX", "WINDOWS".
func (o ComputeImageOutput) OsType() pulumi.StringOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringOutput { return v.OsType }).(pulumi.StringOutput)
}

// Optimize the image to create a disk.
func (o ComputeImageOutput) Pooled() pulumi.BoolOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.BoolOutput { return v.Pooled }).(pulumi.BoolOutput)
}

// License IDs that indicate which licenses are
// attached to this image.
func (o ComputeImageOutput) ProductIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringArrayOutput { return v.ProductIds }).(pulumi.StringArrayOutput)
}

// The size of the image, specified in GB.
func (o ComputeImageOutput) Size() pulumi.Float64Output {
	return o.ApplyT(func(v *ComputeImage) pulumi.Float64Output { return v.Size }).(pulumi.Float64Output)
}

// The ID of a disk to use as the source of the
// image. Changing this ID forces a new resource to be created.
func (o ComputeImageOutput) SourceDisk() pulumi.StringOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringOutput { return v.SourceDisk }).(pulumi.StringOutput)
}

// The name of the family to use as the source of the new image.
// The ID of the latest image is taken from the "standard-images" folder. Changing the family forces
// a new resource to be created.
func (o ComputeImageOutput) SourceFamily() pulumi.StringOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringOutput { return v.SourceFamily }).(pulumi.StringOutput)
}

// The ID of an existing image to use as the source of the
// image. Changing this ID forces a new resource to be created.
func (o ComputeImageOutput) SourceImage() pulumi.StringOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringOutput { return v.SourceImage }).(pulumi.StringOutput)
}

// The ID of a snapshot to use as the source of the
// image. Changing this ID forces a new resource to be created.
func (o ComputeImageOutput) SourceSnapshot() pulumi.StringOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringOutput { return v.SourceSnapshot }).(pulumi.StringOutput)
}

// The URL to use as the source of the
// image. Changing this URL forces a new resource to be created.
func (o ComputeImageOutput) SourceUrl() pulumi.StringOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringOutput { return v.SourceUrl }).(pulumi.StringOutput)
}

// The status of the image.
func (o ComputeImageOutput) Status() pulumi.StringOutput {
	return o.ApplyT(func(v *ComputeImage) pulumi.StringOutput { return v.Status }).(pulumi.StringOutput)
}

func (o ComputeImageOutput) Timeouts() ComputeImageTimeoutsOutput {
	return o.ApplyT(func(v *ComputeImage) ComputeImageTimeoutsOutput { return v.Timeouts }).(ComputeImageTimeoutsOutput)
}

type ComputeImageArrayOutput struct{ *pulumi.OutputState }

func (ComputeImageArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*ComputeImage)(nil)).Elem()
}

func (o ComputeImageArrayOutput) ToComputeImageArrayOutput() ComputeImageArrayOutput {
	return o
}

func (o ComputeImageArrayOutput) ToComputeImageArrayOutputWithContext(ctx context.Context) ComputeImageArrayOutput {
	return o
}

func (o ComputeImageArrayOutput) Index(i pulumi.IntInput) ComputeImageOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *ComputeImage {
		return vs[0].([]*ComputeImage)[vs[1].(int)]
	}).(ComputeImageOutput)
}

type ComputeImageMapOutput struct{ *pulumi.OutputState }

func (ComputeImageMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*ComputeImage)(nil)).Elem()
}

func (o ComputeImageMapOutput) ToComputeImageMapOutput() ComputeImageMapOutput {
	return o
}

func (o ComputeImageMapOutput) ToComputeImageMapOutputWithContext(ctx context.Context) ComputeImageMapOutput {
	return o
}

func (o ComputeImageMapOutput) MapIndex(k pulumi.StringInput) ComputeImageOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *ComputeImage {
		return vs[0].(map[string]*ComputeImage)[vs[1].(string)]
	}).(ComputeImageOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*ComputeImageInput)(nil)).Elem(), &ComputeImage{})
	pulumi.RegisterInputType(reflect.TypeOf((*ComputeImageArrayInput)(nil)).Elem(), ComputeImageArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*ComputeImageMapInput)(nil)).Elem(), ComputeImageMap{})
	pulumi.RegisterOutputType(ComputeImageOutput{})
	pulumi.RegisterOutputType(ComputeImageArrayOutput{})
	pulumi.RegisterOutputType(ComputeImageMapOutput{})
}
