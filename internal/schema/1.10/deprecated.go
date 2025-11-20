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

func patchVariableBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["deprecated"] = &schema.AttributeSchema{
		Constraint:  schema.LiteralType{Type: cty.String},
		Description: lang.Markdown("A message indicating why the variable is deprecated, and what to use instead."),
		IsOptional:  true,
	}
	return bs
}

func patchOutputBlockSchema(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["deprecated"] = &schema.AttributeSchema{
		Constraint:  schema.LiteralType{Type: cty.String},
		Description: lang.Markdown("A message indicating why the output is deprecated, and what to use instead."),
		IsOptional:  true,
	}
	return bs
}
