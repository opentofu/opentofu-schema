// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
)

var rootSchema = &hcl.BodySchema{
	Blocks: []hcl.BlockHeaderSchema{
		{
			Type:       "run",
			LabelNames: []string{"name"},
		},
		{
			Type: "variables",
		},
		{
			Type:       "provider",
			LabelNames: []string{"name"},
		},
	},
}
