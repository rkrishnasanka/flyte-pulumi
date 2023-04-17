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
		_, err = createEKSCluster(ctx, eksClusterRole, eksNodeGroupRole, subnets, clusterSecurityGroup)
		if err != nil {
			panic(err)
		}

		return nil
	})

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

}

func toPulumiStringArray(a []string) pulumi.StringArrayInput {
	var res []pulumi.StringInput
	for _, s := range a {
		res = append(res, pulumi.String(s))
	}
	return pulumi.StringArray(res)
}
