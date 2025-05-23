// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/schema/refscope"
)

var movedBlockSchema = &schema.BlockSchema{
	Description: lang.Markdown("Refactoring declaration to specify what address to move where"),
	Body: &schema.BodySchema{
		HoverURL: "https://opentofu.org/docs/language/modules/develop/refactoring/#moved-block-syntax",
		Attributes: map[string]*schema.AttributeSchema{
			"from": {
				Constraint: schema.OneOf{
					schema.Reference{OfScopeId: refscope.ModuleScope},
					schema.Reference{OfScopeId: refscope.ResourceScope},
				},
				IsRequired:  true,
				Description: lang.Markdown("Source address to move away from"),
			},
			"to": {
				Constraint: schema.OneOf{
					schema.Reference{OfScopeId: refscope.ModuleScope},
					schema.Reference{OfScopeId: refscope.ResourceScope},
				},
				IsRequired:  true,
				Description: lang.Markdown("Destination address to move to"),
			},
		},
	},
}
