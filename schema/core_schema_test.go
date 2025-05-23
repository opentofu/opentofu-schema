// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	mod_v0_12 "github.com/opentofu/opentofu-schema/internal/schema/0.12"
	mod_v0_13 "github.com/opentofu/opentofu-schema/internal/schema/0.13"
	mod_v0_14 "github.com/opentofu/opentofu-schema/internal/schema/0.14"
	mod_v1_1 "github.com/opentofu/opentofu-schema/internal/schema/1.1"
	mod_v1_2 "github.com/opentofu/opentofu-schema/internal/schema/1.2"
	"github.com/zclconf/go-cty-debug/ctydebug"
)

func TestCoreModuleSchemaForVersion_tooOld(t *testing.T) {
	v := version.Must(version.NewVersion("0.11.0"))
	_, err := CoreModuleSchemaForVersion(v)
	if err == nil {
		t.Fatal("expected error for v0.11")
	}
	if !strings.Contains(err.Error(), "no compatible schema") {
		t.Fatalf("error mismatch: %q", err.Error())
	}
}

func TestCoreModuleSchemaForVersion_validate(t *testing.T) {
	versions := []string{
		"0.12.0-alpha1",
		"0.12.0-rc1",
		"0.12.0",
		"0.12.20",
		"0.13.0-alpha1",
		"0.13.0",
		"0.14.0-beta2",
		"0.14.0",
		"0.15.0",
		"1.0.0",
		"1.1.0",
		"1.2.0",
	}

	for _, v := range versions {
		ver, err := version.NewVersion(v)
		if err != nil {
			t.Fatal(err)
		}
		bodySchema, err := CoreModuleSchemaForVersion(ver)
		if err != nil {
			t.Fatal(err)
		}

		err = bodySchema.Validate()
		if err != nil {
			t.Fatalf("%s: %s", v, err)
		}
	}
}

func TestCoreModuleSchemaForVersion_matching(t *testing.T) {
	testCases := []struct {
		version       *version.Version
		matchedSchema versionedBodySchema
	}{
		{
			version.Must(version.NewVersion("0.12.0-alpha1")),
			mod_v0_12.ModuleSchema,
		},
		{
			version.Must(version.NewVersion("0.12.0-rc1")),
			mod_v0_12.ModuleSchema,
		},
		{
			version.Must(version.NewVersion("0.13.0-alpha1")),
			mod_v0_13.ModuleSchema,
		},
		{
			version.Must(version.NewVersion("0.14.0-beta2")),
			mod_v0_14.ModuleSchema,
		},
		{
			version.Must(version.NewVersion("1.1.0")),
			mod_v1_1.ModuleSchema,
		},
		{
			version.Must(version.NewVersion("1.2.0")),
			mod_v1_2.ModuleSchema,
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.version.String()), func(t *testing.T) {
			bodySchema, err := CoreModuleSchemaForVersion(tc.version)
			if err != nil {
				t.Fatal(err)
			}

			expectedSchema := tc.matchedSchema(tc.version)
			if diff := cmp.Diff(expectedSchema, bodySchema, ctydebug.CmpOptions); diff != "" {
				t.Fatalf("schema mismatch: %s", diff)
			}
		})
	}
}

type versionedBodySchema func(*version.Version) *schema.BodySchema
