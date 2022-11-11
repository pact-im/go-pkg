// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Get information about a Yandex Managed MongoDB cluster. For more information, see
// [the official documentation](https://cloud.yandex.com/docs/managed-mongodb/concepts).
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
// 		foo, err := yandex.LookupMdbMongodbCluster(ctx, &GetMdbMongodbClusterArgs{
// 			Name: pulumi.StringRef("test"),
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		ctx.Export("networkId", foo.NetworkId)
// 		return nil
// 	})
// }
// ```
func LookupMdbMongodbCluster(ctx *pulumi.Context, args *LookupMdbMongodbClusterArgs, opts ...pulumi.InvokeOption) (*LookupMdbMongodbClusterResult, error) {
	var rv LookupMdbMongodbClusterResult
	err := ctx.Invoke("yandex:index/getMdbMongodbCluster:getMdbMongodbCluster", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getMdbMongodbCluster.
type LookupMdbMongodbClusterArgs struct {
	// Configuration of the MongoDB cluster. The structure is documented below.
	ClusterConfig *GetMdbMongodbClusterClusterConfig `pulumi:"clusterConfig"`
	// The ID of the MongoDB cluster.
	ClusterId *string `pulumi:"clusterId"`
	// Creation timestamp of the key.
	CreatedAt *string `pulumi:"createdAt"`
	// A database of the MongoDB cluster. The structure is documented below.
	Databases          []GetMdbMongodbClusterDatabase `pulumi:"databases"`
	DeletionProtection *bool                          `pulumi:"deletionProtection"`
	// Description of the MongoDB cluster.
	Description *string `pulumi:"description"`
	// Deployment environment of the MongoDB cluster.
	Environment *string `pulumi:"environment"`
	// Folder that the resource belongs to. If value is omitted, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// The health of the host.
	Health *string `pulumi:"health"`
	// A host of the MongoDB cluster. The structure is documented below.
	Hosts []GetMdbMongodbClusterHost `pulumi:"hosts"`
	// A set of key/value label pairs to assign to the MongoDB cluster.
	Labels            map[string]string                      `pulumi:"labels"`
	MaintenanceWindow *GetMdbMongodbClusterMaintenanceWindow `pulumi:"maintenanceWindow"`
	// The name of the MongoDB cluster.
	Name *string `pulumi:"name"`
	// ID of the network, to which the MongoDB cluster belongs.
	NetworkId *string `pulumi:"networkId"`
	// Resources allocated to hosts of the MongoDB cluster. The structure is documented below.
	Resources *GetMdbMongodbClusterResources `pulumi:"resources"`
	// A set of ids of security groups assigned to hosts of the cluster.
	SecurityGroupIds []string `pulumi:"securityGroupIds"`
	// MongoDB Cluster mode enabled/disabled.
	Sharded *bool `pulumi:"sharded"`
	// Status of the cluster.
	Status *string `pulumi:"status"`
	// A user of the MongoDB cluster. The structure is documented below.
	Users []GetMdbMongodbClusterUser `pulumi:"users"`
}

// A collection of values returned by getMdbMongodbCluster.
type LookupMdbMongodbClusterResult struct {
	// Configuration of the MongoDB cluster. The structure is documented below.
	ClusterConfig *GetMdbMongodbClusterClusterConfig `pulumi:"clusterConfig"`
	ClusterId     string                             `pulumi:"clusterId"`
	// Creation timestamp of the key.
	CreatedAt string `pulumi:"createdAt"`
	// A database of the MongoDB cluster. The structure is documented below.
	Databases          []GetMdbMongodbClusterDatabase `pulumi:"databases"`
	DeletionProtection bool                           `pulumi:"deletionProtection"`
	// Description of the MongoDB cluster.
	Description string `pulumi:"description"`
	// Deployment environment of the MongoDB cluster.
	Environment *string `pulumi:"environment"`
	FolderId    string  `pulumi:"folderId"`
	// The health of the host.
	Health string `pulumi:"health"`
	// A host of the MongoDB cluster. The structure is documented below.
	Hosts []GetMdbMongodbClusterHost `pulumi:"hosts"`
	Id    string                     `pulumi:"id"`
	// A set of key/value label pairs to assign to the MongoDB cluster.
	Labels            map[string]string                      `pulumi:"labels"`
	MaintenanceWindow *GetMdbMongodbClusterMaintenanceWindow `pulumi:"maintenanceWindow"`
	// The name of the database.
	Name *string `pulumi:"name"`
	// ID of the network, to which the MongoDB cluster belongs.
	NetworkId *string `pulumi:"networkId"`
	// Resources allocated to hosts of the MongoDB cluster. The structure is documented below.
	Resources *GetMdbMongodbClusterResources `pulumi:"resources"`
	// A set of ids of security groups assigned to hosts of the cluster.
	SecurityGroupIds []string `pulumi:"securityGroupIds"`
	// MongoDB Cluster mode enabled/disabled.
	Sharded bool `pulumi:"sharded"`
	// Status of the cluster.
	Status string `pulumi:"status"`
	// A user of the MongoDB cluster. The structure is documented below.
	Users []GetMdbMongodbClusterUser `pulumi:"users"`
}

func LookupMdbMongodbClusterOutput(ctx *pulumi.Context, args LookupMdbMongodbClusterOutputArgs, opts ...pulumi.InvokeOption) LookupMdbMongodbClusterResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (LookupMdbMongodbClusterResult, error) {
			args := v.(LookupMdbMongodbClusterArgs)
			r, err := LookupMdbMongodbCluster(ctx, &args, opts...)
			var s LookupMdbMongodbClusterResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(LookupMdbMongodbClusterResultOutput)
}

// A collection of arguments for invoking getMdbMongodbCluster.
type LookupMdbMongodbClusterOutputArgs struct {
	// Configuration of the MongoDB cluster. The structure is documented below.
	ClusterConfig GetMdbMongodbClusterClusterConfigPtrInput `pulumi:"clusterConfig"`
	// The ID of the MongoDB cluster.
	ClusterId pulumi.StringPtrInput `pulumi:"clusterId"`
	// Creation timestamp of the key.
	CreatedAt pulumi.StringPtrInput `pulumi:"createdAt"`
	// A database of the MongoDB cluster. The structure is documented below.
	Databases          GetMdbMongodbClusterDatabaseArrayInput `pulumi:"databases"`
	DeletionProtection pulumi.BoolPtrInput                    `pulumi:"deletionProtection"`
	// Description of the MongoDB cluster.
	Description pulumi.StringPtrInput `pulumi:"description"`
	// Deployment environment of the MongoDB cluster.
	Environment pulumi.StringPtrInput `pulumi:"environment"`
	// Folder that the resource belongs to. If value is omitted, the default provider folder is used.
	FolderId pulumi.StringPtrInput `pulumi:"folderId"`
	// The health of the host.
	Health pulumi.StringPtrInput `pulumi:"health"`
	// A host of the MongoDB cluster. The structure is documented below.
	Hosts GetMdbMongodbClusterHostArrayInput `pulumi:"hosts"`
	// A set of key/value label pairs to assign to the MongoDB cluster.
	Labels            pulumi.StringMapInput                         `pulumi:"labels"`
	MaintenanceWindow GetMdbMongodbClusterMaintenanceWindowPtrInput `pulumi:"maintenanceWindow"`
	// The name of the MongoDB cluster.
	Name pulumi.StringPtrInput `pulumi:"name"`
	// ID of the network, to which the MongoDB cluster belongs.
	NetworkId pulumi.StringPtrInput `pulumi:"networkId"`
	// Resources allocated to hosts of the MongoDB cluster. The structure is documented below.
	Resources GetMdbMongodbClusterResourcesPtrInput `pulumi:"resources"`
	// A set of ids of security groups assigned to hosts of the cluster.
	SecurityGroupIds pulumi.StringArrayInput `pulumi:"securityGroupIds"`
	// MongoDB Cluster mode enabled/disabled.
	Sharded pulumi.BoolPtrInput `pulumi:"sharded"`
	// Status of the cluster.
	Status pulumi.StringPtrInput `pulumi:"status"`
	// A user of the MongoDB cluster. The structure is documented below.
	Users GetMdbMongodbClusterUserArrayInput `pulumi:"users"`
}

func (LookupMdbMongodbClusterOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupMdbMongodbClusterArgs)(nil)).Elem()
}

// A collection of values returned by getMdbMongodbCluster.
type LookupMdbMongodbClusterResultOutput struct{ *pulumi.OutputState }

func (LookupMdbMongodbClusterResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupMdbMongodbClusterResult)(nil)).Elem()
}

func (o LookupMdbMongodbClusterResultOutput) ToLookupMdbMongodbClusterResultOutput() LookupMdbMongodbClusterResultOutput {
	return o
}

func (o LookupMdbMongodbClusterResultOutput) ToLookupMdbMongodbClusterResultOutputWithContext(ctx context.Context) LookupMdbMongodbClusterResultOutput {
	return o
}

// Configuration of the MongoDB cluster. The structure is documented below.
func (o LookupMdbMongodbClusterResultOutput) ClusterConfig() GetMdbMongodbClusterClusterConfigPtrOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) *GetMdbMongodbClusterClusterConfig { return v.ClusterConfig }).(GetMdbMongodbClusterClusterConfigPtrOutput)
}

func (o LookupMdbMongodbClusterResultOutput) ClusterId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) string { return v.ClusterId }).(pulumi.StringOutput)
}

// Creation timestamp of the key.
func (o LookupMdbMongodbClusterResultOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) string { return v.CreatedAt }).(pulumi.StringOutput)
}

// A database of the MongoDB cluster. The structure is documented below.
func (o LookupMdbMongodbClusterResultOutput) Databases() GetMdbMongodbClusterDatabaseArrayOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) []GetMdbMongodbClusterDatabase { return v.Databases }).(GetMdbMongodbClusterDatabaseArrayOutput)
}

func (o LookupMdbMongodbClusterResultOutput) DeletionProtection() pulumi.BoolOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) bool { return v.DeletionProtection }).(pulumi.BoolOutput)
}

// Description of the MongoDB cluster.
func (o LookupMdbMongodbClusterResultOutput) Description() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) string { return v.Description }).(pulumi.StringOutput)
}

// Deployment environment of the MongoDB cluster.
func (o LookupMdbMongodbClusterResultOutput) Environment() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) *string { return v.Environment }).(pulumi.StringPtrOutput)
}

func (o LookupMdbMongodbClusterResultOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) string { return v.FolderId }).(pulumi.StringOutput)
}

// The health of the host.
func (o LookupMdbMongodbClusterResultOutput) Health() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) string { return v.Health }).(pulumi.StringOutput)
}

// A host of the MongoDB cluster. The structure is documented below.
func (o LookupMdbMongodbClusterResultOutput) Hosts() GetMdbMongodbClusterHostArrayOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) []GetMdbMongodbClusterHost { return v.Hosts }).(GetMdbMongodbClusterHostArrayOutput)
}

func (o LookupMdbMongodbClusterResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) string { return v.Id }).(pulumi.StringOutput)
}

// A set of key/value label pairs to assign to the MongoDB cluster.
func (o LookupMdbMongodbClusterResultOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) map[string]string { return v.Labels }).(pulumi.StringMapOutput)
}

func (o LookupMdbMongodbClusterResultOutput) MaintenanceWindow() GetMdbMongodbClusterMaintenanceWindowPtrOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) *GetMdbMongodbClusterMaintenanceWindow {
		return v.MaintenanceWindow
	}).(GetMdbMongodbClusterMaintenanceWindowPtrOutput)
}

// The name of the database.
func (o LookupMdbMongodbClusterResultOutput) Name() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) *string { return v.Name }).(pulumi.StringPtrOutput)
}

// ID of the network, to which the MongoDB cluster belongs.
func (o LookupMdbMongodbClusterResultOutput) NetworkId() pulumi.StringPtrOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) *string { return v.NetworkId }).(pulumi.StringPtrOutput)
}

// Resources allocated to hosts of the MongoDB cluster. The structure is documented below.
func (o LookupMdbMongodbClusterResultOutput) Resources() GetMdbMongodbClusterResourcesPtrOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) *GetMdbMongodbClusterResources { return v.Resources }).(GetMdbMongodbClusterResourcesPtrOutput)
}

// A set of ids of security groups assigned to hosts of the cluster.
func (o LookupMdbMongodbClusterResultOutput) SecurityGroupIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) []string { return v.SecurityGroupIds }).(pulumi.StringArrayOutput)
}

// MongoDB Cluster mode enabled/disabled.
func (o LookupMdbMongodbClusterResultOutput) Sharded() pulumi.BoolOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) bool { return v.Sharded }).(pulumi.BoolOutput)
}

// Status of the cluster.
func (o LookupMdbMongodbClusterResultOutput) Status() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) string { return v.Status }).(pulumi.StringOutput)
}

// A user of the MongoDB cluster. The structure is documented below.
func (o LookupMdbMongodbClusterResultOutput) Users() GetMdbMongodbClusterUserArrayOutput {
	return o.ApplyT(func(v LookupMdbMongodbClusterResult) []GetMdbMongodbClusterUser { return v.Users }).(GetMdbMongodbClusterUserArrayOutput)
}

func init() {
	pulumi.RegisterOutputType(LookupMdbMongodbClusterResultOutput{})
}