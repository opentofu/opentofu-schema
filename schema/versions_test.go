// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"
	"testing"

	"github.com/hashicorp/go-version"
)

func TestResolveVersion(t *testing.T) {
	testCases := []struct {
		installedVersion *version.Version
		constraint       version.Constraints
		expectedVersion  *version.Version
	}{
		{
			nil,
			version.Constraints{},
			LatestAvailableVersion,
		},
		{
			version.Must(version.NewVersion("0.11.0")),
			nil,
			OldestAvailableVersion,
		},
		{
			nil,
			version.MustConstraints(version.NewConstraint("> 999.999.999")),
			LatestAvailableVersion,
		},
		{
			version.Must(version.NewVersion("999.999.999")),
			nil,
			LatestAvailableVersion,
		},
		{
			version.Must(version.NewVersion("1.6.2")),
			nil,
			version.Must(version.NewVersion("1.6.2")),
		},
		{
			version.Must(version.NewVersion("1.6.1")),
			version.MustConstraints(version.NewConstraint("> 999")),
			version.Must(version.NewVersion("1.6.1")),
		},
		{
			version.Must(version.NewVersion("1.6.3")),
			nil,
			version.Must(version.NewVersion("1.6.3")),
		},
		{
			version.Must(version.NewVersion("1.7.0-alpha20231025")),
			nil,
			version.Must(version.NewVersion("1.7.0")),
		},
		{
			version.Must(version.NewVersion("1.6.0-beta2")),
			nil,
			version.Must(version.NewVersion("1.6.0")),
		},
	}

	for i, tc := range testCases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.constraint.String()), func(t *testing.T) {
			resolvedVersion := ResolveVersion(tc.installedVersion, tc.constraint)
			if !tc.expectedVersion.Equal(resolvedVersion) {
				t.Fatalf("unexpected version: %q, expected: %q", resolvedVersion, tc.expectedVersion)
			}
		})
	}
}
