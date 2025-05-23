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
	bs.Body.Blocks["cloud"].Body.Blocks["workspaces"].Body = &schema.BodySchema{
		Attributes: map[string]*schema.AttributeSchema{
			"name": {
				Constraint: schema.LiteralType{Type: cty.String},
				IsOptional: true,
				Description: lang.Markdown("The name of a cloud workspace " +
					"to be used with this configuration. When configured only the specified workspace " +
					"can be used. This option conflicts with `tags`."),
			},
			"project": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.PlainText("The name of a cloud project. Workspaces that need creating will be created within this project."),
			},
			"tags": {
				Constraint: schema.Set{
					Elem: schema.LiteralType{Type: cty.String},
				},
				IsOptional: true,
				Description: lang.Markdown("A set of tags used to select remote cloud workspaces" +
					" to be used for this single configuration. New workspaces will automatically be tagged " +
					"with these tag values. Generally, this is the primary and recommended strategy to use. " +
					"This option conflicts with `name`."),
			},
		},
	}

	return bs
}
