// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	tfmod "github.com/opentofu/opentofu-schema/module"
	tfaddr "github.com/opentofu/registry-address"
)

func (bs *SchemaMerger) mergeResourceSchema(bSchema *schema.BodySchema, rName string, rSchema *schema.BodySchema, providerAddr tfaddr.Provider, localProviderAddr lang.Address, localRef tfmod.ProviderRef) {
	depKeys := schema.DependencyKeys{
		Labels: []schema.LabelDependent{
			{Index: 0, Value: rName},
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

	// In OpenTofu's Search Registry, we don't save the resource prefix on the URL, example:
	// random_uuid becomes uuid on the URL
	if providerAddr.Namespace != "" {
		registryName := rName[len(providerAddr.Type)+1:]
		docsUrl := fmt.Sprintf("https://search.opentofu.org/provider/%s/%s/latest/docs/resources/%s", providerAddr.Namespace, providerAddr.Type, registryName)
		rSchema.DocsLink = &schema.DocsLink{
			URL:     docsUrl,
			Tooltip: fmt.Sprintf("%s/%s/%s Documentation", providerAddr.Namespace, providerAddr.Type, rName),
		}
	}

	bSchema.Blocks["resource"].DependentBody[schema.NewSchemaKey(depKeys)] = rSchema

	// No explicit association is required
	// if the resource prefix matches provider name
	if typeBelongsToProvider(rName, localRef) {
		depKeys := schema.DependencyKeys{
			Labels: []schema.LabelDependent{
				{Index: 0, Value: rName},
			},
		}
		bSchema.Blocks["resource"].DependentBody[schema.NewSchemaKey(depKeys)] = rSchema
	}
}
