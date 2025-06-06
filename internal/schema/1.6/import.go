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

func importBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Import resources into OpenTofu to bring them under OpenTofu's management"),
		Body: &schema.BodySchema{
			HoverURL: "https://opentofu.org/docs/language/import/",
			Attributes: map[string]*schema.AttributeSchema{
				"provider": {
					Constraint:  schema.Reference{OfScopeId: refscope.ProviderScope},
					IsOptional:  true,
					Description: lang.Markdown("Reference to a `provider` configuration block, e.g. `mycloud.west` or `mycloud`"),
				},
				"id": {
					Constraint:  schema.AnyExpression{OfType: cty.String},
					IsRequired:  true,
					Description: lang.Markdown("ID of the resource to be imported. e.g. `i-abcd1234`"),
				},
				"to": {
					Constraint:  schema.Reference{OfScopeId: refscope.ResourceScope},
					IsRequired:  true,
					Description: lang.Markdown("An address of the resource instance to import to. e.g. `aws_instance.example` or `module.foo.aws_instance.bar`"),
				},
			},
		},
	}
}
