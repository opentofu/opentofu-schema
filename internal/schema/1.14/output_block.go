// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"

	v1_12_mod "github.com/opentofu/opentofu-schema/internal/schema/1.12"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_12_mod.ModuleSchema(v)

	// Update the outpu block to include the "type" attribute for version 1.14
	bs.Blocks["output"] = patchOutputBlockWithType(bs.Blocks["output"])

	return bs
}

func patchOutputBlockWithType(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Address.AsTypeOf = &schema.BlockAsTypeOf{
		AttributeExpr: "type",
	}
	bs.Body.Attributes["type"] = &schema.AttributeSchema{
		Constraint:  schema.TypeDeclaration{},
		IsOptional:  true,
		Description: lang.Markdown("Type constraint restricting the type of value to accept, e.g. ``string` or `list(string)`"),
	}

	return bs
}
