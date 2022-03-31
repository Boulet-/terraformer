package azure

import (
	"context"
	"log"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/authentication"

	"github.com/Azure/azure-sdk-for-go/services/apimanagement/mgmt/2021-08-01/apimanagement"
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
)

type APIManagementGenerator struct {
	AzureService
}

func (g APIManagementGenerator) listServices() ([]terraformutils.Resource, error) {
	var resources []terraformutils.Resource
	ctx := context.Background()

	APIManagementServiceClient := apimanagement.NewServiceClient(g.Args["config"].(authentication.Config).SubscriptionID)
	APIManagementServiceClient.Authorizer = g.Args["authorizer"].(autorest.Authorizer)
	var (
		serviceIterator apimanagement.ServiceListResultIterator
		err             error
	)
	if rg := g.Args["resource_group"].(string); rg != "" {
		serviceIterator, err = APIManagementServiceClient.ListByResourceGroupComplete(ctx, rg)
	} else {
		serviceIterator, err = APIManagementServiceClient.ListComplete(ctx)
	}
	if err != nil {
		return nil, err
	}
	for serviceIterator.NotDone() {
		service := serviceIterator.Value()
		resources = append(resources, terraformutils.NewSimpleResource(
			*service.ID,
			*service.Name,
			"azurerm_api_management",
			g.ProviderName,
			[]string{}))

		if err := serviceIterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return resources, err
		}
	}

	return resources, nil
}

func (g *APIManagementGenerator) InitResources() error {
	resources, err := g.listServices()
	if err != nil {
		return err
	}

	g.Resources = append(g.Resources, resources...)

	return nil
}
