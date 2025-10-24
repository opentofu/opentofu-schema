// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/schema/refscope"
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
			"[Read more about ephemeral variables.](https://opentofu.org/docs/v1.11/language/values/variables/#ephemerality)"),
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
			"[Read more about ephemeral variables.](https://opentofu.org/docs/v1.11/language/values/outputs/#ephemerality)"),
	}

	return bs
}

func defaultDependsOnValue() schema.Constraint {
	return schema.Set{
		Elem: schema.OneOf{
			schema.Reference{OfScopeId: refscope.DataScope},
			schema.Reference{OfScopeId: refscope.ModuleScope},
			schema.Reference{OfScopeId: refscope.ResourceScope},
			schema.Reference{OfScopeId: refscope.EphemeralScope},
			schema.Reference{OfScopeId: refscope.VariableScope},
			schema.Reference{OfScopeId: refscope.LocalScope},
		},
	}
}

// patchDependencyScopeConstraintsWithEphemeral Adds Ephemeral scope to depends_on blocks of different blocks
func patchDependencyScopeConstraintsWithEphemeral(bs *schema.BodySchema) {
	// Every block type that support depends_on attribute currently in the schema can also depend on the ephemeral resources and their constraints are identical
	for _, block := range bs.Blocks {
		attrs := block.Body.Attributes
		if _, ok := attrs["depends_on"]; !ok {
			continue
		}
		attrs["depends_on"].Constraint = defaultDependsOnValue()
	}

	// The outlier that we also need to patch is the "data" block inside the "check" block
	// This is configured separately and needs to have it's depends_on updated separately
	bs.Blocks["check"].Body.Blocks["data"].Body.Attributes["depends_on"].Constraint = defaultDependsOnValue()
}
