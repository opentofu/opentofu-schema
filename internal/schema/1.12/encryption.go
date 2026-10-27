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

func patchEncryptionBlockSchema(b *schema.BlockSchema) *schema.BlockSchema {
	keyProviderBlock := b.Body.Blocks["key_provider"]

	if gcpKmsSchema, exists := keyProviderBlock.DependentBody[labelKey("gcp_kms")]; exists {
		gcpKmsSchema.Attributes["additional_authenticated_data"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.String},
			Description: lang.Markdown("Base64-encoded [additional authenticated data (AAD)](https://cloud.google.com/kms/docs/additional-authenticated-data) sent with both encrypt and decrypt calls."),
		}
	}

	if awsKmsSchema, exists := keyProviderBlock.DependentBody[labelKey("aws_kms")]; exists {
		awsKmsSchema.Attributes["encryption_context"] = &schema.AttributeSchema{
			Constraint:  schema.LiteralType{Type: cty.Map(cty.String)},
			Description: lang.Markdown("Optional map of key-value string pairs sent to AWS KMS as [Encryption Context](https://docs.aws.amazon.com/kms/latest/developerguide/encrypt_context.html) with every `GenerateDataKey` and `Decrypt` call."),
		}
	}

	return b
}
