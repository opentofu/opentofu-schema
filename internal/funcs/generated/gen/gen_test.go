// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"testing"

	tfjson "github.com/hashicorp/terraform-json"
)

func TestSignaturesFromTofu(t *testing.T) {
	ctx := context.Background()
	functions, err := signaturesFromTofu(ctx)

	if err != nil {
		t.Error(err)
	}

	if len(functions.Signatures) < 51 {
		t.Error("it should return at least 51 releases")
	}

	functionKeyToTest := "timecmp"
	var sig *tfjson.FunctionSignature
	var ok bool
	if sig, ok = functions.Signatures[functionKeyToTest]; !ok {
		t.Errorf("function %s isn't on the release", functionKeyToTest)
	}

	expectedReturnType := "number"
	if sig.ReturnType.FriendlyName() != expectedReturnType {
		t.Fatalf("the function should have %s as return type, got %s instead", expectedReturnType, sig.ReturnType.FriendlyName())
	}
}
