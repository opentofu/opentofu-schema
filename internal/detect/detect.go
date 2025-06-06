// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package detect

import (
	"fmt"
	"net/url"
	"path/filepath"
)

var RemoteSourceDetectors = []Detector{
	new(GitHubDetector),
	new(GitDetector),
	new(BitBucketDetector),
	new(GCSDetector),
	new(S3Detector),
}

// Detector defines the interface that an invalid URL or a URL with a blank
// scheme is passed through in order to determine if its shorthand for
// something else well-known.
type Detector interface {
	// Detect will detect whether the string matches a known pattern to
	// turn it into a proper URL.
	Detect(string) (string, bool, error)
}

// Detect turns a source string into another source string if it is
// detected to be of a known pattern.
//
// The third parameter should be the list of detectors to use in the
// order to try them. If you don't want to configure this, just use
// the global Detectors variable.
//
// This is safe to be called with an already valid source string: Detect
// will just return it.
func Detect(src string, ds []Detector) (string, error) {
	getForce, getSrc := getForcedGetter(src)

	// Separate out the subdir if there is one, we don't pass that to detect
	getSrc, subDir := SourceDirSubdir(getSrc)

	u, err := url.Parse(getSrc)
	if err == nil && u.Scheme != "" {
		// Valid URL
		return src, nil
	}

	for _, d := range ds {
		result, ok, err := d.Detect(getSrc)
		if err != nil {
			return "", err
		}
		if !ok {
			continue
		}

		var detectForce string
		detectForce, result = getForcedGetter(result)
		result, detectSubdir := SourceDirSubdir(result)

		// If we have a subdir from the detection, then prepend it to our
		// requested subdir.
		if detectSubdir != "" {
			if subDir != "" {
				subDir = filepath.Join(detectSubdir, subDir)
			} else {
				subDir = detectSubdir
			}
		}

		if subDir != "" {
			u, err := url.Parse(result)
			if err != nil {
				return "", fmt.Errorf("error parsing URL: %s", err)
			}
			u.Path += "//" + subDir

			// a subdir may contain wildcards, but in order to support them we
			// have to ensure the path isn't escaped.
			u.RawPath = u.Path

			result = u.String()
		}

		// Preserve the forced getter if it exists. We try to use the
		// original set force first, followed by any force set by the
		// detector.
		if getForce != "" {
			result = fmt.Sprintf("%s::%s", getForce, result)
		} else if detectForce != "" {
			result = fmt.Sprintf("%s::%s", detectForce, result)
		}

		return result, nil
	}

	return "", fmt.Errorf("invalid source string: %s", src)
}
