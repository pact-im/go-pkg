// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Manages a Elasticsearch cluster within the Yandex.Cloud. For more information, see
// [the official documentation](https://cloud.yandex.com/docs/managed-elasticsearch/concepts).
//
// ## Import
//
// A cluster can be imported using the `id` of the resource, e.g.
//
// ```sh
//  $ pulumi import yandex:index/mdbElasticSearchCluster:MdbElasticSearchCluster foo cluster_id
// ```
type MdbElasticSearchCluster struct {
	pulumi.CustomResourceState

	// Configuration of the Elasticsearch cluster. The structure is documented below.
	Config MdbElasticSearchClusterConfigOutput `pulumi:"config"`
	// Creation timestamp of the key.
	CreatedAt pulumi.StringOutput `pulumi:"createdAt"`
	// Inhibits deletion of the cluster.  Can be either `true` or `false`.
	DeletionProtection pulumi.BoolOutput `pulumi:"deletionProtection"`
	// Description of the Elasticsearch cluster.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// Deployment environment of the Elasticsearch cluster. Can be either `PRESTABLE` or `PRODUCTION`.
	Environment pulumi.StringOutput `pulumi:"environment"`
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId pulumi.StringOutput `pulumi:"folderId"`
	// Aggregated health of the cluster. Can be either `ALIVE`, `DEGRADED`, `DEAD` or `HEALTH_UNKNOWN`.
	// For more information see `health` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-elasticsearch/api-ref/Cluster/).
	Health pulumi.StringOutput `pulumi:"health"`
	// A host of the Elasticsearch cluster. The structure is documented below.
	Hosts MdbElasticSearchClusterHostArrayOutput `pulumi:"hosts"`
	// A set of key/value label pairs to assign to the Elasticsearch cluster.
	Labels            pulumi.StringMapOutput                            `pulumi:"labels"`
	MaintenanceWindow MdbElasticSearchClusterMaintenanceWindowPtrOutput `pulumi:"maintenanceWindow"`
	// User defined host name.
	Name pulumi.StringOutput `pulumi:"name"`
	// ID of the network, to which the Elasticsearch cluster belongs.
	NetworkId pulumi.StringOutput `pulumi:"networkId"`
	// A set of ids of security groups assigned to hosts of the cluster.
	SecurityGroupIds pulumi.StringArrayOutput `pulumi:"securityGroupIds"`
	// ID of the service account authorized for this cluster.
	ServiceAccountId pulumi.StringPtrOutput `pulumi:"serviceAccountId"`
	// Status of the cluster. Can be either `CREATING`, `STARTING`, `RUNNING`, `UPDATING`, `STOPPING`, `STOPPED`, `ERROR` or `STATUS_UNKNOWN`.
	// For more information see `status` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-elasticsearch/api-ref/Cluster/).
	Status   pulumi.StringOutput                   `pulumi:"status"`
	Timeouts MdbElasticSearchClusterTimeoutsOutput `pulumi:"timeouts"`
}

// NewMdbElasticSearchCluster registers a new resource with the given unique name, arguments, and options.
func NewMdbElasticSearchCluster(ctx *pulumi.Context,
	name string, args *MdbElasticSearchClusterArgs, opts ...pulumi.ResourceOption) (*MdbElasticSearchCluster, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Config == nil {
		return nil, errors.New("invalid value for required argument 'Config'")
	}
	if args.Environment == nil {
		return nil, errors.New("invalid value for required argument 'Environment'")
	}
	if args.NetworkId == nil {
		return nil, errors.New("invalid value for required argument 'NetworkId'")
	}
	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	var resource MdbElasticSearchCluster
	err := ctx.RegisterResource("yandex:index/mdbElasticSearchCluster:MdbElasticSearchCluster", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetMdbElasticSearchCluster gets an existing MdbElasticSearchCluster resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetMdbElasticSearchCluster(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *MdbElasticSearchClusterState, opts ...pulumi.ResourceOption) (*MdbElasticSearchCluster, error) {
	var resource MdbElasticSearchCluster
	err := ctx.ReadResource("yandex:index/mdbElasticSearchCluster:MdbElasticSearchCluster", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering MdbElasticSearchCluster resources.
type mdbElasticSearchClusterState struct {
	// Configuration of the Elasticsearch cluster. The structure is documented below.
	Config *MdbElasticSearchClusterConfig `pulumi:"config"`
	// Creation timestamp of the key.
	CreatedAt *string `pulumi:"createdAt"`
	// Inhibits deletion of the cluster.  Can be either `true` or `false`.
	DeletionProtection *bool `pulumi:"deletionProtection"`
	// Description of the Elasticsearch cluster.
	Description *string `pulumi:"description"`
	// Deployment environment of the Elasticsearch cluster. Can be either `PRESTABLE` or `PRODUCTION`.
	Environment *string `pulumi:"environment"`
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// Aggregated health of the cluster. Can be either `ALIVE`, `DEGRADED`, `DEAD` or `HEALTH_UNKNOWN`.
	// For more information see `health` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-elasticsearch/api-ref/Cluster/).
	Health *string `pulumi:"health"`
	// A host of the Elasticsearch cluster. The structure is documented below.
	Hosts []MdbElasticSearchClusterHost `pulumi:"hosts"`
	// A set of key/value label pairs to assign to the Elasticsearch cluster.
	Labels            map[string]string                         `pulumi:"labels"`
	MaintenanceWindow *MdbElasticSearchClusterMaintenanceWindow `pulumi:"maintenanceWindow"`
	// User defined host name.
	Name *string `pulumi:"name"`
	// ID of the network, to which the Elasticsearch cluster belongs.
	NetworkId *string `pulumi:"networkId"`
	// A set of ids of security groups assigned to hosts of the cluster.
	SecurityGroupIds []string `pulumi:"securityGroupIds"`
	// ID of the service account authorized for this cluster.
	ServiceAccountId *string `pulumi:"serviceAccountId"`
	// Status of the cluster. Can be either `CREATING`, `STARTING`, `RUNNING`, `UPDATING`, `STOPPING`, `STOPPED`, `ERROR` or `STATUS_UNKNOWN`.
	// For more information see `status` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-elasticsearch/api-ref/Cluster/).
	Status   *string                          `pulumi:"status"`
	Timeouts *MdbElasticSearchClusterTimeouts `pulumi:"timeouts"`
}

type MdbElasticSearchClusterState struct {
	// Configuration of the Elasticsearch cluster. The structure is documented below.
	Config MdbElasticSearchClusterConfigPtrInput
	// Creation timestamp of the key.
	CreatedAt pulumi.StringPtrInput
	// Inhibits deletion of the cluster.  Can be either `true` or `false`.
	DeletionProtection pulumi.BoolPtrInput
	// Description of the Elasticsearch cluster.
	Description pulumi.StringPtrInput
	// Deployment environment of the Elasticsearch cluster. Can be either `PRESTABLE` or `PRODUCTION`.
	Environment pulumi.StringPtrInput
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId pulumi.StringPtrInput
	// Aggregated health of the cluster. Can be either `ALIVE`, `DEGRADED`, `DEAD` or `HEALTH_UNKNOWN`.
	// For more information see `health` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-elasticsearch/api-ref/Cluster/).
	Health pulumi.StringPtrInput
	// A host of the Elasticsearch cluster. The structure is documented below.
	Hosts MdbElasticSearchClusterHostArrayInput
	// A set of key/value label pairs to assign to the Elasticsearch cluster.
	Labels            pulumi.StringMapInput
	MaintenanceWindow MdbElasticSearchClusterMaintenanceWindowPtrInput
	// User defined host name.
	Name pulumi.StringPtrInput
	// ID of the network, to which the Elasticsearch cluster belongs.
	NetworkId pulumi.StringPtrInput
	// A set of ids of security groups assigned to hosts of the cluster.
	SecurityGroupIds pulumi.StringArrayInput
	// ID of the service account authorized for this cluster.
	ServiceAccountId pulumi.StringPtrInput
	// Status of the cluster. Can be either `CREATING`, `STARTING`, `RUNNING`, `UPDATING`, `STOPPING`, `STOPPED`, `ERROR` or `STATUS_UNKNOWN`.
	// For more information see `status` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-elasticsearch/api-ref/Cluster/).
	Status   pulumi.StringPtrInput
	Timeouts MdbElasticSearchClusterTimeoutsPtrInput
}

func (MdbElasticSearchClusterState) ElementType() reflect.Type {
	return reflect.TypeOf((*mdbElasticSearchClusterState)(nil)).Elem()
}

type mdbElasticSearchClusterArgs struct {
	// Configuration of the Elasticsearch cluster. The structure is documented below.
	Config MdbElasticSearchClusterConfig `pulumi:"config"`
	// Inhibits deletion of the cluster.  Can be either `true` or `false`.
	DeletionProtection *bool `pulumi:"deletionProtection"`
	// Description of the Elasticsearch cluster.
	Description *string `pulumi:"description"`
	// Deployment environment of the Elasticsearch cluster. Can be either `PRESTABLE` or `PRODUCTION`.
	Environment string `pulumi:"environment"`
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// A host of the Elasticsearch cluster. The structure is documented below.
	Hosts []MdbElasticSearchClusterHost `pulumi:"hosts"`
	// A set of key/value label pairs to assign to the Elasticsearch cluster.
	Labels            map[string]string                         `pulumi:"labels"`
	MaintenanceWindow *MdbElasticSearchClusterMaintenanceWindow `pulumi:"maintenanceWindow"`
	// User defined host name.
	Name *string `pulumi:"name"`
	// ID of the network, to which the Elasticsearch cluster belongs.
	NetworkId string `pulumi:"networkId"`
	// A set of ids of security groups assigned to hosts of the cluster.
	SecurityGroupIds []string `pulumi:"securityGroupIds"`
	// ID of the service account authorized for this cluster.
	ServiceAccountId *string                         `pulumi:"serviceAccountId"`
	Timeouts         MdbElasticSearchClusterTimeouts `pulumi:"timeouts"`
}

// The set of arguments for constructing a MdbElasticSearchCluster resource.
type MdbElasticSearchClusterArgs struct {
	// Configuration of the Elasticsearch cluster. The structure is documented below.
	Config MdbElasticSearchClusterConfigInput
	// Inhibits deletion of the cluster.  Can be either `true` or `false`.
	DeletionProtection pulumi.BoolPtrInput
	// Description of the Elasticsearch cluster.
	Description pulumi.StringPtrInput
	// Deployment environment of the Elasticsearch cluster. Can be either `PRESTABLE` or `PRODUCTION`.
	Environment pulumi.StringInput
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId pulumi.StringPtrInput
	// A host of the Elasticsearch cluster. The structure is documented below.
	Hosts MdbElasticSearchClusterHostArrayInput
	// A set of key/value label pairs to assign to the Elasticsearch cluster.
	Labels            pulumi.StringMapInput
	MaintenanceWindow MdbElasticSearchClusterMaintenanceWindowPtrInput
	// User defined host name.
	Name pulumi.StringPtrInput
	// ID of the network, to which the Elasticsearch cluster belongs.
	NetworkId pulumi.StringInput
	// A set of ids of security groups assigned to hosts of the cluster.
	SecurityGroupIds pulumi.StringArrayInput
	// ID of the service account authorized for this cluster.
	ServiceAccountId pulumi.StringPtrInput
	Timeouts         MdbElasticSearchClusterTimeoutsInput
}

func (MdbElasticSearchClusterArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*mdbElasticSearchClusterArgs)(nil)).Elem()
}

type MdbElasticSearchClusterInput interface {
	pulumi.Input

	ToMdbElasticSearchClusterOutput() MdbElasticSearchClusterOutput
	ToMdbElasticSearchClusterOutputWithContext(ctx context.Context) MdbElasticSearchClusterOutput
}

func (*MdbElasticSearchCluster) ElementType() reflect.Type {
	return reflect.TypeOf((**MdbElasticSearchCluster)(nil)).Elem()
}

func (i *MdbElasticSearchCluster) ToMdbElasticSearchClusterOutput() MdbElasticSearchClusterOutput {
	return i.ToMdbElasticSearchClusterOutputWithContext(context.Background())
}

func (i *MdbElasticSearchCluster) ToMdbElasticSearchClusterOutputWithContext(ctx context.Context) MdbElasticSearchClusterOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MdbElasticSearchClusterOutput)
}

// MdbElasticSearchClusterArrayInput is an input type that accepts MdbElasticSearchClusterArray and MdbElasticSearchClusterArrayOutput values.
// You can construct a concrete instance of `MdbElasticSearchClusterArrayInput` via:
//
//          MdbElasticSearchClusterArray{ MdbElasticSearchClusterArgs{...} }
type MdbElasticSearchClusterArrayInput interface {
	pulumi.Input

	ToMdbElasticSearchClusterArrayOutput() MdbElasticSearchClusterArrayOutput
	ToMdbElasticSearchClusterArrayOutputWithContext(context.Context) MdbElasticSearchClusterArrayOutput
}

type MdbElasticSearchClusterArray []MdbElasticSearchClusterInput

func (MdbElasticSearchClusterArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*MdbElasticSearchCluster)(nil)).Elem()
}

func (i MdbElasticSearchClusterArray) ToMdbElasticSearchClusterArrayOutput() MdbElasticSearchClusterArrayOutput {
	return i.ToMdbElasticSearchClusterArrayOutputWithContext(context.Background())
}

func (i MdbElasticSearchClusterArray) ToMdbElasticSearchClusterArrayOutputWithContext(ctx context.Context) MdbElasticSearchClusterArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MdbElasticSearchClusterArrayOutput)
}

// MdbElasticSearchClusterMapInput is an input type that accepts MdbElasticSearchClusterMap and MdbElasticSearchClusterMapOutput values.
// You can construct a concrete instance of `MdbElasticSearchClusterMapInput` via:
//
//          MdbElasticSearchClusterMap{ "key": MdbElasticSearchClusterArgs{...} }
type MdbElasticSearchClusterMapInput interface {
	pulumi.Input

	ToMdbElasticSearchClusterMapOutput() MdbElasticSearchClusterMapOutput
	ToMdbElasticSearchClusterMapOutputWithContext(context.Context) MdbElasticSearchClusterMapOutput
}

type MdbElasticSearchClusterMap map[string]MdbElasticSearchClusterInput

func (MdbElasticSearchClusterMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*MdbElasticSearchCluster)(nil)).Elem()
}

func (i MdbElasticSearchClusterMap) ToMdbElasticSearchClusterMapOutput() MdbElasticSearchClusterMapOutput {
	return i.ToMdbElasticSearchClusterMapOutputWithContext(context.Background())
}

func (i MdbElasticSearchClusterMap) ToMdbElasticSearchClusterMapOutputWithContext(ctx context.Context) MdbElasticSearchClusterMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MdbElasticSearchClusterMapOutput)
}

type MdbElasticSearchClusterOutput struct{ *pulumi.OutputState }

func (MdbElasticSearchClusterOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**MdbElasticSearchCluster)(nil)).Elem()
}

func (o MdbElasticSearchClusterOutput) ToMdbElasticSearchClusterOutput() MdbElasticSearchClusterOutput {
	return o
}

func (o MdbElasticSearchClusterOutput) ToMdbElasticSearchClusterOutputWithContext(ctx context.Context) MdbElasticSearchClusterOutput {
	return o
}

// Configuration of the Elasticsearch cluster. The structure is documented below.
func (o MdbElasticSearchClusterOutput) Config() MdbElasticSearchClusterConfigOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) MdbElasticSearchClusterConfigOutput { return v.Config }).(MdbElasticSearchClusterConfigOutput)
}

// Creation timestamp of the key.
func (o MdbElasticSearchClusterOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.StringOutput { return v.CreatedAt }).(pulumi.StringOutput)
}

// Inhibits deletion of the cluster.  Can be either `true` or `false`.
func (o MdbElasticSearchClusterOutput) DeletionProtection() pulumi.BoolOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.BoolOutput { return v.DeletionProtection }).(pulumi.BoolOutput)
}

// Description of the Elasticsearch cluster.
func (o MdbElasticSearchClusterOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// Deployment environment of the Elasticsearch cluster. Can be either `PRESTABLE` or `PRODUCTION`.
func (o MdbElasticSearchClusterOutput) Environment() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.StringOutput { return v.Environment }).(pulumi.StringOutput)
}

// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
func (o MdbElasticSearchClusterOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.StringOutput { return v.FolderId }).(pulumi.StringOutput)
}

// Aggregated health of the cluster. Can be either `ALIVE`, `DEGRADED`, `DEAD` or `HEALTH_UNKNOWN`.
// For more information see `health` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-elasticsearch/api-ref/Cluster/).
func (o MdbElasticSearchClusterOutput) Health() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.StringOutput { return v.Health }).(pulumi.StringOutput)
}

// A host of the Elasticsearch cluster. The structure is documented below.
func (o MdbElasticSearchClusterOutput) Hosts() MdbElasticSearchClusterHostArrayOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) MdbElasticSearchClusterHostArrayOutput { return v.Hosts }).(MdbElasticSearchClusterHostArrayOutput)
}

// A set of key/value label pairs to assign to the Elasticsearch cluster.
func (o MdbElasticSearchClusterOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.StringMapOutput { return v.Labels }).(pulumi.StringMapOutput)
}

func (o MdbElasticSearchClusterOutput) MaintenanceWindow() MdbElasticSearchClusterMaintenanceWindowPtrOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) MdbElasticSearchClusterMaintenanceWindowPtrOutput {
		return v.MaintenanceWindow
	}).(MdbElasticSearchClusterMaintenanceWindowPtrOutput)
}

// User defined host name.
func (o MdbElasticSearchClusterOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// ID of the network, to which the Elasticsearch cluster belongs.
func (o MdbElasticSearchClusterOutput) NetworkId() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.StringOutput { return v.NetworkId }).(pulumi.StringOutput)
}

// A set of ids of security groups assigned to hosts of the cluster.
func (o MdbElasticSearchClusterOutput) SecurityGroupIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.StringArrayOutput { return v.SecurityGroupIds }).(pulumi.StringArrayOutput)
}

// ID of the service account authorized for this cluster.
func (o MdbElasticSearchClusterOutput) ServiceAccountId() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.StringPtrOutput { return v.ServiceAccountId }).(pulumi.StringPtrOutput)
}

// Status of the cluster. Can be either `CREATING`, `STARTING`, `RUNNING`, `UPDATING`, `STOPPING`, `STOPPED`, `ERROR` or `STATUS_UNKNOWN`.
// For more information see `status` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-elasticsearch/api-ref/Cluster/).
func (o MdbElasticSearchClusterOutput) Status() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) pulumi.StringOutput { return v.Status }).(pulumi.StringOutput)
}

func (o MdbElasticSearchClusterOutput) Timeouts() MdbElasticSearchClusterTimeoutsOutput {
	return o.ApplyT(func(v *MdbElasticSearchCluster) MdbElasticSearchClusterTimeoutsOutput { return v.Timeouts }).(MdbElasticSearchClusterTimeoutsOutput)
}

type MdbElasticSearchClusterArrayOutput struct{ *pulumi.OutputState }

func (MdbElasticSearchClusterArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*MdbElasticSearchCluster)(nil)).Elem()
}

func (o MdbElasticSearchClusterArrayOutput) ToMdbElasticSearchClusterArrayOutput() MdbElasticSearchClusterArrayOutput {
	return o
}

func (o MdbElasticSearchClusterArrayOutput) ToMdbElasticSearchClusterArrayOutputWithContext(ctx context.Context) MdbElasticSearchClusterArrayOutput {
	return o
}

func (o MdbElasticSearchClusterArrayOutput) Index(i pulumi.IntInput) MdbElasticSearchClusterOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *MdbElasticSearchCluster {
		return vs[0].([]*MdbElasticSearchCluster)[vs[1].(int)]
	}).(MdbElasticSearchClusterOutput)
}

type MdbElasticSearchClusterMapOutput struct{ *pulumi.OutputState }

func (MdbElasticSearchClusterMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*MdbElasticSearchCluster)(nil)).Elem()
}

func (o MdbElasticSearchClusterMapOutput) ToMdbElasticSearchClusterMapOutput() MdbElasticSearchClusterMapOutput {
	return o
}

func (o MdbElasticSearchClusterMapOutput) ToMdbElasticSearchClusterMapOutputWithContext(ctx context.Context) MdbElasticSearchClusterMapOutput {
	return o
}

func (o MdbElasticSearchClusterMapOutput) MapIndex(k pulumi.StringInput) MdbElasticSearchClusterOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *MdbElasticSearchCluster {
		return vs[0].(map[string]*MdbElasticSearchCluster)[vs[1].(string)]
	}).(MdbElasticSearchClusterOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*MdbElasticSearchClusterInput)(nil)).Elem(), &MdbElasticSearchCluster{})
	pulumi.RegisterInputType(reflect.TypeOf((*MdbElasticSearchClusterArrayInput)(nil)).Elem(), MdbElasticSearchClusterArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*MdbElasticSearchClusterMapInput)(nil)).Elem(), MdbElasticSearchClusterMap{})
	pulumi.RegisterOutputType(MdbElasticSearchClusterOutput{})
	pulumi.RegisterOutputType(MdbElasticSearchClusterArrayOutput{})
	pulumi.RegisterOutputType(MdbElasticSearchClusterMapOutput{})
}
