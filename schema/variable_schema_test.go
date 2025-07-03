// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/schema/refscope"
	"github.com/opentofu/opentofu-schema/module"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
)

func TestSchemaForVariables(t *testing.T) {
	testCases := []struct {
		name           string
		variables      map[string]module.Variable
		expectedSchema *schema.BodySchema
	}{
		{
			"empty schema",
			make(map[string]module.Variable),
			&schema.BodySchema{Attributes: make(map[string]*schema.AttributeSchema)},
		},
		{
			"one attribute schema",
			map[string]module.Variable{
				"name": {
					Description: "name of the module",
					Type:        cty.String,
				},
			},
			&schema.BodySchema{Attributes: map[string]*schema.AttributeSchema{
				"name": {
					Description: lang.MarkupContent{
						Value: "name of the module",
						Kind:  lang.PlainTextKind,
					},
					IsRequired: true,
					Constraint: schema.LiteralType{Type: cty.String},
					OriginForTarget: &schema.PathTarget{
						Address:     schema.Address{schema.StaticStep{Name: "var"}, schema.AttrNameStep{}},
						Path:        lang.Path{Path: "./local", LanguageID: "opentofu"},
						Constraints: schema.Constraints{ScopeId: "variable", Type: cty.String},
					},
				},
			}},
		},
		{
			"two attribute schema",
			map[string]module.Variable{
				"name": {
					Description:  "name of the module",
					Type:         cty.String,
					DefaultValue: cty.StringVal("default"),
				},
				"id": {
					Description: "id of the module",
					Type:        cty.Number,
					IsSensitive: true,
				},
			},
			&schema.BodySchema{Attributes: map[string]*schema.AttributeSchema{
				"name": {
					Description: lang.MarkupContent{
						Value: "name of the module",
						Kind:  lang.PlainTextKind,
					},
					IsOptional: true,
					Constraint: schema.LiteralType{Type: cty.String},
					OriginForTarget: &schema.PathTarget{
						Address:     schema.Address{schema.StaticStep{Name: "var"}, schema.AttrNameStep{}},
						Path:        lang.Path{Path: "./local", LanguageID: "opentofu"},
						Constraints: schema.Constraints{ScopeId: "variable", Type: cty.String},
					},
				},
				"id": {
					Description: lang.MarkupContent{
						Value: "id of the module",
						Kind:  lang.PlainTextKind,
					},
					Constraint:  schema.LiteralType{Type: cty.Number},
					IsSensitive: true,
					IsRequired:  true,
					OriginForTarget: &schema.PathTarget{
						Address:     schema.Address{schema.StaticStep{Name: "var"}, schema.AttrNameStep{}},
						Path:        lang.Path{Path: "./local", LanguageID: "opentofu"},
						Constraints: schema.Constraints{ScopeId: "variable", Type: cty.Number},
					},
				},
			}},
		},
	}

	modPath := "./local"
	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.name), func(t *testing.T) {
			actualSchema, err := SchemaForVariables(tc.variables, modPath)

			if err != nil {
				t.Fatal(err)
			}

			diff := cmp.Diff(tc.expectedSchema, actualSchema, ctydebug.CmpOptions)
			if diff != "" {
				t.Fatalf("unexpected schema %s", diff)
			}
		})
	}
}

func TestGetTargetablesForAddrType(t *testing.T) {
	addr := lang.Address{
		lang.RootStep{Name: "var"},
		lang.AttrStep{Name: "complex_variable_name"},
	}
	rootType := cty.ObjectWithOptionalAttrs(map[string]cty.Type{
		"name": cty.String,
		"type": cty.ObjectWithOptionalAttrs(map[string]cty.Type{
			"nested_name": cty.String,
		}, []string{"nested_name"}),
	}, []string{"name", "type"})

	actualTargetables := getTargetablesForAddrType(addr, rootType)
	expectedTargetables := schema.Targetables{
		&schema.Targetable{
			Address: lang.Address{
				lang.RootStep{Name: "var"},
				lang.AttrStep{Name: "complex_variable_name"},
				lang.AttrStep{Name: "name"},
			},
			ScopeId:           refscope.VariableScope,
			AsType:            cty.String,
			NestedTargetables: nil,
		},
		&schema.Targetable{
			Address: lang.Address{
				lang.RootStep{Name: "var"},
				lang.AttrStep{Name: "complex_variable_name"},
				lang.AttrStep{Name: "type"},
			},
			ScopeId: refscope.VariableScope,
			AsType:  rootType.AttributeTypes()["type"],
			NestedTargetables: schema.Targetables{
				&schema.Targetable{
					Address: lang.Address{
						lang.RootStep{Name: "var"},
						lang.AttrStep{Name: "complex_variable_name"},
						lang.AttrStep{Name: "type"},
						lang.AttrStep{Name: "nested_name"},
					},
					ScopeId:           refscope.VariableScope,
					AsType:            cty.String,
					NestedTargetables: nil,
				},
			},
		},
	}

	diff := cmp.Diff(expectedTargetables, actualTargetables, ctydebug.CmpOptions)
	if diff != "" {
		t.Fatalf("unexpected targetables %s", diff)
	}
}
