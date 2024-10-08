# How to deploy AGIC via Helm using Workload Identity

> **_NOTE:_** [Application Gateway for Containers](https://aka.ms/agc) has been released, which introduces numerous performance, resilience, and feature changes. Please consider leveraging Application Gateway for Containers for your next deployment.

This assumes you have an existing Application Gateway. If not, you can create it with command:

```bash
az network application-gateway create -g myResourceGroup -n myApplicationGateway --sku Standard_v2 --public-ip-address myPublicIP --vnet-name myVnet --subnet mySubnet --priority 100
```

## 1. Set environment variables

```bash
export RESOURCE_GROUP="myResourceGroup"
export APPLICATION_GATEWAY_NAME="myApplicationGateway"
export USER_ASSIGNED_IDENTITY_NAME="myIdentity"
export FEDERATED_IDENTITY_CREDENTIAL_NAME="myFedIdentity"
```

## 2. Create resource group, AKS cluster and identity

```bash
az group create --name "${RESOURCE_GROUP}"  --location eastus
az aks create -g "${RESOURCE_GROUP}" -n myAKSCluster --node-count 1 --enable-oidc-issuer --enable-workload-identity 
az identity create --name "${USER_ASSIGNED_IDENTITY_NAME}" --resource-group "${RESOURCE_GROUP}" 
```

## 3. Export the oidcIssuerProfile.issuerUrl

```bash
export AKS_OIDC_ISSUER="$(az aks show -n myAKSCluster -g "${RESOURCE_GROUP}" --query "oidcIssuerProfile.issuerUrl" -otsv)"
```

## 4. Create federated identity credential

**Note**: the name of the service account that gets created after the helm installation is “ingress-azure” and the following command assumes it will be deployed in “default” namespace. Please change the namespace name in the next command if you deploy the AGIC related Kubernetes resources in other namespace.

```bash
az identity federated-credential create --name ${FEDERATED_IDENTITY_CREDENTIAL_NAME} --identity-name ${USER_ASSIGNED_IDENTITY_NAME} --resource-group ${RESOURCE_GROUP} --issuer ${AKS_OIDC_ISSUER} --subject system:serviceaccount:default:ingress-azure
```

## 5. Obtain the ClientID of the identity created before that is needed for the next step

```bash
az identity show --resource-group "${RESOURCE_GROUP}" --name "${USER_ASSIGNED_IDENTITY_NAME}" --query 'clientId' -otsv
```

## 6. Export the Application Gateway resource ID

```bash
export APP_GW_ID="$(az network application-gateway show --name "${APPLICATION_GATEWAY_NAME}"  --resource-group "${RESOURCE_GROUP}"  --query 'id' --output tsv)"
```

## 7. Add Contributor role for the identity over the Application Gateway

```bash
az role assignment create --assignee <identityClientID> --scope "${APP_GW_ID}" --role Contributor
```

## 8. In helm-config.yaml specify

```yaml
armAuth:
    type: workloadIdentity
    identityClientID: <identityClientID>
```

## 9. Get the AKS cluster credentials

```bash
az aks get-credentials -g "${RESOURCE_GROUP}" -n myAKSCluster
```

## 10. Install the helm chart

```bash
helm install ingress-azure \
  -f helm-config.yaml \
  oci://mcr.microsoft.com/azure-application-gateway/charts/ingress-azure \
  --version 1.7.5
```
