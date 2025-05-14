// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v012_mod "github.com/opentofu/opentofu-schema/internal/schema/0.12"
)

var (
	FileProvisioner       = v012_mod.FileProvisioner
	LocalExecProvisioner  = v012_mod.LocalExecProvisioner
	RemoteExecProvisioner = v012_mod.RemoteExecProvisioner
)

func ConnectionDependentBodies(v *version.Version) map[schema.SchemaKey]*schema.BodySchema {
	return v012_mod.ConnectionDependentBodies(v)
}

func ProvisionerDependentBodies(v *version.Version) map[schema.SchemaKey]*schema.BodySchema {
	m := map[schema.SchemaKey]*schema.BodySchema{
		labelKey("file"):        FileProvisioner,
		labelKey("local-exec"):  LocalExecProvisioner,
		labelKey("remote-exec"): RemoteExecProvisioner,
	}

	// Vendor provisioners are deprecated in 0.13.4+
	// Some of these provisioners have complex schemas
	// but we can at least helpfully list their names
	if v.GreaterThanOrEqual(v0_13_4) {
		m[labelKey("chef")] = &schema.BodySchema{IsDeprecated: true}
		m[labelKey("salt-masterless")] = &schema.BodySchema{IsDeprecated: true}
		m[labelKey("habitat")] = &schema.BodySchema{IsDeprecated: true}
		m[labelKey("puppet")] = &schema.BodySchema{IsDeprecated: true}
	} else {
		m[labelKey("chef")] = &schema.BodySchema{}
		m[labelKey("salt-masterless")] = &schema.BodySchema{}
		m[labelKey("habitat")] = &schema.BodySchema{}
		m[labelKey("puppet")] = &schema.BodySchema{}
	}

	return m
}

func labelKey(value string) schema.SchemaKey {
	return schema.NewSchemaKey(schema.DependencyKeys{
		Labels: []schema.LabelDependent{{Index: 0, Value: value}},
	})
}
