# mackerel-plugin-prometheus

Prometheus (https://prometheus.io/) custom metrics plugin for mackerel.io agent.

## Synopsis

```shell
PROMETHEUS_API_SERVER=<URL of Prometheus API Server> \
  mackerel-plugin-prometheus [-d] [-h] [-n] [-s=<Service Name>] -q=<Query String> -l=<Labels>
```

```shell
mackerel-plugin-prometheus -h
  -d  Debug mode
  -h  Show this help message
  -l string
      Labels to map on metric names
  -n  Not send metrics to Mackerel, just check metrics with -d
  -q string
      Query string to retrieve metrics
  -s string
      Service Name where to send metrics, or Host
```

## Example of mackerel-agent.conf

```
[plugin.metrics.prometheus]
command = "/path/to/mackerel-plugin-prometheus -q 'kube_pod_container_resource_requests_memory_bytes{job=\"kubernetes-pods\"}' -l app.container.memory"
```
