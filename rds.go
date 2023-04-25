package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createRDSCluster(ctx *pulumi.Context, subnets *ec2.GetSubnetsResult) (*rds.Cluster, error) {
	subnetGroup, err := rds.NewSubnetGroup(ctx, "flyteadmin-rds-subnetgroup", &rds.SubnetGroupArgs{
		SubnetIds: toPulumiStringArray(subnets.Ids),
		// Tags: pulumi.StringMap{
		// 	"Name": pulumi.String("My DB subnet group"),
		// },
	})
	if err != nil {
		return nil, err
	}

	rdsCluster, err := rds.NewCluster(ctx, "flyteadmin-rds-cluster", &rds.ClusterArgs{
		ClusterIdentifier: pulumi.String("flyteadmin"),
		Engine:            pulumi.String("aurora-postgresql"),
		EngineVersion:     pulumi.String("15"),
		DatabaseName:      pulumi.String("flyteadmin"),
		MasterUsername:    pulumi.String("flyteadmin"),
		MasterPassword:    pulumi.String("thisisaweakpassword"),
		// DbClusterInstanceClass: pulumi.String("db.r5.large"),
		DbSubnetGroupName: subnetGroup.Name,
		// AllocatedStorage:  pulumi.Int(20),
		// VpcSecurityGroupIds: pulumi.StringArray{
		// 	pulumi.String("sg-0a1b2c3d4e5f6g7h8"),
		// },
		// SkipFinalSnapshot: pulumi.Bool(true),
		// DeletionProtection: pulumi.Bool(false),
		// StorageEncrypted: pulumi.Bool(true),
	})
	if err != nil {
		return nil, err
	}

	var clusterInstances []*rds.ClusterInstance
	for index := 0; index < 1; index++ {
		key0 := index
		val0 := index
		__res, err := rds.NewClusterInstance(ctx, fmt.Sprintf("flyte-rds-clusterInstances-%v", key0), &rds.ClusterInstanceArgs{
			Identifier:        pulumi.String(fmt.Sprintf("flyte-aurora-cluster-instance-%v", val0)),
			ClusterIdentifier: rdsCluster.ID(),
			InstanceClass:     pulumi.String("db.r5.large"),
			Engine:            rdsCluster.Engine,
			EngineVersion:     rdsCluster.EngineVersion,
		})
		if err != nil {
			return nil, err
		}
		clusterInstances = append(clusterInstances, __res)
	}
	for _, clusterInstance := range clusterInstances {
		fmt.Printf("clusterInstance: %v\n", clusterInstance.ID())
	}

	return rdsCluster, nil
}

func getCurrentRDSCluster(ctx *pulumi.Context, vpcID string) error {

	panic("implement me")
	rds.LookupCluster(ctx, &rds.LookupClusterArgs{
		ClusterIdentifier: "database-1",
	})
	return nil
}
