// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfschema "github.com/hashicorp/terraform-schema/schema"
	test_v1_6 "github.com/opentofu/opentofu-schema/internal/schema/tests/1.6"
	test_v1_8 "github.com/opentofu/opentofu-schema/internal/schema/tests/1.8"
)

var (
	v1_6 = version.Must(version.NewVersion("1.6"))
	v1_8 = version.Must(version.NewVersion("1.8"))
)

// CoreTestSchemaForVersion finds a schema for test configuration files
// that is relevant for the given OpenTofu version.
// It will return an error if such schema cannot be found.
func CoreTestSchemaForVersion(v *version.Version) (*schema.BodySchema, error) {
	ver := v.Core()

	if ver.GreaterThanOrEqual(v1_8) {
		return test_v1_8.TestSchema(ver), nil
	}
	if ver.GreaterThanOrEqual(v1_6) {
		return test_v1_6.TestSchema(ver), nil
	}

	return nil, tfschema.NoCompatibleSchemaErr{Version: ver}
}
