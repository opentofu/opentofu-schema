// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"sort"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/schema/refscope"
	"github.com/opentofu/opentofu-schema/module"
	"github.com/zclconf/go-cty/cty"
)

// AnySchemaForVariableCollection returns a schema for collecting all
// variables in a variable file. It doesn't check if a variable has
// been defined in the module or not.
//
// We can use this schema to collect variable references without waiting
// on the module metadata.
func AnySchemaForVariableCollection(modPath string) *schema.BodySchema {
	return &schema.BodySchema{
		AnyAttribute: &schema.AttributeSchema{
			OriginForTarget: &schema.PathTarget{
				Address: schema.Address{
					schema.StaticStep{Name: "var"},
					schema.AttrNameStep{},
				},
				Path: lang.Path{
					Path:       modPath,
					LanguageID: ModuleLanguageID,
				},
				Constraints: schema.Constraints{
					ScopeId: refscope.VariableScope,
					Type:    cty.DynamicPseudoType,
				},
			},
			Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
		},
	}
}

func SchemaForVariables(vars map[string]module.Variable, modPath string) (*schema.BodySchema, error) {
	attributes := make(map[string]*schema.AttributeSchema)

	for name, modVar := range vars {
		aSchema := ModuleVarToAttribute(modVar)
		varType := modVar.Type
		aSchema.Constraint = schema.LiteralType{Type: varType}
		aSchema.OriginForTarget = &schema.PathTarget{
			Address: schema.Address{
				schema.StaticStep{Name: "var"},
				schema.AttrNameStep{},
			},
			Path: lang.Path{
				Path:       modPath,
				LanguageID: ModuleLanguageID,
			},
			Constraints: schema.Constraints{
				ScopeId: refscope.VariableScope,
				Type:    varType,
			},
		}

		attributes[name] = aSchema
	}

	return &schema.BodySchema{
		Attributes: attributes,
	}, nil
}

func ModuleVarToAttribute(modVar module.Variable) *schema.AttributeSchema {
	aSchema := &schema.AttributeSchema{
		IsSensitive: modVar.IsSensitive,
	}

	if modVar.Description != "" {
		aSchema.Description = lang.PlainText(modVar.Description)
	}

	if modVar.DefaultValue == cty.NilVal {
		aSchema.IsRequired = true
	} else {
		aSchema.IsOptional = true
	}

	return aSchema
}

// targetablesForAddrType is used to generate targets for complex object variables. Nested targets are supported as well in case they are objects too.
func targetablesForAddrType(addr lang.Address, rootType cty.Type) schema.Targetables {
	if !rootType.IsObjectType() {
		return nil
	}

	targetables := schema.Targetables{}

	for attrName, attrType := range rootType.AttributeTypes() {
		nestedAddr := addr.Copy()
		nestedAddr = append(nestedAddr, lang.AttrStep{Name: attrName})
		targetables = append(targetables, &schema.Targetable{
			Address:           nestedAddr,
			ScopeId:           refscope.VariableScope,
			AsType:            attrType,
			NestedTargetables: targetablesForAddrType(nestedAddr, attrType),
		})
	}

	sort.Sort(targetables)
	return targetables
}
