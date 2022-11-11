// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Get information about a Yandex Managed Greenplum cluster. For more information, see
// [the official documentation](https://cloud.yandex.com/docs/managed-greenplum/).
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
// 		foo, err := yandex.LookupMdbGreenplumCluster(ctx, &GetMdbGreenplumClusterArgs{
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
func LookupMdbGreenplumCluster(ctx *pulumi.Context, args *LookupMdbGreenplumClusterArgs, opts ...pulumi.InvokeOption) (*LookupMdbGreenplumClusterResult, error) {
	var rv LookupMdbGreenplumClusterResult
	err := ctx.Invoke("yandex:index/getMdbGreenplumCluster:getMdbGreenplumCluster", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getMdbGreenplumCluster.
type LookupMdbGreenplumClusterArgs struct {
	// The ID of the Greenplum cluster.
	ClusterId *string `pulumi:"clusterId"`
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// Greenplum cluster config.
	GreenplumConfig map[string]string `pulumi:"greenplumConfig"`
	// The name of the Greenplum cluster.
	Name *string `pulumi:"name"`
	// Configuration of the connection pooler. The structure is documented below.
	PoolerConfig *GetMdbGreenplumClusterPoolerConfig `pulumi:"poolerConfig"`
}

// A collection of values returned by getMdbGreenplumCluster.
type LookupMdbGreenplumClusterResult struct {
	// Access policy to the Greenplum cluster. The structure is documented below.
	Accesses []GetMdbGreenplumClusterAccess `pulumi:"accesses"`
	// Flag that indicates whether master hosts was created with a public IP.
	AssignPublicIp bool `pulumi:"assignPublicIp"`
	// Time to start the daily backup, in the UTC timezone. The structure is documented below.
	BackupWindowStarts []GetMdbGreenplumClusterBackupWindowStart `pulumi:"backupWindowStarts"`
	ClusterId          string                                    `pulumi:"clusterId"`
	// Timestamp of cluster creation.
	CreatedAt string `pulumi:"createdAt"`
	// Flag to protect the cluster from deletion.
	DeletionProtection bool `pulumi:"deletionProtection"`
	// Description of the Greenplum cluster.
	Description string `pulumi:"description"`
	// Deployment environment of the Greenplum cluster.
	Environment string `pulumi:"environment"`
	FolderId    string `pulumi:"folderId"`
	// Greenplum cluster config.
	GreenplumConfig map[string]string `pulumi:"greenplumConfig"`
	// Aggregated health of the cluster.
	Health string `pulumi:"health"`
	Id     string `pulumi:"id"`
	// A set of key/value label pairs to assign to the Greenplum cluster.
	Labels map[string]string `pulumi:"labels"`
	// Maintenance window settings of the Greenplum cluster. The structure is documented below.
	MaintenanceWindows []GetMdbGreenplumClusterMaintenanceWindow `pulumi:"maintenanceWindows"`
	// Number of hosts in master subcluster.
	MasterHostCount float64 `pulumi:"masterHostCount"`
	// Info about hosts in master subcluster. The structure is documented below.
	MasterHosts []GetMdbGreenplumClusterMasterHost `pulumi:"masterHosts"`
	// Settings for master subcluster. The structure is documented below.
	MasterSubclusters []GetMdbGreenplumClusterMasterSubcluster `pulumi:"masterSubclusters"`
	Name              string                                   `pulumi:"name"`
	// ID of the network, to which the Greenplum cluster belongs.
	NetworkId string `pulumi:"networkId"`
	// Configuration of the connection pooler. The structure is documented below.
	PoolerConfig *GetMdbGreenplumClusterPoolerConfig `pulumi:"poolerConfig"`
	// A set of ids of security groups assigned to hosts of the cluster.
	SecurityGroupIds []string `pulumi:"securityGroupIds"`
	// Number of hosts in segment subcluster.
	SegmentHostCount float64 `pulumi:"segmentHostCount"`
	// Info about hosts in segment subcluster. The structure is documented below.
	SegmentHosts []GetMdbGreenplumClusterSegmentHost `pulumi:"segmentHosts"`
	// Number of segments on segment host.
	SegmentInHost float64 `pulumi:"segmentInHost"`
	// Settings for segment subcluster. The structure is documented below.
	SegmentSubclusters []GetMdbGreenplumClusterSegmentSubcluster `pulumi:"segmentSubclusters"`
	// Status of the cluster.
	Status string `pulumi:"status"`
	// The ID of the subnet, to which the hosts belongs. The subnet must be a part of the network to which the cluster belongs.
	SubnetId string `pulumi:"subnetId"`
	// Greenplum cluster admin user name.
	UserName string `pulumi:"userName"`
	// Version of the Greenplum cluster.
	Version string `pulumi:"version"`
	// The availability zone where the Greenplum hosts will be created.
	Zone string `pulumi:"zone"`
}

func LookupMdbGreenplumClusterOutput(ctx *pulumi.Context, args LookupMdbGreenplumClusterOutputArgs, opts ...pulumi.InvokeOption) LookupMdbGreenplumClusterResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (LookupMdbGreenplumClusterResult, error) {
			args := v.(LookupMdbGreenplumClusterArgs)
			r, err := LookupMdbGreenplumCluster(ctx, &args, opts...)
			var s LookupMdbGreenplumClusterResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(LookupMdbGreenplumClusterResultOutput)
}

// A collection of arguments for invoking getMdbGreenplumCluster.
type LookupMdbGreenplumClusterOutputArgs struct {
	// The ID of the Greenplum cluster.
	ClusterId pulumi.StringPtrInput `pulumi:"clusterId"`
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId pulumi.StringPtrInput `pulumi:"folderId"`
	// Greenplum cluster config.
	GreenplumConfig pulumi.StringMapInput `pulumi:"greenplumConfig"`
	// The name of the Greenplum cluster.
	Name pulumi.StringPtrInput `pulumi:"name"`
	// Configuration of the connection pooler. The structure is documented below.
	PoolerConfig GetMdbGreenplumClusterPoolerConfigPtrInput `pulumi:"poolerConfig"`
}

func (LookupMdbGreenplumClusterOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupMdbGreenplumClusterArgs)(nil)).Elem()
}

// A collection of values returned by getMdbGreenplumCluster.
type LookupMdbGreenplumClusterResultOutput struct{ *pulumi.OutputState }

func (LookupMdbGreenplumClusterResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupMdbGreenplumClusterResult)(nil)).Elem()
}

func (o LookupMdbGreenplumClusterResultOutput) ToLookupMdbGreenplumClusterResultOutput() LookupMdbGreenplumClusterResultOutput {
	return o
}

func (o LookupMdbGreenplumClusterResultOutput) ToLookupMdbGreenplumClusterResultOutputWithContext(ctx context.Context) LookupMdbGreenplumClusterResultOutput {
	return o
}

// Access policy to the Greenplum cluster. The structure is documented below.
func (o LookupMdbGreenplumClusterResultOutput) Accesses() GetMdbGreenplumClusterAccessArrayOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) []GetMdbGreenplumClusterAccess { return v.Accesses }).(GetMdbGreenplumClusterAccessArrayOutput)
}

// Flag that indicates whether master hosts was created with a public IP.
func (o LookupMdbGreenplumClusterResultOutput) AssignPublicIp() pulumi.BoolOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) bool { return v.AssignPublicIp }).(pulumi.BoolOutput)
}

// Time to start the daily backup, in the UTC timezone. The structure is documented below.
func (o LookupMdbGreenplumClusterResultOutput) BackupWindowStarts() GetMdbGreenplumClusterBackupWindowStartArrayOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) []GetMdbGreenplumClusterBackupWindowStart {
		return v.BackupWindowStarts
	}).(GetMdbGreenplumClusterBackupWindowStartArrayOutput)
}

func (o LookupMdbGreenplumClusterResultOutput) ClusterId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.ClusterId }).(pulumi.StringOutput)
}

// Timestamp of cluster creation.
func (o LookupMdbGreenplumClusterResultOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.CreatedAt }).(pulumi.StringOutput)
}

// Flag to protect the cluster from deletion.
func (o LookupMdbGreenplumClusterResultOutput) DeletionProtection() pulumi.BoolOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) bool { return v.DeletionProtection }).(pulumi.BoolOutput)
}

// Description of the Greenplum cluster.
func (o LookupMdbGreenplumClusterResultOutput) Description() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.Description }).(pulumi.StringOutput)
}

// Deployment environment of the Greenplum cluster.
func (o LookupMdbGreenplumClusterResultOutput) Environment() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.Environment }).(pulumi.StringOutput)
}

func (o LookupMdbGreenplumClusterResultOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.FolderId }).(pulumi.StringOutput)
}

// Greenplum cluster config.
func (o LookupMdbGreenplumClusterResultOutput) GreenplumConfig() pulumi.StringMapOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) map[string]string { return v.GreenplumConfig }).(pulumi.StringMapOutput)
}

// Aggregated health of the cluster.
func (o LookupMdbGreenplumClusterResultOutput) Health() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.Health }).(pulumi.StringOutput)
}

func (o LookupMdbGreenplumClusterResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.Id }).(pulumi.StringOutput)
}

// A set of key/value label pairs to assign to the Greenplum cluster.
func (o LookupMdbGreenplumClusterResultOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) map[string]string { return v.Labels }).(pulumi.StringMapOutput)
}

// Maintenance window settings of the Greenplum cluster. The structure is documented below.
func (o LookupMdbGreenplumClusterResultOutput) MaintenanceWindows() GetMdbGreenplumClusterMaintenanceWindowArrayOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) []GetMdbGreenplumClusterMaintenanceWindow {
		return v.MaintenanceWindows
	}).(GetMdbGreenplumClusterMaintenanceWindowArrayOutput)
}

// Number of hosts in master subcluster.
func (o LookupMdbGreenplumClusterResultOutput) MasterHostCount() pulumi.Float64Output {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) float64 { return v.MasterHostCount }).(pulumi.Float64Output)
}

// Info about hosts in master subcluster. The structure is documented below.
func (o LookupMdbGreenplumClusterResultOutput) MasterHosts() GetMdbGreenplumClusterMasterHostArrayOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) []GetMdbGreenplumClusterMasterHost { return v.MasterHosts }).(GetMdbGreenplumClusterMasterHostArrayOutput)
}

// Settings for master subcluster. The structure is documented below.
func (o LookupMdbGreenplumClusterResultOutput) MasterSubclusters() GetMdbGreenplumClusterMasterSubclusterArrayOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) []GetMdbGreenplumClusterMasterSubcluster {
		return v.MasterSubclusters
	}).(GetMdbGreenplumClusterMasterSubclusterArrayOutput)
}

func (o LookupMdbGreenplumClusterResultOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.Name }).(pulumi.StringOutput)
}

// ID of the network, to which the Greenplum cluster belongs.
func (o LookupMdbGreenplumClusterResultOutput) NetworkId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.NetworkId }).(pulumi.StringOutput)
}

// Configuration of the connection pooler. The structure is documented below.
func (o LookupMdbGreenplumClusterResultOutput) PoolerConfig() GetMdbGreenplumClusterPoolerConfigPtrOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) *GetMdbGreenplumClusterPoolerConfig { return v.PoolerConfig }).(GetMdbGreenplumClusterPoolerConfigPtrOutput)
}

// A set of ids of security groups assigned to hosts of the cluster.
func (o LookupMdbGreenplumClusterResultOutput) SecurityGroupIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) []string { return v.SecurityGroupIds }).(pulumi.StringArrayOutput)
}

// Number of hosts in segment subcluster.
func (o LookupMdbGreenplumClusterResultOutput) SegmentHostCount() pulumi.Float64Output {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) float64 { return v.SegmentHostCount }).(pulumi.Float64Output)
}

// Info about hosts in segment subcluster. The structure is documented below.
func (o LookupMdbGreenplumClusterResultOutput) SegmentHosts() GetMdbGreenplumClusterSegmentHostArrayOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) []GetMdbGreenplumClusterSegmentHost { return v.SegmentHosts }).(GetMdbGreenplumClusterSegmentHostArrayOutput)
}

// Number of segments on segment host.
func (o LookupMdbGreenplumClusterResultOutput) SegmentInHost() pulumi.Float64Output {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) float64 { return v.SegmentInHost }).(pulumi.Float64Output)
}

// Settings for segment subcluster. The structure is documented below.
func (o LookupMdbGreenplumClusterResultOutput) SegmentSubclusters() GetMdbGreenplumClusterSegmentSubclusterArrayOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) []GetMdbGreenplumClusterSegmentSubcluster {
		return v.SegmentSubclusters
	}).(GetMdbGreenplumClusterSegmentSubclusterArrayOutput)
}

// Status of the cluster.
func (o LookupMdbGreenplumClusterResultOutput) Status() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.Status }).(pulumi.StringOutput)
}

// The ID of the subnet, to which the hosts belongs. The subnet must be a part of the network to which the cluster belongs.
func (o LookupMdbGreenplumClusterResultOutput) SubnetId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.SubnetId }).(pulumi.StringOutput)
}

// Greenplum cluster admin user name.
func (o LookupMdbGreenplumClusterResultOutput) UserName() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.UserName }).(pulumi.StringOutput)
}

// Version of the Greenplum cluster.
func (o LookupMdbGreenplumClusterResultOutput) Version() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.Version }).(pulumi.StringOutput)
}

// The availability zone where the Greenplum hosts will be created.
func (o LookupMdbGreenplumClusterResultOutput) Zone() pulumi.StringOutput {
	return o.ApplyT(func(v LookupMdbGreenplumClusterResult) string { return v.Zone }).(pulumi.StringOutput)
}

func init() {
	pulumi.RegisterOutputType(LookupMdbGreenplumClusterResultOutput{})
}
