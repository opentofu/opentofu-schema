// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package schema

import (
	"fmt"
	"sort"

	"github.com/hashicorp/hcl-lang/lang"
	"github.com/hashicorp/hcl-lang/schema"
	"github.com/hashicorp/hcl/v2"
	"github.com/opentofu/opentofu-schema/internal/schema/refscope"
	"github.com/opentofu/opentofu-schema/module"
	"github.com/opentofu/opentofu-schema/registry"
	tfaddr "github.com/opentofu/registry-address"
	"github.com/zclconf/go-cty/cty"
)

func schemaForDependentRegistryModuleBlock(module module.DeclaredModuleCall, modMeta *registry.ModuleData) (*schema.BodySchema, error) {
	attributes := make(map[string]*schema.AttributeSchema, 0)

	for _, input := range modMeta.Inputs {
		aSchema := &schema.AttributeSchema{
			Description: input.Description,
		}
		if input.Required {
			aSchema.IsRequired = true
		} else {
			aSchema.IsOptional = true
		}

		typ := input.Type
		if typ == cty.NilType {
			typ = cty.DynamicPseudoType
		}
		aSchema.Constraint = convertAttributeTypeToConstraint(typ)

		attributes[input.Name] = aSchema
	}

	bodySchema := &schema.BodySchema{
		Attributes: attributes,
	}

	if module.LocalName == "" {
		// avoid creating output refs if we don't have reference name
		return bodySchema, nil
	}

	modOutputTypes := make(map[string]cty.Type, 0)
	targetableOutputs := make(schema.Targetables, 0)

	for _, output := range modMeta.Outputs {
		addr := lang.Address{
			lang.RootStep{Name: "module"},
			lang.AttrStep{Name: module.LocalName},
			lang.AttrStep{Name: output.Name},
		}

		targetable := &schema.Targetable{
			Address:     addr,
			AsType:      cty.DynamicPseudoType,
			ScopeId:     refscope.ModuleScope,
			Description: output.Description,
			// The Registry API doesn't tell us anything more about output type structure
			// so we cannot target nested fields within objects, maps or lists
		}

		modOutputTypes[output.Name] = cty.DynamicPseudoType
		targetableOutputs = append(targetableOutputs, targetable)
	}

	sort.Sort(targetableOutputs)

	addr := lang.Address{
		lang.RootStep{Name: "module"},
		lang.AttrStep{Name: module.LocalName},
	}
	bodySchema.TargetableAs = append(bodySchema.TargetableAs, &schema.Targetable{
		Address:           addr,
		ScopeId:           refscope.ModuleScope,
		AsType:            cty.Object(modOutputTypes),
		NestedTargetables: targetableOutputs,
	})

	sourceAddr, ok := module.SourceAddr.(tfaddr.Module)
	if ok && sourceAddr.Package.Host == "registry.opentofu.org" {
		versionStr := ""
		if modMeta.Version == nil {
			versionStr = "latest"
		} else {
			versionStr = fmt.Sprintf("v%s", modMeta.Version.String())
		}

		bodySchema.DocsLink = &schema.DocsLink{
			URL: fmt.Sprintf(
				`https://search.opentofu.org/module/%s/%s`,
				sourceAddr.Package.ForRegistryProtocol(),
				versionStr,
			),
		}
	}

	return bodySchema, nil
}

func schemaForDependentModuleBlock(module module.DeclaredModuleCall, modMeta *module.Meta) (*schema.BodySchema, error) {
	attributes := make(map[string]*schema.AttributeSchema, 0)

	for name, modVar := range modMeta.Variables {
		varType := modVar.Type
		if varType == cty.NilType {
			varType = cty.DynamicPseudoType
		}
		aSchema := moduleVarToAttribute(modVar)
		aSchema.Constraint = convertAttributeTypeToConstraint(varType)
		aSchema.OriginForTarget = &schema.PathTarget{
			Address: schema.Address{
				schema.StaticStep{Name: "var"},
				schema.AttrNameStep{},
			},
			Path: lang.Path{
				Path:       modMeta.Path,
				LanguageID: ModuleLanguageID,
			},
			Constraints: schema.Constraints{
				ScopeId: refscope.VariableScope,
				Type:    varType,
			},
		}

		attributes[name] = aSchema
	}

	bodySchema := &schema.BodySchema{
		Attributes: attributes,
	}

	if module.LocalName == "" {
		// avoid creating output refs if we don't have reference name
		return bodySchema, nil
	}

	modOutputTypes := make(map[string]cty.Type, 0)
	modOutputVals := make(map[string]cty.Value, 0)
	targetableOutputs := make(schema.Targetables, 0)
	impliedOrigins := make(schema.ImpliedOrigins, 0)

	for name, output := range modMeta.Outputs {
		addr := lang.Address{
			lang.RootStep{Name: "module"},
			lang.AttrStep{Name: module.LocalName},
			lang.AttrStep{Name: name},
		}

		typ := cty.DynamicPseudoType
		if !output.Value.IsNull() {
			typ = output.Value.Type()
		}

		targetable := &schema.Targetable{
			Address:           addr,
			ScopeId:           refscope.ModuleScope,
			AsType:            typ,
			IsSensitive:       output.IsSensitive,
			NestedTargetables: schema.NestedTargetablesForValue(addr, refscope.ModuleScope, output.Value),
		}
		if output.Description != "" {
			targetable.Description = lang.PlainText(output.Description)
		}

		targetableOutputs = append(targetableOutputs, targetable)

		modOutputTypes[name] = typ
		modOutputVals[name] = output.Value

		impliedOrigins = append(impliedOrigins, schema.ImpliedOrigin{
			OriginAddress: lang.Address{
				lang.RootStep{Name: "module"},
				lang.AttrStep{Name: module.LocalName},
				lang.AttrStep{Name: name},
			},
			TargetAddress: lang.Address{
				lang.RootStep{Name: "output"},
				lang.AttrStep{Name: name},
			},
			Path: lang.Path{
				Path:       modMeta.Path,
				LanguageID: ModuleLanguageID,
			},
			Constraints: schema.Constraints{
				ScopeId: refscope.OutputScope,
			},
		})
	}

	bodySchema.ImpliedOrigins = impliedOrigins

	sort.Sort(targetableOutputs)

	addr := lang.Address{
		lang.RootStep{Name: "module"},
		lang.AttrStep{Name: module.LocalName},
	}
	bodySchema.TargetableAs = append(bodySchema.TargetableAs, &schema.Targetable{
		Address:           addr,
		ScopeId:           refscope.ModuleScope,
		AsType:            cty.Object(modOutputTypes),
		NestedTargetables: targetableOutputs,
	})

	if len(modMeta.Filenames) > 0 {
		filename := modMeta.Filenames[0]

		// Prioritize main.tf based on best practices as documented at
		if sliceContains(modMeta.Filenames, "main.tf") {
			filename = "main.tf"
		}

		bodySchema.Targets = &schema.Target{
			Path: lang.Path{
				Path:       modMeta.Path,
				LanguageID: "opentofu",
			},
			Range: hcl.Range{
				Filename: filename,
				Start:    hcl.InitialPos,
				End:      hcl.InitialPos,
			},
		}
	}

	registryAddr, ok := module.SourceAddr.(tfaddr.Module)
	if ok && registryAddr.Package.Host == "registry.opentofu.org" {
		versionStr := ""
		if module.Version == nil {
			versionStr = "latest"
		} else {
			versionStr = fmt.Sprintf("v%s", module.Version.String())
		}

		bodySchema.DocsLink = &schema.DocsLink{
			URL: fmt.Sprintf(
				`https://search.opentofu.org/module/%s/%s`,
				registryAddr.Package.ForRegistryProtocol(),
				versionStr,
			),
		}
	}

	return bodySchema, nil
}

func sliceContains(slice []string, value string) bool {
	for _, val := range slice {
		if val == value {
			return true
		}
	}
	return false
}
