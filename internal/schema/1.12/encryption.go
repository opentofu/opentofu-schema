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

func patchAzureVaultEncryption(b *schema.BlockSchema) *schema.BlockSchema {
	// Add azure vault as a key provider in version 1.12
	keyProviderBlock := b.Body.Blocks["key_provider"]

	keyProviderBlock.DependentBody[labelKey("azure_vault")] = azureVaultSchema()

	return b
}

func azureVaultSchema() *schema.BodySchema {
	return &schema.BodySchema{
		Description: lang.Markdown("Azure Vault key provider uses Azure Vault to generate keys"),
		HoverURL:    "https://opentofu.org/docs/language/state/encryption/#azure-vault",
		Attributes: map[string]*schema.AttributeSchema{
			"vault_uri": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("Vault URI in Azure. Format: `https://{vault-name}.vault.azure.net`"),
			},
			"vault_key_name": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The name of the key in the specified Azure Vault."),
			},
			"key_length": {
				Constraint:  schema.LiteralType{Type: cty.Number},
				IsRequired:  true,
				Description: lang.Markdown("Number of bytes to generate as a key. Must be at least 1."),
			},
			"symmetric": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				Description: lang.Markdown("Optional boolean signifier that the provided key is symmetric (HSM only)"),
			},
			"symmetric_key_size": {
				Constraint:  schema.LiteralType{Type: cty.Number},
				Description: lang.Markdown("The size of the symmetric key (128, 192, or 256). Required when `symmetric` is true."),
			},
			// Add encrypted_metadata_alias to azure_vault
			"encrypted_metadata_alias": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("Optional identifier to store metadata in encrypted state/plan files. Allows changing key provider name."),
			},
		},
	}
}
