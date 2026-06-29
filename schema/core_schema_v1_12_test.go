// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"testing"

	"github.com/hashicorp/go-version"
)

func TestCoreModuleSchemaForVersion_v1_12_language(t *testing.T) {
	v := version.Must(version.NewVersion("1.12.0"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	languageBlock, exists := schemaForVersion.Blocks["language"]
	if !exists {
		t.Fatal("expected language block in v1.12 schema")
	}

	if len(languageBlock.Labels) != 0 {
		t.Errorf("expected language block to have 0 labels, got %d", len(languageBlock.Labels))
	}

	if _, exists := languageBlock.Body.Attributes["edition"]; !exists {
		t.Error("expected language block to have edition attribute")
	}

	if _, exists := languageBlock.Body.Attributes["experiments"]; !exists {
		t.Error("expected language block to have experiments attribute")
	}

	compatibleWith, exists := languageBlock.Body.Blocks["compatible_with"]
	if !exists {
		t.Fatal("expected language block to have compatible_with block")
	}

	if _, exists := compatibleWith.Body.Attributes["opentofu"]; !exists {
		t.Error("expected compatible_with block to have opentofu attribute")
	}
}

func TestCoreModuleSchemaForVersion_v1_11_noLanguage(t *testing.T) {
	v := version.Must(version.NewVersion("1.11.0"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	if _, exists := schemaForVersion.Blocks["language"]; exists {
		t.Error("did not expect language block in v1.11 schema")
	}
}
