// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
)

func importIdentitySchema(bs *schema.BlockSchema) *schema.BlockSchema {
	// FIXME: the schema is obtained via resourceSchema.IdentitySchema,
	// getting the resource schema from the provider
	// Can we "patch" this in here?
	bs.Body.Attributes["identity"] = &schema.AttributeSchema{
		IsRequired:  false,
		IsOptional:  true,
		Description: lang.Markdown("An expression, often a map, uniquely identifying a resource of the type being imported."),
	}
	bs.Body.Attributes["id"].IsRequired = false

	return bs
}
