// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/zclconf/go-cty/cty"

	v0_15_mod "github.com/opentofu/opentofu-schema/internal/schema/0.15"
)

var (
	FileProvisioner      = v0_15_mod.FileProvisioner
	LocalExecProvisioner = func() *schema.BodySchema {
		bodySchema := v0_15_mod.LocalExecProvisioner
		bodySchema.Attributes["quiet"] = &schema.AttributeSchema{
			Constraint:   schema.LiteralType{Type: cty.Bool},
			DefaultValue: schema.DefaultValue{Value: cty.False},
			IsOptional:   true,
			Description:  lang.Markdown("Whether to suppress script output"),
		}
		return bodySchema
	}()
	RemoteExecProvisioner = v0_15_mod.RemoteExecProvisioner
)

func ProvisionerDependentBodies(v *version.Version) map[schema.SchemaKey]*schema.BodySchema {
	return map[schema.SchemaKey]*schema.BodySchema{
		labelKey("file"):        FileProvisioner,
		labelKey("local-exec"):  LocalExecProvisioner,
		labelKey("remote-exec"): RemoteExecProvisioner,
	}
}

func labelKey(value string) schema.SchemaKey {
	return schema.NewSchemaKey(schema.DependencyKeys{
		Labels: []schema.LabelDependent{{Index: 0, Value: value}},
	})
}
