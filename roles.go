package main

import (
	"fmt"

	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

// Create the roles for the EKS Cluster
func createEKSRoles(ctx *pulumi.Context) (*iam.Role, *iam.Role, error) {

	// Create EKS Cluster Role
	eksClusterRole, err := iam.NewRole(ctx, "eks-iam-eksClusterRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(
			`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": {
						"Service": [
							"eks.amazonaws.com"
						]
					},
					"Action": "sts:AssumeRole"
				}
			]
		}`),
		// PermissionsBoundary: pulumi.String("arn:aws:iam::aws:policy/AdministratorAccess"),
	})

	if err != nil {
		return nil, nil, err
	}

	eksPolicies := []string{
		"arn:aws:iam::aws:policy/AmazonEKSServicePolicy",
		"arn:aws:iam::aws:policy/AmazonEKSClusterPolicy",
		"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
	}
	for i, eksPolicy := range eksPolicies {
		_, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("rpa-%d", i), &iam.RolePolicyAttachmentArgs{
			PolicyArn: pulumi.String(eksPolicy),
			Role:      eksClusterRole.Name,
		})
		if err != nil {
			return nil, nil, err
		}
	}

	// Create the EC2 NodeGroup Role
	nodeGroupRole, err := iam.NewRole(ctx, "eks-iam-nodeGroupRole", &iam.RoleArgs{
		AssumeRolePolicy: pulumi.String(
			`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Action": [
					"sts:AssumeRole"
				],
				"Principal": {
					"Service": [
						"ec2.amazonaws.com"
					]
				}
			}
		]
	}`),
	})
	if err != nil {
		return nil, nil, err
	}
	nodeGroupPolicies := []string{
		"arn:aws:iam::aws:policy/AmazonEKSWorkerNodePolicy",
		"arn:aws:iam::aws:policy/AmazonEKS_CNI_Policy",
		"arn:aws:iam::aws:policy/AmazonEC2ContainerRegistryReadOnly",
	}
	for i, nodeGroupPolicy := range nodeGroupPolicies {
		_, err := iam.NewRolePolicyAttachment(ctx, fmt.Sprintf("ngpa-%d", i), &iam.RolePolicyAttachmentArgs{
			Role:      nodeGroupRole.Name,
			PolicyArn: pulumi.String(nodeGroupPolicy),
		})
		if err != nil {
			return nil, nil, err
		}
	}

	return eksClusterRole, nodeGroupRole, nil
}

func createFlyteSystemRole(ctx *pulumi.Context) (*iam.Role, error) {
	panic("Not Implemented !")
	// flyte_system_role, err := iam.NewRole(ctx, "flyte-system", nil)
	// if err != nil {
	// 	return nil, err
	// }
	// _, err = iam.NewRolePolicyAttachment(ctx, "flyte-system-attachment", &iam.RolePolicyAttachmentArgs{
	// 	Role:      flyte_system_role.Name,
	// 	PolicyArn: pulumi.String("arn:aws:iam::aws:policy/AmazonS3FullAccess"),
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// return flyte_system_role, nil
}
