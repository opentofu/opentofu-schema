// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl-lang/schema"
	tfmod "github.com/opentofu/opentofu-schema/module"
)

type FunctionsMerger struct {
	coreFunctions map[string]schema.FunctionSignature
	tofuVersion   *version.Version
	stateReader   StateReader
}

func NewFunctionsMerger(coreFunctions map[string]schema.FunctionSignature) *FunctionsMerger {
	return &FunctionsMerger{
		coreFunctions: coreFunctions,
	}
}

func (m *FunctionsMerger) SetStateReader(mr StateReader) {
	m.stateReader = mr
}

func (m *FunctionsMerger) SetTofuVersion(v *version.Version) {
	m.tofuVersion = v
}

func (m *FunctionsMerger) FunctionsForModule(meta *tfmod.Meta) (map[string]schema.FunctionSignature, error) {
	if m.coreFunctions == nil {
		return nil, coreFunctionsRequiredErr{}
	}

	if meta == nil {
		return m.coreFunctions, nil
	}

	if m.stateReader == nil {
		return m.coreFunctions, nil
	}

	if m.tofuVersion.LessThan(v1_8) {
		return m.coreFunctions, nil
	}

	mergedFunctions := make(map[string]schema.FunctionSignature, len(m.coreFunctions))
	for fName, fSig := range m.coreFunctions {
		mergedFunctions[fName] = *fSig.Copy()
	}

	providerRefs := ProviderReferences(meta.ProviderReferences)

	for pAddr, pVersionCons := range meta.ProviderRequirements {
		pSchema, err := m.stateReader.ProviderSchema(meta.Path, pAddr, pVersionCons)
		if err != nil {
			continue
		}

		refs := providerRefs.ReferencesOfProvider(pAddr)

		for _, localRef := range refs {
			for fName, fSig := range pSchema.Functions {
				mergedFunctions[fmt.Sprintf("provider::%s::%s", localRef.LocalName, fName)] = *fSig.Copy()
			}
		}
	}

	return mergedFunctions, nil
}
