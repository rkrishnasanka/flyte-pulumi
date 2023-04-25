package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/eks"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createRDSCluster(ctx *pulumi.Context, subnets *ec2.GetSubnetsResult, eksCluster *eks.Cluster) (*rds.Instance, error) {
	subnetGroup, err := rds.NewSubnetGroup(ctx, "flyteadmin-rds-subnetgroup", &rds.SubnetGroupArgs{
		SubnetIds: toPulumiStringArray(subnets.Ids),
		// Tags: pulumi.StringMap{
		// 	"Name": pulumi.String("My DB subnet group"),
		// },
	})
	if err != nil {
		return nil, err
	}

	// rdsCluster, err := rds.NewCluster(ctx, "flyteadmin-rds-cluster", &rds.ClusterArgs{
	// 	ClusterIdentifier: pulumi.String("flyteadmin"),
	// 	Engine:            pulumi.String("aurora-postgresql"),
	// 	EngineVersion:     pulumi.String("15"),
	// 	DatabaseName:      pulumi.String("flyteadmin"),
	// 	MasterUsername:    pulumi.String("flyteadmin"),
	// 	MasterPassword:    pulumi.String("thisisaweakpassword"),
	// 	// DbClusterInstanceClass: pulumi.String("db.r5.large"),
	// 	DbSubnetGroupName: subnetGroup.Name,
	// 	// AllocatedStorage:  pulumi.Int(20),
	// 	// VpcSecurityGroupIds: pulumi.StringArray{
	// 	// 	pulumi.String("sg-0a1b2c3d4e5f6g7h8"),
	// 	// },
	// 	// SkipFinalSnapshot: pulumi.Bool(true),
	// 	// DeletionProtection: pulumi.Bool(false),
	// 	// StorageEncrypted: pulumi.Bool(true),
	// })
	// _, err = rds.NewInstance(ctx, "flyteadmin-rds-instance", &rds.InstanceArgs{
	// 	InstanceClass:     pulumi.String("db.r5.large"),
	// 	ClusterIdentifier: rdsCluster.ID(),
	// 	// InstanceClass:     pulumi.String("db.r5.large"),
	// 	Engine:            pulumi.String("aurora-postgresql"),
	// 	EngineVersion:     pulumi.String("15"),
	// 	DbSubnetGroupName: subnetGroup.Name,
	// 	InstanceClass:     pulumi.String("db.t3.micro"),
	// })

	rdsInstance, err := rds.NewInstance(ctx, "default", &rds.InstanceArgs{
		AllocatedStorage: pulumi.Int(20),
		DbName:           pulumi.String("flyteadmin"),
		Engine:           pulumi.String("aurora-postgresql"),
		EngineVersion:    pulumi.String("15"),
		InstanceClass:    pulumi.String("db.t3.micro"),
		// ParameterGroupName: pulumi.String("default.mysql5.7"),
		Password:          pulumi.String("thisisaweakpassword"),
		SkipFinalSnapshot: pulumi.Bool(true),
		Username:          pulumi.String("flyteadmin"),
		DbSubnetGroupName: subnetGroup.Name,
	})
	if err != nil {
		return nil, err
	}

	return rdsInstance, nil
}

func getCurrentRDSCluster(ctx *pulumi.Context, vpcID string) error {

	panic("implement me")
	rds.LookupCluster(ctx, &rds.LookupClusterArgs{
		ClusterIdentifier: "database-1",
	})
	return nil
}
