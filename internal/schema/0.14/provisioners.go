// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"

	v013_mod "github.com/opentofu/opentofu-schema/internal/schema/0.13"
)

var (
	FileProvisioner       = v013_mod.FileProvisioner
	LocalExecProvisioner  = v013_mod.LocalExecProvisioner
	RemoteExecProvisioner = v013_mod.RemoteExecProvisioner
)

func ConnectionDependentBodies(v *version.Version) map[schema.SchemaKey]*schema.BodySchema {
	return v013_mod.ConnectionDependentBodies(v)
}

var ProvisionerDependentBodies = map[schema.SchemaKey]*schema.BodySchema{
	labelKey("file"):        FileProvisioner,
	labelKey("local-exec"):  LocalExecProvisioner,
	labelKey("remote-exec"): RemoteExecProvisioner,

	// Vendor provisioners are deprecated in 0.13.4+
	// Some of these provisioners have complex schemas
	// but we can at least helpfully list their names
	labelKey("chef"):            {IsDeprecated: true},
	labelKey("salt-masterless"): {IsDeprecated: true},
	labelKey("habitat"):         {IsDeprecated: true},
	labelKey("puppet"):          {IsDeprecated: true},
}

func labelKey(value string) schema.SchemaKey {
	return schema.NewSchemaKey(schema.DependencyKeys{
		Labels: []schema.LabelDependent{{Index: 0, Value: value}},
	})
}
