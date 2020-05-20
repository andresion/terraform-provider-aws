# Label Dictionary

<!-- non breaking spaces are to ensure that the badges are consistent. -->

| Label | Description | Automation |
|---------|-------------|----------|
| [![breaking-change][breaking-change-badge]][breaking-change]&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp; | Introduces a breaking change in current functionality; breaking changes are usually deferred to the next major release. | None |
| [![bug][bug-badge]][bug] | Addresses a defect in current functionality. | None |
| [![crash][crash-badge]][crash] | Results from or addresses a Terraform crash or kernel panic. | None |
| [![dependencies][dependencies-badge]][dependencies] | Used to indicate dependency or vendoring changes. | Added by Hashibot. |
| [![documentation][documentation-badge]][documentation] | Introduces or discusses updates to documentation. | None |
| [![enhancement][enhancement-badge]][enhancement] | Requests to existing resources that expand the functionality or scope. | None |
| [![good first issue][good-first-issue-badge]][good-first-issue] | Call to action for new contributors looking for a place to start. Smaller or straightforward issues. | None |
| [![hacktoberfest][hacktoberfest-badge]][hacktoberfest] | Call to action for Hacktoberfest (OSS Initiative). | None |
| [![hashibot ignore][hashibot-ignore-badge]][hashibot-ignore] | Issues or PRs labelled with this are ignored by Hashibot. | None |
| [![help wanted][help-wanted-badge]][help-wanted] | Call to action for contributors. Indicates an area of the codebase we’d like to expand/work on but don’t have the bandwidth inside the team. | None |
| [![needs-triage][needs-triage-badge]][needs-triage] | Waiting for first response or review from a maintainer. | Added to all new issues or PRs by GitHub action in `.github/workflows/issues.yml` or PRs by Hashibot in `.hashibot.hcl` unless they were submitted by a maintainer. |
| [![new-data-source][new-data-source-badge]][new-data-source] | Introduces a new data source. | None |
| [![new-resource][new-resource-badge]][new-resource] | Introduces a new resrouce. | None |
| [![proposal][proposal-badge]][proposal] | Proposes new design or functionality. | None |
| [![provider][provider-badge]][provider] | Pertains to the provider itself, rather than any interaction with AWS. | Added by Hashibot when the code change is in an area configured in `.hashibot.hcl` |
| [![question][question-badge]][question] | Includes a question about existing functionality; most questions will be re-routed to discuss.hashicorp.com. | None |
| [![regression][regression-badge]][regression] | Pertains to a degraded workflow resulting from an upstream patch or internal enhancement; usually categorized as a bug. | None |
| [![reinvent][reinvent-badge]][reinvent] | Pertains to a service or feature announced at reinvent. | None |
| ![service <*>][service-badge] | Indicates the service that is covered or introduced (i.e. service/s3) | Added by Hashibot when the code change matches a service definition in `.hashibot.hcl`.
| ![size%2F<*>][size-badge] | Managed by automation to categorize the size of a PR | Added by Hashibot to indicate the size of the PR. |
| [![stale][stale-badge]][stale] | Old or inactive issues managed by automation, if no further action taken these will get closed. | Added by a Github Action, configuration is found: `.github/workflows/stale.yml`. |
| [![technical-debt][technical-debt-badge]][technical-debt] | Addresses areas of the codebase that need refactoring or redesign. |  None |
| [![tests][tests-badge]][tests] | On a PR this indicates expanded test coverage. On an Issue this proposes expanded coverage or enhancement to test infrastructure. | None |
| [![thinking][thinking-badge]][thinking] | Requires additional research by the maintainers. | None |
| [![upstream-terraform][upstream-terraform-badge]][upstream-terraform] | Addresses functionality related to the Terraform core binary. | None |
| [![upstream][upstream-badge]][upstream] | Addresses functionality related to the cloud provider. | None |
| [![waiting-response][waiting-response-badge]][waiting-response] | Maintainers are waiting on response from community or contributor. | None |

[breaking-change-badge]: https://img.shields.io/badge/breaking--change-d93f0b
[breaking-change]: https://github.com/terraform-providers/terraform-provider-aws/labels/breaking-change
[bug-badge]: https://img.shields.io/badge/bug-f7c6c7
[bug]: https://github.com/terraform-providers/terraform-provider-aws/labels/bug
[crash-badge]: https://img.shields.io/badge/crash-e11d21
[crash]: https://github.com/terraform-providers/terraform-provider-aws/labels/crash
[dependencies-badge]: https://img.shields.io/badge/dependencies-fad8c7
[dependencies]: https://github.com/terraform-providers/terraform-provider-aws/labels/dependencies
[documentation-badge]: https://img.shields.io/badge/documentation-fef2c0
[documentation]: https://github.com/terraform-providers/terraform-provider-aws/labels/documentation
[enhancement-badge]: https://img.shields.io/badge/enhancement-d4c5f9
[enhancement]: https://github.com/terraform-providers/terraform-provider-aws/labels/enhancement
[good-first-issue-badge]: https://img.shields.io/badge/good%20first%20issue-128A0C
[good-first-issue]: https://github.com/terraform-providers/terraform-provider-aws/labels/good%20first%20issue
[hacktoberfest-badge]: https://img.shields.io/badge/hacktoberfest-2c0fad
[hacktoberfest]: https://github.com/terraform-providers/terraform-provider-aws/labels/hacktoberfest
[hashibot-ignore-badge]: https://img.shields.io/badge/hashibot%2Fignore-2c0fad
[hashibot-ignore]: https://github.com/terraform-providers/terraform-provider-aws/labels/hashibot-ignore
[help-wanted-badge]: https://img.shields.io/badge/help%20wanted-128A0C
[help-wanted]: https://github.com/terraform-providers/terraform-provider-aws/labels/help-wanted
[needs-triage-badge]: https://img.shields.io/badge/needs--triage-e236d7
[needs-triage]: https://github.com/terraform-providers/terraform-provider-aws/labels/needs-triage
[new-data-source-badge]: https://img.shields.io/badge/new--data--source-d4c5f9
[new-data-source]: https://github.com/terraform-providers/terraform-provider-aws/labels/new-data-source
[new-resource-badge]: https://img.shields.io/badge/new--resource-d4c5f9
[new-resource]: https://github.com/terraform-providers/terraform-provider-aws/labels/new-resource
[proposal-badge]: https://img.shields.io/badge/proposal-fbca04
[proposal]: https://github.com/terraform-providers/terraform-provider-aws/labels/proposal
[provider-badge]: https://img.shields.io/badge/provider-bfd4f2
[provider]: https://github.com/terraform-providers/terraform-provider-aws/labels/provider
[question-badge]: https://img.shields.io/badge/question-d4c5f9
[question]: https://github.com/terraform-providers/terraform-provider-aws/labels/question
[regression-badge]: https://img.shields.io/badge/regression-e11d21
[regression]: https://github.com/terraform-providers/terraform-provider-aws/labels/regression
[reinvent-badge]: https://img.shields.io/badge/reinvent-c5def5
[reinvent]: https://github.com/terraform-providers/terraform-provider-aws/labels/reinvent
[service-badge]: https://img.shields.io/badge/service%2F<*>-bfd4f2
[size-badge]: https://img.shields.io/badge/size%2F<*>-ffffff
[stale-badge]: https://img.shields.io/badge/stale-e11d21
[stale]: https://github.com/terraform-providers/terraform-provider-aws/labels/stale
[technical-debt-badge]: https://img.shields.io/badge/technical--debt-1d76db
[technical-debt]: https://github.com/terraform-providers/terraform-provider-aws/labels/technical-debt
[tests-badge]: https://img.shields.io/badge/tests-DDDDDD
[tests]: https://github.com/terraform-providers/terraform-provider-aws/labels/tests
[thinking-badge]: https://img.shields.io/badge/thinking-bfd4f2
[thinking]: https://github.com/terraform-providers/terraform-provider-aws/labels/thinking
[upstream-terraform-badge]: https://img.shields.io/badge/upstream--terraform-CCCCCC
[upstream-terraform]: https://github.com/terraform-providers/terraform-provider-aws/labels/upstream-terraform
[upstream-badge]: https://img.shields.io/badge/upstream-fad8c7
[upstream]: https://github.com/terraform-providers/terraform-provider-aws/labels/upstream
[waiting-response-badge]: https://img.shields.io/badge/waiting--response-5319e7
[waiting-response]: https://github.com/terraform-providers/terraform-provider-aws/labels/waiting-response
