// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package detect

import (
	"testing"
)

func TestGitHubDetector(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		// HTTP
		{"github.com/hashicorp/foo", "git::https://github.com/hashicorp/foo.git"},
		{"github.com/hashicorp/foo.git", "git::https://github.com/hashicorp/foo.git"},
		{
			"github.com/hashicorp/foo/bar",
			"git::https://github.com/hashicorp/foo.git//bar",
		},
		{
			"github.com/hashicorp/foo?foo=bar",
			"git::https://github.com/hashicorp/foo.git?foo=bar",
		},
		{
			"github.com/hashicorp/foo.git?foo=bar",
			"git::https://github.com/hashicorp/foo.git?foo=bar",
		},
	}

	f := new(GitHubDetector)
	for i, tc := range cases {
		output, ok, err := f.Detect(tc.Input)
		if err != nil {
			t.Fatalf("err: %s", err)
		}
		if !ok {
			t.Fatal("not ok")
		}

		if output != tc.Output {
			t.Fatalf("%d: bad: %#v", i, output)
		}
	}
}
