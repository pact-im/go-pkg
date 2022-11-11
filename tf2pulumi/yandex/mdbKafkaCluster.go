// Code generated by the Pulumi Terraform Bridge (tfgen) Tool DO NOT EDIT.
// *** WARNING: Do not edit by hand unless you're certain you know what you are doing! ***

package yandex

import (
	"context"
	"errors"
	"reflect"

	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Manages a Kafka cluster within the Yandex.Cloud. For more information, see
// [the official documentation](https://cloud.yandex.com/docs/managed-kafka/concepts).
//
// ## Example Usage
//
// Example of creating a Single Node Kafka.
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
// 		fooVpcNetwork, err := yandex.NewVpcNetwork(ctx, "fooVpcNetwork", nil)
// 		if err != nil {
// 			return err
// 		}
// 		fooVpcSubnet, err := yandex.NewVpcSubnet(ctx, "fooVpcSubnet", &yandex.VpcSubnetArgs{
// 			NetworkId: fooVpcNetwork.ID(),
// 			V4CidrBlocks: pulumi.StringArray{
// 				pulumi.String("10.5.0.0/24"),
// 			},
// 			Zone: pulumi.String("ru-central1-a"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		_, err = yandex.NewMdbKafkaCluster(ctx, "fooMdbKafkaCluster", &yandex.MdbKafkaClusterArgs{
// 			Config: &MdbKafkaClusterConfigArgs{
// 				AssignPublicIp: pulumi.Bool(false),
// 				BrokersCount:   pulumi.Float64(1),
// 				Kafka: &MdbKafkaClusterConfigKafkaArgs{
// 					KafkaConfig: &MdbKafkaClusterConfigKafkaKafkaConfigArgs{
// 						CompressionType:             pulumi.String("COMPRESSION_TYPE_ZSTD"),
// 						DefaultReplicationFactor:    pulumi.String("1"),
// 						LogFlushIntervalMessages:    pulumi.String("1024"),
// 						LogFlushIntervalMs:          pulumi.String("1000"),
// 						LogFlushSchedulerIntervalMs: pulumi.String("1000"),
// 						LogPreallocate:              pulumi.Bool(true),
// 						LogRetentionBytes:           pulumi.String("1073741824"),
// 						LogRetentionHours:           pulumi.String("168"),
// 						LogRetentionMinutes:         pulumi.String("10080"),
// 						LogRetentionMs:              pulumi.String("86400000"),
// 						LogSegmentBytes:             pulumi.String("134217728"),
// 						MessageMaxBytes:             pulumi.String("1048588"),
// 						NumPartitions:               pulumi.String("10"),
// 						OffsetsRetentionMinutes:     pulumi.String("10080"),
// 						ReplicaFetchMaxBytes:        pulumi.String("1048576"),
// 						SslCipherSuites: pulumi.StringArray{
// 							pulumi.String("TLS_DHE_RSA_WITH_AES_128_CBC_SHA"),
// 							pulumi.String("TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256"),
// 						},
// 					},
// 					Resources: &MdbKafkaClusterConfigKafkaResourcesArgs{
// 						DiskSize:         pulumi.Float64(32),
// 						DiskTypeId:       pulumi.String("network-ssd"),
// 						ResourcePresetId: pulumi.String("s2.micro"),
// 					},
// 				},
// 				SchemaRegistry:  pulumi.Bool(false),
// 				UnmanagedTopics: pulumi.Bool(false),
// 				Version:         pulumi.String("2.8"),
// 				Zones: pulumi.StringArray{
// 					pulumi.String("ru-central1-a"),
// 				},
// 			},
// 			Environment: pulumi.String("PRESTABLE"),
// 			NetworkId:   fooVpcNetwork.ID(),
// 			SubnetIds: pulumi.StringArray{
// 				fooVpcSubnet.ID(),
// 			},
// 			Users: MdbKafkaClusterUserArray{
// 				&MdbKafkaClusterUserArgs{
// 					Name:     pulumi.String("producer-application"),
// 					Password: pulumi.String("password"),
// 					Permissions: MdbKafkaClusterUserPermissionArray{
// 						&MdbKafkaClusterUserPermissionArgs{
// 							Role:      pulumi.String("ACCESS_ROLE_PRODUCER"),
// 							TopicName: pulumi.String("input"),
// 						},
// 					},
// 				},
// 				&MdbKafkaClusterUserArgs{
// 					Name:     pulumi.String("worker"),
// 					Password: pulumi.String("password"),
// 					Permissions: MdbKafkaClusterUserPermissionArray{
// 						&MdbKafkaClusterUserPermissionArgs{
// 							Role:      pulumi.String("ACCESS_ROLE_CONSUMER"),
// 							TopicName: pulumi.String("input"),
// 						},
// 						&MdbKafkaClusterUserPermissionArgs{
// 							Role:      pulumi.String("ACCESS_ROLE_PRODUCER"),
// 							TopicName: pulumi.String("output"),
// 						},
// 					},
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
// Example of creating a HA Kafka Cluster with two brokers per AZ (6 brokers + 3 zk)
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
// 		fooVpcNetwork, err := yandex.NewVpcNetwork(ctx, "fooVpcNetwork", nil)
// 		if err != nil {
// 			return err
// 		}
// 		fooVpcSubnet, err := yandex.NewVpcSubnet(ctx, "fooVpcSubnet", &yandex.VpcSubnetArgs{
// 			NetworkId: fooVpcNetwork.ID(),
// 			V4CidrBlocks: pulumi.StringArray{
// 				pulumi.String("10.1.0.0/24"),
// 			},
// 			Zone: pulumi.String("ru-central1-a"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		bar, err := yandex.NewVpcSubnet(ctx, "bar", &yandex.VpcSubnetArgs{
// 			NetworkId: fooVpcNetwork.ID(),
// 			V4CidrBlocks: pulumi.StringArray{
// 				pulumi.String("10.2.0.0/24"),
// 			},
// 			Zone: pulumi.String("ru-central1-b"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		baz, err := yandex.NewVpcSubnet(ctx, "baz", &yandex.VpcSubnetArgs{
// 			NetworkId: fooVpcNetwork.ID(),
// 			V4CidrBlocks: pulumi.StringArray{
// 				pulumi.String("10.3.0.0/24"),
// 			},
// 			Zone: pulumi.String("ru-central1-c"),
// 		})
// 		if err != nil {
// 			return err
// 		}
// 		_, err = yandex.NewMdbKafkaCluster(ctx, "fooMdbKafkaCluster", &yandex.MdbKafkaClusterArgs{
// 			Config: &MdbKafkaClusterConfigArgs{
// 				AssignPublicIp: pulumi.Bool(true),
// 				BrokersCount:   pulumi.Float64(2),
// 				Kafka: &MdbKafkaClusterConfigKafkaArgs{
// 					KafkaConfig: &MdbKafkaClusterConfigKafkaKafkaConfigArgs{
// 						CompressionType:             pulumi.String("COMPRESSION_TYPE_ZSTD"),
// 						DefaultReplicationFactor:    pulumi.String("6"),
// 						LogFlushIntervalMessages:    pulumi.String("1024"),
// 						LogFlushIntervalMs:          pulumi.String("1000"),
// 						LogFlushSchedulerIntervalMs: pulumi.String("1000"),
// 						LogPreallocate:              pulumi.Bool(true),
// 						LogRetentionBytes:           pulumi.String("1073741824"),
// 						LogRetentionHours:           pulumi.String("168"),
// 						LogRetentionMinutes:         pulumi.String("10080"),
// 						LogRetentionMs:              pulumi.String("86400000"),
// 						LogSegmentBytes:             pulumi.String("134217728"),
// 						MessageMaxBytes:             pulumi.String("1048588"),
// 						NumPartitions:               pulumi.String("10"),
// 						OffsetsRetentionMinutes:     pulumi.String("10080"),
// 						ReplicaFetchMaxBytes:        pulumi.String("1048576"),
// 						SslCipherSuites: pulumi.StringArray{
// 							pulumi.String("TLS_DHE_RSA_WITH_AES_128_CBC_SHA"),
// 							pulumi.String("TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305_SHA256"),
// 						},
// 					},
// 					Resources: &MdbKafkaClusterConfigKafkaResourcesArgs{
// 						DiskSize:         pulumi.Float64(128),
// 						DiskTypeId:       pulumi.String("network-ssd"),
// 						ResourcePresetId: pulumi.String("s2.medium"),
// 					},
// 				},
// 				SchemaRegistry:  pulumi.Bool(false),
// 				UnmanagedTopics: pulumi.Bool(false),
// 				Version:         pulumi.String("2.8"),
// 				Zones: pulumi.StringArray{
// 					pulumi.String("ru-central1-a"),
// 					pulumi.String("ru-central1-b"),
// 					pulumi.String("ru-central1-c"),
// 				},
// 				Zookeeper: &MdbKafkaClusterConfigZookeeperArgs{
// 					Resources: &MdbKafkaClusterConfigZookeeperResourcesArgs{
// 						DiskSize:         pulumi.Float64(20),
// 						DiskTypeId:       pulumi.String("network-ssd"),
// 						ResourcePresetId: pulumi.String("s2.micro"),
// 					},
// 				},
// 			},
// 			Environment: pulumi.String("PRESTABLE"),
// 			NetworkId:   fooVpcNetwork.ID(),
// 			SubnetIds: pulumi.StringArray{
// 				fooVpcSubnet.ID(),
// 				bar.ID(),
// 				baz.ID(),
// 			},
// 			Users: MdbKafkaClusterUserArray{
// 				&MdbKafkaClusterUserArgs{
// 					Name:     pulumi.String("producer-application"),
// 					Password: pulumi.String("password"),
// 					Permissions: MdbKafkaClusterUserPermissionArray{
// 						&MdbKafkaClusterUserPermissionArgs{
// 							Role:      pulumi.String("ACCESS_ROLE_PRODUCER"),
// 							TopicName: pulumi.String("input"),
// 						},
// 					},
// 				},
// 				&MdbKafkaClusterUserArgs{
// 					Name:     pulumi.String("worker"),
// 					Password: pulumi.String("password"),
// 					Permissions: MdbKafkaClusterUserPermissionArray{
// 						&MdbKafkaClusterUserPermissionArgs{
// 							Role:      pulumi.String("ACCESS_ROLE_CONSUMER"),
// 							TopicName: pulumi.String("input"),
// 						},
// 						&MdbKafkaClusterUserPermissionArgs{
// 							Role:      pulumi.String("ACCESS_ROLE_PRODUCER"),
// 							TopicName: pulumi.String("output"),
// 						},
// 					},
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
// A cluster can be imported using the `id` of the resource, e.g.
//
// ```sh
//  $ pulumi import yandex:index/mdbKafkaCluster:MdbKafkaCluster foo cluster_id
// ```
type MdbKafkaCluster struct {
	pulumi.CustomResourceState

	// Configuration of the Kafka cluster. The structure is documented below.
	Config MdbKafkaClusterConfigOutput `pulumi:"config"`
	// Timestamp of cluster creation.
	CreatedAt pulumi.StringOutput `pulumi:"createdAt"`
	// Inhibits deletion of the cluster.  Can be either `true` or `false`.
	DeletionProtection pulumi.BoolOutput `pulumi:"deletionProtection"`
	// Description of the Kafka cluster.
	Description pulumi.StringPtrOutput `pulumi:"description"`
	// Deployment environment of the Kafka cluster. Can be either `PRESTABLE` or `PRODUCTION`.
	// The default is `PRODUCTION`.
	Environment pulumi.StringPtrOutput `pulumi:"environment"`
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId pulumi.StringOutput `pulumi:"folderId"`
	// Health of the host.
	Health pulumi.StringOutput `pulumi:"health"`
	// A list of IDs of the host groups to place VMs of the cluster on.
	HostGroupIds pulumi.StringArrayOutput `pulumi:"hostGroupIds"`
	// A host of the Kafka cluster. The structure is documented below.
	Hosts MdbKafkaClusterHostArrayOutput `pulumi:"hosts"`
	// A set of key/value label pairs to assign to the Kafka cluster.
	Labels pulumi.StringMapOutput `pulumi:"labels"`
	// Maintenance policy of the Kafka cluster. The structure is documented below.
	MaintenanceWindow MdbKafkaClusterMaintenanceWindowPtrOutput `pulumi:"maintenanceWindow"`
	// The name of the topic.
	Name pulumi.StringOutput `pulumi:"name"`
	// ID of the network, to which the Kafka cluster belongs.
	NetworkId pulumi.StringOutput `pulumi:"networkId"`
	// Security group ids, to which the Kafka cluster belongs.
	SecurityGroupIds pulumi.StringArrayOutput `pulumi:"securityGroupIds"`
	// Status of the cluster. Can be either `CREATING`, `STARTING`, `RUNNING`, `UPDATING`, `STOPPING`, `STOPPED`, `ERROR` or `STATUS_UNKNOWN`.
	// For more information see `status` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-kafka/api-ref/Cluster/).
	Status pulumi.StringOutput `pulumi:"status"`
	// IDs of the subnets, to which the Kafka cluster belongs.
	SubnetIds pulumi.StringArrayOutput      `pulumi:"subnetIds"`
	Timeouts  MdbKafkaClusterTimeoutsOutput `pulumi:"timeouts"`
	// To manage topics, please switch to using a separate resource type `MdbKafkaTopic`.
	//
	// Deprecated: topic is deprecated
	Topics MdbKafkaClusterTopicArrayOutput `pulumi:"topics"`
	// A user of the Kafka cluster. The structure is documented below.
	Users MdbKafkaClusterUserArrayOutput `pulumi:"users"`
}

// NewMdbKafkaCluster registers a new resource with the given unique name, arguments, and options.
func NewMdbKafkaCluster(ctx *pulumi.Context,
	name string, args *MdbKafkaClusterArgs, opts ...pulumi.ResourceOption) (*MdbKafkaCluster, error) {
	if args == nil {
		return nil, errors.New("missing one or more required arguments")
	}

	if args.Config == nil {
		return nil, errors.New("invalid value for required argument 'Config'")
	}
	if args.NetworkId == nil {
		return nil, errors.New("invalid value for required argument 'NetworkId'")
	}
	if args.Timeouts == nil {
		return nil, errors.New("invalid value for required argument 'Timeouts'")
	}
	var resource MdbKafkaCluster
	err := ctx.RegisterResource("yandex:index/mdbKafkaCluster:MdbKafkaCluster", name, args, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// GetMdbKafkaCluster gets an existing MdbKafkaCluster resource's state with the given name, ID, and optional
// state properties that are used to uniquely qualify the lookup (nil if not required).
func GetMdbKafkaCluster(ctx *pulumi.Context,
	name string, id pulumi.IDInput, state *MdbKafkaClusterState, opts ...pulumi.ResourceOption) (*MdbKafkaCluster, error) {
	var resource MdbKafkaCluster
	err := ctx.ReadResource("yandex:index/mdbKafkaCluster:MdbKafkaCluster", name, id, state, &resource, opts...)
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

// Input properties used for looking up and filtering MdbKafkaCluster resources.
type mdbKafkaClusterState struct {
	// Configuration of the Kafka cluster. The structure is documented below.
	Config *MdbKafkaClusterConfig `pulumi:"config"`
	// Timestamp of cluster creation.
	CreatedAt *string `pulumi:"createdAt"`
	// Inhibits deletion of the cluster.  Can be either `true` or `false`.
	DeletionProtection *bool `pulumi:"deletionProtection"`
	// Description of the Kafka cluster.
	Description *string `pulumi:"description"`
	// Deployment environment of the Kafka cluster. Can be either `PRESTABLE` or `PRODUCTION`.
	// The default is `PRODUCTION`.
	Environment *string `pulumi:"environment"`
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// Health of the host.
	Health *string `pulumi:"health"`
	// A list of IDs of the host groups to place VMs of the cluster on.
	HostGroupIds []string `pulumi:"hostGroupIds"`
	// A host of the Kafka cluster. The structure is documented below.
	Hosts []MdbKafkaClusterHost `pulumi:"hosts"`
	// A set of key/value label pairs to assign to the Kafka cluster.
	Labels map[string]string `pulumi:"labels"`
	// Maintenance policy of the Kafka cluster. The structure is documented below.
	MaintenanceWindow *MdbKafkaClusterMaintenanceWindow `pulumi:"maintenanceWindow"`
	// The name of the topic.
	Name *string `pulumi:"name"`
	// ID of the network, to which the Kafka cluster belongs.
	NetworkId *string `pulumi:"networkId"`
	// Security group ids, to which the Kafka cluster belongs.
	SecurityGroupIds []string `pulumi:"securityGroupIds"`
	// Status of the cluster. Can be either `CREATING`, `STARTING`, `RUNNING`, `UPDATING`, `STOPPING`, `STOPPED`, `ERROR` or `STATUS_UNKNOWN`.
	// For more information see `status` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-kafka/api-ref/Cluster/).
	Status *string `pulumi:"status"`
	// IDs of the subnets, to which the Kafka cluster belongs.
	SubnetIds []string                 `pulumi:"subnetIds"`
	Timeouts  *MdbKafkaClusterTimeouts `pulumi:"timeouts"`
	// To manage topics, please switch to using a separate resource type `MdbKafkaTopic`.
	//
	// Deprecated: topic is deprecated
	Topics []MdbKafkaClusterTopic `pulumi:"topics"`
	// A user of the Kafka cluster. The structure is documented below.
	Users []MdbKafkaClusterUser `pulumi:"users"`
}

type MdbKafkaClusterState struct {
	// Configuration of the Kafka cluster. The structure is documented below.
	Config MdbKafkaClusterConfigPtrInput
	// Timestamp of cluster creation.
	CreatedAt pulumi.StringPtrInput
	// Inhibits deletion of the cluster.  Can be either `true` or `false`.
	DeletionProtection pulumi.BoolPtrInput
	// Description of the Kafka cluster.
	Description pulumi.StringPtrInput
	// Deployment environment of the Kafka cluster. Can be either `PRESTABLE` or `PRODUCTION`.
	// The default is `PRODUCTION`.
	Environment pulumi.StringPtrInput
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId pulumi.StringPtrInput
	// Health of the host.
	Health pulumi.StringPtrInput
	// A list of IDs of the host groups to place VMs of the cluster on.
	HostGroupIds pulumi.StringArrayInput
	// A host of the Kafka cluster. The structure is documented below.
	Hosts MdbKafkaClusterHostArrayInput
	// A set of key/value label pairs to assign to the Kafka cluster.
	Labels pulumi.StringMapInput
	// Maintenance policy of the Kafka cluster. The structure is documented below.
	MaintenanceWindow MdbKafkaClusterMaintenanceWindowPtrInput
	// The name of the topic.
	Name pulumi.StringPtrInput
	// ID of the network, to which the Kafka cluster belongs.
	NetworkId pulumi.StringPtrInput
	// Security group ids, to which the Kafka cluster belongs.
	SecurityGroupIds pulumi.StringArrayInput
	// Status of the cluster. Can be either `CREATING`, `STARTING`, `RUNNING`, `UPDATING`, `STOPPING`, `STOPPED`, `ERROR` or `STATUS_UNKNOWN`.
	// For more information see `status` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-kafka/api-ref/Cluster/).
	Status pulumi.StringPtrInput
	// IDs of the subnets, to which the Kafka cluster belongs.
	SubnetIds pulumi.StringArrayInput
	Timeouts  MdbKafkaClusterTimeoutsPtrInput
	// To manage topics, please switch to using a separate resource type `MdbKafkaTopic`.
	//
	// Deprecated: topic is deprecated
	Topics MdbKafkaClusterTopicArrayInput
	// A user of the Kafka cluster. The structure is documented below.
	Users MdbKafkaClusterUserArrayInput
}

func (MdbKafkaClusterState) ElementType() reflect.Type {
	return reflect.TypeOf((*mdbKafkaClusterState)(nil)).Elem()
}

type mdbKafkaClusterArgs struct {
	// Configuration of the Kafka cluster. The structure is documented below.
	Config MdbKafkaClusterConfig `pulumi:"config"`
	// Inhibits deletion of the cluster.  Can be either `true` or `false`.
	DeletionProtection *bool `pulumi:"deletionProtection"`
	// Description of the Kafka cluster.
	Description *string `pulumi:"description"`
	// Deployment environment of the Kafka cluster. Can be either `PRESTABLE` or `PRODUCTION`.
	// The default is `PRODUCTION`.
	Environment *string `pulumi:"environment"`
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId *string `pulumi:"folderId"`
	// A list of IDs of the host groups to place VMs of the cluster on.
	HostGroupIds []string `pulumi:"hostGroupIds"`
	// A set of key/value label pairs to assign to the Kafka cluster.
	Labels map[string]string `pulumi:"labels"`
	// Maintenance policy of the Kafka cluster. The structure is documented below.
	MaintenanceWindow *MdbKafkaClusterMaintenanceWindow `pulumi:"maintenanceWindow"`
	// The name of the topic.
	Name *string `pulumi:"name"`
	// ID of the network, to which the Kafka cluster belongs.
	NetworkId string `pulumi:"networkId"`
	// Security group ids, to which the Kafka cluster belongs.
	SecurityGroupIds []string `pulumi:"securityGroupIds"`
	// IDs of the subnets, to which the Kafka cluster belongs.
	SubnetIds []string                `pulumi:"subnetIds"`
	Timeouts  MdbKafkaClusterTimeouts `pulumi:"timeouts"`
	// To manage topics, please switch to using a separate resource type `MdbKafkaTopic`.
	//
	// Deprecated: topic is deprecated
	Topics []MdbKafkaClusterTopic `pulumi:"topics"`
	// A user of the Kafka cluster. The structure is documented below.
	Users []MdbKafkaClusterUser `pulumi:"users"`
}

// The set of arguments for constructing a MdbKafkaCluster resource.
type MdbKafkaClusterArgs struct {
	// Configuration of the Kafka cluster. The structure is documented below.
	Config MdbKafkaClusterConfigInput
	// Inhibits deletion of the cluster.  Can be either `true` or `false`.
	DeletionProtection pulumi.BoolPtrInput
	// Description of the Kafka cluster.
	Description pulumi.StringPtrInput
	// Deployment environment of the Kafka cluster. Can be either `PRESTABLE` or `PRODUCTION`.
	// The default is `PRODUCTION`.
	Environment pulumi.StringPtrInput
	// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
	FolderId pulumi.StringPtrInput
	// A list of IDs of the host groups to place VMs of the cluster on.
	HostGroupIds pulumi.StringArrayInput
	// A set of key/value label pairs to assign to the Kafka cluster.
	Labels pulumi.StringMapInput
	// Maintenance policy of the Kafka cluster. The structure is documented below.
	MaintenanceWindow MdbKafkaClusterMaintenanceWindowPtrInput
	// The name of the topic.
	Name pulumi.StringPtrInput
	// ID of the network, to which the Kafka cluster belongs.
	NetworkId pulumi.StringInput
	// Security group ids, to which the Kafka cluster belongs.
	SecurityGroupIds pulumi.StringArrayInput
	// IDs of the subnets, to which the Kafka cluster belongs.
	SubnetIds pulumi.StringArrayInput
	Timeouts  MdbKafkaClusterTimeoutsInput
	// To manage topics, please switch to using a separate resource type `MdbKafkaTopic`.
	//
	// Deprecated: topic is deprecated
	Topics MdbKafkaClusterTopicArrayInput
	// A user of the Kafka cluster. The structure is documented below.
	Users MdbKafkaClusterUserArrayInput
}

func (MdbKafkaClusterArgs) ElementType() reflect.Type {
	return reflect.TypeOf((*mdbKafkaClusterArgs)(nil)).Elem()
}

type MdbKafkaClusterInput interface {
	pulumi.Input

	ToMdbKafkaClusterOutput() MdbKafkaClusterOutput
	ToMdbKafkaClusterOutputWithContext(ctx context.Context) MdbKafkaClusterOutput
}

func (*MdbKafkaCluster) ElementType() reflect.Type {
	return reflect.TypeOf((**MdbKafkaCluster)(nil)).Elem()
}

func (i *MdbKafkaCluster) ToMdbKafkaClusterOutput() MdbKafkaClusterOutput {
	return i.ToMdbKafkaClusterOutputWithContext(context.Background())
}

func (i *MdbKafkaCluster) ToMdbKafkaClusterOutputWithContext(ctx context.Context) MdbKafkaClusterOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MdbKafkaClusterOutput)
}

// MdbKafkaClusterArrayInput is an input type that accepts MdbKafkaClusterArray and MdbKafkaClusterArrayOutput values.
// You can construct a concrete instance of `MdbKafkaClusterArrayInput` via:
//
//          MdbKafkaClusterArray{ MdbKafkaClusterArgs{...} }
type MdbKafkaClusterArrayInput interface {
	pulumi.Input

	ToMdbKafkaClusterArrayOutput() MdbKafkaClusterArrayOutput
	ToMdbKafkaClusterArrayOutputWithContext(context.Context) MdbKafkaClusterArrayOutput
}

type MdbKafkaClusterArray []MdbKafkaClusterInput

func (MdbKafkaClusterArray) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*MdbKafkaCluster)(nil)).Elem()
}

func (i MdbKafkaClusterArray) ToMdbKafkaClusterArrayOutput() MdbKafkaClusterArrayOutput {
	return i.ToMdbKafkaClusterArrayOutputWithContext(context.Background())
}

func (i MdbKafkaClusterArray) ToMdbKafkaClusterArrayOutputWithContext(ctx context.Context) MdbKafkaClusterArrayOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MdbKafkaClusterArrayOutput)
}

// MdbKafkaClusterMapInput is an input type that accepts MdbKafkaClusterMap and MdbKafkaClusterMapOutput values.
// You can construct a concrete instance of `MdbKafkaClusterMapInput` via:
//
//          MdbKafkaClusterMap{ "key": MdbKafkaClusterArgs{...} }
type MdbKafkaClusterMapInput interface {
	pulumi.Input

	ToMdbKafkaClusterMapOutput() MdbKafkaClusterMapOutput
	ToMdbKafkaClusterMapOutputWithContext(context.Context) MdbKafkaClusterMapOutput
}

type MdbKafkaClusterMap map[string]MdbKafkaClusterInput

func (MdbKafkaClusterMap) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*MdbKafkaCluster)(nil)).Elem()
}

func (i MdbKafkaClusterMap) ToMdbKafkaClusterMapOutput() MdbKafkaClusterMapOutput {
	return i.ToMdbKafkaClusterMapOutputWithContext(context.Background())
}

func (i MdbKafkaClusterMap) ToMdbKafkaClusterMapOutputWithContext(ctx context.Context) MdbKafkaClusterMapOutput {
	return pulumi.ToOutputWithContext(ctx, i).(MdbKafkaClusterMapOutput)
}

type MdbKafkaClusterOutput struct{ *pulumi.OutputState }

func (MdbKafkaClusterOutput) ElementType() reflect.Type {
	return reflect.TypeOf((**MdbKafkaCluster)(nil)).Elem()
}

func (o MdbKafkaClusterOutput) ToMdbKafkaClusterOutput() MdbKafkaClusterOutput {
	return o
}

func (o MdbKafkaClusterOutput) ToMdbKafkaClusterOutputWithContext(ctx context.Context) MdbKafkaClusterOutput {
	return o
}

// Configuration of the Kafka cluster. The structure is documented below.
func (o MdbKafkaClusterOutput) Config() MdbKafkaClusterConfigOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) MdbKafkaClusterConfigOutput { return v.Config }).(MdbKafkaClusterConfigOutput)
}

// Timestamp of cluster creation.
func (o MdbKafkaClusterOutput) CreatedAt() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringOutput { return v.CreatedAt }).(pulumi.StringOutput)
}

// Inhibits deletion of the cluster.  Can be either `true` or `false`.
func (o MdbKafkaClusterOutput) DeletionProtection() pulumi.BoolOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.BoolOutput { return v.DeletionProtection }).(pulumi.BoolOutput)
}

// Description of the Kafka cluster.
func (o MdbKafkaClusterOutput) Description() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringPtrOutput { return v.Description }).(pulumi.StringPtrOutput)
}

// Deployment environment of the Kafka cluster. Can be either `PRESTABLE` or `PRODUCTION`.
// The default is `PRODUCTION`.
func (o MdbKafkaClusterOutput) Environment() pulumi.StringPtrOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringPtrOutput { return v.Environment }).(pulumi.StringPtrOutput)
}

// The ID of the folder that the resource belongs to. If it is not provided, the default provider folder is used.
func (o MdbKafkaClusterOutput) FolderId() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringOutput { return v.FolderId }).(pulumi.StringOutput)
}

// Health of the host.
func (o MdbKafkaClusterOutput) Health() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringOutput { return v.Health }).(pulumi.StringOutput)
}

// A list of IDs of the host groups to place VMs of the cluster on.
func (o MdbKafkaClusterOutput) HostGroupIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringArrayOutput { return v.HostGroupIds }).(pulumi.StringArrayOutput)
}

// A host of the Kafka cluster. The structure is documented below.
func (o MdbKafkaClusterOutput) Hosts() MdbKafkaClusterHostArrayOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) MdbKafkaClusterHostArrayOutput { return v.Hosts }).(MdbKafkaClusterHostArrayOutput)
}

// A set of key/value label pairs to assign to the Kafka cluster.
func (o MdbKafkaClusterOutput) Labels() pulumi.StringMapOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringMapOutput { return v.Labels }).(pulumi.StringMapOutput)
}

// Maintenance policy of the Kafka cluster. The structure is documented below.
func (o MdbKafkaClusterOutput) MaintenanceWindow() MdbKafkaClusterMaintenanceWindowPtrOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) MdbKafkaClusterMaintenanceWindowPtrOutput { return v.MaintenanceWindow }).(MdbKafkaClusterMaintenanceWindowPtrOutput)
}

// The name of the topic.
func (o MdbKafkaClusterOutput) Name() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringOutput { return v.Name }).(pulumi.StringOutput)
}

// ID of the network, to which the Kafka cluster belongs.
func (o MdbKafkaClusterOutput) NetworkId() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringOutput { return v.NetworkId }).(pulumi.StringOutput)
}

// Security group ids, to which the Kafka cluster belongs.
func (o MdbKafkaClusterOutput) SecurityGroupIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringArrayOutput { return v.SecurityGroupIds }).(pulumi.StringArrayOutput)
}

// Status of the cluster. Can be either `CREATING`, `STARTING`, `RUNNING`, `UPDATING`, `STOPPING`, `STOPPED`, `ERROR` or `STATUS_UNKNOWN`.
// For more information see `status` field of JSON representation in [the official documentation](https://cloud.yandex.com/docs/managed-kafka/api-ref/Cluster/).
func (o MdbKafkaClusterOutput) Status() pulumi.StringOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringOutput { return v.Status }).(pulumi.StringOutput)
}

// IDs of the subnets, to which the Kafka cluster belongs.
func (o MdbKafkaClusterOutput) SubnetIds() pulumi.StringArrayOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) pulumi.StringArrayOutput { return v.SubnetIds }).(pulumi.StringArrayOutput)
}

func (o MdbKafkaClusterOutput) Timeouts() MdbKafkaClusterTimeoutsOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) MdbKafkaClusterTimeoutsOutput { return v.Timeouts }).(MdbKafkaClusterTimeoutsOutput)
}

// To manage topics, please switch to using a separate resource type `MdbKafkaTopic`.
//
// Deprecated: topic is deprecated
func (o MdbKafkaClusterOutput) Topics() MdbKafkaClusterTopicArrayOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) MdbKafkaClusterTopicArrayOutput { return v.Topics }).(MdbKafkaClusterTopicArrayOutput)
}

// A user of the Kafka cluster. The structure is documented below.
func (o MdbKafkaClusterOutput) Users() MdbKafkaClusterUserArrayOutput {
	return o.ApplyT(func(v *MdbKafkaCluster) MdbKafkaClusterUserArrayOutput { return v.Users }).(MdbKafkaClusterUserArrayOutput)
}

type MdbKafkaClusterArrayOutput struct{ *pulumi.OutputState }

func (MdbKafkaClusterArrayOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*[]*MdbKafkaCluster)(nil)).Elem()
}

func (o MdbKafkaClusterArrayOutput) ToMdbKafkaClusterArrayOutput() MdbKafkaClusterArrayOutput {
	return o
}

func (o MdbKafkaClusterArrayOutput) ToMdbKafkaClusterArrayOutputWithContext(ctx context.Context) MdbKafkaClusterArrayOutput {
	return o
}

func (o MdbKafkaClusterArrayOutput) Index(i pulumi.IntInput) MdbKafkaClusterOutput {
	return pulumi.All(o, i).ApplyT(func(vs []interface{}) *MdbKafkaCluster {
		return vs[0].([]*MdbKafkaCluster)[vs[1].(int)]
	}).(MdbKafkaClusterOutput)
}

type MdbKafkaClusterMapOutput struct{ *pulumi.OutputState }

func (MdbKafkaClusterMapOutput) ElementType() reflect.Type {
	return reflect.TypeOf((*map[string]*MdbKafkaCluster)(nil)).Elem()
}

func (o MdbKafkaClusterMapOutput) ToMdbKafkaClusterMapOutput() MdbKafkaClusterMapOutput {
	return o
}

func (o MdbKafkaClusterMapOutput) ToMdbKafkaClusterMapOutputWithContext(ctx context.Context) MdbKafkaClusterMapOutput {
	return o
}

func (o MdbKafkaClusterMapOutput) MapIndex(k pulumi.StringInput) MdbKafkaClusterOutput {
	return pulumi.All(o, k).ApplyT(func(vs []interface{}) *MdbKafkaCluster {
		return vs[0].(map[string]*MdbKafkaCluster)[vs[1].(string)]
	}).(MdbKafkaClusterOutput)
}

func init() {
	pulumi.RegisterInputType(reflect.TypeOf((*MdbKafkaClusterInput)(nil)).Elem(), &MdbKafkaCluster{})
	pulumi.RegisterInputType(reflect.TypeOf((*MdbKafkaClusterArrayInput)(nil)).Elem(), MdbKafkaClusterArray{})
	pulumi.RegisterInputType(reflect.TypeOf((*MdbKafkaClusterMapInput)(nil)).Elem(), MdbKafkaClusterMap{})
	pulumi.RegisterOutputType(MdbKafkaClusterOutput{})
	pulumi.RegisterOutputType(MdbKafkaClusterArrayOutput{})
	pulumi.RegisterOutputType(MdbKafkaClusterMapOutput{})
}