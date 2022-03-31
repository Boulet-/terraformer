package azure

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/profiles/latest/servicebus/mgmt/servicebus"
	"github.com/Azure/go-autorest/autorest"
	"github.com/hashicorp/go-azure-helpers/authentication"

	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
)

type ServiceBusGenerator struct {
	AzureService
}

func (g ServiceBusGenerator) listNamespaces() ([]servicebus.SBNamespace, error) {
	var resources []servicebus.SBNamespace
	ctx := context.Background()

	namespaceClient := servicebus.NewNamespacesClient(g.Args["config"].(authentication.Config).SubscriptionID)
	namespaceClient.Authorizer = g.Args["authorizer"].(autorest.Authorizer)
	var (
		serviceIterator servicebus.SBNamespaceListResultIterator
		err             error
	)
	if rg := g.Args["resource_group"].(string); rg != "" {
		serviceIterator, err = namespaceClient.ListByResourceGroupComplete(ctx, rg)
	} else {
		serviceIterator, err = namespaceClient.ListComplete(ctx)
	}
	if err != nil {
		return nil, err
	}
	for serviceIterator.NotDone() {
		service := serviceIterator.Value()
		resources = append(resources, service)

		if err := serviceIterator.NextWithContext(ctx); err != nil {
			log.Println(err)
			return resources, err
		}
	}

	return resources, nil
}

// func (g ServiceBusGenerator) listQueues(namespace *servicebus.SBNamespace, namespaceRG *ResourceID) ([]terraformutils.Resource, error) {
// 	var resources []terraformutils.Resource
// 	ctx := context.Background()

// 	queueClient := servicebus.NewQueuesClient(g.Args["config"].(authentication.Config).SubscriptionID)
// 	queueClient.Authorizer = g.Args["authorizer"].(autorest.Authorizer)
// 	var (
// 		serviceIterator servicebus.SBQueueListResultIterator
// 		err             error
// 	)
// 	serviceIterator, err = queueClient.ListByNamespaceComplete(ctx, namespaceRG.ResourceGroup, *namespace.Name, nil, nil)

// 	if err != nil {
// 		return nil, err
// 	}
// 	for serviceIterator.NotDone() {
// 		service := serviceIterator.Value()
// 		resources = append(resources, terraformutils.NewSimpleResource(
// 			*service.ID,
// 			*service.Name,
// 			"azurerm_servicebus_queue",
// 			g.ProviderName,
// 			[]string{}))

// 		if err := serviceIterator.NextWithContext(ctx); err != nil {
// 			log.Println(err)
// 			return resources, err
// 		}
// 	}

// 	return resources, nil
// }

func (g *ServiceBusGenerator) InitResources() error {
	namespaces, err := g.listNamespaces()
	if err != nil {
		return err
	}

	for _, namespace := range namespaces {
		g.Resources = append(g.Resources, terraformutils.NewSimpleResource(
			*namespace.ID,
			*namespace.Name,
			"azurerm_servicebus_queue",
			g.ProviderName,
			[]string{}))
		// namespaceRg, err := ParseAzureResourceID(*namespace.ID)
		// if err != nil {
		// 	return err
		// }

		// queues, err := g.listQueues(&namespace, namespaceRg)
		// if err != nil {
		// 	return err
		// }
		// g.Resources = append(g.Resources, queues...)
	}
	return nil
}
