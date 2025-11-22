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

	// Update variable block to support ephemeral attribute
	bs.Blocks["variable"] = patchVariableBlockSchema(bs.Blocks["variable"])

	// Update output block to support ephemeral attribute
	bs.Blocks["output"] = patchOutputBlockSchema(bs.Blocks["output"])

	// Assign new set of identical contraints, including ephemeral resource as a viable dependency type to all block types
	patchDependencyScopeConstraintsWithEphemeral(bs)

	// Update lifecycle blocks to include "enabled" field for version 1.11
	bs.Blocks["resource"].Body.Blocks["lifecycle"] = patchResourceLifecycleBlockV1_11(bs.Blocks["resource"].Body.Blocks["lifecycle"])
	bs.Blocks["data"].Body.Blocks["lifecycle"] = patchDataLifecycleBlockV1_11(bs.Blocks["data"].Body.Blocks["lifecycle"])

	bs.Blocks["module"].Body.Blocks = moduleBlocks()

	return bs
}
