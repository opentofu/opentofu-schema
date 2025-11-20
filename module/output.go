// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package module

import (
	"github.com/zclconf/go-cty/cty"
)

type Output struct {
	Description string
	IsSensitive bool
	Value       cty.Value

	// Deprecated is a string to mark an output as deprecated with instructions to end users
	// of the module.
	Deprecated string
}
