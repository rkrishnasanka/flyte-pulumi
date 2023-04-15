package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Find default VPC for the EKS cluster
func findDefaultVPC(ctx *pulumi.Context) (*ec2.LookupVpcResult, error) {

	// Find default VPC
	vpcID := "vpc-009f4986f9d95a645"
	vpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{Id: &vpcID})
	if err != nil {
		return nil, err
	}

	return vpc, nil
}

// Find the default Relational Database Service (RDS)
func findDefaultRDS(ctx *pulumi.Context) (*rds.LookupInstanceResult, error) {
	// // Find default RDS
	// rdsName := "Whatever"
	// rdsID := "rds-1234567890"
	// rds, err := rds.GetInstance(ctx, rdsName, rdsID, nil)

	// if err != nil {
	// 	panic("Could not find RDS")
	// }
	instanceLookup, err := rds.LookupInstance(ctx, &rds.LookupInstanceArgs{
		DbInstanceIdentifier: "database-1",
	})
	if err != nil {
		return nil, fmt.Errorf("could not find RDS: %w", err)
	}
	return instanceLookup, nil
}
