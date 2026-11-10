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

func getKeyProviderBody(t *testing.T, ver string, provider string) *schema.BodySchema {
	t.Helper()
	s, err := CoreModuleSchemaForVersion(version.Must(version.NewVersion(ver)))
	if err != nil {
		t.Fatal(err)
	}
	key := schema.NewSchemaKey(schema.DependencyKeys{
		Labels: []schema.LabelDependent{{Index: 0, Value: provider}},
	})
	body, exists := s.Blocks["terraform"].Body.Blocks["encryption"].Body.Blocks["key_provider"].DependentBody[key]
	if !exists {
		t.Fatalf("%s dependent body not found", provider)
	}
	return body
}

func TestCoreModuleSchemaForVersion_v1_12_encryption_patches(t *testing.T) {
	tests := []struct {
		name      string
		provider  string
		attribute string
		wantType  cty.Type
	}{
		{"gcp_kms_aad", "gcp_kms", "additional_authenticated_data", cty.String},
		{"aws_kms_encryption_context", "aws_kms", "encryption_context", cty.Map(cty.String)},
	}

	for _, tt := range tests {
		t.Run(tt.name+"_present_in_1.12.4", func(t *testing.T) {
			body := getKeyProviderBody(t, "1.12.4", tt.provider)
			attr, exists := body.Attributes[tt.attribute]
			if !exists {
				t.Fatalf("expected %s to have %s in v1.12.4", tt.provider, tt.attribute)
			}
			if attr.IsRequired {
				t.Errorf("expected %s to be optional", tt.attribute)
			}
			lt, ok := attr.Constraint.(schema.LiteralType)
			if !ok {
				t.Fatalf("expected LiteralType constraint for %s", tt.attribute)
			}
			if !lt.Type.Equals(tt.wantType) {
				t.Errorf("expected %s, got %s", tt.wantType.FriendlyName(), lt.Type.FriendlyName())
			}
			if attr.Description.Value == "" {
				t.Errorf("expected %s to have a description", tt.attribute)
			}
		})

		t.Run(tt.name+"_absent_before_1.12.4", func(t *testing.T) {
			body := getKeyProviderBody(t, "1.12.0", tt.provider)
			if _, exists := body.Attributes[tt.attribute]; exists {
				t.Errorf("expected %s to NOT have %s before v1.12.4", tt.provider, tt.attribute)
			}
		})
	}
}
