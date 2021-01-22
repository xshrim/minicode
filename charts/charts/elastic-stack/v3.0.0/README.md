# elastic-stack Helm Chart

This chart installs an elasticsearch cluster with filebeat and kibana by default.
You can optionally enable elasticsearch-curator and disable kibana.

It is a fork of the original elastic-stack chart https://github.com/helm/charts/tree/master/stable/elastic-stack
that is now deprecated.

## Prerequisites Details

* Kubernetes 1.10+
* PV dynamic provisioning support on the underlying infrastructure

## Chart Details

By default, this chart will deploy a 3-node elasticsearch cluster.
Each node requests 1 CPU and 2GB ram and can't be collocated.
This can be configured to suit your needs.

## Installing the Chart

To install the chart with the release name `my-release`:

```bash
$ helm repo add gitlab https://charts.gitlab.io
$ helm install --name my-release gitlab/elastic-stack
```

## Deleting the Charts

Delete the Helm deployment as normal

```
$ helm delete my-release
```

Deletion of the StatefulSet doesn't cascade to deleting associated PVCs. To delete them:

```
$ kubectl delete pvc --selector app=elastic-stack-elasticsearch-master
```

## Configuration

Each requirement is configured with the options provided by that Chart.
Please consult the relevant charts for their configuration options.
