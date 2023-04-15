package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/rds"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Find default VPC for the EKS cluster
func findDefaultVPC(ctx *pulumi.Context) *ec2.LookupVpcResult {

	// Find default VPC
	vpcID := "vpc-009f4986f9d95a645"
	vpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{Id: &vpcID})
	if err != nil {
		panic("Could not find VPC")
	}

	return vpc
}

// Find the default Relational Database Service (RDS)
func findDefaultRDS(ctx *pulumi.Context) *rds.Instance {
	// // Find default RDS
	// rdsName := "Whatever"
	// rdsID := "rds-1234567890"
	// rds, err := rds.GetInstance(ctx, rdsName, rdsID, nil)

	// if err != nil {
	// 	panic("Could not find RDS")
	// }

	// return rds
	panic("Not implemented")
}
