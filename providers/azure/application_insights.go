package azure

import (
	"context"
	"log"

	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/authentication"

	"github.com/Azure/azure-sdk-for-go/services/appinsights/mgmt/2020-02-02/insights"
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
)

type AppInsightsGenerator struct {
	AzureService
}

func (g AppInsightsGenerator) listServices() ([]terraformutils.Resource, error) {
	var resources []terraformutils.Resource
	ctx := context.Background()

	appInsightServiceClient := insights.NewComponentsClient(g.Args["config"].(authentication.Config).SubscriptionID)
	appInsightServiceClient.Authorizer = g.Args["authorizer"].(autorest.Authorizer)
	var (
		serviceIterator insights.ApplicationInsightsComponentListResultIterator
		err             error
	)
	if rg := g.Args["resource_group"].(string); rg != "" {
		serviceIterator, err = appInsightServiceClient.ListByResourceGroupComplete(ctx, rg)
	} else {
		serviceIterator, err = appInsightServiceClient.ListComplete(ctx)
	}
	if err != nil {
		return nil, err
	}
	for serviceIterator.NotDone() {
		service := serviceIterator.Value()
		resources = append(resources, terraformutils.NewSimpleResource(
			*service.ID,
			*service.Name,
			"azurerm_application_insights",
			g.ProviderName,
			[]string{}))

		if err := serviceIterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return resources, err
		}
	}

	return resources, nil
}

func (g *AppInsightsGenerator) InitResources() error {
	resources, err := g.listServices()
	if err != nil {
		return err
	}

	g.Resources = append(g.Resources, resources...)

	return nil
}
