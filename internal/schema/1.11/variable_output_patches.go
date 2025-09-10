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

// patchVariableBlockSchema adds the ephemeral attribute to variable blocks
func patchVariableBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	// Add the ephemeral attribute
	bs.Body.Attributes["ephemeral"] = &schema.AttributeSchema{
		Constraint:   schema.LiteralType{Type: cty.Bool},
		DefaultValue: schema.DefaultValue{Value: cty.BoolVal(false)},
		IsOptional:   true,
		Description: lang.Markdown("Marks variable as ephemeral. OpenTofu will not store ephemeral variable in state at all" +
			"and will store them only their name in plan. \n Ephemeral variables can only be used in the limited context where ephemerals are allowed. " +
			"[Read more about ephemeral variables.](https://opentofu.org/docs/language/values/variables/#ephemerality)"),
	}

	return bs
}

// patchOutputBlockSchema adds the ephemeral attribute to output blocks
func patchOutputBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	// Add the ephemeral attribute
	bs.Body.Attributes["ephemeral"] = &schema.AttributeSchema{
		Constraint:   schema.LiteralType{Type: cty.Bool},
		DefaultValue: schema.DefaultValue{Value: cty.BoolVal(false)},
		IsOptional:   true,
		Description: lang.Markdown("Marks output as ephemeral. " +
			"Ephemeral outputs can only be used in the limited context where ephemerals are allowed.\n" +
			"[Read more about ephemeral variables.](https://opentofu.org/docs/language/values/outputs/#ephemerality)"),
	}

	return bs
}
