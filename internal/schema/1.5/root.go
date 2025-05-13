// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_4_mod "github.com/opentofu/opentofu-schema/internal/schema/1.4"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_4_mod.ModuleSchema(v)
	bs.Blocks["import"] = importBlock()
	bs.Blocks["check"] = checkBlock()

	return bs
}
