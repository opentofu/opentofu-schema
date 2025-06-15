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
	methodBlock := b.Body.Blocks["method"]

	// Add chain option to pbkdf2 and make passphrase optional
	if pbkdf2Schema, exists := keyProviderBlock.DependentBody[labelKey("pbkdf2")]; exists {
		pbkdf2Schema.Attributes["chain"] = &schema.AttributeSchema{
			Constraint:  schema.Reference{OfType: cty.DynamicPseudoType},
			Description: lang.Markdown("Receive the passphrase from another key provider. Required if `passphrase` is not specified."),
		}
		// Make passphrase optional since chain can be used instead
		pbkdf2Schema.Attributes["passphrase"].IsRequired = false
		pbkdf2Schema.Attributes["passphrase"].Description = lang.Markdown("Enter a long and complex passphrase. Required if `chain` is not specified. Minimum 16 characters.")
	}

	// Remove experimental status from OpenBao
	if openBaoSchema, exists := keyProviderBlock.DependentBody[labelKey("openbao")]; exists {
		openBaoSchema.Description = lang.Markdown("OpenBao key provider uses the OpenBao Transit Secret Engine to generate data keys")
	}

	// Add external key provider (experimental)
	keyProviderBlock.DependentBody[labelKey("external")] = &schema.BodySchema{
		Description: lang.Markdown("External key provider allows you to use external commands to generate keys (experimental)"),
		HoverURL:    "https://opentofu.org/docs/language/state/encryption/#external",
		Attributes: map[string]*schema.AttributeSchema{
			"command": {
				Constraint:  schema.List{Elem: schema.LiteralType{Type: cty.String}},
				IsRequired:  true,
				Description: lang.Markdown("External command to run in an array format, each parameter being an item in an array."),
			},
		},
	}

	// Add external method (experimental)
	methodBlock.DependentBody[labelKey("external")] = &schema.BodySchema{
		Description: lang.Markdown("External method allows you to use external commands for encryption and decryption (experimental)"),
		HoverURL:    "https://opentofu.org/docs/language/state/encryption/#external",
		Attributes: map[string]*schema.AttributeSchema{
			"encrypt_command": {
				Constraint:  schema.List{Elem: schema.LiteralType{Type: cty.String}},
				IsRequired:  true,
				Description: lang.Markdown("External command to run for encryption in an array format, each parameter being an item in an array."),
			},
			"decrypt_command": {
				Constraint:  schema.List{Elem: schema.LiteralType{Type: cty.String}},
				IsRequired:  true,
				Description: lang.Markdown("External command to run for decryption in an array format, each parameter being an item in an array."),
			},
			"keys": {
				Constraint:  schema.Reference{OfType: cty.DynamicPseudoType},
				Description: lang.Markdown("Reference to a key provider if the external command requires keys."),
			},
		},
	}

	return b
}
