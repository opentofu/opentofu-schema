# Generated OpenTofu function signatures

This package can generate function signature files for OpenTofu >= 1.6 automatically.

It is intended to run whenever OpenTofu releases a new OpenTofu version. If the OpenTofu version contains updated function signatures, it generates a new file for that version. When no changes are detected, one should commit the version bump in `gen/gen.go`.

Pre-releases are accepted with the following caveats:

 - The given pre-release version of OpenTofu is downloaded
 - The pre-release part of the version is omitted in all other contexts (e.g. function name, file name, constraint etc.)

... meaning that e.g. `1.6.0-alpha1` downloads OpenTofu `1.6.0-alpha1` but assumes compatibility with `1.6.0`. Therefore, generating signatures based on pre-releases should be used with the assumption that those signatures won't change before the final release.

## Running

Update the `tofuVersion` in `gen/gen.go` with the version of the new OpenTofu release.

Run `go generate ./internal/funcs/generated`. This command will:

1. Install the specified OpenTofu version.
1. Run `tofu metadata functions -json` to obtain the function signatures.
1. Compare a hash of the JSON string to `functionSignatureHash` in `gen/gen.go`.
   1. If there are no changes, we will stop here.
1. Create a new Go file with the function signatures of that OpenTofu version.
1. Regenerate the function signature selection in `functions.go`.

**If everything looks solid, update the `functionSignatureHash` with the one from the output, and commit all changes.**
