// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Get information about a Yandex Compute instance group.
func LookupComputeInstanceGroup(ctx *pulumi.Context, args *LookupComputeInstanceGroupArgs, opts ...pulumi.InvokeOption) (*LookupComputeInstanceGroupResult, error) {
	var rv LookupComputeInstanceGroupResult
	err := ctx.Invoke("yandex:index/getComputeInstanceGroup:getComputeInstanceGroup", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getComputeInstanceGroup.
type LookupComputeInstanceGroupArgs struct {
	// The ID of a specific instance group.
	InstanceGroupId string `pulumi:"instanceGroupId"`
}

// A collection of values returned by getComputeInstanceGroup.
type LookupComputeInstanceGroupResult struct {
	// The allocation policy of the instance group by zone and region. The structure is documented below.
	AllocationPolicies        []GetComputeInstanceGroupAllocationPolicy         `pulumi:"allocationPolicies"`
	ApplicationBalancerStates []GetComputeInstanceGroupApplicationBalancerState `pulumi:"applicationBalancerStates"`
	// Application Load balancing (L7) specifications. The structure is documented below.
	ApplicationLoadBalancers []GetComputeInstanceGroupApplicationLoadBalancer `pulumi:"applicationLoadBalancers"`
	// The instance group creation timestamp.
	CreatedAt string `pulumi:"createdAt"`
	// Flag that protects the instance group from accidental deletion.
	DeletionProtection bool `pulumi:"deletionProtection"`
	// The deployment policy of the instance group. The structure is documented below.
	DeployPolicies []GetComputeInstanceGroupDeployPolicy `pulumi:"deployPolicies"`
	// A description of the boot disk.
	Description string `pulumi:"description"`
	// Folder ID of custom metric in Yandex Monitoring that should be used for scaling.
	FolderId string `pulumi:"folderId"`
	// Health check specification. The structure is documented below.
	HealthChecks    []GetComputeInstanceGroupHealthCheck `pulumi:"healthChecks"`
	Id              string                               `pulumi:"id"`
	InstanceGroupId string                               `pulumi:"instanceGroupId"`
	// The instance template that the instance group belongs to. The structure is documented below.
	InstanceTemplates []GetComputeInstanceGroupInstanceTemplate `pulumi:"instanceTemplates"`
	// A list of instances in the specified instance group. The structure is documented below.
	Instances []GetComputeInstanceGroupInstance `pulumi:"instances"`
	// A map of labels applied to this instance.
	// * `resources.0.memory` - The memory size allocated to the instance.
	// * `resources.0.cores` - Number of CPU cores allocated to the instance.
	// * `resources.0.core_fraction` - Baseline core performance as a percent.
	// * `resources.0.gpus` - Number of GPU cores allocated to the instance.
	Labels map[string]string `pulumi:"labels"`
	// Information about which entities can be attached to this load balancer. The structure is documented below.
	LoadBalancerStates []GetComputeInstanceGroupLoadBalancerState `pulumi:"loadBalancerStates"`
	// Load balancing specification. The structure is documented below.
	LoadBalancers []GetComputeInstanceGroupLoadBalancer `pulumi:"loadBalancers"`
	// Timeout for waiting for the VM to become healthy. If the timeout is exceeded, the VM will be turned off based on the deployment policy. Specified in seconds.
	MaxCheckingHealthDuration float64 `pulumi:"maxCheckingHealthDuration"`
	// The name of the managed instance.
	Name string `pulumi:"name"`
	// The scaling policy of the instance group. The structure is documented below.
	ScalePolicies []GetComputeInstanceGroupScalePolicy `pulumi:"scalePolicies"`
	// The service account ID for the instance.
	ServiceAccountId string `pulumi:"serviceAccountId"`
	// The status of the instance.
	Status string `pulumi:"status"`
	// A set of key/value  variables pairs to assign to the instance group.
	Variables map[string]string `pulumi:"variables"`
}

func LookupComputeInstanceGroupOutput(ctx *pulumi.Context, args LookupComputeInstanceGroupOutputArgs, opts ...pulumi.InvokeOption) LookupComputeInstanceGroupResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (LookupComputeInstanceGroupResult, error) {
			args := v.(LookupComputeInstanceGroupArgs)
			r, err := LookupComputeInstanceGroup(ctx, &args, opts...)
			var s LookupComputeInstanceGroupResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(LookupComputeInstanceGroupResultOutput)
}

// A collection of arguments for invoking getComputeInstanceGroup.
type LookupComputeInstanceGroupOutputArgs struct {
	// The ID of a specific instance group.
	InstanceGroupId pulumi.StringInput `pulumi:"instanceGroupId"`
}

func (LookupComputeInstanceGroupOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupComputeInstanceGroupArgs)(nil)).Elem()
}

// A collection of values returned by getComputeInstanceGroup.
type LookupComputeInstanceGroupResultOutput struct{ *pulumi.OutputState }

func (LookupComputeInstanceGroupResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupComputeInstanceGroupResult)(nil)).Elem()
}

func (o LookupComputeInstanceGroupResultOutput) ToLookupComputeInstanceGroupResultOutput() LookupComputeInstanceGroupResultOutput {
	return o
}

func (o LookupComputeInstanceGroupResultOutput) ToLookupComputeInstanceGroupResultOutputWithContext(ctx context.Context) LookupComputeInstanceGroupResultOutput {
	return o
}

// The allocation policy of the instance group by zone and region. The structure is documented below.
func (o LookupComputeInstanceGroupResultOutput) AllocationPolicies() GetComputeInstanceGroupAllocationPolicyArrayOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) []GetComputeInstanceGroupAllocationPolicy {
		return v.AllocationPolicies
	}).(GetComputeInstanceGroupAllocationPolicyArrayOutput)
}

func (o LookupComputeInstanceGroupResultOutput) ApplicationBalancerStates() GetComputeInstanceGroupApplicationBalancerStateArrayOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) []GetComputeInstanceGroupApplicationBalancerState {
		return v.ApplicationBalancerStates
	}).(GetComputeInstanceGroupApplicationBalancerStateArrayOutput)
}

// Application Load balancing (L7) specifications. The structure is documented below.
func (o LookupComputeInstanceGroupResultOutput) ApplicationLoadBalancers() GetComputeInstanceGroupApplicationLoadBalancerArrayOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) []GetComputeInstanceGroupApplicationLoadBalancer {
		return v.ApplicationLoadBalancers
	}).(GetComputeInstanceGroupApplicationLoadBalancerArrayOutput)
}

// The instance group creation timestamp.
func (o LookupComputeInstanceGroupResultOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) string { return v.CreatedAt }).(pulumi.StringOutput)
}

// Flag that protects the instance group from accidental deletion.
func (o LookupComputeInstanceGroupResultOutput) DeletionProtection() pulumi.BoolOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) bool { return v.DeletionProtection }).(pulumi.BoolOutput)
}

// The deployment policy of the instance group. The structure is documented below.
func (o LookupComputeInstanceGroupResultOutput) DeployPolicies() GetComputeInstanceGroupDeployPolicyArrayOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) []GetComputeInstanceGroupDeployPolicy {
		return v.DeployPolicies
	}).(GetComputeInstanceGroupDeployPolicyArrayOutput)
}

// A description of the boot disk.
func (o LookupComputeInstanceGroupResultOutput) Description() pulumi.StringOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) string { return v.Description }).(pulumi.StringOutput)
}

// Folder ID of custom metric in Yandex Monitoring that should be used for scaling.
func (o LookupComputeInstanceGroupResultOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) string { return v.FolderId }).(pulumi.StringOutput)
}

// Health check specification. The structure is documented below.
func (o LookupComputeInstanceGroupResultOutput) HealthChecks() GetComputeInstanceGroupHealthCheckArrayOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) []GetComputeInstanceGroupHealthCheck { return v.HealthChecks }).(GetComputeInstanceGroupHealthCheckArrayOutput)
}

func (o LookupComputeInstanceGroupResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) string { return v.Id }).(pulumi.StringOutput)
}

func (o LookupComputeInstanceGroupResultOutput) InstanceGroupId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) string { return v.InstanceGroupId }).(pulumi.StringOutput)
}

// The instance template that the instance group belongs to. The structure is documented below.
func (o LookupComputeInstanceGroupResultOutput) InstanceTemplates() GetComputeInstanceGroupInstanceTemplateArrayOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) []GetComputeInstanceGroupInstanceTemplate {
		return v.InstanceTemplates
	}).(GetComputeInstanceGroupInstanceTemplateArrayOutput)
}

// A list of instances in the specified instance group. The structure is documented below.
func (o LookupComputeInstanceGroupResultOutput) Instances() GetComputeInstanceGroupInstanceArrayOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) []GetComputeInstanceGroupInstance { return v.Instances }).(GetComputeInstanceGroupInstanceArrayOutput)
}

// A map of labels applied to this instance.
// * `resources.0.memory` - The memory size allocated to the instance.
// * `resources.0.cores` - Number of CPU cores allocated to the instance.
// * `resources.0.core_fraction` - Baseline core performance as a percent.
// * `resources.0.gpus` - Number of GPU cores allocated to the instance.
func (o LookupComputeInstanceGroupResultOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) map[string]string { return v.Labels }).(pulumi.StringMapOutput)
}

// Information about which entities can be attached to this load balancer. The structure is documented below.
func (o LookupComputeInstanceGroupResultOutput) LoadBalancerStates() GetComputeInstanceGroupLoadBalancerStateArrayOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) []GetComputeInstanceGroupLoadBalancerState {
		return v.LoadBalancerStates
	}).(GetComputeInstanceGroupLoadBalancerStateArrayOutput)
}

// Load balancing specification. The structure is documented below.
func (o LookupComputeInstanceGroupResultOutput) LoadBalancers() GetComputeInstanceGroupLoadBalancerArrayOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) []GetComputeInstanceGroupLoadBalancer { return v.LoadBalancers }).(GetComputeInstanceGroupLoadBalancerArrayOutput)
}

// Timeout for waiting for the VM to become healthy. If the timeout is exceeded, the VM will be turned off based on the deployment policy. Specified in seconds.
func (o LookupComputeInstanceGroupResultOutput) MaxCheckingHealthDuration() pulumi.Float64Output {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) float64 { return v.MaxCheckingHealthDuration }).(pulumi.Float64Output)
}

// The name of the managed instance.
func (o LookupComputeInstanceGroupResultOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) string { return v.Name }).(pulumi.StringOutput)
}

// The scaling policy of the instance group. The structure is documented below.
func (o LookupComputeInstanceGroupResultOutput) ScalePolicies() GetComputeInstanceGroupScalePolicyArrayOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) []GetComputeInstanceGroupScalePolicy { return v.ScalePolicies }).(GetComputeInstanceGroupScalePolicyArrayOutput)
}

// The service account ID for the instance.
func (o LookupComputeInstanceGroupResultOutput) ServiceAccountId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) string { return v.ServiceAccountId }).(pulumi.StringOutput)
}

// The status of the instance.
func (o LookupComputeInstanceGroupResultOutput) Status() pulumi.StringOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) string { return v.Status }).(pulumi.StringOutput)
}

// A set of key/value  variables pairs to assign to the instance group.
func (o LookupComputeInstanceGroupResultOutput) Variables() pulumi.StringMapOutput {
	return o.ApplyT(func(v LookupComputeInstanceGroupResult) map[string]string { return v.Variables }).(pulumi.StringMapOutput)
}

func init() {
	pulumi.RegisterOutputType(LookupComputeInstanceGroupResultOutput{})
}