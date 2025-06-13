package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/schema/refscope"
	"github.com/opentofu/opentofu-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

func labelKey(value string) schema.SchemaKey {
	return schema.NewSchemaKey(schema.DependencyKeys{
		Labels: []schema.LabelDependent{{Index: 0, Value: value}},
	})
}

// keyProviderTypes with their markdown descriptions
func keyProviderTypes() map[schema.SchemaKey]*schema.BodySchema {
	return map[schema.SchemaKey]*schema.BodySchema{
		labelKey("pbkdf2"): &schema.BodySchema{
			Description: lang.Markdown("PBKDF2 key provider allows you to use a long passphrase to generate a key for encryption methods such as AES-GCM"),
			HoverURL:    "https://opentofu.org/docs/language/state/encryption/#pbkdf2",
			Attributes: map[string]*schema.AttributeSchema{
				"passphrase": {
					Constraint:  schema.LiteralType{Type: cty.String},
					Description: lang.Markdown("Enter a long and complex passphrase. Required if `chain` is not specified. Minimum 16 characters."),
				},
				"chain": {
					Constraint:  schema.LiteralType{Type: cty.String},
					Description: lang.Markdown("Receive the passphrase from another key provider. Required if `passphrase` is not specified."),
				},
				"key_length": {
					Constraint:  schema.LiteralType{Type: cty.Number},
					Description: lang.Markdown("Number of bytes to generate as a key. Minimum 1. Default: 32"),
				},
				"iterations": {
					Constraint:  schema.LiteralType{Type: cty.Number},
					Description: lang.Markdown("Number of iterations. Minimum 200,000. Default: 600,000"),
				},
				"salt_length": {
					Constraint:  schema.LiteralType{Type: cty.Number},
					Description: lang.Markdown("Length of the salt for the key derivation. Minimum 1. Default: 32"),
				},
				"hash_function": {
					Constraint:  schema.LiteralType{Type: cty.String},
					Description: lang.Markdown("Hash function to use: `sha256` or `sha512`. Default: sha512"),
				},
				"encrypted_metadata_alias": {
					Constraint:  schema.LiteralType{Type: cty.String},
					Description: lang.Markdown("Optional identifier to store metadata in encrypted state/plan files. Allows changing key provider name."),
				},
			},
		},
		labelKey("aws_kms"): &schema.BodySchema{
			Description: lang.Markdown("AWS KMS key provider uses Amazon Web Services Key Management Service to generate keys"),
			HoverURL:    "https://opentofu.org/docs/language/state/encryption/#aws-kms",
			Attributes: map[string]*schema.AttributeSchema{
				"kms_key_id": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsRequired:  true,
					Description: lang.Markdown("Key ID for AWS KMS"),
				},
				"key_spec": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsRequired:  true,
					Description: lang.Markdown("Key spec for AWS KMS. Adapt this to your encryption method (e.g. `AES_256`)"),
				},
				"encrypted_metadata_alias": {
					Constraint:  schema.LiteralType{Type: cty.String},
					Description: lang.Markdown("Optional identifier to store metadata in encrypted state/plan files. Allows changing key provider name."),
				},
			},
		},
		labelKey("gcp_kms"): &schema.BodySchema{
			Description: lang.Markdown("GCP KMS key provider uses Google Cloud Key Management Service to generate keys"),
			HoverURL:    "https://opentofu.org/docs/language/state/encryption/#gcp-kms",
			Attributes: map[string]*schema.AttributeSchema{
				"kms_encryption_key": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsRequired:  true,
					Description: lang.Markdown("Key ID for GCP KMS"),
				},
				"key_length": {
					Constraint:  schema.LiteralType{Type: cty.Number},
					IsRequired:  true,
					Description: lang.Markdown("Number of bytes to generate as a key. Must be in range from 1 to 1024 bytes."),
				},
				"encrypted_metadata_alias": {
					Constraint:  schema.LiteralType{Type: cty.String},
					Description: lang.Markdown("Optional identifier to store metadata in encrypted state/plan files. Allows changing key provider name."),
				},
			},
		},
		labelKey("openbao"): &schema.BodySchema{
			Description: lang.Markdown("OpenBao key provider uses the OpenBao Transit Secret Engine to generate data keys"),
			HoverURL:    "https://opentofu.org/docs/language/state/encryption/#openbao",
			Attributes: map[string]*schema.AttributeSchema{
				"key_name": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsRequired:  true,
					Description: lang.Markdown("Name of the transit encryption key to use to encrypt/decrypt the datakey. Pre-configure it in your OpenBao server."),
				},
				"token": {
					Constraint:  schema.LiteralType{Type: cty.String},
					Description: lang.Markdown("Authorization Token to use when accessing OpenBao API. Can be read from `BAO_TOKEN` environment variable."),
				},
				"address": {
					Constraint:  schema.LiteralType{Type: cty.String},
					Description: lang.Markdown("OpenBao server address to access the API. Can be read from `BAO_ADDR` environment variable. Default: https://127.0.0.1:8200"),
				},
				"transit_engine_path": {
					Constraint:  schema.LiteralType{Type: cty.String},
					Description: lang.Markdown("Path at which the Transit Secret Engine is enabled in OpenBao. Default: /transit"),
				},
				"key_length": {
					Constraint:  schema.LiteralType{Type: cty.Number},
					Description: lang.Markdown("Number of bytes to generate as a key. Available options are 16, 32 or 64 bytes. Default: 32"),
				},
				"encrypted_metadata_alias": {
					Constraint:  schema.LiteralType{Type: cty.String},
					Description: lang.Markdown("Optional identifier to store metadata in encrypted state/plan files. Allows changing key provider name."),
				},
			},
		},
	}
}

func keyProviderBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "encryption"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName: "encryption",
			ScopeId:      refscope.EncryptionScope,
			AsReference:  true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Encryption},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.PlainText("key_provider type"),
				Completable:            true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type},
				Description:            lang.Markdown("key_provider name"),
			},
		},
		DependentBody: keyProviderTypes(),

		Description: lang.Markdown("TODO"),
		Body: &schema.BodySchema{
			HoverURL: "https://opentofu.org/docs/language/state/encryption/#key-providers",
			//TODO
		},
		MinItems: 1,
		MaxItems: 1,
	}
}

func encryptionBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("State and Plan encryption configuration block"),
		Body: &schema.BodySchema{
			HoverURL: "https://opentofu.org/docs/language/state/encryption/#configuration",
			Blocks: map[string]*schema.BlockSchema{
				"key_provider": keyProviderBlock(),
			},
		},
	}
}
