// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Manages a single Secuirity Group Rule within the Yandex.Cloud. For more information, see the official documentation
// of [security groups](https://cloud.yandex.com/docs/vpc/concepts/security-groups)
// and [security group rules](https://cloud.yandex.com/docs/vpc/concepts/security-groups#rules).
//
// > **NOTE:** There is another way to manage security group rules by `ingress` and `egress` arguments in yandex_vpc_security_group. Both ways are equivalent but not compatible now. Using in-line rules of VpcSecurityGroup with Security Group Rule resource at the same time will cause a conflict of rules configuration.
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
// 		_, err := yandex.NewVpcNetwork(ctx, "lab-net", nil)
// 		if err != nil {
// 			return err
// 		}
// 		group1, err := yandex.NewVpcSecurityGroup(ctx, "group1", &yandex.VpcSecurityGroupArgs{
// 			Description: pulumi.String("description for my security group"),
// 			NetworkId:   lab_net.ID(),
// 			Labels: pulumi.StringMap{
// 				"my-label": pulumi.String("my-label-value"),
// 			},
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		_, err = yandex.NewVpcSecurityGroupRule(ctx, "rule1", &yandex.VpcSecurityGroupRuleArgs{
// 			SecurityGroupBinding: group1.ID(),
// 			Direction:            pulumi.String("ingress"),
// 			Description:          pulumi.String("rule1 description"),
// 			V4CidrBlocks: pulumi.StringArray{
// 				pulumi.String("10.0.1.0/24"),
// 				pulumi.String("10.0.2.0/24"),
// 			},
// 			Port:     pulumi.Float64(8080),
// 			Protocol: pulumi.String("TCP"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		_, err = yandex.NewVpcSecurityGroupRule(ctx, "rule2", &yandex.VpcSecurityGroupRuleArgs{
// 			SecurityGroupBinding: group1.ID(),
// 			Direction:            pulumi.String("egress"),
// 			Description:          pulumi.String("rule2 description"),
// 			V4CidrBlocks: pulumi.StringArray{
// 				pulumi.String("10.0.1.0/24"),
// 			},
// 			FromPort: pulumi.Float64(8090),
// 			ToPort:   pulumi.Float64(8099),
// 			Protocol: pulumi.String("UDP"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
type VpcSecurityGroupRule struct {
	pulumi.CustomResourceState

	// Description of the rule.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// direction of the rule. Can be `ingress` (inbound) or `egress` (outbound).
	Direction pulumi.StringOutput `pulumi:"direction"`
	// Minimum port number.
	FromPort pulumi.Float64PtrOutput `pulumi:"fromPort"`
	// Labels to assign to this rule.
	Labels pulumi.StringMapOutput `pulumi:"labels"`
	// Port number (if applied to a single port).
	Port pulumi.Float64PtrOutput `pulumi:"port"`
	// Special-purpose targets such as "selfSecurityGroup". [See docs](https://cloud.yandex.com/docs/vpc/concepts/security-groups) for possible options.
	PredefinedTarget pulumi.StringPtrOutput `pulumi:"predefinedTarget"`
	// One of `ANY`, `TCP`, `UDP`, `ICMP`, `IPV6_ICMP`.
	Protocol pulumi.StringPtrOutput `pulumi:"protocol"`
	// ID of the security group this rule belongs to.
	SecurityGroupBinding pulumi.StringOutput `pulumi:"securityGroupBinding"`
	// Target security group ID for this rule.
	SecurityGroupId pulumi.StringPtrOutput             `pulumi:"securityGroupId"`
	Timeouts        VpcSecurityGroupRuleTimeoutsOutput `pulumi:"timeouts"`
	// Maximum port number.
	ToPort pulumi.Float64PtrOutput `pulumi:"toPort"`
	// The blocks of IPv4 addresses for this rule.
	V4CidrBlocks pulumi.StringArrayOutput `pulumi:"v4CidrBlocks"`
	// The blocks of IPv6 addresses for this rule. `v6CidrBlocks` argument is currently not supported. It will be available in the future.
	V6CidrBlocks pulumi.StringArrayOutput `pulumi:"v6CidrBlocks"`
}

// NewVpcSecurityGroupRule registers a new resource with the given unique name, arguments, and options.
func NewVpcSecurityGroupRule(ctx *pulumi.Context,
	name string, args *VpcSecurityGroupRuleArgs, opts ...pulumi.ResourceOption) (*VpcSecurityGroupRule, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Direction == nil {
		return nil, errors.New("invalid value for required argument 'Direction'")
	}
	if args.SecurityGroupBinding == nil {
		return nil, errors.New("invalid value for required argument 'SecurityGroupBinding'")
	}
	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	var resource VpcSecurityGroupRule
	err := ctx.RegisterResource("yandex:index/vpcSecurityGroupRule:VpcSecurityGroupRule", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetVpcSecurityGroupRule gets an existing VpcSecurityGroupRule resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetVpcSecurityGroupRule(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *VpcSecurityGroupRuleState, opts ...pulumi.ResourceOption) (*VpcSecurityGroupRule, error) {
	var resource VpcSecurityGroupRule
	err := ctx.ReadResource("yandex:index/vpcSecurityGroupRule:VpcSecurityGroupRule", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering VpcSecurityGroupRule resources.
type vpcSecurityGroupRuleState struct {
	// Description of the rule.
	Description *string `pulumi:"description"`
	// direction of the rule. Can be `ingress` (inbound) or `egress` (outbound).
	Direction *string `pulumi:"direction"`
	// Minimum port number.
	FromPort *float64 `pulumi:"fromPort"`
	// Labels to assign to this rule.
	Labels map[string]string `pulumi:"labels"`
	// Port number (if applied to a single port).
	Port *float64 `pulumi:"port"`
	// Special-purpose targets such as "selfSecurityGroup". [See docs](https://cloud.yandex.com/docs/vpc/concepts/security-groups) for possible options.
	PredefinedTarget *string `pulumi:"predefinedTarget"`
	// One of `ANY`, `TCP`, `UDP`, `ICMP`, `IPV6_ICMP`.
	Protocol *string `pulumi:"protocol"`
	// ID of the security group this rule belongs to.
	SecurityGroupBinding *string `pulumi:"securityGroupBinding"`
	// Target security group ID for this rule.
	SecurityGroupId *string                       `pulumi:"securityGroupId"`
	Timeouts        *VpcSecurityGroupRuleTimeouts `pulumi:"timeouts"`
	// Maximum port number.
	ToPort *float64 `pulumi:"toPort"`
	// The blocks of IPv4 addresses for this rule.
	V4CidrBlocks []string `pulumi:"v4CidrBlocks"`
	// The blocks of IPv6 addresses for this rule. `v6CidrBlocks` argument is currently not supported. It will be available in the future.
	V6CidrBlocks []string `pulumi:"v6CidrBlocks"`
}

type VpcSecurityGroupRuleState struct {
	// Description of the rule.
	Description pulumi.StringPtrInput
	// direction of the rule. Can be `ingress` (inbound) or `egress` (outbound).
	Direction pulumi.StringPtrInput
	// Minimum port number.
	FromPort pulumi.Float64PtrInput
	// Labels to assign to this rule.
	Labels pulumi.StringMapInput
	// Port number (if applied to a single port).
	Port pulumi.Float64PtrInput
	// Special-purpose targets such as "selfSecurityGroup". [See docs](https://cloud.yandex.com/docs/vpc/concepts/security-groups) for possible options.
	PredefinedTarget pulumi.StringPtrInput
	// One of `ANY`, `TCP`, `UDP`, `ICMP`, `IPV6_ICMP`.
	Protocol pulumi.StringPtrInput
	// ID of the security group this rule belongs to.
	SecurityGroupBinding pulumi.StringPtrInput
	// Target security group ID for this rule.
	SecurityGroupId pulumi.StringPtrInput
	Timeouts        VpcSecurityGroupRuleTimeoutsPtrInput
	// Maximum port number.
	ToPort pulumi.Float64PtrInput
	// The blocks of IPv4 addresses for this rule.
	V4CidrBlocks pulumi.StringArrayInput
	// The blocks of IPv6 addresses for this rule. `v6CidrBlocks` argument is currently not supported. It will be available in the future.
	V6CidrBlocks pulumi.StringArrayInput
}

func (VpcSecurityGroupRuleState) ElementType() reflect.Type {
	return reflect.TypeOf((*vpcSecurityGroupRuleState)(nil)).Elem()
}

type vpcSecurityGroupRuleArgs struct {
	// Description of the rule.
	Description *string `pulumi:"description"`
	// direction of the rule. Can be `ingress` (inbound) or `egress` (outbound).
	Direction string `pulumi:"direction"`
	// Minimum port number.
	FromPort *float64 `pulumi:"fromPort"`
	// Labels to assign to this rule.
	Labels map[string]string `pulumi:"labels"`
	// Port number (if applied to a single port).
	Port *float64 `pulumi:"port"`
	// Special-purpose targets such as "selfSecurityGroup". [See docs](https://cloud.yandex.com/docs/vpc/concepts/security-groups) for possible options.
	PredefinedTarget *string `pulumi:"predefinedTarget"`
	// One of `ANY`, `TCP`, `UDP`, `ICMP`, `IPV6_ICMP`.
	Protocol *string `pulumi:"protocol"`
	// ID of the security group this rule belongs to.
	SecurityGroupBinding string `pulumi:"securityGroupBinding"`
	// Target security group ID for this rule.
	SecurityGroupId *string                      `pulumi:"securityGroupId"`
	Timeouts        VpcSecurityGroupRuleTimeouts `pulumi:"timeouts"`
	// Maximum port number.
	ToPort *float64 `pulumi:"toPort"`
	// The blocks of IPv4 addresses for this rule.
	V4CidrBlocks []string `pulumi:"v4CidrBlocks"`
	// The blocks of IPv6 addresses for this rule. `v6CidrBlocks` argument is currently not supported. It will be available in the future.
	V6CidrBlocks []string `pulumi:"v6CidrBlocks"`
}

// The set of arguments for constructing a VpcSecurityGroupRule resource.
type VpcSecurityGroupRuleArgs struct {
	// Description of the rule.
	Description pulumi.StringPtrInput
	// direction of the rule. Can be `ingress` (inbound) or `egress` (outbound).
	Direction pulumi.StringInput
	// Minimum port number.
	FromPort pulumi.Float64PtrInput
	// Labels to assign to this rule.
	Labels pulumi.StringMapInput
	// Port number (if applied to a single port).
	Port pulumi.Float64PtrInput
	// Special-purpose targets such as "selfSecurityGroup". [See docs](https://cloud.yandex.com/docs/vpc/concepts/security-groups) for possible options.
	PredefinedTarget pulumi.StringPtrInput
	// One of `ANY`, `TCP`, `UDP`, `ICMP`, `IPV6_ICMP`.
	Protocol pulumi.StringPtrInput
	// ID of the security group this rule belongs to.
	SecurityGroupBinding pulumi.StringInput
	// Target security group ID for this rule.
	SecurityGroupId pulumi.StringPtrInput
	Timeouts        VpcSecurityGroupRuleTimeoutsInput
	// Maximum port number.
	ToPort pulumi.Float64PtrInput
	// The blocks of IPv4 addresses for this rule.
	V4CidrBlocks pulumi.StringArrayInput
	// The blocks of IPv6 addresses for this rule. `v6CidrBlocks` argument is currently not supported. It will be available in the future.
	V6CidrBlocks pulumi.StringArrayInput
}

func (VpcSecurityGroupRuleArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*vpcSecurityGroupRuleArgs)(nil)).Elem()
}

type VpcSecurityGroupRuleInput interface {
	pulumi.Input

	ToVpcSecurityGroupRuleOutput() VpcSecurityGroupRuleOutput
	ToVpcSecurityGroupRuleOutputWithContext(ctx context.Context) VpcSecurityGroupRuleOutput
}

func (*VpcSecurityGroupRule) ElementType() reflect.Type {
	return reflect.TypeOf((**VpcSecurityGroupRule)(nil)).Elem()
}

func (i *VpcSecurityGroupRule) ToVpcSecurityGroupRuleOutput() VpcSecurityGroupRuleOutput {
	return i.ToVpcSecurityGroupRuleOutputWithContext(context.Background())
}

func (i *VpcSecurityGroupRule) ToVpcSecurityGroupRuleOutputWithContext(ctx context.Context) VpcSecurityGroupRuleOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VpcSecurityGroupRuleOutput)
}

// VpcSecurityGroupRuleArrayInput is an input type that accepts VpcSecurityGroupRuleArray and VpcSecurityGroupRuleArrayOutput values.
// You can construct a concrete instance of `VpcSecurityGroupRuleArrayInput` via:
//
//          VpcSecurityGroupRuleArray{ VpcSecurityGroupRuleArgs{...} }
type VpcSecurityGroupRuleArrayInput interface {
	pulumi.Input

	ToVpcSecurityGroupRuleArrayOutput() VpcSecurityGroupRuleArrayOutput
	ToVpcSecurityGroupRuleArrayOutputWithContext(context.Context) VpcSecurityGroupRuleArrayOutput
}

type VpcSecurityGroupRuleArray []VpcSecurityGroupRuleInput

func (VpcSecurityGroupRuleArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*VpcSecurityGroupRule)(nil)).Elem()
}

func (i VpcSecurityGroupRuleArray) ToVpcSecurityGroupRuleArrayOutput() VpcSecurityGroupRuleArrayOutput {
	return i.ToVpcSecurityGroupRuleArrayOutputWithContext(context.Background())
}

func (i VpcSecurityGroupRuleArray) ToVpcSecurityGroupRuleArrayOutputWithContext(ctx context.Context) VpcSecurityGroupRuleArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VpcSecurityGroupRuleArrayOutput)
}

// VpcSecurityGroupRuleMapInput is an input type that accepts VpcSecurityGroupRuleMap and VpcSecurityGroupRuleMapOutput values.
// You can construct a concrete instance of `VpcSecurityGroupRuleMapInput` via:
//
//          VpcSecurityGroupRuleMap{ "key": VpcSecurityGroupRuleArgs{...} }
type VpcSecurityGroupRuleMapInput interface {
	pulumi.Input

	ToVpcSecurityGroupRuleMapOutput() VpcSecurityGroupRuleMapOutput
	ToVpcSecurityGroupRuleMapOutputWithContext(context.Context) VpcSecurityGroupRuleMapOutput
}

type VpcSecurityGroupRuleMap map[string]VpcSecurityGroupRuleInput

func (VpcSecurityGroupRuleMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*VpcSecurityGroupRule)(nil)).Elem()
}

func (i VpcSecurityGroupRuleMap) ToVpcSecurityGroupRuleMapOutput() VpcSecurityGroupRuleMapOutput {
	return i.ToVpcSecurityGroupRuleMapOutputWithContext(context.Background())
}

func (i VpcSecurityGroupRuleMap) ToVpcSecurityGroupRuleMapOutputWithContext(ctx context.Context) VpcSecurityGroupRuleMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(VpcSecurityGroupRuleMapOutput)
}

type VpcSecurityGroupRuleOutput struct{ *pulumi.OutputState }

func (VpcSecurityGroupRuleOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**VpcSecurityGroupRule)(nil)).Elem()
}

func (o VpcSecurityGroupRuleOutput) ToVpcSecurityGroupRuleOutput() VpcSecurityGroupRuleOutput {
	return o
}

func (o VpcSecurityGroupRuleOutput) ToVpcSecurityGroupRuleOutputWithContext(ctx context.Context) VpcSecurityGroupRuleOutput {
	return o
}

// Description of the rule.
func (o VpcSecurityGroupRuleOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// direction of the rule. Can be `ingress` (inbound) or `egress` (outbound).
func (o VpcSecurityGroupRuleOutput) Direction() pulumi.StringOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.StringOutput { return v.Direction }).(pulumi.StringOutput)
}

// Minimum port number.
func (o VpcSecurityGroupRuleOutput) FromPort() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.Float64PtrOutput { return v.FromPort }).(pulumi.Float64PtrOutput)
}

// Labels to assign to this rule.
func (o VpcSecurityGroupRuleOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.StringMapOutput { return v.Labels }).(pulumi.StringMapOutput)
}

// Port number (if applied to a single port).
func (o VpcSecurityGroupRuleOutput) Port() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.Float64PtrOutput { return v.Port }).(pulumi.Float64PtrOutput)
}

// Special-purpose targets such as "selfSecurityGroup". [See docs](https://cloud.yandex.com/docs/vpc/concepts/security-groups) for possible options.
func (o VpcSecurityGroupRuleOutput) PredefinedTarget() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.StringPtrOutput { return v.PredefinedTarget }).(pulumi.StringPtrOutput)
}

// One of `ANY`, `TCP`, `UDP`, `ICMP`, `IPV6_ICMP`.
func (o VpcSecurityGroupRuleOutput) Protocol() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.StringPtrOutput { return v.Protocol }).(pulumi.StringPtrOutput)
}

// ID of the security group this rule belongs to.
func (o VpcSecurityGroupRuleOutput) SecurityGroupBinding() pulumi.StringOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.StringOutput { return v.SecurityGroupBinding }).(pulumi.StringOutput)
}

// Target security group ID for this rule.
func (o VpcSecurityGroupRuleOutput) SecurityGroupId() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.StringPtrOutput { return v.SecurityGroupId }).(pulumi.StringPtrOutput)
}

func (o VpcSecurityGroupRuleOutput) Timeouts() VpcSecurityGroupRuleTimeoutsOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) VpcSecurityGroupRuleTimeoutsOutput { return v.Timeouts }).(VpcSecurityGroupRuleTimeoutsOutput)
}

// Maximum port number.
func (o VpcSecurityGroupRuleOutput) ToPort() pulumi.Float64PtrOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.Float64PtrOutput { return v.ToPort }).(pulumi.Float64PtrOutput)
}

// The blocks of IPv4 addresses for this rule.
func (o VpcSecurityGroupRuleOutput) V4CidrBlocks() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.StringArrayOutput { return v.V4CidrBlocks }).(pulumi.StringArrayOutput)
}

// The blocks of IPv6 addresses for this rule. `v6CidrBlocks` argument is currently not supported. It will be available in the future.
func (o VpcSecurityGroupRuleOutput) V6CidrBlocks() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *VpcSecurityGroupRule) pulumi.StringArrayOutput { return v.V6CidrBlocks }).(pulumi.StringArrayOutput)
}

type VpcSecurityGroupRuleArrayOutput struct{ *pulumi.OutputState }

func (VpcSecurityGroupRuleArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*VpcSecurityGroupRule)(nil)).Elem()
}

func (o VpcSecurityGroupRuleArrayOutput) ToVpcSecurityGroupRuleArrayOutput() VpcSecurityGroupRuleArrayOutput {
	return o
}

func (o VpcSecurityGroupRuleArrayOutput) ToVpcSecurityGroupRuleArrayOutputWithContext(ctx context.Context) VpcSecurityGroupRuleArrayOutput {
	return o
}

func (o VpcSecurityGroupRuleArrayOutput) Index(i pulumi.IntInput) VpcSecurityGroupRuleOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *VpcSecurityGroupRule {
		return vs[0].([]*VpcSecurityGroupRule)[vs[1].(int)]
	}).(VpcSecurityGroupRuleOutput)
}

type VpcSecurityGroupRuleMapOutput struct{ *pulumi.OutputState }

func (VpcSecurityGroupRuleMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*VpcSecurityGroupRule)(nil)).Elem()
}

func (o VpcSecurityGroupRuleMapOutput) ToVpcSecurityGroupRuleMapOutput() VpcSecurityGroupRuleMapOutput {
	return o
}

func (o VpcSecurityGroupRuleMapOutput) ToVpcSecurityGroupRuleMapOutputWithContext(ctx context.Context) VpcSecurityGroupRuleMapOutput {
	return o
}

func (o VpcSecurityGroupRuleMapOutput) MapIndex(k pulumi.StringInput) VpcSecurityGroupRuleOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *VpcSecurityGroupRule {
		return vs[0].(map[string]*VpcSecurityGroupRule)[vs[1].(string)]
	}).(VpcSecurityGroupRuleOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*VpcSecurityGroupRuleInput)(nil)).Elem(), &VpcSecurityGroupRule{})
	pulumi.RegisterInputType(reflect.TypeOf((*VpcSecurityGroupRuleArrayInput)(nil)).Elem(), VpcSecurityGroupRuleArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*VpcSecurityGroupRuleMapInput)(nil)).Elem(), VpcSecurityGroupRuleMap{})
	pulumi.RegisterOutputType(VpcSecurityGroupRuleOutput{})
	pulumi.RegisterOutputType(VpcSecurityGroupRuleArrayOutput{})
	pulumi.RegisterOutputType(VpcSecurityGroupRuleMapOutput{})
}
