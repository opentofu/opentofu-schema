// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"github.com/hashicorp/hcl/v2"
	"github.com/opentofu/opentofu-schema/backend"
	"github.com/zclconf/go-cty/cty"
)

func decodeBackendsBlock(block *hcl.Block) (backend.BackendData, hcl.Diagnostics) {
	bType := block.Labels[0]
	attrs, diags := block.Body.JustAttributes()

	switch bType {
	case "remote":
		if attr, ok := attrs["hostname"]; ok {
			val, vDiags := attr.Expr.Value(nil)
			diags = append(diags, vDiags...)
			if val.IsWhollyKnown() && val.Type() == cty.String {
				return &backend.Remote{
					Hostname: val.AsString(),
				}, nil
			}
		}

		return &backend.Remote{}, nil
	}

	return &backend.UnknownBackendData{}, diags
}

func decodeCloudBlock(block *hcl.Block) (*backend.Cloud, hcl.Diagnostics) {
	attrs, _ := block.Body.JustAttributes()
	// Ignore diagnostics which may complain about unknown blocks

	if attr, ok := attrs["hostname"]; ok {
		val, vDiags := attr.Expr.Value(nil)
		if val.IsWhollyKnown() && val.Type() == cty.String {
			return &backend.Cloud{
				Hostname: val.AsString(),
			}, vDiags
		}
	}

	return &backend.Cloud{}, nil
}
