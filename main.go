package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/s3"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type FlyteConfig struct {
	// The name of the VPC to use for the EKS cluster
	VpcName string
	// The name of the RDS instance to use for the Flyte database
	RdsName string
	// The name of the EKS cluster to create
	EksClusterName string
	// The name of the EKS node group to create
	PrimaryNodeGroupScalingConfig *ScalingConfig
	// RDS password
	RdsPassword string
}

type ScalingConfig struct {
	// The minimum number of nodes to run in the node group
	MinSize int
	// The maximum number of nodes to run in the node group
	MaxSize int
	// The number of nodes to run in the node group
	DesiredSize int
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// Read back the default VPC and public subnets, which we will use.
		// t := true

		vpcID := "vpc-012c395322b57c628"

		// if vpcID == "" {
		// 	vpc, err := findDefaultVPC(ctx)
		// 	if err != nil {
		// 		panic(err)
		// 	}
		// }

		vpc, err := ec2.LookupVpc(ctx, &ec2.LookupVpcArgs{Id: &vpcID})
		if err != nil {
			return err
		}

		// subnet := []string{"subnet-0063ed129917d3c44"}

		subnets, err := ec2.GetSubnets(ctx, &ec2.GetSubnetsArgs{
			Filters: []ec2.GetSubnetsFilter{
				{
					Name: "subnet-id",
					Values: []string{
						"subnet-0063ed129917d3c44",
						"subnet-0dfa77319d1ada651",
						"subnet-083462886f234b559",
					},
				},
			},
		})
		if err != nil {
			panic(err)
		}

		// // Get the subnets from the VPC
		// subnets, err := ec2.GetSubnets(ctx, &ec2.GetSubnetsArgs{
		// 	Filters: []ec2.GetSubnetsFilter{
		// 		{
		// 			Name:   "vpc-id",
		// 			Values: []string{vpc.Id},
		// 		},
		// 	},
		// })
		// if err != nil {
		// 	panic(err)
		// }

		// Create the S3 bucket for storing the flyte data
		flyteBucket, err := s3.NewBucket(ctx, "flyte-admin-bucket", nil)
		if err != nil {
			return err
		}

		// TODO - Figure out if the vpc has access to the bucket

		// Export the name of the bucket
		ctx.Export("flyte-admin-bucketName", flyteBucket.ID())

		// Create the roles for the EKS Cluster
		// eksClusterRole, eksNodeGroupRole, err := createEKSRoles(ctx)
		eksClusterRole, eksNodeGroupRole, err := createEKSRoles(ctx)
		if err != nil {
			panic(err)
		}

		// Create the Security Groups for the EKS Cluster
		clusterSecurityGroup, err := createEKSSecurityGroup(ctx, vpc)
		if err != nil {
			panic(err)
		}

		// Create the EKS Cluster
		eksCluster, err := createEKSCluster(ctx, eksClusterRole, eksNodeGroupRole, subnets, clusterSecurityGroup)
		if err != nil {
			panic(err)
		}

		// Create the RDS instance
		_, err = createRDSCluster(ctx, subnets, eksCluster)
		if err != nil {
			panic(err)
		}

		return nil
	})

}

func toPulumiStringArray(a []string) pulumi.StringArrayInput {
	var res []pulumi.StringInput
	for _, s := range a {
		res = append(res, pulumi.String(s))
	}
	return pulumi.StringArray(res)
}
