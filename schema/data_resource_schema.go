// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0
package schema

import (
	"fmt"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/schema/backends"
	tfmod "github.com/opentofu/opentofu-schema/module"
	tfaddr "github.com/opentofu/registry-address"
)

func (sm *SchemaMerger) mergeDataSourceSchema(bSchema *schema.BodySchema, dsName string, dsSchema *schema.BodySchema, providerAddr tfaddr.Provider, localProviderAddr lang.Address, localRef tfmod.ProviderRef) {
	depKeys := schema.DependencyKeys{
		Labels: []schema.LabelDependent{
			{Index: 0, Value: dsName},
		},
		Attributes: []schema.AttributeDependent{
			{
				Name: "provider",
				Expr: schema.ExpressionValue{
					Address: localProviderAddr,
				},
			},
		},
	}

	// Add backend-related core bits of schema
	if isRemoteStateDataSource(providerAddr, dsName) {
		remoteStateDs := dsSchema.Copy()

		addRemoteStateDataSourceDocsURL(providerAddr, dsName, remoteStateDs)

		remoteStateDs.Attributes["backend"].IsDepKey = true
		remoteStateDs.Attributes["backend"].SemanticTokenModifiers = lang.SemanticTokenModifiers{lang.TokenModifierDependent}
		remoteStateDs.Attributes["backend"].Constraint = backends.BackendTypesAsOneOfConstraint(sm.tofuVersion)
		delete(remoteStateDs.Attributes, "config")

		depBodies := sm.dependentBodyForRemoteStateDataSource(remoteStateDs, localProviderAddr, localRef)
		for key, depBody := range depBodies {
			bSchema.Blocks["data"].DependentBody[key] = depBody
			if _, ok := bSchema.Blocks["check"]; ok {
				bSchema.Blocks["check"].Body.Blocks["data"].DependentBody[key] = depBody
			}
		}

		dsSchema = remoteStateDs
	} else {
		addDataSourceDocsURL(providerAddr, dsName, dsSchema)
	}

	bSchema.Blocks["data"].DependentBody[schema.NewSchemaKey(depKeys)] = dsSchema

	if _, ok := bSchema.Blocks["check"]; ok {
		bSchema.Blocks["check"].Body.Blocks["data"].DependentBody[schema.NewSchemaKey(depKeys)] = dsSchema
	}

	// No explicit association is required
	// if the resource prefix matches provider name
	if typeBelongsToProvider(dsName, localRef) {
		depKeys := schema.DependencyKeys{
			Labels: []schema.LabelDependent{
				{Index: 0, Value: dsName},
			},
		}
		bSchema.Blocks["data"].DependentBody[schema.NewSchemaKey(depKeys)] = dsSchema
		if _, ok := bSchema.Blocks["check"]; ok {
			bSchema.Blocks["check"].Body.Blocks["data"].DependentBody[schema.NewSchemaKey(depKeys)] = dsSchema
		}
	}
}

func addDataSourceDocsURL(providerAddr tfaddr.Provider, dsName string, dsSchema *schema.BodySchema) {
	namespace := providerAddr.Namespace
	if providerAddr.IsLegacy() {
		// When namespaces are legacy, we assume their namespace is hashicorp
		namespace = "hashicorp"
	}

	if namespace == "" {
		return
	}

	// In OpenTofu's Search Registry, we don't save the data source prefix on the URL, example:
	// random_uuid becomes uuid on the URL
	registryDataSourceName := dsName
	if len(providerAddr.Type)+1 <= len(dsName) {
		registryDataSourceName = dsName[len(providerAddr.Type)+1:]
	}

	docsURL := fmt.Sprintf("https://search.opentofu.org/provider/%s/%s/latest/docs/datasources/%s", namespace, providerAddr.Type, registryDataSourceName)
	dsSchema.DocsLink = &schema.DocsLink{
		URL:     docsURL,
		Tooltip: fmt.Sprintf("%s/%s/%s Documentation", namespace, providerAddr.Type, dsName),
	}
	dsSchema.HoverURL = docsURL
}

func addRemoteStateDataSourceDocsURL(providerAddr tfaddr.Provider, dsName string, dsSchema *schema.BodySchema) {
	docsURL := "https://opentofu.org/docs/language/state/remote-state-data/"
	dsSchema.DocsLink = &schema.DocsLink{
		URL:     docsURL,
		Tooltip: fmt.Sprintf("%s Documentation", dsName),
	}
	dsSchema.HoverURL = docsURL
}
