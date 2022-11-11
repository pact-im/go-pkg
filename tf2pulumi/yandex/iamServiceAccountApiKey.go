// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Allows management of a [Yandex.Cloud IAM service account API key](https://cloud.yandex.com/docs/iam/concepts/authorization/api-key).
// The API key is a private key used for simplified authorization in the Yandex.Cloud API. API keys are only used for [service accounts](https://cloud.yandex.com/docs/iam/concepts/users/service-accounts).
//
// API keys do not expire. This means that this authentication method is simpler, but less secure. Use it if you can't automatically request an [IAM token](https://cloud.yandex.com/docs/iam/concepts/authorization/iam-token).
//
// ## Example Usage
//
// This snippet creates an API key.
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
// 		_, err := yandex.NewIamServiceAccountApiKey(ctx, "sa-api-key", &yandex.IamServiceAccountApiKeyArgs{
// 			Description:      pulumi.String("api key for authorization"),
// 			PgpKey:           pulumi.String("keybase:keybaseusername"),
// 			ServiceAccountId: pulumi.String("some_sa_id"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		return nil
// 	})
// }
// ```
type IamServiceAccountApiKey struct {
	pulumi.CustomResourceState

	// Creation timestamp of the static access key.
	CreatedAt pulumi.StringOutput `pulumi:"createdAt"`
	// The description of the key.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// The encrypted secret key, base64 encoded. This is only populated when `pgpKey` is supplied.
	EncryptedSecretKey pulumi.StringOutput `pulumi:"encryptedSecretKey"`
	// The fingerprint of the PGP key used to encrypt the secret key. This is only populated when `pgpKey` is supplied.
	KeyFingerprint pulumi.StringOutput `pulumi:"keyFingerprint"`
	// An optional PGP key to encrypt the resulting secret key material. May either be a base64-encoded public key or a keybase username in the form `keybase:keybaseusername`.
	PgpKey pulumi.StringPtrOutput `pulumi:"pgpKey"`
	// The secret key. This is only populated when no `pgpKey` is provided.
	SecretKey pulumi.StringOutput `pulumi:"secretKey"`
	// ID of the service account to an API key for.
	ServiceAccountId pulumi.StringOutput `pulumi:"serviceAccountId"`
}

// NewIamServiceAccountApiKey registers a new resource with the given unique name, arguments, and options.
func NewIamServiceAccountApiKey(ctx *pulumi.Context,
	name string, args *IamServiceAccountApiKeyArgs, opts ...pulumi.ResourceOption) (*IamServiceAccountApiKey, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.ServiceAccountId == nil {
		return nil, errors.New("invalid value for required argument 'ServiceAccountId'")
	}
	secrets := pulumi.AdditionalSecretOutputs([]string{
		"secretKey",
	})
	opts = append(opts, secrets)
	var resource IamServiceAccountApiKey
	err := ctx.RegisterResource("yandex:index/iamServiceAccountApiKey:IamServiceAccountApiKey", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetIamServiceAccountApiKey gets an existing IamServiceAccountApiKey resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetIamServiceAccountApiKey(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *IamServiceAccountApiKeyState, opts ...pulumi.ResourceOption) (*IamServiceAccountApiKey, error) {
	var resource IamServiceAccountApiKey
	err := ctx.ReadResource("yandex:index/iamServiceAccountApiKey:IamServiceAccountApiKey", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering IamServiceAccountApiKey resources.
type iamServiceAccountApiKeyState struct {
	// Creation timestamp of the static access key.
	CreatedAt *string `pulumi:"createdAt"`
	// The description of the key.
	Description *string `pulumi:"description"`
	// The encrypted secret key, base64 encoded. This is only populated when `pgpKey` is supplied.
	EncryptedSecretKey *string `pulumi:"encryptedSecretKey"`
	// The fingerprint of the PGP key used to encrypt the secret key. This is only populated when `pgpKey` is supplied.
	KeyFingerprint *string `pulumi:"keyFingerprint"`
	// An optional PGP key to encrypt the resulting secret key material. May either be a base64-encoded public key or a keybase username in the form `keybase:keybaseusername`.
	PgpKey *string `pulumi:"pgpKey"`
	// The secret key. This is only populated when no `pgpKey` is provided.
	SecretKey *string `pulumi:"secretKey"`
	// ID of the service account to an API key for.
	ServiceAccountId *string `pulumi:"serviceAccountId"`
}

type IamServiceAccountApiKeyState struct {
	// Creation timestamp of the static access key.
	CreatedAt pulumi.StringPtrInput
	// The description of the key.
	Description pulumi.StringPtrInput
	// The encrypted secret key, base64 encoded. This is only populated when `pgpKey` is supplied.
	EncryptedSecretKey pulumi.StringPtrInput
	// The fingerprint of the PGP key used to encrypt the secret key. This is only populated when `pgpKey` is supplied.
	KeyFingerprint pulumi.StringPtrInput
	// An optional PGP key to encrypt the resulting secret key material. May either be a base64-encoded public key or a keybase username in the form `keybase:keybaseusername`.
	PgpKey pulumi.StringPtrInput
	// The secret key. This is only populated when no `pgpKey` is provided.
	SecretKey pulumi.StringPtrInput
	// ID of the service account to an API key for.
	ServiceAccountId pulumi.StringPtrInput
}

func (IamServiceAccountApiKeyState) ElementType() reflect.Type {
	return reflect.TypeOf((*iamServiceAccountApiKeyState)(nil)).Elem()
}

type iamServiceAccountApiKeyArgs struct {
	// The description of the key.
	Description *string `pulumi:"description"`
	// An optional PGP key to encrypt the resulting secret key material. May either be a base64-encoded public key or a keybase username in the form `keybase:keybaseusername`.
	PgpKey *string `pulumi:"pgpKey"`
	// ID of the service account to an API key for.
	ServiceAccountId string `pulumi:"serviceAccountId"`
}

// The set of arguments for constructing a IamServiceAccountApiKey resource.
type IamServiceAccountApiKeyArgs struct {
	// The description of the key.
	Description pulumi.StringPtrInput
	// An optional PGP key to encrypt the resulting secret key material. May either be a base64-encoded public key or a keybase username in the form `keybase:keybaseusername`.
	PgpKey pulumi.StringPtrInput
	// ID of the service account to an API key for.
	ServiceAccountId pulumi.StringInput
}

func (IamServiceAccountApiKeyArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*iamServiceAccountApiKeyArgs)(nil)).Elem()
}

type IamServiceAccountApiKeyInput interface {
	pulumi.Input

	ToIamServiceAccountApiKeyOutput() IamServiceAccountApiKeyOutput
	ToIamServiceAccountApiKeyOutputWithContext(ctx context.Context) IamServiceAccountApiKeyOutput
}

func (*IamServiceAccountApiKey) ElementType() reflect.Type {
	return reflect.TypeOf((**IamServiceAccountApiKey)(nil)).Elem()
}

func (i *IamServiceAccountApiKey) ToIamServiceAccountApiKeyOutput() IamServiceAccountApiKeyOutput {
	return i.ToIamServiceAccountApiKeyOutputWithContext(context.Background())
}

func (i *IamServiceAccountApiKey) ToIamServiceAccountApiKeyOutputWithContext(ctx context.Context) IamServiceAccountApiKeyOutput {
	return pulumi.ToOutputWithContext(ctx, i).(IamServiceAccountApiKeyOutput)
}

// IamServiceAccountApiKeyArrayInput is an input type that accepts IamServiceAccountApiKeyArray and IamServiceAccountApiKeyArrayOutput values.
// You can construct a concrete instance of `IamServiceAccountApiKeyArrayInput` via:
//
//          IamServiceAccountApiKeyArray{ IamServiceAccountApiKeyArgs{...} }
type IamServiceAccountApiKeyArrayInput interface {
	pulumi.Input

	ToIamServiceAccountApiKeyArrayOutput() IamServiceAccountApiKeyArrayOutput
	ToIamServiceAccountApiKeyArrayOutputWithContext(context.Context) IamServiceAccountApiKeyArrayOutput
}

type IamServiceAccountApiKeyArray []IamServiceAccountApiKeyInput

func (IamServiceAccountApiKeyArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*IamServiceAccountApiKey)(nil)).Elem()
}

func (i IamServiceAccountApiKeyArray) ToIamServiceAccountApiKeyArrayOutput() IamServiceAccountApiKeyArrayOutput {
	return i.ToIamServiceAccountApiKeyArrayOutputWithContext(context.Background())
}

func (i IamServiceAccountApiKeyArray) ToIamServiceAccountApiKeyArrayOutputWithContext(ctx context.Context) IamServiceAccountApiKeyArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(IamServiceAccountApiKeyArrayOutput)
}

// IamServiceAccountApiKeyMapInput is an input type that accepts IamServiceAccountApiKeyMap and IamServiceAccountApiKeyMapOutput values.
// You can construct a concrete instance of `IamServiceAccountApiKeyMapInput` via:
//
//          IamServiceAccountApiKeyMap{ "key": IamServiceAccountApiKeyArgs{...} }
type IamServiceAccountApiKeyMapInput interface {
	pulumi.Input

	ToIamServiceAccountApiKeyMapOutput() IamServiceAccountApiKeyMapOutput
	ToIamServiceAccountApiKeyMapOutputWithContext(context.Context) IamServiceAccountApiKeyMapOutput
}

type IamServiceAccountApiKeyMap map[string]IamServiceAccountApiKeyInput

func (IamServiceAccountApiKeyMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*IamServiceAccountApiKey)(nil)).Elem()
}

func (i IamServiceAccountApiKeyMap) ToIamServiceAccountApiKeyMapOutput() IamServiceAccountApiKeyMapOutput {
	return i.ToIamServiceAccountApiKeyMapOutputWithContext(context.Background())
}

func (i IamServiceAccountApiKeyMap) ToIamServiceAccountApiKeyMapOutputWithContext(ctx context.Context) IamServiceAccountApiKeyMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(IamServiceAccountApiKeyMapOutput)
}

type IamServiceAccountApiKeyOutput struct{ *pulumi.OutputState }

func (IamServiceAccountApiKeyOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**IamServiceAccountApiKey)(nil)).Elem()
}

func (o IamServiceAccountApiKeyOutput) ToIamServiceAccountApiKeyOutput() IamServiceAccountApiKeyOutput {
	return o
}

func (o IamServiceAccountApiKeyOutput) ToIamServiceAccountApiKeyOutputWithContext(ctx context.Context) IamServiceAccountApiKeyOutput {
	return o
}

// Creation timestamp of the static access key.
func (o IamServiceAccountApiKeyOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v *IamServiceAccountApiKey) pulumi.StringOutput { return v.CreatedAt }).(pulumi.StringOutput)
}

// The description of the key.
func (o IamServiceAccountApiKeyOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *IamServiceAccountApiKey) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// The encrypted secret key, base64 encoded. This is only populated when `pgpKey` is supplied.
func (o IamServiceAccountApiKeyOutput) EncryptedSecretKey() pulumi.StringOutput {
	return o.ApplyT(func(v *IamServiceAccountApiKey) pulumi.StringOutput { return v.EncryptedSecretKey }).(pulumi.StringOutput)
}

// The fingerprint of the PGP key used to encrypt the secret key. This is only populated when `pgpKey` is supplied.
func (o IamServiceAccountApiKeyOutput) KeyFingerprint() pulumi.StringOutput {
	return o.ApplyT(func(v *IamServiceAccountApiKey) pulumi.StringOutput { return v.KeyFingerprint }).(pulumi.StringOutput)
}

// An optional PGP key to encrypt the resulting secret key material. May either be a base64-encoded public key or a keybase username in the form `keybase:keybaseusername`.
func (o IamServiceAccountApiKeyOutput) PgpKey() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *IamServiceAccountApiKey) pulumi.StringPtrOutput { return v.PgpKey }).(pulumi.StringPtrOutput)
}

// The secret key. This is only populated when no `pgpKey` is provided.
func (o IamServiceAccountApiKeyOutput) SecretKey() pulumi.StringOutput {
	return o.ApplyT(func(v *IamServiceAccountApiKey) pulumi.StringOutput { return v.SecretKey }).(pulumi.StringOutput)
}

// ID of the service account to an API key for.
func (o IamServiceAccountApiKeyOutput) ServiceAccountId() pulumi.StringOutput {
	return o.ApplyT(func(v *IamServiceAccountApiKey) pulumi.StringOutput { return v.ServiceAccountId }).(pulumi.StringOutput)
}

type IamServiceAccountApiKeyArrayOutput struct{ *pulumi.OutputState }

func (IamServiceAccountApiKeyArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*IamServiceAccountApiKey)(nil)).Elem()
}

func (o IamServiceAccountApiKeyArrayOutput) ToIamServiceAccountApiKeyArrayOutput() IamServiceAccountApiKeyArrayOutput {
	return o
}

func (o IamServiceAccountApiKeyArrayOutput) ToIamServiceAccountApiKeyArrayOutputWithContext(ctx context.Context) IamServiceAccountApiKeyArrayOutput {
	return o
}

func (o IamServiceAccountApiKeyArrayOutput) Index(i pulumi.IntInput) IamServiceAccountApiKeyOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *IamServiceAccountApiKey {
		return vs[0].([]*IamServiceAccountApiKey)[vs[1].(int)]
	}).(IamServiceAccountApiKeyOutput)
}

type IamServiceAccountApiKeyMapOutput struct{ *pulumi.OutputState }

func (IamServiceAccountApiKeyMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*IamServiceAccountApiKey)(nil)).Elem()
}

func (o IamServiceAccountApiKeyMapOutput) ToIamServiceAccountApiKeyMapOutput() IamServiceAccountApiKeyMapOutput {
	return o
}

func (o IamServiceAccountApiKeyMapOutput) ToIamServiceAccountApiKeyMapOutputWithContext(ctx context.Context) IamServiceAccountApiKeyMapOutput {
	return o
}

func (o IamServiceAccountApiKeyMapOutput) MapIndex(k pulumi.StringInput) IamServiceAccountApiKeyOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *IamServiceAccountApiKey {
		return vs[0].(map[string]*IamServiceAccountApiKey)[vs[1].(string)]
	}).(IamServiceAccountApiKeyOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*IamServiceAccountApiKeyInput)(nil)).Elem(), &IamServiceAccountApiKey{})
	pulumi.RegisterInputType(reflect.TypeOf((*IamServiceAccountApiKeyArrayInput)(nil)).Elem(), IamServiceAccountApiKeyArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*IamServiceAccountApiKeyMapInput)(nil)).Elem(), IamServiceAccountApiKeyMap{})
	pulumi.RegisterOutputType(IamServiceAccountApiKeyOutput{})
	pulumi.RegisterOutputType(IamServiceAccountApiKeyArrayOutput{})
	pulumi.RegisterOutputType(IamServiceAccountApiKeyMapOutput{})
}
