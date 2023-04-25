```go

// OIDC issuer URL for the cluster.
// Add the OIDC provider to the cluster so that we can create IAM roles that can be assumed by Kubernetes Service Accounts.
_, err = eks.NewIdentityProviderConfig(ctx, "eks-flyte-cluster-oidc", &eks.IdentityProviderConfigArgs{
		ClusterName: eksCluster.Name,
		Oidc: &eks.IdentityProviderConfigOidcArgs{
			IdentityProviderConfigName: pulumi.String("oidc"),
			IssuerUrl:                  eksCluster.Identities.Index().Oidcs().Index(0).Issuer().Elem().ToStringOutput(), // eksCluster.Identities[0].Oidcs[0].Issuer,
    },
})
if err != nil {
    return nil, err
}


```

```go

//EKS Tutorial code
namespace, err := corev1.NewNamespace(ctx, "app-ns", &corev1.NamespaceArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Name: pulumi.String("joe-duffy"),
    },
}, pulumi.Provider(k8sProvider))
if err != nil {
    return err
}

appLabels := pulumi.StringMap{
    "app": pulumi.String("iac-workshop"),
}
_, err = appsv1.NewDeployment(ctx, "app-dep", &appsv1.DeploymentArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Namespace: namespace.Metadata.Elem().Name(),
    },
    Spec: appsv1.DeploymentSpecArgs{
        Selector: &metav1.LabelSelectorArgs{
            MatchLabels: appLabels,
        },
        Replicas: pulumi.Int(3),
        Template: &corev1.PodTemplateSpecArgs{
            Metadata: &metav1.ObjectMetaArgs{
                Labels: appLabels,
            },
            Spec: &corev1.PodSpecArgs{
                Containers: corev1.ContainerArray{
                    corev1.ContainerArgs{
                        Name:  pulumi.String("iac-workshop"),
                        Image: pulumi.String("jocatalin/kubernetes-bootcamp:v2"),
                    },
                },
            },
        },
    },
}, pulumi.Provider(k8sProvider))
if err != nil {
    return err
}

service, err := corev1.NewService(ctx, "app-service", &corev1.ServiceArgs{
    Metadata: &metav1.ObjectMetaArgs{
        Namespace: namespace.Metadata.Elem().Name(),
        Labels:    appLabels,
    },
    Spec: &corev1.ServiceSpecArgs{
        Ports: corev1.ServicePortArray{
            corev1.ServicePortArgs{
                Port:       pulumi.Int(80),
                TargetPort: pulumi.Int(8080),
            },
        },
        Selector: appLabels,
        Type:     pulumi.String("LoadBalancer"),
    },
}, pulumi.Provider(k8sProvider))
if err != nil {
    return err
}

ctx.Export("url", service.Status.ApplyT(func(status *corev1.ServiceStatus) *string {
    ingress := status.LoadBalancer.Ingress[0]
    if ingress.Hostname != nil {
        return ingress.Hostname
    }
    return ingress.Ip
}))



```

```go

// Attaches COREDNS to the EKS Cluster. For now this is uncessary since the default eks
// seems to have CoreDNS already installed.
func attachCoreDNS(ctx *pulumi.Context, eksCluster *eks.Cluster) error {
	_, err := eks.NewAddon(ctx, "eks-core-dns-addon", &eks.AddonArgs{
		ClusterName:      eksCluster.Name,
		AddonName:        pulumi.String("coredns"),
		AddonVersion:     pulumi.String("v1.8.7-eksbuild.3"),
		ResolveConflicts: pulumi.String("OVERWRITE"),
	})
	if err != nil {
		return err
	}

	return nil
}

```

```go

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

```