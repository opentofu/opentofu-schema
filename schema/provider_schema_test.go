// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfaddr "github.com/opentofu/registry-address"
	"github.com/zclconf/go-cty-debug/ctydebug"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/function"
)

func TestProviderSchema_SetProviderVersion(t *testing.T) {
	ps := &ProviderSchema{
		Provider: &schema.BodySchema{},
		Resources: map[string]*schema.BodySchema{
			"foo": {
				Attributes: map[string]*schema.AttributeSchema{
					"str": {
						Constraint: schema.LiteralType{Type: cty.String},
						IsOptional: true,
					},
				},
			},
		},
		EphemeralResources: map[string]*schema.BodySchema{
			"foo": {
				Attributes: map[string]*schema.AttributeSchema{
					"str": {
						Constraint: schema.LiteralType{Type: cty.String},
						IsOptional: true,
					},
				},
			},
		},
		DataSources: map[string]*schema.BodySchema{
			"bar": {
				Attributes: map[string]*schema.AttributeSchema{
					"num": {
						Constraint: schema.LiteralType{Type: cty.Number},
						IsOptional: true,
					},
				},
			},
		},
		Functions: map[string]*schema.FunctionSignature{
			"baz": {
				Params: []function.Parameter{
					{
						Name:        "a",
						Type:        cty.String,
						Description: "first parameter",
					},
				},
				Description: "baz",
			},
		},
	}
	expectedSchema := &ProviderSchema{
		Provider: &schema.BodySchema{
			Detail:   "hashicorp/aws 3.76.1",
			HoverURL: "https://search.opentofu.org/provider/hashicorp/aws/v3.76.1/",
			DocsLink: &schema.DocsLink{
				URL:     "https://search.opentofu.org/provider/hashicorp/aws/v3.76.1/",
				Tooltip: "hashicorp/aws Documentation",
			},
		},
		Resources: map[string]*schema.BodySchema{
			"foo": {
				Detail: "hashicorp/aws 3.76.1",
				Attributes: map[string]*schema.AttributeSchema{
					"str": {
						Constraint: schema.LiteralType{Type: cty.String},
						IsOptional: true,
					},
				},
			},
		},
		EphemeralResources: map[string]*schema.BodySchema{
			"foo": {
				Detail: "hashicorp/aws 3.76.1",
				Attributes: map[string]*schema.AttributeSchema{
					"str": {
						Constraint: schema.LiteralType{Type: cty.String},
						IsOptional: true,
					},
				},
			},
		},
		DataSources: map[string]*schema.BodySchema{
			"bar": {
				Detail: "hashicorp/aws 3.76.1",
				Attributes: map[string]*schema.AttributeSchema{
					"num": {
						Constraint: schema.LiteralType{Type: cty.Number},
						IsOptional: true,
					},
				},
			},
		},
		Functions: map[string]*schema.FunctionSignature{
			"baz": {
				Description: "baz",
				Detail:      "hashicorp/aws 3.76.1",
				Params: []function.Parameter{
					{Name: "a", Type: cty.String, Description: "first parameter"},
				},
			},
		},
	}

	pAddr := tfaddr.Provider{
		Hostname:  tfaddr.DefaultProviderRegistryHost,
		Namespace: "hashicorp",
		Type:      "aws",
	}
	pv := version.Must(version.NewVersion("3.76.1"))

	ps.SetProviderVersion(pAddr, pv)

	if diff := cmp.Diff(expectedSchema, ps, ctydebug.CmpOptions); diff != "" {
		t.Fatalf("unexpected schema: %s", diff)
	}
}
