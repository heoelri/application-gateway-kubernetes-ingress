// -------------------------------------------------------------------------------------------
// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// --------------------------------------------------------------------------------------------

package azure

import (
	"fmt"
	"regexp"
	"strings"

	"k8s.io/klog/v2"

	"github.com/Azure/application-gateway-kubernetes-ingress/pkg/controllererrors"
)

type (
	// SubscriptionID is the subscription of the resource in the resourceID
	SubscriptionID string

	// ResourceGroup is the resource group in which resource is deployed in the resourceID
	ResourceGroup string

	// ResourceName is the resource name in the resourceID
	ResourceName string
)

var (
	operationIDRegex = regexp.MustCompile(`/operations/(.+)\?api-version`)
)

// ParseResourceID gets subscriptionId, resource group, resource name from resourceID
func ParseResourceID(ID string) (SubscriptionID, ResourceGroup, ResourceName) {
	split := strings.Split(ID, "/")
	if len(split) < 9 {
		klog.Errorf("resourceID %s is invalid. There should be atleast 9 segments in resourceID", ID)
		return "", "", ""
	}

	return SubscriptionID(split[2]), ResourceGroup(split[4]), ResourceName(split[8])
}

// ParseSubResourceID gets subscriptionId, resource group, resource name, sub resource name from resourceID
func ParseSubResourceID(ID string) (SubscriptionID, ResourceGroup, ResourceName, ResourceName) {
	split := strings.Split(ID, "/")
	if len(split) < 11 {
		klog.Errorf("resourceID %s is invalid. There should be atleast 9 segments in resourceID", ID)
		return "", "", "", ""
	}

	return SubscriptionID(split[2]), ResourceGroup(split[4]), ResourceName(split[8]), ResourceName(split[10])
}

// ResourceID generates a resource id
func ResourceID(subscriptionID SubscriptionID, resourceGroup ResourceGroup, provider string, resourceKind string, resourcePath string) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s", subscriptionID, resourceGroup, provider, resourceKind, resourcePath)
}

// RouteTableID generates a route table resource id
func RouteTableID(subscriptionID SubscriptionID, resourceGroup ResourceGroup, routeTableName ResourceName) string {
	return ResourceID(subscriptionID, resourceGroup, "Microsoft.Network", "routeTables", string(routeTableName))
}

// ApplicationGatewayID generates a application gateway resource id
func ApplicationGatewayID(subscriptionID SubscriptionID, resourceGroup ResourceGroup, applicationGatewayName ResourceName) string {
	return ResourceID(subscriptionID, resourceGroup, "Microsoft.Network", "applicationGateways", string(applicationGatewayName))
}

// ResourceGroupID generates a resource group resource id
func ResourceGroupID(subscriptionID SubscriptionID, resourceGroup ResourceGroup) string {
	return fmt.Sprintf("/subscriptions/%s/resourceGroups/%s", subscriptionID, resourceGroup)
}

// ConvertToClusterResourceGroup converts infra resource group to aks cluster ID
func ConvertToClusterResourceGroup(subscriptionID SubscriptionID, resourceGroup ResourceGroup, err error) (string, error) {
	if err != nil {
		return "", err
	}

	split := strings.Split(string(resourceGroup), "_")
	if len(split) != 4 || strings.ToUpper(split[0]) != "MC" {
		err := controllererrors.NewErrorf(
			controllererrors.ErrorMissingResourceGroup,
			"infrastructure resource group name: %s is expected to be of format MC_ResourceGroup_ResourceName_Location", string(resourceGroup),
		)
		klog.V(3).Info(err.Error())
		return "", err
	}

	return fmt.Sprintf("/subscriptions/%s/resourcegroups/%s/providers/Microsoft.ContainerService/managedClusters/%s", subscriptionID, split[1], split[2]), nil
}

// GetOperationIDFromPollingURL extracts operationID from pollingURL
func GetOperationIDFromPollingURL(pollingURL string) string {
	operationID := ""
	if pollingURL != "" {
		matchGroup := operationIDRegex.FindStringSubmatch(pollingURL)
		if len(matchGroup) == 2 {
			operationID = matchGroup[1]
		}
	}

	return operationID
}
