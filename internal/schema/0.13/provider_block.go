// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	"github.com/opentofu/opentofu-schema/internal/schema/refscope"
	"github.com/opentofu/opentofu-schema/internal/schema/tokmod"
)

func providerBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.LabelStep{Index: 0},
				schema.AttrValueStep{Name: "alias", IsOptional: true},
			},
			FriendlyName: "provider",
			ScopeId:      refscope.ProviderScope,
			AsReference:  true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Provider},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name, lang.TokenModifierDependent},
				Description:            lang.PlainText("Provider Name"),
				IsDepKey:               true,
				Completable:            true,
			},
		},
		Description: lang.PlainText("A provider block is used to specify a provider configuration"),
		Body: &schema.BodySchema{
			Extensions: &schema.BodyExtensions{
				DynamicBlocks: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"alias": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("Alias for using the same provider with different configurations for different resources, e.g. `eu-west`"),
				},
				"version": {
					Constraint:   schema.LiteralType{Type: cty.String},
					IsOptional:   true,
					IsDeprecated: true,
					Description: lang.Markdown("Specifies a version constraint for the provider. e.g. `~> 1.0`.\n" +
						"**DEPRECATED:** Use `required_providers` block to manage provider version instead."),
				},
			},
		},
	}
}
