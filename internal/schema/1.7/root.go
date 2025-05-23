// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v1_6_mod "github.com/opentofu/opentofu-schema/internal/schema/1.6"
)

func ModuleSchema(v *version.Version) *schema.BodySchema {
	bs := v1_6_mod.ModuleSchema(v)
	bs.Blocks["removed"] = removedBlock()
	bs.Blocks["import"].Body.Extensions = &schema.BodyExtensions{
		ForEach: true,
	}
	return bs
}
