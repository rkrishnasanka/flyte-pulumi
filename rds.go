package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createRDSCluster(ctx *pulumi.Context, subnets *ec2.GetSubnetsResult) (*rds.Instance, error) {
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

	// var clusterInstances []*rds.ClusterInstance
	// for index := 0; index < 2; index++ {
	// 	key0 := index
	// 	val0 := index
	// 	__res, err := rds.NewClusterInstance(ctx, fmt.Sprintf("clusterInstances-%v", key0), &rds.ClusterInstanceArgs{
	// 		Identifier:        pulumi.String(fmt.Sprintf("aurora-cluster-demo-%v", val0)),
	// 		ClusterIdentifier: _default.ID(),
	// 		InstanceClass:     pulumi.String("db.r4.large"),
	// 		Engine:            _default.Engine,
	// 		EngineVersion:     _default.EngineVersion,
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	clusterInstances = append(clusterInstances, __res)
	// }

	rdsInstance, err := rds.NewInstance(ctx, "flyte-db", &rds.InstanceArgs{
		AllocatedStorage: pulumi.Int(20),
		DbName:           pulumi.String("flyteadmin"),
		Engine:           pulumi.String("postgres"),
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
