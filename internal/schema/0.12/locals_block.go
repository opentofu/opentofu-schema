// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/schema/refscope"
	"github.com/opentofu/opentofu-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func localsBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Locals},
		Description: lang.Markdown("Local values assigning names to expressions, so you can use these multiple times without repetition\n" +
			"e.g. `service_name = \"forum\"`"),
		Body: &schema.BodySchema{
			AnyAttribute: &schema.AttributeSchema{
				Address: &schema.AttributeAddrSchema{
					Steps: []schema.AddrStep{
						schema.StaticStep{Name: "local"},
						schema.AttrNameStep{},
					},
					ScopeId:     refscope.LocalScope,
					AsExprType:  true,
					AsReference: true,
				},
				Constraint: schema.AnyExpression{OfType: cty.DynamicPseudoType},
			},
		},
	}
}
