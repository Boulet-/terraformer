// Copyright 2019 The Terraformer Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package azure

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/services/storage/mgmt/2019-04-01/storage"
	"github.com/Azure/go-autorest/autorest"
	"github.com/GoogleCloudPlatform/terraformer/terraformutils"
	"github.com/hashicorp/go-azure-helpers/authentication"
)

type StorageAccountGenerator struct {
	AzureService
}

func (g StorageAccountGenerator) createResourcesByResourceGroup(ctx context.Context, client storage.AccountsClient, rg string) ([]terraformutils.Resource, error) {
	accountListResult, err := client.ListByResourceGroup(ctx, rg)
	if err != nil {
		return nil, err
	}
	var resources []terraformutils.Resource
	if accounts := accountListResult.Value; accounts != nil {
		for _, account := range *accounts {
			resources = append(resources, terraformutils.NewSimpleResource(
				*account.ID,
				*account.Name,
				"azurerm_storage_account",
				"azurerm",
				[]string{}))
		}
	}
	return resources, nil
}
func (g StorageAccountGenerator) createResources(ctx context.Context, client storage.AccountsClient) ([]terraformutils.Resource, error) {
	accountListResultIterator, err := client.ListComplete(ctx)
	if err != nil {
		return nil, err
	}
	var resources []terraformutils.Resource
	for accountListResultIterator.NotDone() {
		account := accountListResultIterator.Value()
		resources = append(resources, terraformutils.NewSimpleResource(
			*account.ID,
			*account.Name,
			"azurerm_storage_account",
			"azurerm",
			[]string{}))
		if err := accountListResultIterator.Next(); err != nil {
			log.Println(err)
			return resources, err
		}
	}
	return resources, nil
}

func (g *StorageAccountGenerator) InitResources() error {
	ctx := context.Background()
	accountsClient := storage.NewAccountsClient(g.Args["config"].(authentication.Config).SubscriptionID)
	accountsClient.Authorizer = g.Args["authorizer"].(autorest.Authorizer)
	if rg := g.Args["resource_group"].(string); rg != "" {
		output, err := g.createResourcesByResourceGroup(ctx, accountsClient, rg)
		g.Resources = output
		return err
	}
	output, err := g.createResources(ctx, accountsClient)
	g.Resources = output
	return err
}

func (g *StorageAccountGenerator) PostConvertHook() error {
	for _, resource := range g.Resources {
		if resource.InstanceInfo.Type != "azurerm_storage_account" {
			continue
		}

		// Remove default value to 0 that are against the constraint [1-365]
		if resource.InstanceState.Attributes["queue_properties.0.hour_metrics.0.retention_policy_days"] == "0" {
			delete(resource.Item["queue_properties"].([]interface{})[0].(map[string]interface{})["hour_metrics"].([]interface{})[0].(map[string]interface{}), "retention_policy_days")
		}
		if resource.InstanceState.Attributes["queue_properties.0.logging.0.retention_policy_days"] == "0" {
			delete(resource.Item["queue_properties"].([]interface{})[0].(map[string]interface{})["logging"].([]interface{})[0].(map[string]interface{}), "retention_policy_days")
		}
		if resource.InstanceState.Attributes["queue_properties.0.minute_metrics.0.retention_policy_days"] == "0" {
			delete(resource.Item["queue_properties"].([]interface{})[0].(map[string]interface{})["minute_metrics"].([]interface{})[0].(map[string]interface{}), "retention_policy_days")
		}
	}
	return nil
}
