package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/eks"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func createEKSCluster(ctx *pulumi.Context, eksClusterRole *iam.Role, nodeGroupRole *iam.Role, subnets *ec2.GetSubnetsResult, securityGroup *ec2.SecurityGroup) (*eks.Cluster, error) {
	// Create an EKS cluster
	eksCluster, err := eks.NewCluster(ctx, "eks-flyte-cluster", &eks.ClusterArgs{
		RoleArn: pulumi.StringInput(eksClusterRole.Arn),
		VpcConfig: &eks.ClusterVpcConfigArgs{
			PublicAccessCidrs: pulumi.StringArray{
				pulumi.String("0.0.0.0/0"),
			},
			SecurityGroupIds: pulumi.StringArray{
				securityGroup.ID().ToStringOutput(),
			},
			SubnetIds: toPulumiStringArray(subnets.Ids),
		},
	})
	if err != nil {
		return nil, err
	}

	// Create the EKS Node Group
	// TODO - We need to update the scaling capabilities as a new argument to the function and make it user definable

	nodeGroupName := "fyte-eks-nodegroup-primary"
	_, err = eks.NewNodeGroup(ctx, nodeGroupName, &eks.NodeGroupArgs{
		ClusterName:   eksCluster.Name,
		NodeGroupName: pulumi.String(nodeGroupName),
		NodeRoleArn:   pulumi.StringInput(nodeGroupRole.Arn),
		SubnetIds:     toPulumiStringArray(subnets.Ids),
		ScalingConfig: &eks.NodeGroupScalingConfigArgs{
			DesiredSize: pulumi.Int(5),
			MaxSize:     pulumi.Int(5),
			MinSize:     pulumi.Int(2),
		},

		// TODO - Figure out how we need to setup the instance sizes
		// InstanceTypes: pulumi.StringArray{
		// 	pulumi.String("t2.small"), // Replace with your desired instance type(s)
		// },

		// TODO - Add SSH Key
		// RemoteAccess: &eks.NodeGroupRemoteAccessArgs{
		// 	Ec2SshKey: pulumi.String("my-ssh-key"), // Replace with your desired SSH key name
		// },
	})
	if err != nil {
		return nil, err
	}

	ctx.Export("kubeconfig", generateKubeconfig(eksCluster.Endpoint,
		eksCluster.CertificateAuthority.Data().Elem(), eksCluster.Name))

	return eksCluster, nil
}
