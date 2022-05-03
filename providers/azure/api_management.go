package azure

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
)

type APIManagementGenerator struct {
	AzureService
}

func (g *APIManagementGenerator) listServices() ([]apimanagement.ServiceResource, error) {
	subscriptionID, resourceGroup, authorizer := g.getClientArgs()
	client := apimanagement.NewServiceClient(subscriptionID)
	client.Authorizer = authorizer
	var (
		iterator apimanagement.ServiceListResultIterator
		err      error
	)
	ctx := context.Background()

	if resourceGroup != "" {
		iterator, err = client.ListByResourceGroupComplete(ctx, resourceGroup)
	} else {
		iterator, err = client.ListComplete(ctx)
	}
	if err != nil {
		return nil, err
	}

	var resources []apimanagement.ServiceResource
	for iterator.NotDone() {
		item := iterator.Value()
		resources = append(resources, item)
		if err := iterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return resources, err
		}
	}

	return resources, nil
}

func (g *APIManagementGenerator) appendAPIs(service *apimanagement.ServiceResource, serviceRg *ResourceID) error {
	subscriptionID, _, authorizer := g.getClientArgs()
	client := apimanagement.NewAPIClient(subscriptionID)
	client.Authorizer = authorizer

	ctx := context.Background()
	iterator, err := client.ListByServiceComplete(ctx, serviceRg.ResourceGroup, *service.Name, "", nil, nil, "", nil)
	if err != nil {
		return err
	}
	for iterator.NotDone() {
		item := iterator.Value()

		g.AppendSimpleResource(*item.ID, *item.Name, "azurerm_api_management_api")
		err = g.appendAPIDiagnostics(service, serviceRg, *item.Name)
		if err != nil {
			return err
		}
		// err = g.appendAPIOperations(service, serviceRg, *item.Name)
		// if err != nil {
		// 	return err
		// }
		// err = g.appendAPIPolicies(service, serviceRg, *item.Name)
		// if err != nil {
		// 	return err
		// }
		err = g.appendAPIReleases(service, serviceRg, *item.Name)
		if err != nil {
			return err
		}
		err = g.appendAPISchemas(service, serviceRg, *item.Name)
		if err != nil {
			return err
		}
		err = g.appendAPITags(service, serviceRg, *item.Name)
		if err != nil {
			return err
		}

		if err := iterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
func (g *APIManagementGenerator) appendAPIDiagnostics(service *apimanagement.ServiceResource, serviceRg *ResourceID, apiid string) error {
	subscriptionID, _, authorizer := g.getClientArgs()
	client := apimanagement.NewAPIDiagnosticClient(subscriptionID)
	client.Authorizer = authorizer

	ctx := context.Background()
	iterator, err := client.ListByServiceComplete(ctx, serviceRg.ResourceGroup, *service.Name, apiid, "", nil, nil)
	if err != nil {
		return err
	}
	for iterator.NotDone() {
		item := iterator.Value()
		g.AppendSimpleResource(*item.ID, *item.Name, "azurerm_api_management_api_diagnostic")

		if err := iterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

// func (g *APIManagementGenerator) appendAPIOperations(service *apimanagement.ServiceResource, serviceRg *ResourceID, apiid string) error {
// 	subscriptionID, _, authorizer := g.getClientArgs()
// 	client := apimanagement.NewAPIOperationClient(subscriptionID)
// 	client.Authorizer = authorizer

// 	ctx := context.Background()
// 	iterator, err := client.ListByAPIComplete(ctx, serviceRg.ResourceGroup, *service.Name, apiid, "", nil, nil, "")
// 	if err != nil {
// 		return err
// 	}
// 	for iterator.NotDone() {
// 		item := iterator.Value()
// 		g.AppendSimpleResource(*item.ID, fmt.Sprintf("%s-%s", apiid, *item.Name), "azurerm_api_management_api_operation")
// 		// err = g.appendAPIOperationPolicies(service, serviceRg, apiid, *item.Name)
// 		// if err != nil {
// 		// 	return err
// 		// }

// 		if err := iterator.NextWithContext(ctx); err != nil {
// 			log.Println(err)
// 			return err
// 		}
// 	}

// 	return nil
// }
// func (g *APIManagementGenerator) appendAPIOperationPolicies(service *apimanagement.ServiceResource, serviceRg *ResourceID, appid string, operationID string) error {
// 	subscriptionID, _, authorizer := g.getClientArgs()
// 	client := apimanagement.NewAPIOperationPolicyClient(subscriptionID)
// 	client.Authorizer = authorizer

// 	ctx := context.Background()
// 	collection, err := client.ListByOperation(ctx, serviceRg.ResourceGroup, *service.Name, appid, operationID)
// 	if err != nil {
// 		return err
// 	}

// 	for _, item := range *collection.Value {
// 		g.AppendSimpleResource(*item.ID, *item.Name, "azurerm_api_management_api_operation_policy")
// 	}

// 	return nil
// }
// func (g *APIManagementGenerator) appendAPIPolicies(service *apimanagement.ServiceResource, serviceRg *ResourceID, apiid string) error {
// 	subscriptionID, _, authorizer := g.getClientArgs()
// 	client := apimanagement.NewAPIPolicyClient(subscriptionID)
// 	client.Authorizer = authorizer

// 	ctx := context.Background()
// 	collection, err := client.ListByAPI(ctx, serviceRg.ResourceGroup, *service.Name, apiid)
// 	if err != nil {
// 		return err
// 	}
// 	for _, item := range *collection.Value {
// 		fmt.Println(*item.Name)
// 		g.AppendSimpleResource(*item.ID, *item.Name, "azurerm_api_management_api_policy")
// 	}

// 	return nil
// }
func (g *APIManagementGenerator) appendAPIReleases(service *apimanagement.ServiceResource, serviceRg *ResourceID, apiid string) error {
	subscriptionID, _, authorizer := g.getClientArgs()
	client := apimanagement.NewAPIReleaseClient(subscriptionID)
	client.Authorizer = authorizer

	ctx := context.Background()
	iterator, err := client.ListByServiceComplete(ctx, serviceRg.ResourceGroup, *service.Name, apiid, "", nil, nil)
	if err != nil {
		return err
	}
	for iterator.NotDone() {
		item := iterator.Value()
		g.AppendSimpleResource(*item.ID, *item.Name, "azurerm_api_management_api_release")

		if err := iterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
func (g *APIManagementGenerator) appendAPISchemas(service *apimanagement.ServiceResource, serviceRg *ResourceID, apiid string) error {
	subscriptionID, _, authorizer := g.getClientArgs()
	client := apimanagement.NewAPISchemaClient(subscriptionID)
	client.Authorizer = authorizer

	ctx := context.Background()
	iterator, err := client.ListByAPIComplete(ctx, serviceRg.ResourceGroup, *service.Name, apiid, "", nil, nil)
	if err != nil {
		return err
	}
	for iterator.NotDone() {
		item := iterator.Value()
		g.AppendSimpleResource(*item.ID, *item.Name, "azurerm_api_management_api_schema")

		if err := iterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
func (g *APIManagementGenerator) appendAPITags(service *apimanagement.ServiceResource, serviceRg *ResourceID, apiid string) error {
	subscriptionID, _, authorizer := g.getClientArgs()
	client := apimanagement.NewAPITagDescriptionClient(subscriptionID)
	client.Authorizer = authorizer

	ctx := context.Background()
	iterator, err := client.ListByServiceComplete(ctx, serviceRg.ResourceGroup, *service.Name, apiid, "", nil, nil)
	if err != nil {
		return err
	}
	for iterator.NotDone() {
		item := iterator.Value()
		g.AppendSimpleResource(*item.ID, *item.Name, "azurerm_api_management_api_tag")

		if err := iterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (g *APIManagementGenerator) appendBackends(service *apimanagement.ServiceResource, serviceRg *ResourceID) error {
	subscriptionID, _, authorizer := g.getClientArgs()
	client := apimanagement.NewBackendClient(subscriptionID)
	client.Authorizer = authorizer

	ctx := context.Background()
	iterator, err := client.ListByServiceComplete(ctx, serviceRg.ResourceGroup, *service.Name, "", nil, nil)
	if err != nil {
		return err
	}
	for iterator.NotDone() {
		item := iterator.Value()
		g.AppendSimpleResource(*item.ID, *item.Name, "azurerm_api_management_backend")

		if err := iterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (g *APIManagementGenerator) appendGateways(service *apimanagement.ServiceResource, serviceRg *ResourceID) error {
	subscriptionID, _, authorizer := g.getClientArgs()
	client := apimanagement.NewBackendClient(subscriptionID)
	client.Authorizer = authorizer

	ctx := context.Background()
	iterator, err := client.ListByServiceComplete(ctx, serviceRg.ResourceGroup, *service.Name, "", nil, nil)
	if err != nil {
		return err
	}
	for iterator.NotDone() {
		item := iterator.Value()
		g.AppendSimpleResource(*item.ID, *item.Name, "azurerm_api_management_gateway")
		err = g.appendGatewayAPIs(service, serviceRg, *item.Name)
		if err != nil {
			return err
		}

		if err := iterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}
func (g *APIManagementGenerator) appendGatewayAPIs(service *apimanagement.ServiceResource, serviceRg *ResourceID, gatewayid string) error {
	subscriptionID, _, authorizer := g.getClientArgs()
	client := apimanagement.NewGatewayAPIClient(subscriptionID)
	client.Authorizer = authorizer

	ctx := context.Background()
	iterator, err := client.ListByServiceComplete(ctx, serviceRg.ResourceGroup, *service.Name, gatewayid, "", nil, nil)
	if err != nil {
		return err
	}
	for iterator.NotDone() {
		item := iterator.Value()
		g.AppendSimpleResource(*item.ID, *item.Name, "azurerm_api_management_gateway_api")

		if err := iterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

func (g *APIManagementGenerator) InitResources() error {
	services, err := g.listServices()
	if err != nil {
		return err
	}

	for _, service := range services {
		g.AppendSimpleResource(*service.ID, *service.Name, "azurerm_api_management")
		serviceRg, err := ParseAzureResourceID(*service.ID)
		if err != nil {
			return err
		}

		err = g.appendAPIs(&service, serviceRg)
		if err != nil {
			return err
		}
		err = g.appendBackends(&service, serviceRg)
		if err != nil {
			return err
		}
		err = g.appendGateways(&service, serviceRg)
		if err != nil {
			return err
		}
	}
	return nil
}
