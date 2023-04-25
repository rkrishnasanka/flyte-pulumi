# flyte-pulumi


## Manual Testing 


```sh
aws eks update-kubeconfig --name <CLUSTER-NAME> --region us-east-1
```

```sh
kubectl config current-context <CLUSTER-ARM>
```

```sh
kubectl get pods
```

```sh
aws eks describe-cluster --region us-east-1 --name <CLUSTER-NAME> --query "cluster.identity.oidc.issuer" --output text
```

```sh
kubectl run pgsql-postgresql-client --rm --tty -i --restart='Never' --namespace testdb --image docker.io/bitnami/postgresql:11.7.0-debian-10-r9 --env='PGPASSWORD=thisisaweakpassword' --command -- psql testdb --host <RDS-ENDPOINT> -U flyteadmin -d flyteadmin -p 5432
```

```sh
kubectl get pods -n kube-system | grep coredns
```


# Debugging CoreDNS



```
kubectl describe coredns-7975d6fb9b-cnwqf -n kube-system
```
