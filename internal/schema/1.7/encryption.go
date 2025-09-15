// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

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

func awsKmsSchema() *schema.BodySchema {
	return &schema.BodySchema{
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
			"access_key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("AWS access key ID. Can also be sourced from the `AWS_ACCESS_KEY_ID` environment variable"),
			},
			"secret_key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("AWS secret access key. Can also be sourced from the `AWS_SECRET_ACCESS_KEY` environment variable"),
			},
			"region": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("AWS region. Can also be sourced from the `AWS_REGION` or `AWS_DEFAULT_REGION` environment variables"),
			},
			"profile": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("AWS profile name to use from the shared credentials file"),
			},
			"token": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("AWS session token. Can also be sourced from the `AWS_SESSION_TOKEN` environment variable"),
			},
			"max_retries": {
				Constraint:  schema.LiteralType{Type: cty.Number},
				Description: lang.Markdown("Maximum number of times to retry API calls. Default is 5"),
			},
			"skip_credentials_validation": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				Description: lang.Markdown("Skip the credentials validation via the STS API"),
			},
			"skip_requesting_account_id": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				Description: lang.Markdown("Skip requesting the account ID"),
			},
			"sts_region": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("AWS STS region to use for assuming roles"),
			},
			"http_proxy": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("HTTP proxy URL to use for API requests"),
			},
			"https_proxy": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("HTTPS proxy URL to use for API requests"),
			},
			"no_proxy": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("Comma-separated list of hosts to exclude from proxy"),
			},
			"insecure": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				Description: lang.Markdown("Whether to explicitly allow insecure HTTPS requests"),
			},
			"use_dualstack_endpoint": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				Description: lang.Markdown("Whether to use the dual-stack endpoint for AWS services"),
			},
			"use_fips_endpoint": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				Description: lang.Markdown("Whether to use FIPS-compliant endpoints"),
			},
			"custom_ca_bundle": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("Path to a custom CA bundle to use for HTTPS requests"),
			},
			"ec2_metadata_service_endpoint": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("Address of the EC2 metadata service endpoint"),
			},
			"ec2_metadata_service_endpoint_mode": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("Protocol to use with EC2 metadata (IPv4 or IPv6)"),
			},
			"skip_metadata_api_check": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				Description: lang.Markdown("Skip the AWS metadata API check"),
			},
			"shared_credentials_files": {
				Constraint:  schema.List{Elem: schema.LiteralType{Type: cty.String}},
				Description: lang.Markdown("List of paths to shared credentials files"),
			},
			"shared_config_files": {
				Constraint:  schema.List{Elem: schema.LiteralType{Type: cty.String}},
				Description: lang.Markdown("List of paths to shared configuration files"),
			},
			"allowed_account_ids": {
				Constraint:  schema.List{Elem: schema.LiteralType{Type: cty.String}},
				Description: lang.Markdown("List of allowed AWS account IDs"),
			},
			"forbidden_account_ids": {
				Constraint:  schema.List{Elem: schema.LiteralType{Type: cty.String}},
				Description: lang.Markdown("List of forbidden AWS account IDs"),
			},
			"retry_mode": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("Specifies how many times to retry API calls. Valid values are `standard` and `adaptive`"),
			},
			"assume_role": {
				Description: lang.Markdown("Configuration for assuming an IAM role"),
				Constraint: schema.Object{
					Name: "Assume an IAM Role Configuration",
					Attributes: map[string]*schema.AttributeSchema{
						"role_arn": {
							Constraint:  schema.LiteralType{Type: cty.String},
							IsRequired:  true,
							Description: lang.Markdown("ARN of the IAM role to assume"),
						},
						"duration": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("Duration of the assumed role session (e.g. `1h`, `30m`)"),
						},
						"external_id": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("External ID to use when assuming the role"),
						},
						"policy": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("IAM policy in JSON format to apply to the assumed role session"),
						},
						"policy_arns": {
							Constraint:  schema.List{Elem: schema.LiteralType{Type: cty.String}},
							Description: lang.Markdown("List of ARNs of IAM policies to apply to the assumed role session"),
						},
						"session_name": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("Name to use for the assumed role session"),
						},
						"tags": {
							Constraint:  schema.Map{Elem: schema.LiteralType{Type: cty.String}},
							Description: lang.Markdown("Map of tags to apply to the assumed role session"),
						},
						"transitive_tag_keys": {
							Constraint:  schema.List{Elem: schema.LiteralType{Type: cty.String}},
							Description: lang.Markdown("List of tag keys to pass to subsequent assumed roles"),
						},
					},
				},
			},
			"assume_role_with_web_identity": {
				Description: lang.Markdown("Configuration for assuming an IAM role using web identity"),
				Constraint: schema.Object{
					Name: "`assume_role_with_web_identity` configuration block",
					Attributes: map[string]*schema.AttributeSchema{
						"role_arn": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("ARN of the IAM role to assume"),
						},
						"duration": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("Duration of the assumed role session (e.g. `1h`, `30m`)"),
						},
						"policy": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("IAM policy in JSON format to apply to the assumed role session"),
						},
						"policy_arns": {
							Constraint:  schema.List{Elem: schema.LiteralType{Type: cty.String}},
							Description: lang.Markdown("List of ARNs of IAM policies to apply to the assumed role session"),
						},
						"session_name": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("Name to use for the assumed role session"),
						},
						"web_identity_token": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("OAuth 2.0 access token or OpenID Connect ID token"),
						},
						"web_identity_token_file": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("Path to file containing web identity token"),
						},
					},
				},
			},
		},
		Blocks: map[string]*schema.BlockSchema{
			"endpoints": {
				Description: lang.Markdown("Configuration for custom AWS service endpoints"),
				Body: &schema.BodySchema{
					Attributes: map[string]*schema.AttributeSchema{
						"iam": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("Custom endpoint URL for AWS IAM API"),
						},
						"sts": {
							Constraint:  schema.LiteralType{Type: cty.String},
							Description: lang.Markdown("Custom endpoint URL for AWS STS API"),
						},
					},
				},
			},
		},
	}
}

func gcpKmsSchema() *schema.BodySchema {
	return &schema.BodySchema{
		Description: lang.Markdown("GCP KMS key provider uses Google Cloud Key Management Service to generate keys"),
		HoverURL:    "https://opentofu.org/docs/language/state/encryption/#gcp-kms",
		Attributes: map[string]*schema.AttributeSchema{
			"kms_encryption_key": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("[Key ID for GCP KMS](https://cloud.google.com/kms/docs/create-key#kms-create-symmetric-encrypt-decrypt-console).                          | N/A  | -                                  |"),
			},
			"key_length": {
				Constraint:  schema.LiteralType{Type: cty.Number},
				IsRequired:  true,
				Description: lang.Markdown("Number of bytes to generate as a key. Must be in range from 1 to 1024 bytes."),
			},
			"credentials": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("Local path to Google Cloud Platform account credentials in JSON format. If unset, the path uses [Google Application Default Credentials](https://developers.google.com/identity/protocols/application-default-credentials).  The provided credentials must have the Storage Object Admin role on the bucket. **Warning**: if using the Google Cloud Platform provider as well, it will also pick up the `GOOGLE_CREDENTIALS` environment variable."),
			},
			"access_token": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("A temporary **OAuth 2.0 access token** obtained from the Google Authorization server, i.e. the `Authorization: Bearer` token used to authenticate HTTP requests to GCP APIs. This is an alternative to `credentials`. If both are specified, `access_token` will be used over the `credentials` field."),
			},
			"impersonate_service_account": {
				Constraint:  schema.LiteralType{Type: cty.String},
				Description: lang.Markdown("The service account to impersonate for accessing the State Bucket. You must have `roles/iam.serviceAccountTokenCreator` role on that account for the impersonation to succeed. If you are using a delegation chain, you can specify that using the `impersonate_service_account_delegates` field. Can also be sourced from the `GOOGLE_IMPERSONATE_SERVICE_ACCOUNT` or `GOOGLE_BACKEND_IMPERSONATE_SERVICE_ACCOUNT` environment variables"),
			},
			"impersonate_service_account_delegates": {
				Constraint:  schema.List{Elem: schema.LiteralType{Type: cty.String}},
				Description: lang.Markdown("The delegation chain for an impersonating a service account as described [here](https://cloud.google.com/iam/docs/creating-short-lived-service-account-credentials#sa-credentials-delegated)."),
			},
		},
	}
}

// keyProviderTypes with their markdown descriptions for 1.7
func keyProviderTypes() map[schema.SchemaKey]*schema.BodySchema {
	return map[schema.SchemaKey]*schema.BodySchema{
		labelKey("pbkdf2"): {
			Description: lang.Markdown("PBKDF2 key provider allows you to use a long passphrase to generate a key for encryption methods such as AES-GCM"),
			HoverURL:    "https://opentofu.org/docs/language/state/encryption/#pbkdf2",
			Attributes: map[string]*schema.AttributeSchema{
				"passphrase": {
					Constraint:  schema.LiteralType{Type: cty.String},
					IsRequired:  true,
					Description: lang.Markdown("Enter a long and complex passphrase. Minimum 16 characters."),
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
			},
		},
		labelKey("aws_kms"): awsKmsSchema(),
		labelKey("gcp_kms"): gcpKmsSchema(),
		labelKey("openbao"): {
			Description: lang.Markdown("OpenBao key provider uses the OpenBao Transit Secret Engine to generate data keys (experimental)"),
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
			},
		},
	}
}

func keyProviderBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "key_provider"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName: "key_provider",
			ScopeId:      refscope.EncryptionKeyProviderScope,
			AsReference:  true,
		},
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Encryption},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type, lang.TokenModifierDependent},
				Description:            lang.PlainText("key_provider type"),
				Completable:            true,
				IsDepKey:               true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.Markdown("key_provider name"),
			},
		},
		DependentBody: keyProviderTypes(),
		Description:   lang.Markdown("Key provider configuration for encryption"),
		Body: &schema.BodySchema{
			HoverURL: "https://opentofu.org/docs/language/state/encryption/#key-providers",
		},
	}
}

// methodTypes with their markdown descriptions for 1.7
func methodTypes() map[schema.SchemaKey]*schema.BodySchema {
	return map[schema.SchemaKey]*schema.BodySchema{
		labelKey("aes_gcm"): {
			Description: lang.Markdown("AES-GCM encryption method"),
			HoverURL:    "https://opentofu.org/docs/language/state/encryption/#aes-gcm",
			Attributes: map[string]*schema.AttributeSchema{
				"keys": {
					Constraint:  schema.Reference{OfScopeId: refscope.EncryptionKeyProviderScope},
					IsRequired:  true,
					Description: lang.Markdown("Reference to a key provider"),
				},
			},
		},
		labelKey("unencrypted"): {
			Description: lang.Markdown("Unencrypted method for migration purposes"),
			HoverURL:    "https://opentofu.org/docs/language/state/encryption/#unencrypted",
			Attributes:  map[string]*schema.AttributeSchema{},
		},
	}
}

func methodBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Address: &schema.BlockAddrSchema{
			Steps: []schema.AddrStep{
				schema.StaticStep{Name: "method"},
				schema.LabelStep{Index: 0},
				schema.LabelStep{Index: 1},
			},
			FriendlyName: "method",
			ScopeId:      refscope.EncryptionMethodScope,
			AsReference:  true,
		},
		Labels: []*schema.LabelSchema{
			{
				Name:                   "type",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Type},
				Description:            lang.PlainText("method type"),
				Completable:            true,
				IsDepKey:               true,
			},
			{
				Name:                   "name",
				SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Name},
				Description:            lang.Markdown("method name"),
			},
		},
		DependentBody: methodTypes(),
		Description:   lang.Markdown("Encryption method configuration"),
		Body: &schema.BodySchema{
			HoverURL: "https://opentofu.org/docs/language/state/encryption/#methods",
		},
	}
}

func stateBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("State encryption configuration"),
		Body: &schema.BodySchema{
			HoverURL: "https://opentofu.org/docs/language/state/encryption/#configuration",
			Attributes: map[string]*schema.AttributeSchema{
				"method": {
					Constraint:  schema.Reference{OfScopeId: refscope.EncryptionMethodScope},
					IsRequired:  true,
					Description: lang.Markdown("Reference to an encryption method"),
				},
				"enforced": {
					Constraint:  schema.LiteralType{Type: cty.Bool},
					Description: lang.Markdown("Whether encryption is enforced"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"fallback": fallbackSchema(10),
			},
		},
		MaxItems: 1,
	}
}

func planBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Plan encryption configuration"),
		Body: &schema.BodySchema{
			HoverURL: "https://opentofu.org/docs/language/state/encryption/#configuration",
			Attributes: map[string]*schema.AttributeSchema{
				"method": {
					Constraint:  schema.Reference{OfScopeId: refscope.EncryptionMethodScope},
					IsRequired:  true,
					Description: lang.Markdown("Reference to an encryption method"),
				},
				"enforced": {
					Constraint:  schema.LiteralType{Type: cty.Bool},
					Description: lang.Markdown("Whether encryption is enforced"),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"fallback": fallbackSchema(10),
			},
		},
		MaxItems: 1,
	}
}

func fallbackSchema(recursiveDepth int) *schema.BlockSchema {
	block := &schema.BlockSchema{
		Description: lang.Markdown("Fallback method for reading existing encrypted data"),
		Body: &schema.BodySchema{
			Attributes: map[string]*schema.AttributeSchema{
				"method": {
					Constraint: schema.Reference{
						OfScopeId: refscope.EncryptionMethodScope,
					},
					IsRequired:  true,
					Description: lang.Markdown("Reference to a fallback encryption method"),
				},
			},
		},
		MaxItems: 1,
	}
	if recursiveDepth == 0 {
		return block
	}

	block.Body.Blocks["fallback"] = fallbackSchema(recursiveDepth - 1)
	return block
}

func remoteStateDataSourcesBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description: lang.Markdown("Remote state data sources encryption configuration"),
		Body: &schema.BodySchema{
			HoverURL: "https://opentofu.org/docs/language/state/encryption/#remote-state-data-sources",
			Blocks: map[string]*schema.BlockSchema{
				"default": {
					Description: lang.Markdown("Default encryption method for remote state data sources"),
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"method": {
								Constraint: schema.Reference{
									OfScopeId: refscope.EncryptionMethodScope,
								},
								IsRequired:  true,
								Description: lang.Markdown("Reference to an encryption method"),
							},
						},
						Blocks: map[string]*schema.BlockSchema{
							"fallback": fallbackSchema(10),
						},
					},
					MaxItems: 1,
				},
				"remote_state_data_source": {
					Description: lang.Markdown("Specific remote state data source encryption configuration"),
					Labels: []*schema.LabelSchema{
						{
							Name:        "name",
							Description: lang.Markdown("Name of the remote state data source of type *terraform_remote_state* <br> examples: *myname*, *mymodule.myname*, *mymodule.myname[0]* "),
							IsDepKey:    true,
							Completable: true, // This will be set based on the available datasource during early evaluation
						},
					},
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"method": {
								Constraint:  schema.Reference{OfScopeId: refscope.EncryptionMethodScope},
								IsRequired:  true,
								Description: lang.Markdown("Reference to an encryption method"),
							},
						},
						Blocks: map[string]*schema.BlockSchema{
							"fallback": fallbackSchema(10),
						},
					},
				},
			},
		},
		MaxItems: 1,
	}
}

func encryptionBlock() *schema.BlockSchema {
	return &schema.BlockSchema{
		Description:            lang.Markdown("State and Plan encryption configuration block"),
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Encryption},
		Body: &schema.BodySchema{
			HoverURL: "https://opentofu.org/docs/language/state/encryption/#configuration",
			Blocks: map[string]*schema.BlockSchema{
				"key_provider":              keyProviderBlock(),
				"method":                    methodBlock(),
				"state":                     stateBlock(),
				"plan":                      planBlock(),
				"remote_state_data_sources": remoteStateDataSourcesBlock(),
			},
		},
		MaxItems: 1,
	}
}
