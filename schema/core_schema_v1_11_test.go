// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"testing"

	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/schema/refscope"

	"github.com/hashicorp/go-version"
	"github.com/zclconf/go-cty/cty"
)

// Generic test scenarios to check if the ephemeral schema is setup correctly
func TestCoreModuleSchemaForVersion_v1_11_ephemeral(t *testing.T) {
	v := version.Must(version.NewVersion("1.11.0-beta1"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	// Test that ephemeral block exists
	ephemeralBlock, exists := schemaForVersion.Blocks["ephemeral"]
	if !exists {
		t.Fatal("expected ephemeral block in v1.11 schema")
	}

	if ephemeralBlock.Address.FriendlyName != "ephemeral" {
		t.Errorf("expected ephemeral block friendly name to be 'ephemeral', got %q", ephemeralBlock.Address.FriendlyName)
	}

	// Test ephemeral block has two labels (type and name)
	if len(ephemeralBlock.Labels) != 2 {
		t.Errorf("expected ephemeral block to have 2 labels, got %d", len(ephemeralBlock.Labels))
	}

	if ephemeralBlock.Labels[0].Name != "type" {
		t.Errorf("expected first label to be 'type', got %q", ephemeralBlock.Labels[0].Name)
	}

	if ephemeralBlock.Labels[1].Name != "name" {
		t.Errorf("expected second label to be 'name', got %q", ephemeralBlock.Labels[1].Name)
	}

	if !ephemeralBlock.Body.Extensions.Count {
		t.Error("expected ephemeral block to support count")
	}

	if !ephemeralBlock.Body.Extensions.ForEach {
		t.Error("expected ephemeral block to support for_each")
	}

	// Test ephemeral block has provider and depends_on attributes
	if _, exists := ephemeralBlock.Body.Attributes["provider"]; !exists {
		t.Error("expected ephemeral block to have provider attribute")
	}

	if _, exists := ephemeralBlock.Body.Attributes["depends_on"]; !exists {
		t.Error("expected ephemeral block to have depends_on attribute")
	}

	// Test ephemeral block has lifecycle block
	if _, exists := ephemeralBlock.Body.Blocks["lifecycle"]; !exists {
		t.Error("expected ephemeral block to have lifecycle block")
	} else {
		// Ephemeral only supports pre and post conditions
		if _, exists := ephemeralBlock.Body.Blocks["lifecycle"].Body.Blocks["precondition"]; !exists {
			t.Error("expected ephemeral lifecycle block to support precondition")
		}
		if _, exists := ephemeralBlock.Body.Blocks["lifecycle"].Body.Blocks["postcondition"]; !exists {
			t.Error("expected ephemeral lifecycle block to support postcondition")
		}
	}

	// Test variable block has ephemeral attribute
	variableBlock, exists := schemaForVersion.Blocks["variable"]
	if !exists {
		t.Fatal("expected variable block in v1.11 schema")
	}

	ephemeralAttr, exists := variableBlock.Body.Attributes["ephemeral"]
	if !exists {
		t.Error("expected variable block to have ephemeral attribute in v1.11")
	} else {
		// Test ephemeral attribute is optional boolean
		if !ephemeralAttr.IsOptional {
			t.Error("expected variable ephemeral attribute to be optional")
		}
	}

	// Test output block has an ephemeral attribute
	outputBlock, exists := schemaForVersion.Blocks["output"]
	if !exists {
		t.Fatal("expected output block in v1.11 schema")
	}

	ephemeralAttr, exists = outputBlock.Body.Attributes["ephemeral"]
	if !exists {
		t.Error("expected output block to have ephemeral attribute in v1.11")
	} else {
		// Test ephemeral attribute is optional (boolean)
		if !ephemeralAttr.IsOptional {
			t.Error("expected output ephemeral attribute to be optional")
		}
	}
}

// Since we need to overwrite parts of the schema in order to update it
// there have been instance, where we overwrote an existing block and lost the definition
// this is a naive test and simply check if every other block still exists in the schema we return for 1.11
func TestCoreModuleSchemaForVersion_v1_11_block_integrity(t *testing.T) {
	v := version.Must(version.NewVersion("1.11.0-beta1"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	// Test that all other expected blocks still exist after merging the 1.11 schema, no accidental overrides
	expectedBlocks := []string{"resource", "terraform", "moved", "removed", "data", "locals", "module", "output", "provider", "variable", "import", "check"}

	for _, blockName := range expectedBlocks {
		if _, exists := schemaForVersion.Blocks[blockName]; !exists {
			t.Errorf("expected %s block to exist in v1.11 schema", blockName)
		}
	}
}

// Check that the correct blocks have a way to depend on the ephemeral block
func TestCoreModuleSchemaForVersion_v1_11_ephemeral_depends_on(t *testing.T) {
	v := version.Must(version.NewVersion("1.11.0-beta1"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}
	// Check that all expected blocks can depend_on ephemeral
	dependantBlocks := []string{"resource", "data", "module", "output"}
	for _, blockName := range dependantBlocks {
		if !blockCanDependOnEphemeral(schemaForVersion.Blocks[blockName]) {
			t.Errorf("expected the `%s` block to be able to depend_on ephemeral", blockName)
		}
	}
	// Additional case for nested "data" block under the "check" block
	if !blockCanDependOnEphemeral(schemaForVersion.Blocks["check"].Body.Blocks["data"]) {
		t.Errorf("expected the `data` block under the `check` block to be able to depend_on ephemeral")
	}
}

// blockCanDependOnEphemeral is a test helper that checks if the given block's depends_on attribute
// can reference ephemeral blocks.
func blockCanDependOnEphemeral(block *schema.BlockSchema) bool {
	do, ok := block.Body.Attributes["depends_on"]
	if !ok {
		return false
	}
	refsCons, ok := do.Constraint.(schema.Set)
	if !ok {
		return false
	}
	oneOf, ok := refsCons.Elem.(schema.OneOf)
	if !ok {
		return false
	}
	for _, con := range oneOf {
		ref, ok := con.(schema.Reference)
		if ok && ref.OfScopeId == refscope.EphemeralScope {
			return true
		}
	}
	return false
}

// Test that lifecycle blocks have the enabled attribute in v1.11
func TestCoreModuleSchemaForVersion_v1_11_lifecycle_enabled(t *testing.T) {
	v := version.Must(version.NewVersion("1.11.0-beta1"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	// Test resource lifecycle block has enabled attribute
	resourceBlock, exists := schemaForVersion.Blocks["resource"]
	if !exists {
		t.Fatal("expected resource block in v1.11 schema")
	}

	resourceLifecycle, exists := resourceBlock.Body.Blocks["lifecycle"]
	if !exists {
		t.Fatal("expected resource block to have lifecycle block")
	}

	enabledAttr, exists := resourceLifecycle.Body.Attributes["enabled"]
	if !exists {
		t.Error("expected resource lifecycle block to have enabled attribute in v1.11")
	} else {
		if !enabledAttr.IsOptional {
			t.Error("expected resource lifecycle enabled attribute to be optional")
		}
	}

	// Test resource lifecycle block has enabled attribute
	dataBlock, exists := schemaForVersion.Blocks["data"]
	if !exists {
		t.Fatal("expected resource block in v1.11 schema")
	}

	dataLifecycle, exists := dataBlock.Body.Blocks["lifecycle"]
	if !exists {
		t.Fatal("expected resource block to have lifecycle block")
	}

	enabledDataAttr, exists := dataLifecycle.Body.Attributes["enabled"]
	if !exists {
		t.Error("expected resource lifecycle block to have enabled attribute in v1.11")
	} else {
		if !enabledDataAttr.IsOptional {
			t.Error("expected resource lifecycle enabled attribute to be optional")
		}
	}

	// Test resource lifecycle block has enabled attribute
	moduleBlock, exists := schemaForVersion.Blocks["module"]
	if !exists {
		t.Fatal("expected resource block in v1.11 schema")
	}

	moduleLifecycle, exists := moduleBlock.Body.Blocks["lifecycle"]
	if !exists {
		t.Fatal("expected resource block to have lifecycle block")
	}

	enabledModuleAttr, exists := moduleLifecycle.Body.Attributes["enabled"]
	if !exists {
		t.Error("expected resource lifecycle block to have enabled attribute in v1.11")
	} else {
		if !enabledModuleAttr.IsOptional {
			t.Error("expected resource lifecycle enabled attribute to be optional")
		}
	}
}

// Test that ephemeral lifecycle block has the enabled block and proper structure
func TestCoreModuleSchemaForVersion_v1_11_ephemeral_lifecycle(t *testing.T) {
	v := version.Must(version.NewVersion("1.11.0-beta1"))
	schemaForVersion, err := CoreModuleSchemaForVersion(v)
	if err != nil {
		t.Fatal(err)
	}

	ephemeralBlock, exists := schemaForVersion.Blocks["ephemeral"]
	if !exists {
		t.Fatal("expected ephemeral block in v1.11 schema")
	}

	ephemeralLifecycle, exists := ephemeralBlock.Body.Blocks["lifecycle"]
	if !exists {
		t.Fatal("expected ephemeral block to have lifecycle block")
	}

	enabledAttr, exists := ephemeralLifecycle.Body.Attributes["enabled"]
	if !exists {
		t.Error("expected ephemeral lifecycle block to have enabled block in v1.11")
	} else {
		if !enabledAttr.IsOptional {
			t.Error("expected enabled block to be optional")
		}
	}

	// Test ephemeral lifecycle has precondition block
	preconditionBlock, exists := ephemeralLifecycle.Body.Blocks["precondition"]
	if !exists {
		t.Error("expected ephemeral lifecycle block to have precondition block")
	} else {
		conditionAttr, exists := preconditionBlock.Body.Attributes["condition"]
		if !exists {
			t.Error("expected precondition block to have condition attribute")
		} else {
			if !conditionAttr.IsRequired {
				t.Error("expected precondition condition attribute to be required")
			}
		}

		errorMessageAttr, exists := preconditionBlock.Body.Attributes["error_message"]
		if !exists {
			t.Error("expected precondition block to have error_message attribute")
		} else {
			if !errorMessageAttr.IsRequired {
				t.Error("expected precondition error_message attribute to be required")
			}
		}
	}

	// Test ephemeral lifecycle has postcondition block
	postconditionBlock, exists := ephemeralLifecycle.Body.Blocks["postcondition"]
	if !exists {
		t.Error("expected ephemeral lifecycle block to have postcondition block")
	} else {
		conditionAttr, exists := postconditionBlock.Body.Attributes["condition"]
		if !exists {
			t.Error("expected postcondition block to have condition attribute")
		} else {
			if !conditionAttr.IsRequired {
				t.Error("expected postcondition condition attribute to be required")
			}
			// Postcondition condition should be boolean type
			anyExpr, ok := conditionAttr.Constraint.(schema.AnyExpression)
			if !ok {
				t.Errorf("expected postcondition condition constraint to be AnyExpression, got %T", conditionAttr.Constraint)
			} else {
				if anyExpr.OfType != cty.Bool {
					t.Errorf("expected postcondition condition to be boolean type, got %v", anyExpr.OfType)
				}
			}
		}

		errorMessageAttr, exists := postconditionBlock.Body.Attributes["error_message"]
		if !exists {
			t.Error("expected postcondition block to have error_message attribute")
		} else {
			if !errorMessageAttr.IsRequired {
				t.Error("expected postcondition error_message attribute to be required")
			}
			// Postcondition error_message should be string type
			anyExpr, ok := errorMessageAttr.Constraint.(schema.AnyExpression)
			if !ok {
				t.Errorf("expected postcondition error_message constraint to be AnyExpression, got %T", errorMessageAttr.Constraint)
			} else {
				if anyExpr.OfType != cty.String {
					t.Errorf("expected postcondition error_message to be string type, got %v", anyExpr.OfType)
				}
			}
		}

		// Postcondition should support self references
		if postconditionBlock.Body.Extensions == nil {
			t.Error("expected postcondition block to have extensions")
		} else {
			if !postconditionBlock.Body.Extensions.SelfRefs {
				t.Error("expected postcondition block to support self references")
			}
		}
	}
}
