# Installation Guide

This guide provides detailed instructions for installing Cloud Sentinel K8s in a Kubernetes environment.

## Prerequisites

- `kubectl` with cluster administrator privileges
- Helm v3 (recommended for Helm installation)
- MySQL/PostgreSQL database, or local storage for sqlite

## Installation Methods

### Method 1: Helm Chart (Recommended)

Using Helm provides flexibility for configuration and upgrades:

```bash
# Add Cloud Sentinel K8s repository
helm repo add cloud-sentinel-k8s https://pixelvide.github.io/cloud-sentinel-k8s

# Update repository information
helm repo update

# Install with default configuration
helm install cloud-sentinel-k8s cloud-sentinel-k8s/cloud-sentinel-k8s -n cloud-sentinel-k8s-system --create-namespace
```

#### Custom Installation

You can adjust installation parameters by customizing the values file:

For complete configuration, refer to [Chart Values](../config/chart-values).

Install with custom values:

```bash
helm install cloud-sentinel-k8s cloud-sentinel-k8s/cloud-sentinel-k8s -n cloud-sentinel-k8s-system -f values.yaml
```

### Method 2: YAML Manifest

For quick deployment, you can directly apply the official installation YAML:

```bash
kubectl apply -f https://raw.githubusercontent.com/pixelvide/cloud-sentinel-k8s/main/deploy/install.yaml
```

This method will install Cloud Sentinel K8s with default configuration. For advanced customization, it's recommended to use the Helm Chart.

## Accessing Cloud Sentinel K8s

### Port Forwarding (Testing Environment)

During testing, you can quickly access Cloud Sentinel K8s through port forwarding:

```bash
kubectl port-forward -n cloud-sentinel-k8s-system svc/cloud-sentinel-k8s 8080:8080
```

### LoadBalancer Service

If the cluster supports LoadBalancer, you can directly expose the Cloud Sentinel K8s service:

```bash
kubectl patch svc cloud-sentinel-k8s -n cloud-sentinel-k8s-system -p '{"spec": {"type": "LoadBalancer"}}'
```

Get the assigned IP:

```bash
kubectl get svc cloud-sentinel-k8s -n cloud-sentinel-k8s-system
```

### Ingress (Recommended for Production)

For production environments, it's recommended to expose Cloud Sentinel K8s through an Ingress controller with TLS enabled:

::: warning
Cloud Sentinel K8s's log and web terminal features require websocket support.
Some Ingress controllers may require additional configuration to handle websockets correctly.
:::

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: cloud-sentinel-k8s
  namespace: cloud-sentinel-k8s-system
spec:
  ingressClassName: nginx
  rules:
    - host: cloud-sentinel-k8s.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: cloud-sentinel-k8s
                port:
                  number: 8080
  tls:
    - hosts:
        - cloud-sentinel-k8s.example.com
      secretName: cloud-sentinel-k8s-tls
```

## Serving under a subpath (basePath)

If you want to serve Cloud Sentinel K8s under a subpath (for example `https://example.com/cloud-sentinel-k8s`), use the Helm chart `basePath` value.

How to set it:

- In `values.yaml`:

```yaml
basePath: "/cloud-sentinel-k8s"
```

- Or with Helm CLI:

```fish
helm install cloud-sentinel-k8s cloud-sentinel-k8s/cloud-sentinel-k8s -n cloud-sentinel-k8s-system --create-namespace --set basePath="/cloud-sentinel-k8s"
```

Important notes:

- Ingress configuration: make sure your Ingress `paths` match the subpath and use a matching pathType (e.g., `Prefix`). Example:

```yaml
ingress:
  enabled: true
  hosts:
    - host: cloud-sentinel-k8s.example.com
      paths:
        - path: /cloud-sentinel-k8s
          pathType: Prefix
```

- OAuth / redirects: if you enable OAuth (or any external redirect flows), update the redirect URLs in your OAuth provider to include the base path, e.g. `https://cloud-sentinel-k8s.example.com/cloud-sentinel-k8s/oauth/callback`.
- Environment overrides: if you provide environment variables via `extraEnvs` or an existing secret, ensure `CLOUD_SENTINEL_K8S_BASE` is set consistently with the `basePath` value (otherwise behavior may differ).

## Verifying Installation

After installation, you can access the dashboard to verify that Cloud Sentinel K8s is deployed successfully. The expected interface is as follows:

::: tip
If you need to configure Cloud Sentinel K8s through environment variables, please refer to [Environment Variables](../config/env).
:::

![setup](../screenshots/setup.png)

![setup](../screenshots/setup2.png)

You can complete cluster setup according to the page prompts.

### Quick Setup with In-Cluster Mode

For the simplest setup, select **`in-cluster`** as the cluster type. This option automatically uses the service account credentials that Cloud Sentinel K8s is running with inside the cluster, requiring no additional configuration:

- **No kubeconfig needed**: Cloud Sentinel K8s will use its own service account to access the Kubernetes API
- **Automatic authentication**: Works out of the box with the default RBAC permissions
- **Ideal for single-cluster deployments**: Perfect when Cloud Sentinel K8s is managing the same cluster it's running in

This is the recommended option for getting started quickly, especially in development or when Cloud Sentinel K8s only needs to manage its own cluster.

## Uninstalling Cloud Sentinel K8s

### Helm Uninstall

```bash
helm uninstall cloud-sentinel-k8s -n cloud-sentinel-k8s-system
```

### YAML Uninstall

```bash
kubectl delete -f https://raw.githubusercontent.com/pixelvide/cloud-sentinel-k8s/main/deploy/install.yaml
```

## Next Steps

After Cloud Sentinel K8s installation is complete, you can continue with:

- [Adding Users](../config/user-management)
- [Configuring RBAC](../config/rbac-config)
- [Configuring OAuth Authentication](../config/oauth-setup)
- [Setting up Prometheus Monitoring](../config/prometheus-setup)
