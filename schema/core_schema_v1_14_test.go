// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"testing"

	"github.com/hashicorp/go-version"
)

func TestCoreModuleSchemaForVersion_v1_14_language(t *testing.T) {
	v := version.Must(version.NewVersion("1.14.0"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	outputBlock, exists := schemaForVersion.Blocks["output"]
	if !exists {
		t.Fatal("expected output block in v1.14 schema")
	}

	if _, exists := outputBlock.Body.Attributes["type"]; !exists {
		t.Error("expected output block to have type attribute")
	}
}

func TestCoreModuleSchemaForVersion_v1_12_noLanguage(t *testing.T) {
	v := version.Must(version.NewVersion("1.12.0"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	if _, exists := schemaForVersion.Blocks["output"].Body.Attributes["type"]; exists {
		t.Error("did not expect output type attribute in v1.12 schema")
	}
}
