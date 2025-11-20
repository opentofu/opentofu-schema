// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

// Test that the deprecated attribute exists in variable blocks for v1.10
func TestCoreModuleSchemaForVersion_v1_10_variable_deprecated(t *testing.T) {
	v := version.Must(version.NewVersion("1.10.0"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	// Test that variable block exists
	variableBlock, exists := schemaForVersion.Blocks["variable"]
	if !exists {
		t.Fatal("expected variable block in v1.10 schema")
	}

	// Test variable block has deprecated attribute
	deprecatedAttr, exists := variableBlock.Body.Attributes["deprecated"]
	if !exists {
		t.Fatal("expected variable block to have deprecated attribute in v1.10")
	}

	// Test deprecated attribute is optional string
	if !deprecatedAttr.IsOptional {
		t.Error("expected variable deprecated attribute to be optional")
	}

	// Test that it's a string type
	literalType, ok := deprecatedAttr.Constraint.(schema.LiteralType)
	if !ok {
		t.Fatal("expected deprecated attribute to have LiteralType constraint")
	}

	if !literalType.Type.Equals(cty.String) {
		t.Errorf("expected deprecated attribute to be string type, got %s", literalType.Type.FriendlyName())
	}

	// Test that description exists
	if deprecatedAttr.Description.Value == "" {
		t.Error("expected deprecated attribute to have a description")
	}
}

// Test that all other expected blocks still exist in v1.10
func TestCoreModuleSchemaForVersion_v1_10_block_integrity(t *testing.T) {
	v := version.Must(version.NewVersion("1.10.0"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	// Test that all other expected blocks still exist, no accidental overrides
	expectedBlocks := []string{"resource", "terraform", "moved", "removed", "data", "locals", "module", "output", "provider", "variable", "import", "check"}

	for _, blockName := range expectedBlocks {
		if _, exists := schemaForVersion.Blocks[blockName]; !exists {
			t.Errorf("expected %s block to exist in v1.10 schema", blockName)
		}
	}
}

// Test that other variable attributes still exist
func TestCoreModuleSchemaForVersion_v1_10_variable_attributes(t *testing.T) {
	v := version.Must(version.NewVersion("1.10.0"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	variableBlock := schemaForVersion.Blocks["variable"]

	// Test that all expected variable attributes exist
	// Note for developers in the future, we do not check for "default" here as it's handled in the core schema and not needed to be done here.
	expectedAttrs := []string{"type", "description", "sensitive", "nullable", "deprecated"}

	for _, attrName := range expectedAttrs {
		if _, exists := variableBlock.Body.Attributes[attrName]; !exists {
			t.Errorf("expected variable block to have %s attribute", attrName)
		}
	}

	// Test that validation block exists (it's a block, not an attribute)
	if _, exists := variableBlock.Body.Blocks["validation"]; !exists {
		t.Error("expected variable block to have validation block")
	}
}

// Test that the deprecated attribute exists in output blocks for v1.10
func TestCoreModuleSchemaForVersion_v1_10_output_deprecated(t *testing.T) {
	v := version.Must(version.NewVersion("1.10.0"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	// Test that output block exists
	outputBlock, exists := schemaForVersion.Blocks["output"]
	if !exists {
		t.Fatal("expected output block in v1.10 schema")
	}

	// Test output block has deprecated attribute
	deprecatedAttr, exists := outputBlock.Body.Attributes["deprecated"]
	if !exists {
		t.Fatal("expected output block to have deprecated attribute in v1.10")
	}

	// Test deprecated attribute is optional string
	if !deprecatedAttr.IsOptional {
		t.Error("expected output deprecated attribute to be optional")
	}

	// Test that it's a string type
	literalType, ok := deprecatedAttr.Constraint.(schema.LiteralType)
	if !ok {
		t.Fatal("expected deprecated attribute to have LiteralType constraint")
	}

	if !literalType.Type.Equals(cty.String) {
		t.Errorf("expected deprecated attribute to be string type, got %s", literalType.Type.FriendlyName())
	}

	// Test that description exists
	if deprecatedAttr.Description.Value == "" {
		t.Error("expected deprecated attribute to have a description")
	}
}

// Test that other output attributes still exist
func TestCoreModuleSchemaForVersion_v1_10_output_attributes(t *testing.T) {
	v := version.Must(version.NewVersion("1.10.0"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	outputBlock := schemaForVersion.Blocks["output"]

	// Test that all expected output attributes exist
	// Note for developers in the future, we do not check for "value" here as it's handled differently in the schema
	expectedAttrs := []string{"description", "sensitive", "deprecated"}

	for _, attrName := range expectedAttrs {
		if _, exists := outputBlock.Body.Attributes[attrName]; !exists {
			t.Errorf("expected output block to have %s attribute", attrName)
		}
	}
}
