# Updating the Default Golang Version

This project typically upgrades its Go version for development and testing shortly after release to get the latest and greatest Go functionality.

Create an issue to cover the update.

Before beginning the update process, ensure that you review the release notes to look for any areas of possible friction when updating. Note this down within the issue.

Ensure that the following steps are completed (This best works as a checklist on the issue):

- Verify all formatting, linting, and testing works as expected
- Verify `gox` builds for all currently supported architectures
- Verify `goenv` support for the new version
- Update `README.md`
- Update `.travis.yml`
- Update CHANGELOG.md with any notes practioners need to be aware of.