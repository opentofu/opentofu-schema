// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package earlydecoder

import (
	"maps"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hcl/v2"
	"github.com/opentofu/opentofu-schema/module"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

// resolveStaticModuleCalls evaluates the source and version of module calls
// that reference variables or locals.
func resolveStaticModuleCalls(mod *decodedModule) {
	// If we have none of the stored expressions to evaluate in this "second pass"
	// then we can just exit early
	if len(mod.moduleSourceExprs) == 0 && len(mod.moduleVersionExprs) == 0 {
		return
	}

	vars := make(map[string]cty.Value, len(mod.Variables))
	for name, v := range mod.Variables {
		if v.DefaultValue != cty.NilVal {
			vars[name] = v.DefaultValue
		}
	}

	locals := resolveLocals(mod.localExprs, vars)

	evalCtx := &hcl.EvalContext{
		Variables: map[string]cty.Value{
			"var":   objectOrEmpty(vars),
			"local": objectOrEmpty(locals),
		},
	}

	for name, expr := range mod.moduleSourceExprs {
		mc, ok := mod.ModuleCalls[name]
		if !ok || mc.RawSourceAddr != "" {
			continue
		}
		if source, ok := evalStaticString(expr, evalCtx); ok && source != "" {
			mc.RawSourceAddr = source
			mc.SourceAddr = module.ParseModuleSourceAddr(source)
		}
	}

	for name, expr := range mod.moduleVersionExprs {
		mc, ok := mod.ModuleCalls[name]
		if !ok || len(mc.Version) > 0 {
			continue
		}
		if versionStr, ok := evalStaticString(expr, evalCtx); ok && versionStr != "" {
			if vc, err := version.NewConstraint(versionStr); err == nil {
				mc.Version = vc
			}
		}
	}
}

func resolveLocals(exprs map[string]hcl.Expression, vars map[string]cty.Value) map[string]cty.Value {
	resolved := make(map[string]cty.Value, len(exprs))
	if len(exprs) == 0 {
		return resolved
	}

	remaining := make(map[string]hcl.Expression, len(exprs))
	maps.Copy(remaining, exprs)

	for len(remaining) > 0 {
		ctx := &hcl.EvalContext{
			Variables: map[string]cty.Value{
				"var":   objectOrEmpty(vars),
				"local": objectOrEmpty(resolved),
			},
		}

		progressed := false
		for name, expr := range remaining {
			val, diags := expr.Value(ctx)
			if diags.HasErrors() || val.IsNull() || !val.IsWhollyKnown() {
				continue
			}
			resolved[name] = val
			delete(remaining, name)
			progressed = true
		}
		// Keep going until we cant progress anymore
		if !progressed {
			break
		}
	}

	return resolved
}

// evalStaticString evaluates a string, returns false for second return val if it cant resolve
func evalStaticString(expr hcl.Expression, evalCtx *hcl.EvalContext) (string, bool) {
	val, diags := expr.Value(evalCtx)
	if diags.HasErrors() || val.IsNull() || !val.IsWhollyKnown() {
		return "", false
	}

	val, err := convert.Convert(val, cty.String)
	if err != nil {
		return "", false
	}

	return val.AsString(), true
}

func objectOrEmpty(m map[string]cty.Value) cty.Value {
	if len(m) == 0 {
		return cty.EmptyObjectVal
	}
	return cty.ObjectVal(m)
}
