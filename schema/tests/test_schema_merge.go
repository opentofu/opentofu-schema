// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfmod "github.com/opentofu/opentofu-schema/module"
	tfschema "github.com/opentofu/opentofu-schema/schema"
	tftest "github.com/opentofu/opentofu-schema/test"
	tfaddr "github.com/opentofu/registry-address"
)

type TestSchemaMerger struct {
	coreSchema  *schema.BodySchema
	stateReader StateReader
}

// StateReader exposes a set of methods to read data from the internal language server state
type StateReader interface {
	// ProviderSchema returns the schema for a provider we have stored in memory. The can come
	// from different sources.
	ProviderSchema(modPath string, addr tfaddr.Provider, vc version.Constraints) (*tfschema.ProviderSchema, error)

	// LocalModuleMeta returns the module meta data for a local module. This is the result
	// of the [earlydecoder] when processing module files
	LocalModuleMeta(modPath string) (*tfmod.Meta, error)
}

func NewTestSchemaMerger(coreSchema *schema.BodySchema) *TestSchemaMerger {
	return &TestSchemaMerger{
		coreSchema: coreSchema,
	}
}

func (m *TestSchemaMerger) SetStateReader(mr StateReader) {
	m.stateReader = mr
}

func (m *TestSchemaMerger) SchemaForTest(meta *tftest.Meta) (*schema.BodySchema, error) {
	if m.coreSchema == nil {
		return nil, tfschema.CoreSchemaRequiredErr{}
	}

	if meta == nil {
		return m.coreSchema, nil
	}

	if m.stateReader == nil {
		return m.coreSchema, nil
	}

	mergedSchema := m.coreSchema.Copy()

	// TODO merge mock_provider blocks - use the label as dependency key AND the source if defined
	// TODO merge nested mock_resource blocks - use the label as dependency key
	// TODO merge nested mock_data blocks - use the label as dependency key
	// TODO merge run - module blocks - use the source as dependency key
	// TODO merge variables - source them from the OpenTofu module meta
	// TODO merge provider - source them from the OpenTofu module meta requirements

	return mergedSchema, nil
}
