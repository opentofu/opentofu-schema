// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_10_mod "github.com/opentofu/opentofu-schema/internal/schema/1.10"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_10_mod.ModuleSchema(v)

	// Add the new ephemeral block
	bs.Blocks["ephemeral"] = ephemeralBlockSchema(v)

	return bs
}
