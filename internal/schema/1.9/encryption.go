// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func labelKey(value string) schema.SchemaKey {
	return schema.NewSchemaKey(schema.DependencyKeys{
		Labels: []schema.LabelDependent{{Index: 0, Value: value}},
	})
}

func addEncryptedMetadataAlias(sch *schema.BodySchema) {
	sch.Attributes["encrypted_metadata_alias"] = &schema.AttributeSchema{
		Constraint:  schema.LiteralType{Type: cty.String},
		Description: lang.Markdown("Optional identifier to store metadata in encrypted state/plan files. Allows changing key provider name."),
	}
}

func patchEncryptionBlockSchema(b *schema.BlockSchema) *schema.BlockSchema {
	// Add encrypted_metadata_alias to all key providers in version 1.9
	keyProviderBlock := b.Body.Blocks["key_provider"]

	// Add encrypted_metadata_alias to pbkdf2
	addEncryptedMetadataAlias(keyProviderBlock.DependentBody[labelKey("pbkdf2")])

	// Add encrypted_metadata_alias to aws_kms
	addEncryptedMetadataAlias(keyProviderBlock.DependentBody[labelKey("aws_kms")])

	// Add encrypted_metadata_alias to gcp_kms
	addEncryptedMetadataAlias(keyProviderBlock.DependentBody[labelKey("gcp_kms")])

	// Add encrypted_metadata_alias to openbao
	addEncryptedMetadataAlias(keyProviderBlock.DependentBody[labelKey("openbao")])

	return b
}
