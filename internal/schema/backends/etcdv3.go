// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package backends

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"
)

func etcdv3Backend(v *version.Version) *schema.BodySchema {
	if v.GreaterThanOrEqual(v1_2_0) {
		return &schema.BodySchema{
			IsDeprecated: true,
			Description:  lang.Markdown("etcdv2 backend is deprecated since v1.2.0."),
		}
	}
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/etcdv3.html" //We do not have a page for this, leaving it here
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("etcd v3"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"endpoints": {
				Constraint: schema.List{
					Elem:     schema.LiteralType{Type: cty.String},
					MinItems: 1,
				},
				IsRequired:  true,
				Description: lang.Markdown("Endpoints for the etcd cluster."),
			},

			"username": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Username used to connect to the etcd cluster."),
			},

			"password": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Password used to connect to the etcd cluster."),
			},

			"prefix": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("An optional prefix to be added to keys when to storing state in etcd."),
			},

			"lock": {
				Constraint:  schema.LiteralType{Type: cty.Bool},
				IsOptional:  true,
				Description: lang.Markdown("Whether to lock state access."),
			},

			"cacert_path": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The path to a PEM-encoded CA bundle with which to verify certificates of TLS-enabled etcd servers."),
			},

			"cert_path": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The path to a PEM-encoded certificate to provide to etcd for secure client identification."),
			},

			"key_path": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("The path to a PEM-encoded key to provide to etcd for secure client identification."),
			},
		},
	}

	return bodySchema
}
