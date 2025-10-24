// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import "fmt"

type ephemeralResource struct {
	resource
}

func (r *ephemeralResource) MapKey() string {
	return fmt.Sprintf("ephemeral.%s.%s", r.Type, r.Name)
}
