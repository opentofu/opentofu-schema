// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
)

// ResolveVersion returns OpenTofu version for which we have schema available
// based on either given version and/or constraint.
// Lack of constraint and version implies latest known version.
//
//go:generate go run ../internal/versiongen -w ./versions_gen.go
func ResolveVersion(tfVersion *version.Version, tfCons version.Constraints) *version.Version {
	if tfVersion != nil {
		coreVersion := tfVersion.Core()
		if coreVersion.LessThan(OldestAvailableVersion) {
			return OldestAvailableVersion
		}
		// Specified version core is greater than anything we have currently including prereleases
		if coreVersion.GreaterThan(LatestAvailableVersionIncludingPrereleases) {
			// There is no need to look for the tfVersion in available version
			// since it is greater than any on those. And iteration over the known versions happens below otherwise.
			// Hence, we need to return one of the latest versions known to us.

			// User is asking for a stable release
			if tfVersion.Prerelease() == "" {
				return LatestAvailableVersion
			}

			// User provided version with prerelease label in it, we can assume higher risk tolerance and return the latest non-stable version
			// although the prerelease part inside the patch release won't matter to us
			return LatestAvailableVersionIncludingPrereleases.Core()
		}
		if tfCons.Check(coreVersion) {
			return coreVersion
		}
	}

	for _, v := range tofuVersions {
		if len(tfCons) > 0 && tfCons.Check(v) && v.LessThan(OldestAvailableVersion) {
			return OldestAvailableVersion
		}
		// Check if the version in it's core matches the version core user provided.
		// This schema, as of writing this, doesn't care about differences
		// between sub-release under the patch release (different pre-releases and the final release)
		if tfVersion != nil && tfVersion.Core().Equal(v.Core()) {
			return tfVersion.Core()
		}
		if len(tfCons) > 0 && tfCons.Check(v) {
			return v
		}
	}

	return LatestAvailableVersion
}
