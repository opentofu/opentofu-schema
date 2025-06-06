// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"strings"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/addr"
	"github.com/opentofu/opentofu-schema/internal/schema/backends"
	"github.com/opentofu/opentofu-schema/module"
	tfaddr "github.com/opentofu/registry-address"
	"github.com/zclconf/go-cty/cty"
)

const remoteStateDsName = "terraform_remote_state"

func isRemoteStateDataSource(pAddr tfaddr.Provider, dsName string) bool {
	return (pAddr.Equals(addr.NewBuiltInProvider("terraform")) ||
		pAddr.Equals(addr.NewDefaultProvider("terraform")) ||
		pAddr.Equals(addr.NewLegacyProvider("terraform"))) &&
		dsName == remoteStateDsName
}

func (sm *SchemaMerger) dependentBodyForRemoteStateDataSource(originalBodySchema *schema.BodySchema, providerAddr lang.Address, localRef module.ProviderRef) map[schema.SchemaKey]*schema.BodySchema {
	m := make(map[schema.SchemaKey]*schema.BodySchema, 0)
	backendsAsCons := backends.ConfigsAsObjectConstraint(sm.tofuVersion)

	for backendType, objConstraint := range backendsAsCons {
		depKeys := schema.DependencyKeys{
			Labels: []schema.LabelDependent{
				{Index: 0, Value: remoteStateDsName},
			},
			Attributes: []schema.AttributeDependent{
				{
					Name: "provider",
					Expr: schema.ExpressionValue{
						Address: providerAddr,
					},
				},
				{
					Name: "backend",
					Expr: schema.ExpressionValue{
						Static: cty.StringVal(backendType),
					},
				},
			},
		}

		dsSchema := originalBodySchema.Copy()
		dsSchema.Attributes["backend"] = &schema.AttributeSchema{
			Constraint:             backends.BackendTypesAsOneOfConstraint(sm.tofuVersion),
			IsRequired:             true,
			SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
		}
		dsSchema.Attributes["config"] = &schema.AttributeSchema{
			Constraint: objConstraint,
			IsOptional: true,
		}

		m[schema.NewSchemaKey(depKeys)] = dsSchema

		// No explicit association is required
		// if the resource prefix matches provider name
		if strings.HasPrefix(remoteStateDsName, localRef.LocalName+"_") {
			depKeys := schema.DependencyKeys{
				Labels: []schema.LabelDependent{
					{Index: 0, Value: remoteStateDsName},
				},
				Attributes: []schema.AttributeDependent{
					{
						Name: "backend",
						Expr: schema.ExpressionValue{
							Static: cty.StringVal(backendType),
						},
					},
				},
			}
			m[schema.NewSchemaKey(depKeys)] = dsSchema
		}
	}

	return m
}
