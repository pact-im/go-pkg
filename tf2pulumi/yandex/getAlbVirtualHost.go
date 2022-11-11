// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Get information about a Yandex ALB Virtual Host. For more information, see
// [Yandex.Cloud Application Load Balancer](https://cloud.yandex.com/en/docs/application-load-balancer/quickstart).
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
// 		_, err = yandex.LookupAlbVirtualHost(ctx, &GetAlbVirtualHostArgs{
// 			Name:         pulumi.StringRef(yandex_alb_virtual_host.MyVh.Name),
// 			HttpRouterId: pulumi.StringRef(yandex_alb_virtual_host.MyRouter.Id),
// 		}, nil)
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
//
// This data source is used to define [Application Load Balancer Virtual Host] that can be used by other resources.
func LookupAlbVirtualHost(ctx *pulumi.Context, args *LookupAlbVirtualHostArgs, opts ...pulumi.InvokeOption) (*LookupAlbVirtualHostResult, error) {
	var rv LookupAlbVirtualHostResult
	err := ctx.Invoke("yandex:index/getAlbVirtualHost:getAlbVirtualHost", args, &rv, opts...)
	if err != nil {
		return nil, err
	}
	return &rv, nil
}

// A collection of arguments for invoking getAlbVirtualHost.
type LookupAlbVirtualHostArgs struct {
	// HTTP Router that the resource belongs to.
	HttpRouterId *string `pulumi:"httpRouterId"`
	// Name of the Virtual Host.
	Name *string `pulumi:"name"`
	// The ID of a specific Virtual Host. Virtual Host ID is concatenation of HTTP Router ID
	// and Virtual Host name with `/` symbol between them. For Example, "http_router_id/vhost_name".
	VirtualHostId *string `pulumi:"virtualHostId"`
}

// A collection of values returned by getAlbVirtualHost.
type LookupAlbVirtualHostResult struct {
	// A list of domains (host/authority header) that will be matched to this virtual host. Wildcard hosts are
	// supported in the form of '*.foo.com' or '*-bar.foo.com'. If not specified, all domains will be matched.
	Authorities  []string `pulumi:"authorities"`
	HttpRouterId string   `pulumi:"httpRouterId"`
	Id           string   `pulumi:"id"`
	// Apply the following modifications to the request headers. The structure is documented
	// below.
	ModifyRequestHeaders []GetAlbVirtualHostModifyRequestHeader `pulumi:"modifyRequestHeaders"`
	// Apply the following modifications to the response headers. The structure is documented
	// below.
	ModifyResponseHeaders []GetAlbVirtualHostModifyResponseHeader `pulumi:"modifyResponseHeaders"`
	// name of the route.
	Name         string                         `pulumi:"name"`
	RouteOptions []GetAlbVirtualHostRouteOption `pulumi:"routeOptions"`
	// A Route resource. Routes are matched *in-order*. Be careful when adding them to the end. For instance,
	// having http '/' match first makes all other routes unused. The structure is documented below.
	Routes        []GetAlbVirtualHostRoute `pulumi:"routes"`
	VirtualHostId string                   `pulumi:"virtualHostId"`
}

func LookupAlbVirtualHostOutput(ctx *pulumi.Context, args LookupAlbVirtualHostOutputArgs, opts ...pulumi.InvokeOption) LookupAlbVirtualHostResultOutput {
	return pulumi.ToOutputWithContext(context.Background(), args).
		ApplyT(func(v interface{}) (LookupAlbVirtualHostResult, error) {
			args := v.(LookupAlbVirtualHostArgs)
			r, err := LookupAlbVirtualHost(ctx, &args, opts...)
			var s LookupAlbVirtualHostResult
			if r != nil {
				s = *r
			}
			return s, err
		}).(LookupAlbVirtualHostResultOutput)
}

// A collection of arguments for invoking getAlbVirtualHost.
type LookupAlbVirtualHostOutputArgs struct {
	// HTTP Router that the resource belongs to.
	HttpRouterId pulumi.StringPtrInput `pulumi:"httpRouterId"`
	// Name of the Virtual Host.
	Name pulumi.StringPtrInput `pulumi:"name"`
	// The ID of a specific Virtual Host. Virtual Host ID is concatenation of HTTP Router ID
	// and Virtual Host name with `/` symbol between them. For Example, "http_router_id/vhost_name".
	VirtualHostId pulumi.StringPtrInput `pulumi:"virtualHostId"`
}

func (LookupAlbVirtualHostOutputArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupAlbVirtualHostArgs)(nil)).Elem()
}

// A collection of values returned by getAlbVirtualHost.
type LookupAlbVirtualHostResultOutput struct{ *pulumi.OutputState }

func (LookupAlbVirtualHostResultOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*LookupAlbVirtualHostResult)(nil)).Elem()
}

func (o LookupAlbVirtualHostResultOutput) ToLookupAlbVirtualHostResultOutput() LookupAlbVirtualHostResultOutput {
	return o
}

func (o LookupAlbVirtualHostResultOutput) ToLookupAlbVirtualHostResultOutputWithContext(ctx context.Context) LookupAlbVirtualHostResultOutput {
	return o
}

// A list of domains (host/authority header) that will be matched to this virtual host. Wildcard hosts are
// supported in the form of '*.foo.com' or '*-bar.foo.com'. If not specified, all domains will be matched.
func (o LookupAlbVirtualHostResultOutput) Authorities() pulumi.StringArrayOutput {
	return o.ApplyT(func(v LookupAlbVirtualHostResult) []string { return v.Authorities }).(pulumi.StringArrayOutput)
}

func (o LookupAlbVirtualHostResultOutput) HttpRouterId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupAlbVirtualHostResult) string { return v.HttpRouterId }).(pulumi.StringOutput)
}

func (o LookupAlbVirtualHostResultOutput) Id() pulumi.StringOutput {
	return o.ApplyT(func(v LookupAlbVirtualHostResult) string { return v.Id }).(pulumi.StringOutput)
}

// Apply the following modifications to the request headers. The structure is documented
// below.
func (o LookupAlbVirtualHostResultOutput) ModifyRequestHeaders() GetAlbVirtualHostModifyRequestHeaderArrayOutput {
	return o.ApplyT(func(v LookupAlbVirtualHostResult) []GetAlbVirtualHostModifyRequestHeader {
		return v.ModifyRequestHeaders
	}).(GetAlbVirtualHostModifyRequestHeaderArrayOutput)
}

// Apply the following modifications to the response headers. The structure is documented
// below.
func (o LookupAlbVirtualHostResultOutput) ModifyResponseHeaders() GetAlbVirtualHostModifyResponseHeaderArrayOutput {
	return o.ApplyT(func(v LookupAlbVirtualHostResult) []GetAlbVirtualHostModifyResponseHeader {
		return v.ModifyResponseHeaders
	}).(GetAlbVirtualHostModifyResponseHeaderArrayOutput)
}

// name of the route.
func (o LookupAlbVirtualHostResultOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v LookupAlbVirtualHostResult) string { return v.Name }).(pulumi.StringOutput)
}

func (o LookupAlbVirtualHostResultOutput) RouteOptions() GetAlbVirtualHostRouteOptionArrayOutput {
	return o.ApplyT(func(v LookupAlbVirtualHostResult) []GetAlbVirtualHostRouteOption { return v.RouteOptions }).(GetAlbVirtualHostRouteOptionArrayOutput)
}

// A Route resource. Routes are matched *in-order*. Be careful when adding them to the end. For instance,
// having http '/' match first makes all other routes unused. The structure is documented below.
func (o LookupAlbVirtualHostResultOutput) Routes() GetAlbVirtualHostRouteArrayOutput {
	return o.ApplyT(func(v LookupAlbVirtualHostResult) []GetAlbVirtualHostRoute { return v.Routes }).(GetAlbVirtualHostRouteArrayOutput)
}

func (o LookupAlbVirtualHostResultOutput) VirtualHostId() pulumi.StringOutput {
	return o.ApplyT(func(v LookupAlbVirtualHostResult) string { return v.VirtualHostId }).(pulumi.StringOutput)
}

func init() {
	pulumi.RegisterOutputType(LookupAlbVirtualHostResultOutput{})
}
