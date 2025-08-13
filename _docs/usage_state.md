---
title: "State"
path: "usage/state"
---

All state managed by Ocuroot is stored as JSON documents and organized by Reference, a URI-compatible string
that describes config files within source repos, Releases, Deployments, Environments and Custom State.

State can be queried and manipulated using the `ocuroot state` commands. Run `ocuroot state --help` for more
information

References are of the form:

```
[repo]/-/[path]/@[release]/[subpath]#[fragment]
```

* [repo]: Is the URL or alias of a Git repo.
* [path]: Is the path to a file within the repo, usually a *.ocu.star file.
* [release]: Is a release identifier. If blank, the most recent release is implied.
* [subpath]: A path to a document within the release, such as a deployment to a specific environment.
* [fragment]: An optional path to a field within the document.

For example, `github.com/ocuroot/example/-/frontend/release.ocu.star/@1.0.0/call/build#output/image` would
refer to the container image for the 1.0.0 release of the frontend in an example repo.

Intent References are denoted by the use of `+` instead of `@` for the release. So
`github.com/ocuroot/example/-/frontend/release.ocu.star/+/deploy/production` would
refer to the desired state for deploying the frontend to the production environment.