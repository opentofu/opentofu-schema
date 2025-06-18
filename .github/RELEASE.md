## Releasing

Releases are made on a reasonably regular basis by the maintainers, by using the Github UI. The following notes are only relevant to maintainers.

Release process:

1. On your computer, make sure you have checked out the correct branch:
   * `main` for `alpha` and `beta` releases;
   * `vX.Y` for any other releases (assuming you are releasing version `X.Y.Z`)
2. Make sure the branch is up-to-date by running `git pull`;
3. Create the correct tag: `git tag -m "X.Y.Z" vX.Y.Z` (assuming you are releasing version `X.Y.Z`)
   * If you have a GPG key, consider adding the `-s` option to create a GPG-signed tag
4. Push the tag: `git push origin vX.Y.Z`;
5. Go to the [new release page](https://github.com/opentofu/opentofu-schema/releases/new) and choose the pushed tag. On the release title, put "vX.Y.Z", then select the previous tag, clicking on "Generate release notes" in order to add an automatic description.
6. Click on "Publish Release" to finish the process.
