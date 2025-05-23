// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v015_mod "github.com/opentofu/opentofu-schema/internal/schema/0.15"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v015_mod.ModuleSchema(v)
	bs.Blocks["moved"] = movedBlockSchema
	bs.Blocks["terraform"] = patchTerraformBlockSchema(bs.Blocks["terraform"])
	bs.Blocks["variable"] = patchVariableBlockSchema(bs.Blocks["variable"])

	return bs
}
