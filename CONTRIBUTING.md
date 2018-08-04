# Contributing to [webwire-go](https://github.com/qbeon/webwire-go)

There are many ways to contribute to the project:
- **Reporting bugs** (ideally with solution proposals)
- **Proposing missing features** (ideally with solution proposals)
- **Identifying performance bottlenecks** (ideally with fixes).
- **Providing performance comparisons** to alternative solutions
	(such as comparing common properties like req/s, latency and throughput
	with the [standard HTTP library](https://golang.org/pkg/net/http/)
	or third-party solutions such as [valyala/fasthttp](https://github.com/valyala/fasthttp)).
- **Reviewing and improving the [documentation](https://godoc.org/github.com/qbeon/webwire-go)**
	as well as additional documentation materials.
- **Adding or proposing missing examples** to make the library easier to get started with.
- **Reviewing code** to improve code quality.
- **Spreading the word** to make more developer aware of the library and its features.


## Contribution Guidelines

### Reporting Issues
If you identify a reproducible problem in the library, unclear or missing documentation or a missing feature
then please feel free to post a new issue in the [issues section](https://github.com/qbeon/webwire-go/issues)
following the [guidelines](#writing-good-bug-reports-and-feature-requests).

Before you create a new issue, please ensure there are no similar
[open](https://github.com/qbeon/webwire-go/issues?q=is%3Aissue+is%3Aopen+)
or [closed](https://github.com/qbeon/webwire-go/issues?q=is%3Aissue+is%3Aclosed+) issues.
If you find your issue already exists, feel free to make relevant comments and add your reactions.
Please use a reaction in place of a "+1" comment:
	- 👍 for up vote
	- 👎 for down vote
	- 😕 for confusion
	- 🎉 for celebration

Once submitted, your report will be marked with relevant tags by core repository maintainers as soon as possible.
When a maintainer begins working on your issue - the issue will be assigned to him/her.

Your issue may also be scheduled to certain [milestones](https://github.com/qbeon/webwire-go/milestones)
in the process of its resolution.

### Writing Good Bug Reports and Feature Requests
- File a single issue per problem and feature request.
- Do not enumerate multiple bugs or feature requests in the same issue.
- Do not add your issue as a comment to an existing issue unless it's for the identical input.
	 Many issues look similar, but have different causes.
- Provide as much information as you can to increase the likeliness
	of someone successfully reproducing the issue and elaborating a fix.
- Include the following with each issue:
	- **Expected Behavior**: describe how you expect the library to behave eventually.
	- **Actual Behavior**: describe how the library actually behaves.
	- **Proposed Solution (optional)**: optionally describe how you would resolve this issue.
- Provide either a code snippet that demonstrates the issue or a link to a code repository
	the developers can easily pull down to recreate the issue locally.
- Use [markdown formatting](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet)
	and [code blocks](https://github.com/adam-p/markdown-here/wiki/Markdown-Cheatsheet#code-and-syntax-highlighting)
	to improve readability.

If the developers can't reproduce the issue right away they will ask for more information in the comments section.

### Posting Pull Requests
Please consider discussing your changes in an issue before submitting any pull requests!
Submit pull requests that are not attached to any issue only in those cases,
where an issue would be an obvious, unnecessary overhead.

To enable the core maintainers to quickly review and accept your pull requests,
please ensure that the following requirements are met:
- Create only one pull request per issue.
- Link the related issue.
- Never merge multiple requests in one unless they have the same root cause.
- Avoid introducing new external dependencies as much as possible.
- Keep the changes as small as possible.
- Don't mix cosmetic and behavioral code changes!
- Provide cosmetic code changes in separate requests and denote them as such.
- Add tests if existing tests don't cover the new code.
- Properly document your changes by providing a description about **what** and **why** was changed,
	even if the answers to these questions seem obvious to you.
- Properly document the changed code parts.
- Ensure code quality:
	- Ensure that the tests pass with `GOCACHE=off go test -race ./...`
	- Ensure that [go vet](https://golang.org/cmd/vet/) passes with `go vet ./...`
	- Ensure that [megacheck](https://github.com/dominikh/go-tools/tree/master/cmd/megacheck)
		passes with `megacheck ./...`
- Make sure you're following the [conventions](#conventions).


## Conventions
### Semantic Versioning
[webwire-go](https://github.com/qbeon/webwire-go) follows [semantic versioning](https://semver.org/) since v1.0.0.
We release patch versions for bug fixes, minor versions for new features,
and major versions for any breaking changes.
When we make breaking changes, we also introduce deprecation warnings in a minor version
so that our users learn about the upcoming changes and migrate their code in advance.
Every significant change is documented in the
[changelog](https://github.com/qbeon/webwire-go/blob/master/CHANGELOG.md).

### Git Commit Messages
[webwire-go](https://github.com/qbeon/webwire-go) follows the
[How to Write a Git Commit Message](https://chris.beams.io/posts/git-commit/) guidelines
by [Chris Beams](https://github.com/cbeams).

### File Naming
[webwire-go](https://github.com/qbeon/webwire-go) uses [camelCase](https://en.wikipedia.org/wiki/Camel_case)
as the file naming convention with the files names starting with a lowercase letter.

### Tools
[webwire-go](https://github.com/qbeon/webwire-go) uses [gofmt](https://golang.org/cmd/gofmt/),
[go vet](https://golang.org/cmd/vet/) and [megacheck](https://github.com/dominikh/go-tools/tree/master/cmd/megacheck)
to ensure high code quality.


## License
By contributing to [webwire-go](https://github.com/qbeon/webwire-go),
you agree that your contributions will be licensed under its MIT license.


## Contact
Feel free to contact the core repository maintainers:
- [Roman Sharkov](https://github.com/romshark)
	- **Email**: [roman.sharkov@qbeon.com](mailto:roman.sharkov@qbeon.com), [sharkov@fitcat.com](mailto:sharkov@fitcat.pro)
	- [**Gophers@Slack:** Roman Sharkov (romshark)](https://gophers.slack.com)
	- [**GolangRussian@Slack:** Roman Sharkov (romshark)](https://golang-ru.slack.com)
	- [**Telegram**: @Romshark](t.me/Romshark)
----

If you find something incorrect or missing in this document,
please [provide a pull request](#posting-pull-requests) for minor changes such as typos, or [file an issue](#reporting-issues).
