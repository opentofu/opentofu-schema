// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func patchTerraformBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Blocks["cloud"] = &schema.BlockSchema{
		Description: lang.PlainText("Cloud configuration"),
		MaxItems:    1,
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"hostname": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.Markdown("The hostname to connect to"),
				},
				"organization": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsOptional:  true,
					Description: lang.PlainText("The name of the organization containing the targeted workspace(s)."),
				},
				"token": {
					Constraint: schema.LiteralType{Type: cty.String},
					IsOptional: true,
					Description: lang.Markdown("The token used to authenticate with the cloud platform. " +
						"Typically this argument should not be set, and `tofu login` used instead; " +
						"your credentials will then be fetched from your CLI configuration file " +
						"or configured credential helper."),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"workspaces": {
					Description: lang.Markdown("Workspace mapping strategy, either workspace `tags` or `name` is required."),
					MaxItems:    1,
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"name": {
								Constraint: schema.LiteralType{Type: cty.String},
								IsOptional: true,
								Description: lang.Markdown("The name of a single HCP Terraform workspace " +
									"to be used with this configuration. When configured only the specified workspace " +
									"can be used. This option conflicts with `tags`."),
							},
							"tags": {
								Constraint: schema.Set{
									Elem: schema.LiteralType{Type: cty.String},
								},
								IsOptional: true,
								Description: lang.Markdown("A set of tags used to select remote HCP Terraform workspaces" +
									" to be used for this single configuration. New workspaces will automatically be tagged " +
									"with these tag values. Generally, this is the primary and recommended strategy to use. " +
									"This option conflicts with `name`."),
							},
						},
					},
				},
			},
		},
	}

	return bs
}
