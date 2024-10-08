## Adding Health Probes to your service

.. note::
    [Application Gateway for Containers](https://aka.ms/agc) has been released, which introduces numerous performance, resilience, and feature changes. Please consider leveraging Application Gateway for Containers for your next deployment.

    Details on custom health probes in Application Gateway for Containers [may be found here](https://learn.microsoft.com/azure/application-gateway/for-containers/custom-health-probe).

By default, Ingress controller will provision an HTTP GET probe for the exposed pods.
The probe properties can be customized by adding a [Readiness or Liveness Probe](https://kubernetes.io/docs/tasks/configure-pod-container/configure-liveness-readiness-probes/) to your `deployment`/`pod` spec.

### With `readinessProbe` or `livenessProbe`

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: aspnetapp
spec:
  replicas: 3
  template:
    metadata:
      labels:
        service: site
    spec:
      containers:
      - name: aspnetapp
        image: mcr.microsoft.com/dotnet/samples:aspnetapp
        imagePullPolicy: IfNotPresent
        ports:
        - containerPort: 80
        readinessProbe:
          httpGet:
            path: /
            port: 80
          periodSeconds: 3
          timeoutSeconds: 1
```

Kubernetes API Reference:

* [Container Probes](https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#container-probes)
* [HttpGet Action](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.14/#httpgetaction-v1-core)

**Note:**

1. `readinessProbe` and `livenessProbe` are supported when configured with `httpGet`.
1. Probing on a port other than the one exposed on the pod is currently not supported.
1. `HttpHeaders`, `InitialDelaySeconds`, `SuccessThreshold` are not supported.

### Without `readinessProbe` or `livenessProbe`

If the above probes are not provided, then Ingress Controller make an assumption that the service is reachable on `Path` specified for `backend-path-prefix` annotation or the `path` specified in the `ingress` definition for the service.

### Default Values for Health Probe

For any property that can not be inferred by the readiness/liveness probe, Default values are set.

| Application Gateway Probe Property | Default Value |
|-|-|
| `Path` | / |
| `Host` | localhost |
| `Protocol` | HTTP |
| `Timeout` | 30 |
| `Interval` | 30 |
| `UnhealthyThreshold` | 3 |
