// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/schema/backends"
	"github.com/opentofu/opentofu-schema/internal/schema/refscope"
	"github.com/opentofu/opentofu-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func terraformBlockSchema(v *version.Version) *schema.BlockSchema {
	return &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Terraform},
		Description:            lang.Markdown("`terraform` block used to configure some high-level behaviors of OpenTofu"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"required_version": {
					Constraint: schema.LiteralType{Type: cty.String},
					IsOptional: true,
					Description: lang.Markdown("Constraint to specify which versions of Terraform can be used " +
						"with this configuration, e.g. `~> 0.12`"),
				},
				"experiments": {
					Constraint: schema.Set{
						Elem: schema.OneOf{
							schema.Keyword{
								Keyword: "module_variable_optional_attrs",
								Name:    "feature",
							},
							schema.Keyword{
								Keyword: "provider_sensitive_attrs",
								Name:    "feature",
							},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("A set of experimental language features to enable"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"backend": {
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Backend},
					Description: lang.Markdown("Backend configuration which defines exactly where and how " +
						"operations are performed, where state snapshots are stored, etc."),
					Labels: []*schema.LabelSchema{
						{
							Name:                   "type",
							Description:            lang.Markdown("Backend Type"),
							IsDepKey:               true,
							Completable:            true,
							SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
						},
					},
					DependentBody: backends.ConfigsAsDependentBodies(v),
				},
				"provider_meta": {
					Description: lang.Markdown("Metadata to pass into a provider which supports this"),
					Labels: []*schema.LabelSchema{
						{
							Name:                   "name",
							Description:            lang.Markdown("Provider Name"),
							IsDepKey:               true,
							SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
						},
					},
				},
				"required_providers": {
					SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.RequiredProviders},
					Description: lang.Markdown("What provider version to use within this configuration " +
						"and where to source it from"),
					Body: &schema.BodySchema{
						AnyAttribute: &schema.AttributeSchema{
							Constraint: schema.OneOf{
								schema.Object{
									Attributes: schema.ObjectAttributes{
										"source": &schema.AttributeSchema{
											Constraint: schema.LiteralType{Type: cty.String},
											IsRequired: true,
											Description: lang.Markdown("The global source address for the provider " +
												"you intend to use, such as `hashicorp/aws`"),
										},
										"version": &schema.AttributeSchema{
											Constraint: schema.LiteralType{Type: cty.String},
											IsOptional: true,
											Description: lang.Markdown("Version constraint specifying which subset of " +
												"available provider versions the module is compatible with, e.g. `~> 1.0`"),
										},
									},
								},
								schema.LiteralType{Type: cty.String},
							},
							Description: lang.Markdown("Provider source and version constraint"),
							Address: &schema.AttributeAddrSchema{
								Steps: []schema.AddrStep{
									schema.AttrNameStep{},
								},
								AsReference:  true,
								FriendlyName: "provider",
								ScopeId:      refscope.ProviderScope,
							},
						},
					},
				},
			},
		},
	}
}
