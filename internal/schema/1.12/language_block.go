// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/opentofu/opentofu-schema/internal/schema/tokmod"
	"github.com/zclconf/go-cty/cty"
)

// languageBlockSchema returns the schema for the `language` block introduced
// in OpenTofu v1.12. It is used to configure language-level behaviors such as
// the language edition and compatibility constraints.
//
// https://opentofu.org/docs/language/settings/
func languageBlockSchema() *schema.BlockSchema {
	return &schema.BlockSchema{
		SemanticTokenModifiers: lang.SemanticTokenModifiers{tokmod.Language},
		Description: lang.Markdown("`language` block used to configure language-level behaviors of OpenTofu, " +
			"such as the language edition and compatibility constraints. Accepted only in OpenTofu v1.12 and later."),
		Body: &schema.BodySchema{
			DocsLink: &schema.DocsLink{
				URL: "https://opentofu.org/docs/language/settings/",
			},
			Attributes: map[string]*schema.AttributeSchema{
				"edition": {
					Constraint: schema.Keyword{
						Keyword: "tofu2024",
						Name:    "edition",
					},
					IsOptional: true,
					Description: lang.Markdown("Selects the language edition the module is written for. " +
						"Modules written for OpenTofu should typically not set this argument at all. " +
						"The only currently-available edition is `tofu2024`."),
				},
				"experiments": {
					Constraint:  schema.Set{},
					IsOptional:  true,
					Description: lang.Markdown("A set of experimental language features to enable. "),
				},
			},
			Blocks: map[string]*schema.BlockSchema{
				"compatible_with": {
					Description: lang.Markdown("Declares the version of OpenTofu this module is compatible with, " +
						"causing OpenTofu to return an error early if the current version does not match."),
					Body: &schema.BodySchema{
						Attributes: map[string]*schema.AttributeSchema{
							"opentofu": {
								Constraint: schema.LiteralType{Type: cty.String},
								IsOptional: true,
								Description: lang.Markdown("Version constraint specifying which versions of OpenTofu " +
									"this module is compatible with, e.g. `>= 1.12`"),
							},
						},
					},
				},
			},
		},
	}
}
