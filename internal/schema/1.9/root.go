// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v012_mod "github.com/opentofu/opentofu-schema/internal/schema/0.12"
	v1_3_mod "github.com/opentofu/opentofu-schema/internal/schema/1.3"
	v1_8_mod "github.com/opentofu/opentofu-schema/internal/schema/1.8"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_8_mod.ModuleSchema(v)

	bs.Blocks["removed"].Body.Blocks["connection"] = v012_mod.ConnectionBlock(v)
	bs.Blocks["removed"].Body.Blocks["connection"].DependentBody = v1_3_mod.ConnectionDependentBodies(v)

	bs.Blocks["provider"].Body.Extensions = &schema.BodyExtensions{
		ForEach: true,
	}

	bs.Blocks["terraform"].Body.Blocks["encryption"] = patchEncryptionBlockSchema(bs.Blocks["terraform"].Body.Blocks["encryption"])

	return bs
}
