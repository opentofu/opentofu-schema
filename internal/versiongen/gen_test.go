// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/go-version"
)

func TestGetTofuReleases(t *testing.T) {
	releases, err := GetTofuReleases()
	if err != nil {
		t.Fatal(err)
	}

	minExpectedLength := 51
	if len(releases) < minExpectedLength {
		t.Fatalf("expected >= %d releases, %d given", minExpectedLength, len(releases))
	}

	// The oldest release should really be 1.6.0-alpha1. We're however getting
	// releases sorted by dates and those dates were backfilled as part
	// of some older data migrations where the original dates were lost.
	expectedOldestRelease := release{
		Version: version.Must(version.NewVersion("1.6.0-alpha1")),
	}
	oldestRelease := releases[len(releases)-1]
	if diff := cmp.Diff(expectedOldestRelease, oldestRelease); diff != "" {
		t.Fatalf("unexpected oldest release: %s", diff)
	}
}
