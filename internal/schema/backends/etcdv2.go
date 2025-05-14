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

func etcdv2Backend(v *version.Version) *schema.BodySchema {
	if v.GreaterThanOrEqual(v1_2_0) {
		return &schema.BodySchema{
			IsDeprecated: true,
			Description:  lang.Markdown("etcdv2 backend is deprecated since v1.2.0."),
		}
	}
	docsUrl := "https://www.terraform.io/docs/language/settings/backends/etcd.html" //We do not have a page for this, leaving it here
	bodySchema := &schema.BodySchema{
		Description: lang.Markdown("etcd v2.x"),
		HoverURL:    docsUrl,
		DocsLink: &schema.DocsLink{
			URL: docsUrl,
		},
		Attributes: map[string]*schema.AttributeSchema{
			"path": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("The path where to store the state"),
			},
			"endpoints": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsRequired:  true,
				Description: lang.Markdown("A space-separated list of the etcd endpoints"),
			},
			"username": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Username"),
			},
			"password": {
				Constraint:  schema.LiteralType{Type: cty.String},
				IsOptional:  true,
				Description: lang.Markdown("Password"),
			},
		},
	}

	return bodySchema
}
