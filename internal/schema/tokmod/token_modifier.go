// Copyright (c) The OpenTofu Authors
// SPDX-License-Identifier: MPL-2.0
// Copyright (c) 2024 HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package tokmod

import (
	"github.com/hashicorp/hcl-lang/lang"
)

var (
	Data              = lang.SemanticTokenModifier("opentofu-data")
	Locals            = lang.SemanticTokenModifier("opentofu-locals")
	Module            = lang.SemanticTokenModifier("opentofu-module")
	Output            = lang.SemanticTokenModifier("opentofu-output")
	Provider          = lang.SemanticTokenModifier("opentofu-provider")
	Resource          = lang.SemanticTokenModifier("opentofu-resource")
	Ephemeral         = lang.SemanticTokenModifier("opentofu-ephemeral")
	Provisioner       = lang.SemanticTokenModifier("opentofu-provisioner")
	Connection        = lang.SemanticTokenModifier("opentofu-connection")
	Variable          = lang.SemanticTokenModifier("opentofu-variable")
	Terraform         = lang.SemanticTokenModifier("opentofu-terraform")
	Backend           = lang.SemanticTokenModifier("opentofu-backend")
	Name              = lang.SemanticTokenModifier("opentofu-name")
	Type              = lang.SemanticTokenModifier("opentofu-type")
	RequiredProviders = lang.SemanticTokenModifier("opentofu-requiredProviders")
	Encryption        = lang.SemanticTokenModifier("opentofu-encryption")
	Run               = lang.SemanticTokenModifier("opentofu-run")
	Variables         = lang.SemanticTokenModifier("opentofu-variables")
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
	Run,
	Terraform,
	Type,
	Variable,
	Encryption,
	Variables,
}
