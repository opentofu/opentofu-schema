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

func removedBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Declaration to specify what resource or module to remove from the state"),
		Body: &schema.BodySchema{
			HoverURL: "https://opentofu.org/docs/language/resources/syntax/#removing-resources",
			Attributes: map[string]*schema.AttributeSchema{
				"from": {
					Constraint: schema.OneOf{
						schema.Reference{OfScopeId: refscope.ModuleScope},
						schema.Reference{OfScopeId: refscope.ResourceScope},
					},
					IsRequired:  true,
					Description: lang.Markdown("Address of the module or resource to be removed"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{},
		},
	}
}
