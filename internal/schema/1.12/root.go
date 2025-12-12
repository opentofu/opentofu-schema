// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	v1_11_mod "github.com/opentofu/opentofu-schema/internal/schema/1.11"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_11_mod.ModuleSchema(v)

	// Update the lifecycle block to include the "destroy" attribute for version 1.12 retaining resources on destruction
	bs.Blocks["resource"].Body.Blocks["lifecycle"] = patchResourceLifecycleBlockWithDestroy(bs.Blocks["resource"].Body.Blocks["lifecycle"])

	return bs
}

func patchResourceLifecycleBlockWithDestroy(bs *schema.BlockSchema) *schema.BlockSchema {
	bs.Body.Attributes["destroy"] = &schema.AttributeSchema{
		Constraint:   schema.LiteralType{Type: cty.Bool},
		DefaultValue: schema.DefaultValue{Value: cty.True},
		IsOptional:   true,
		Description:  lang.Markdown(fmt.Sprintf("Setting `destroy` to `false` changes the OpenTofu's default behavior when destroying or replacing the resource. \n OpenTofu will 'forget' the resource instance, removing it from the state without destroying the actual infrastructure object. \n Read more about the `destroy` lifecycle argument in the [documentation](https://opentofu.org/docs/language/resources/behavior/#lifecycle-customizations).")),
	}
	return bs
}
