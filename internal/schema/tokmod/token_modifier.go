// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tokmod

import (
	"github.com/hashicorp/hcl-lang/lang"
)

var (
	Data              = lang.SemanticTokenModifier("tofu-data")
	Locals            = lang.SemanticTokenModifier("tofu-locals")
	Module            = lang.SemanticTokenModifier("tofu-module")
	Output            = lang.SemanticTokenModifier("tofu-output")
	Provider          = lang.SemanticTokenModifier("tofu-provider")
	Resource          = lang.SemanticTokenModifier("tofu-resource")
	Provisioner       = lang.SemanticTokenModifier("tofu-provisioner")
	Connection        = lang.SemanticTokenModifier("tofu-connection")
	Variable          = lang.SemanticTokenModifier("tofu-variable")
	Terraform         = lang.SemanticTokenModifier("tofu-terraform")
	Backend           = lang.SemanticTokenModifier("tofu-backend")
	Name              = lang.SemanticTokenModifier("tofu-name")
	Type              = lang.SemanticTokenModifier("tofu-type")
	RequiredProviders = lang.SemanticTokenModifier("tofu-requiredProviders")
)

var SupportedModifiers = []lang.SemanticTokenModifier{
	Backend,
	Connection,
	Data,
	Locals,
	Module,
	Name,
	Output,
	Provider,
	Provisioner,
	RequiredProviders,
	Resource,
	Terraform,
	Type,
	Variable,
}
