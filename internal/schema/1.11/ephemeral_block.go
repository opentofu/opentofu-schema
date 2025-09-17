// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/schema/refscope"
	"github.com/opentofu/opentofu-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func ephemeralBlockSchema(_ *version.Version) *schema.BlockSchema {
	bs := &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "ephemeral"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName:         "ephemeral",
			ScopeId:              refscope.EphemeralScope,
			AsReference:          true,
			DependentBodyAsData:  true,
			InferDependentBody:   true,
			DependentBodySelfRef: true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Ephemeral},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				Description:            lang.PlainText("Ephemeral Resource Type"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				IsDepKey:               true,
				Completable:            true,
			},
			{
				Name:                   "name",
				Description:            lang.PlainText("Reference Name"),
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
			},
		},
		Description: lang.PlainText("An ephemeral block declares an ephemeral resource of a given type with a given local name. " +
			"Ephemeral resources are temporary values that are not stored in state or plan, designed for handling " +
			"transient and confidential information during OpenTofu execution."),
		Body: &schema.BodySchema{
			DocsLink: &schema.DocsLink{
				URL: "https://opentofu.org/docs/language/ephemerality/ephemeral-resources/",
			},
			Extensions: &schema.BodyExtensions{
				Count:         true,
				ForEach:       true,
				DynamicBlocks: true,
			},
			Attributes: map[string]*schema.AttributeSchema{
				"provider": {
					Constraint:             schema.Reference{OfScopeId: refscope.ProviderScope},
					IsOptional:             true,
					Description:            lang.Markdown("Reference to a `provider` configuration block, e.g. `mycloud.west` or `mycloud`"),
					IsDepKey:               true,
					SemanticTokenModifiers: lang.SemanticTokenModifiers{lang.TokenModifierDependent},
				},
				"depends_on": {
					Constraint: schema.Set{
						Elem: schema.OneOf{
							schema.Reference{OfScopeId: refscope.DataScope},
							schema.Reference{OfScopeId: refscope.ModuleScope},
							schema.Reference{OfScopeId: refscope.ResourceScope},
							schema.Reference{OfScopeId: refscope.VariableScope},
							schema.Reference{OfScopeId: refscope.LocalScope},
							schema.Reference{OfScopeId: refscope.EphemeralScope},
						},
					},
					IsOptional:  true,
					Description: lang.Markdown("Set of references to hidden dependencies, e.g. other resources, data sources, or ephemeral resources"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"lifecycle": ephemeralLifecycleBlock(),
			},
		},
	}

	return bs
}

func ephemeralLifecycleBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Lifecycle customizations for ephemeral resources (only precondition and postcondition are supported)"),
		Body: &schema.BodySchema{
			DocsLink: &schema.DocsLink{
				URL: "https://opentofu.org/docs/language/expressions/custom-conditions/#preconditions-and-postconditions",
			},
			Blocks: map[string]*schema.BlockSchema{
				"precondition": {
					Description: lang.Markdown("Custom condition to check before opening the ephemeral resource"),
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"condition": {
								Constraint:  schema.AnyExpression{},
								IsRequired:  true,
								Description: lang.Markdown("Condition that must be true before the ephemeral resource is opened"),
							},
							"error_message": {
								Constraint:  schema.AnyExpression{},
								IsRequired:  true,
								Description: lang.Markdown("Error message to show when the condition fails"),
							},
						},
					},
				},
				"postcondition": {
					Description: lang.Markdown("Custom condition to check after opening the ephemeral resource"),
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"condition": {
								Constraint:  schema.AnyExpression{OfType: cty.Bool},
								IsRequired:  true,
								Description: lang.Markdown("Condition that must be true after the ephemeral resource is opened"),
							},
							"error_message": {
								Constraint:  schema.AnyExpression{OfType: cty.String},
								IsRequired:  true,
								Description: lang.Markdown("Error message to show when the condition fails"),
							},
						},
					},
				},
			},
		},
	}
}
