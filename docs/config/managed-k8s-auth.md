---
title: Managed Kubernetes Cluster Configuration
---

# Managed Kubernetes Cluster Configuration

## Problem Description

Managed Kubernetes clusters like AKS (Azure Kubernetes Service), EKS (Amazon Elastic Kubernetes Service), etc., typically use `exec` plugins in their default kubeconfig to dynamically obtain authentication credentials. For example:

- **AKS** uses the `kubelogin` command
- **EKS** uses the `aws` CLI
- **GKE** uses the `gcloud` command
- **GitLab Agent** uses the `glab` command

This authentication method works well in local client environments, but can be challenging in server-side environments like Cloud Sentinel K8s because:

1. These CLI tools may not be installed on the server
2. Even if installed, the server environment may not have the corresponding authentication configuration
3. Managing different user credentials in multi-tenant scenarios is difficult

Cloud Sentinel K8s provides two ways to solve this: **Managed Authentication Support** (for AWS and GitLab) and **Service Account Tokens** (for all others).

## Managed Authentication Support [NEW]

Cloud Sentinel K8s natively supports authentication for specific managed Kubernetes providers by securely managing your credentials and injecting them into the CLI tools.

### AWS EKS Authentication

For EKS clusters, Cloud Sentinel K8s supports authentication via `aws` or `aws-iam-authenticator`.

1. **Configure AWS Credentials**: Navigate to **Settings > AWS Settings** and paste your AWS credentials file content (typically found at `~/.aws/credentials`).
2. **Add Cluster**: Import your EKS kubeconfig. Cloud Sentinel K8s will detect the `aws` exec command.
3. **Secure Injection**: The system automatically injects the `AWS_SHARED_CREDENTIALS_FILE` environment variable for your requests, ensuring you only use your own credentials.

### GitLab Agent Authentication

For clusters managed via GitLab Agent, Cloud Sentinel K8s supports authentication using the `glab` CLI.

1. **Configure GitLab Token**: Navigate to **Settings > GitLab Settings**, add your GitLab host (e.g., `gitlab.com`), and provide a Personal Access Token (PAT).
2. **Validate**: Click **Validate** to initialize your `glab` session.
3. **Add Cluster**: Import your cluster kubeconfig that uses `glab` for authentication.
4. **Context Management**: Cloud Sentinel K8s automatically manages the `GLAB_CONFIG_DIR` to use your validated session.

> [!NOTE]
> Support for other managed providers like **AKS (Azure)** and **GKE (Google)** is coming soon.

---

## Alternative: Using Service Account Token

If your provider is not yet natively supported, or if you prefer a more generic approach, you can create a dedicated Service Account for Cloud Sentinel K8s and use its token for authentication.

Cloud Sentinel K8s provides a helper script for creation:

```sh
wget https://raw.githubusercontent.com/pixelvide/cloud-sentinel-k8s/refs/heads/main/scripts/generate-cloud-sentinel-k8s-kubeconfig.sh -O generate-cloud-sentinel-k8s-kubeconfig.sh
chmod +x generate-cloud-sentinel-k8s-kubeconfig.sh
./generate-cloud-sentinel-k8s-kubeconfig.sh
```

### Manual Steps:

1. **Create Service Account and RBAC permissions**:

```yaml
apiVersion: v1
kind: ServiceAccount
metadata:
  name: cloud-sentinel-k8s-admin
  namespace: kube-system
---
apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: cloud-sentinel-k8s-admin
roleRef:
  apiGroup: rbac.authorization.k8s.io
  kind: ClusterRole
  name: cluster-admin
subjects:
  - kind: ServiceAccount
    name: cloud-sentinel-k8s-admin
    namespace: kube-system
```

2. **Create Long-lived Token Secret (Kubernetes 1.24+)**:

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: cloud-sentinel-k8s-admin-token
  namespace: kube-system
  annotations:
    kubernetes.io/service-account.name: cloud-sentinel-k8s-admin
type: kubernetes.io/service-account-token
```

3. **Get token and cluster information**:

```bash
# Get token
TOKEN=$(kubectl get secret cloud-sentinel-k8s-admin-token -n kube-system -o jsonpath='{.data.token}' | base64 -d)

# Get CA certificate
CA_CERT=$(kubectl get secret cloud-sentinel-k8s-admin-token -n kube-system -o jsonpath='{.data.ca\.crt}')

# Get API Server address
API_SERVER=$(kubectl config view --minify -o jsonpath='{.clusters[0].cluster.server}')
```

4. **Generate kubeconfig**:

```bash
cat > cloud-sentinel-k8s-kubeconfig.yaml <<EOF
apiVersion: v1
kind: Config
clusters:
- cluster:
    certificate-authority-data: ${CA_CERT}
    server: ${API_SERVER}
  name: cloud-sentinel-k8s-cluster
contexts:
- context:
    cluster: cloud-sentinel-k8s-cluster
    user: cloud-sentinel-k8s-admin
  name: cloud-sentinel-k8s-context
current-context: cloud-sentinel-k8s-context
users:
- name: cloud-sentinel-k8s-admin
  user:
    token: ${TOKEN}
EOF
```

## Related Documentation

- [Kubernetes Service Account Tokens](https://kubernetes.io/docs/reference/access-authn-authz/service-accounts-admin/)
- [AKS Authentication](https://learn.microsoft.com/en-us/azure/aks/control-kubeconfig-access)
- [EKS Authentication](https://docs.aws.amazon.com/eks/latest/userguide/cluster-auth.html)
- [GKE Authentication](https://cloud.google.com/kubernetes-engine/docs/how-to/api-server-authentication)
