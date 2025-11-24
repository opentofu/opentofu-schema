// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func patchResourceLifecycleBlockV1_11(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["enabled"] = enabledAttribute("resource")
	return bs
}

func patchDataLifecycleBlockV1_11(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["enabled"] = enabledAttribute("data source")
	return bs
}

func moduleLifecycleBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Lifecycle customizations for module blocks"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"enabled": enabledAttribute("module"),
			},
		},
	}
}

func enabledAttribute(blockType string) *schema.AttributeSchema {
	return &schema.AttributeSchema{
		Constraint:   schema.AnyExpression{OfType: cty.Bool},
		DefaultValue: schema.DefaultValue{Value: cty.True},
		IsOptional:   true,
		Description:  lang.Markdown(fmt.Sprintf("Whether the %s is enabled. When set to `false`, the %s will be skipped during plan and apply", blockType, blockType)),
	}
}
