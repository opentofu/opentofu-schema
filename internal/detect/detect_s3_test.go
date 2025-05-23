// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package detect

import (
	"testing"
)

func TestS3Detector(t *testing.T) {
	cases := []struct {
		Input  string
		Output string
	}{
		// Virtual hosted style
		{
			"bucket.s3.amazonaws.com/foo",
			"s3::https://s3.amazonaws.com/bucket/foo",
		},
		{
			"bucket.s3.amazonaws.com/foo/bar",
			"s3::https://s3.amazonaws.com/bucket/foo/bar",
		},
		{
			"bucket.s3.amazonaws.com/foo/bar.baz",
			"s3::https://s3.amazonaws.com/bucket/foo/bar.baz",
		},
		{
			"bucket.s3-eu-west-1.amazonaws.com/foo",
			"s3::https://s3-eu-west-1.amazonaws.com/bucket/foo",
		},
		{
			"bucket.s3-eu-west-1.amazonaws.com/foo/bar",
			"s3::https://s3-eu-west-1.amazonaws.com/bucket/foo/bar",
		},
		{
			"bucket.s3-eu-west-1.amazonaws.com/foo/bar.baz",
			"s3::https://s3-eu-west-1.amazonaws.com/bucket/foo/bar.baz",
		},
		// 5 parts Virtual hosted-style
		{
			"bucket.s3.eu-west-1.amazonaws.com/foo/bar.baz",
			"s3::https://s3.eu-west-1.amazonaws.com/bucket/foo/bar.baz",
		},
		// Path style
		{
			"s3.amazonaws.com/bucket/foo",
			"s3::https://s3.amazonaws.com/bucket/foo",
		},
		{
			"s3.amazonaws.com/bucket/foo/bar",
			"s3::https://s3.amazonaws.com/bucket/foo/bar",
		},
		{
			"s3.amazonaws.com/bucket/foo/bar.baz",
			"s3::https://s3.amazonaws.com/bucket/foo/bar.baz",
		},
		{
			"s3-eu-west-1.amazonaws.com/bucket/foo",
			"s3::https://s3-eu-west-1.amazonaws.com/bucket/foo",
		},
		{
			"s3-eu-west-1.amazonaws.com/bucket/foo/bar",
			"s3::https://s3-eu-west-1.amazonaws.com/bucket/foo/bar",
		},
		{
			"s3-eu-west-1.amazonaws.com/bucket/foo/bar.baz",
			"s3::https://s3-eu-west-1.amazonaws.com/bucket/foo/bar.baz",
		},
		// Misc tests
		{
			"s3-eu-west-1.amazonaws.com/bucket/foo/bar.baz?version=1234",
			"s3::https://s3-eu-west-1.amazonaws.com/bucket/foo/bar.baz?version=1234",
		},
	}

	f := new(S3Detector)
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
