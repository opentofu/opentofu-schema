// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	v012_mod "github.com/opentofu/opentofu-schema/internal/schema/0.12"
	v1_3_mod "github.com/opentofu/opentofu-schema/internal/schema/1.3"
	v1_4_mod "github.com/opentofu/opentofu-schema/internal/schema/1.4"
	"github.com/zclconf/go-cty/cty"
)

func patchRemovedBlock1_10(v *version.Version, bs *schema.BodySchema) {
	bs.Blocks["removed"].Body.Blocks["lifecycle"] = &schema.BlockSchema{
		Description: lang.Markdown("Lifecycle customizations controlling the removal"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"destroy": {
					Constraint:  schema.LiteralType{Type: cty.Bool},
					IsRequired:  true,
					Description: lang.Markdown("Whether OpenTofu will attempt to destroy the objects (`true`) or not, i.e. just remove from state (`false`)."),
				},
			},
		},
		MinItems: 0,
		MaxItems: 1,
	}

	bs.Blocks["removed"].Body.Blocks["provisioner"] = v012_mod.ProvisionerBlock(v)
	bs.Blocks["removed"].Body.Blocks["provisioner"].DependentBody = v1_4_mod.ProvisionerDependentBodies(v)
	bs.Blocks["removed"].Body.Blocks["provisioner"].Body.Blocks["connection"].DependentBody = v1_3_mod.ConnectionDependentBodies(v)
	bs.Blocks["removed"].Body.Blocks["provisioner"].Body.Attributes["when"] = &schema.AttributeSchema{
		Constraint: schema.OneOf{
			schema.Keyword{
				Keyword:     "destroy",
				Description: lang.Markdown("Run the provisioner when the resource is destroyed"),
			},
		},
		IsOptional:  true,
		Description: lang.Markdown("When to run the provisioner - `removed` resources can only be destroyed."),
	}
}
