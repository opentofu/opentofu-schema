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

func patchResourceLifecycleBlockV1_11(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["enabled"] = &schema.AttributeSchema{
		Constraint:   schema.LiteralType{Type: cty.Bool},
		DefaultValue: schema.DefaultValue{Value: cty.True},
		IsOptional:   true,
		Description:  lang.Markdown("Whether the resource is enabled. When set to `false`, the resource will be skipped during plan and apply"),
	}
	return bs
}

func patchDataLifecycleBlockV1_11(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["enabled"] = &schema.AttributeSchema{
		Constraint:   schema.LiteralType{Type: cty.Bool},
		DefaultValue: schema.DefaultValue{Value: cty.True},
		IsOptional:   true,
		Description:  lang.Markdown("Whether the data source is enabled. When set to `false`, the data source will be skipped during plan and apply"),
	}
	return bs
}

func patchEphemeralLifecycleBlockV1_11(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["enabled"] = &schema.AttributeSchema{
		Constraint:   schema.LiteralType{Type: cty.Bool},
		DefaultValue: schema.DefaultValue{Value: cty.True},
		IsOptional:   true,
		Description:  lang.Markdown("Whether the ephemeral resource is enabled. When set to `false`, the ephemeral resource will be skipped during plan and apply"),
	}
	return bs
}

func moduleLifecycleBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Lifecycle customizations for module blocks"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"enabled": {
					Constraint:   schema.LiteralType{Type: cty.Bool},
					DefaultValue: schema.DefaultValue{Value: cty.True},
					IsOptional:   true,
					Description:  lang.Markdown("Whether the module is enabled. When set to `false`, the module will be skipped during plan and apply"),
				},
			},
		},
	}
}
