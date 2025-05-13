// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"fmt"

	"github.com/opentofu/opentofu-schema/module"
)

type dataSource struct {
	Type     string
	Name     string
	Provider module.ProviderRef
}

// MapKey returns a string that can be used to uniquely identify the receiver
// in a map[string]*dataSource.
func (r *dataSource) MapKey() string {
	return fmt.Sprintf("data.%s.%s", r.Type, r.Name)
}
