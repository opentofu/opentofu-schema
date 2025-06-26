// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_9_mod "github.com/opentofu/opentofu-schema/internal/schema/1.9"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_9_mod.ModuleSchema(v)

	bs.Blocks["terraform"].Body.Blocks["encryption"] = patchEncryptionBlockSchema(
		bs.Blocks["terraform"].Body.Blocks["encryption"],
	)

	patchRemovedBlock1_10(v, bs)
	return bs
}
