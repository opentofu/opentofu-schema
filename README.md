# opentofu-schema [WIP]
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopentofu%2Fopentofu-schema.svg?type=shield)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopentofu%2Fopentofu-schema?ref=badge_shield)


This library helps assembling a complete [`hcl-lang`](https://github.com/hashicorp/hcl-lang)
schema for decoding OpenTofu config based on static OpenTofu core schema
and relevant provider schemas.

**There is more than one schema?**

Yes.

 - OpenTofu Core defines top-level schema
   - `provider`, `resource` or `data` blocks incl. meta attributes, such as `alias` or `count`
   - `variable`, `output` blocks etc.
 - each OpenTofu provider defines its own schema for the body of some of these blocks
   - attributes and nested blocks inside `resource`, `data` or `provider` blocks

Each of these can also differ between (core / provider) version.

## Current Status

This project is in use by the OpenTofu Language Server and could _in theory_
be used by other projects which need to decode the _whole_ configuration.

However it has not been tested in any other scenarios.

Please note that this library depends on [`hcl-lang`](https://github.com/hashicorp/hcl-lang)
which itself is not considered stable yet.

**Breaking changes may be introduced.**

## How It Works

### Usage

```go
import (
	tfschema "github.com/opentofu/opentofu-schema/schema"
	"github.com/hashicorp/terraform-json"
)

// parse files e.g. via hclsyntax
parsedFiles := map[string]*hcl.File{ /* ... */ }

// obtain relevant core schema
coreSchema := tfschema.UniversalCoreModuleSchema()

// obtain relevant provider schemas e.g. via tofu-exec
// and marshal them into terraform-json type
providerSchemas := &tfjson.ProviderSchemas{ /* ... */ }

mergedSchema, err := tfschema.MergeCoreWithJsonProviderSchemas(parsedFiles, coreSchema, providerSchemas)
if err != nil {
	// ...
}

```

### Provider Schemas

The only reliable way of obtaining provider schemas at the time of writing is via
OpenTofu CLI by running `tofu providers schema -json`.

[`tofu-exec`](https://github.com/opentofu/tofu-exec) can help automating
the process of obtaining the schema.

[`terraform-json`](https://github.com/hashicorp/terraform-json) provides types
that the JSON output can be marshalled into, which also used by `tofu-exec`
and is considered as standard way of representing the output.


#### Known Issues

At the time of writing there is a known issue affecting the above command
where it requires the following to be true in order to produce schemas:

 - configuration is valid (e.g. contains no incomplete blocks)
 - authentication is provided for any remote backend

Read more at [hashicorp/terraform#24261](https://github.com/hashicorp/terraform/issues/24261).

Other ways of obtaining schemas are also being explored.



## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fopentofu%2Fopentofu-schema.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fopentofu%2Fopentofu-schema?ref=badge_large)