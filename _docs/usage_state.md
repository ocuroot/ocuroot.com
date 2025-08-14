---
title: "State"
path: "usage/state"
---

All state managed by Ocuroot is stored in JSON documents that are organized by Reference.

## Reference format

References are a URI-compatible string that serve as a path to a document within Ocuroot state.
They are of the form:

```
[repo]/-/[path]/@[release]/[subpath]#[fragment]
```

* `[repo]`: The URL or alias of a Git repo.
* `[path]`: The path to a file within the repo, usually a *.ocu.star file.
* `[release]`: A release identifier. If blank, the most recent release is implied.
* `[subpath]`: A path to a document within the release, such as a deployment to a specific environment.
* `[fragment]`: An optional path to a field within the document.

For example, the following ref would refer to the container image for the first release of the frontend 
in an example repo:

```
github.com/ocuroot/example/-/frontend/release.ocu.star/@1/call/build#output/image
```

## Working with state

State can be queried and manipulated using the `ocuroot state` commands. They are summarized here, but you
can also run `ocuroot state --help` for more information.

### ocuroot state get

The command `ocuroot state get` does what it says on the tin, it retrives the document at a specific ref.

Example:

```bash
$ ocuroot state get github.com/ocuroot/ocuroot/-/release.ocu.star/@/call/build_darwin_amd64
{
  "entrypoint": "github.com/ocuroot/ocuroot/-/release.ocu.star/@8/call/build_darwin_amd64/1/functions/1",
  "output": {
    "bucket_path": "ocuroot_binaries:client-binaries/ocuroot/0.3.9-3/darwin-amd64",
    "download_url": "https://downloads.ocuroot.com/ocuroot/0.3.9-3/darwin-amd64/ocuroot"
  },
  "release": "github.com/ocuroot/ocuroot/-/release.ocu.star/@8",
}

$ ocuroot state get github.com/ocuroot/ocuroot/-/release.ocu.star/@/call/build_darwin_amd64#output/download_url
"https://downloads.ocuroot.com/ocuroot/0.3.9-3/darwin-amd64/ocuroot"
```

### ocuroot state match

### ocuroot state view

### ocuroot state set

### ocuroot state delete

## Types of state

### Releases

Releases represent a point-in-time snapshot of a particular package and how it should be built and deployed.

A ref for a Release specifies its repo, path and release. For example:

```
github.com/ocuroot/example/-/path/to/package/release.ocu.star/@1
```    

The document at this ref defines the release process defined in the file `path/to/package/release.ocu.star` in the
repo `github.com/ocuroot/example` at release `1`. Release identifiers are monotonically increasing integers.

The commit hash for a specific release can be found in the commit subpath for the release:

```
github.com/ocuroot/example/-/path/to/package/release.ocu.star/@1/commit/[hash]
```

The hash is stored in the ref itself for quick lookup.
This allows you to look up the commit for a specific release:

```bash
$ ocuroot state match github.com/ocuroot/example/-/release.ocu.star/@1/commit/*
```

Or even find all releases at a specific commit:

```bash
$ ocuroot state match github.com/ocuroot/example/-/**/@*/commit/fa56a23554a75a7ab334f841c5f61f952e52930c
github.com/ocuroot/example/-/frontend/release.ocu.star/@1/commit/fa56a23554a75a7ab334f841c5f61f952e52930c
github.com/ocuroot/example/-/backend/release.ocu.star/@1/commit/fa56a23554a75a7ab334f841c5f61f952e52930c
github.com/ocuroot/example/-/frontend/release.ocu.star/@2/commit/fa56a23554a75a7ab334f841c5f61f952e52930c
```

Note that there can be multiple releases at a single commit.

### Calls



### Deployments

### Environments

### Custom State

## State vs. Intent

Intent References are denoted by the use of `+` instead of `@` for the release. So the below ref would refer
to the desired state for deploying the frontend to the production environment.

```
github.com/ocuroot/example/-/frontend/release.ocu.star/+/deploy/production
```

All manual changes to state must be made by first modifying intent. Ocuroot will then identify the appropriate
actions to take to apply this intent to current state.