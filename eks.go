package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/ec2"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/eks"
	"github.com/pulumi/pulumi-aws/sdk/v5/go/aws/iam"
	"github.com/pulumi/pulumi-kubernetes/sdk/v3/go/kubernetes"
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

	nodeGroupName := "flyte-eks-nodegroup-primary"
	primaryNodeGroup, err := eks.NewNodeGroup(ctx, nodeGroupName, &eks.NodeGroupArgs{
		ClusterName:   eksCluster.Name,
		NodeGroupName: pulumi.String(nodeGroupName),
		NodeRoleArn:   pulumi.StringInput(nodeGroupRole.Arn),
		SubnetIds:     toPulumiStringArray(subnets.Ids),
		ScalingConfig: &eks.NodeGroupScalingConfigArgs{
			DesiredSize: pulumi.Int(5),
			MaxSize:     pulumi.Int(5),
			MinSize:     pulumi.Int(2),
		},
		// Currently fixing the AMI to the latest Amazon Linux 2 AMI
		AmiType: pulumi.String("AL2_x86_64"),

		// TODO - Figure out how we need to setup the instance sizes
		InstanceTypes: pulumi.StringArray{
			pulumi.String("t2.nano"), // Replace with your desired instance type(s)
		},

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
	// k8sProvider
	_, err = kubernetes.NewProvider(ctx, "k8sprovider", &kubernetes.ProviderArgs{
		Kubeconfig: generateKubeconfig(eksCluster.Endpoint,
			eksCluster.CertificateAuthority.Data().Elem(), eksCluster.Name),
	}, pulumi.DependsOn([]pulumi.Resource{primaryNodeGroup}))
	if err != nil {
		return nil, err
	}

	// Retrieve the OIDC issuer URL
	// aws eks describe-cluster --region us-east-1 --name eks-flyte-cluster-12345 --query "cluster.identity.oidc.issuer" --output text

	// // Print the OIDC issuer URL
	// ctx.Export("OIDCIssuerURL", pulumi.Sprintf("OIDC issuer URL: %s", *oidcIssuer))

	// Create an IAM OpenID Connect (OIDC) provider
	// eksctl utils associate-iam-oidc-provider --cluster <Name-EKS-Cluster> --approve
	// _, err = eks.NewIdentityProviderConfig(ctx, "flyte-system-oidc-provider", &eks.IdentityProviderConfigArgs{
	// 	ClusterName: eksCluster.Name,
	// 	// Oidc:        &eks.IdentityProviderConfigOidcArgs{},
	// })
	// if err != nil {
	// 	return nil, err
	// }

	// // Create an IAM role for the service account
	// _, err = createFlyteSystemRole(ctx)
	// if err != nil {
	// 	return nil, err
	// }

	return eksCluster, nil
}

//  EKS Tutorial code
// 	namespace, err := corev1.NewNamespace(ctx, "app-ns", &corev1.NamespaceArgs{
// 		Metadata: &metav1.ObjectMetaArgs{
// 			Name: pulumi.String("joe-duffy"),
// 		},
// 	}, pulumi.Provider(k8sProvider))
// 	if err != nil {
// 		return err
// 	}

// 	appLabels := pulumi.StringMap{
// 		"app": pulumi.String("iac-workshop"),
// 	}
// 	_, err = appsv1.NewDeployment(ctx, "app-dep", &appsv1.DeploymentArgs{
// 		Metadata: &metav1.ObjectMetaArgs{
// 			Namespace: namespace.Metadata.Elem().Name(),
// 		},
// 		Spec: appsv1.DeploymentSpecArgs{
// 			Selector: &metav1.LabelSelectorArgs{
// 				MatchLabels: appLabels,
// 			},
// 			Replicas: pulumi.Int(3),
// 			Template: &corev1.PodTemplateSpecArgs{
// 				Metadata: &metav1.ObjectMetaArgs{
// 					Labels: appLabels,
// 				},
// 				Spec: &corev1.PodSpecArgs{
// 					Containers: corev1.ContainerArray{
// 						corev1.ContainerArgs{
// 							Name:  pulumi.String("iac-workshop"),
// 							Image: pulumi.String("jocatalin/kubernetes-bootcamp:v2"),
// 						},
// 					},
// 				},
// 			},
// 		},
// 	}, pulumi.Provider(k8sProvider))
// 	if err != nil {
// 		return err
// 	}

// 	service, err := corev1.NewService(ctx, "app-service", &corev1.ServiceArgs{
// 		Metadata: &metav1.ObjectMetaArgs{
// 			Namespace: namespace.Metadata.Elem().Name(),
// 			Labels:    appLabels,
// 		},
// 		Spec: &corev1.ServiceSpecArgs{
// 			Ports: corev1.ServicePortArray{
// 				corev1.ServicePortArgs{
// 					Port:       pulumi.Int(80),
// 					TargetPort: pulumi.Int(8080),
// 				},
// 			},
// 			Selector: appLabels,
// 			Type:     pulumi.String("LoadBalancer"),
// 		},
// 	}, pulumi.Provider(k8sProvider))
// 	if err != nil {
// 		return err
// 	}

// 	ctx.Export("url", service.Status.ApplyT(func(status *corev1.ServiceStatus) *string {
// 		ingress := status.LoadBalancer.Ingress[0]
// 		if ingress.Hostname != nil {
// 			return ingress.Hostname
// 		}
// 		return ingress.Ip
// 	}))

// 	return nil
// })

// Create the KubeConfig Structure as per https://docs.aws.amazon.com/eks/latest/userguide/create-kubeconfig.html
func generateKubeconfig(clusterEndpoint pulumi.StringOutput, certData pulumi.StringOutput, clusterName pulumi.StringOutput) pulumi.StringOutput {
	return pulumi.Sprintf(`{
        "apiVersion": "v1",
        "clusters": [{
            "cluster": {
                "server": "%s",
                "certificate-authority-data": "%s"
            },
            "name": "kubernetes",
        }],
        "contexts": [{
            "context": {
                "cluster": "kubernetes",
                "user": "aws",
            },
            "name": "aws",
        }],
        "current-context": "aws",
        "kind": "Config",
        "users": [{
            "name": "aws",
            "user": {
                "exec": {
                    "apiVersion": "client.authentication.k8s.io/v1beta1",
                    "command": "aws-iam-authenticator",
                    "args": [
                        "token",
                        "-i",
                        "%s",
                    ],
                },
            },
        }],
    }`, clusterEndpoint, certData, clusterName)
}
