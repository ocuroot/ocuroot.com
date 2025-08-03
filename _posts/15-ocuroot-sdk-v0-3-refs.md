---
title: "SDK v0.3: Supercharge state management with the new refs model"
slug: ocuroot-sdk-v0-3-refs
excerpt: "Introducing the all-new refs model in Ocuroot SDK v0.3, enabling flexible pipeline connections based on any attribute."
coverImage:
  src: "/assets/blog/ocuroot-sdk-v0-3-refs/cover.jpg"
  alt: "A library filing cabinet with one open drawer containing index cards"
  credit: "Photo by Tima Miroshnichenko"
  creditURL: "https://www.pexels.com/photo/brown-wooden-wall-mounted-rack-6550462/"
date: "2025-06-12T12:00:00-04:00"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/ocuroot-sdk-v0-3-refs/cover.jpg"
---

As we saw in the [previous blog post](/blog/ocuroot-sdk-v0-3-dependencies), Ocuroot really comes
into its own when you start connecting your release pipelines together. This allows you to create a full environment from the ground up without manual intervention.

The initial implementation relied on being able to construct a path to an output value using
the builder pattern:

```python
inputs={
    "binary_url": input.self().call("build").output("binary_url"),
    "db_conn": input.repo("github.com/my_org/infra").package("db").deploy(staging_env).output("conn"),
},
```

This had two major issues. Firstly, it created some pretty verbose code even in simple cases. And
more importantly, it limited the data you could use as inputs to just the outputs of other packages.
Anything more complex would require me to write and document a whole set of functions to support them - plus more to learn for the end user.

I wanted something more flexible, so I've spent the last couple of weeks working on an alternative
model, which I'm calling "refs".

A ref in Ocuroot is a URI-style string that points to a specific value in state. It encodes the repository,
package, specific version (or intent) and a subpath to a document. 

Refs are structured as:

```
<repository url>/<package path>/<@version/+intent>/<subpath>#<json ref>
```

We'll get into the specifics of what each of these mean shortly, but for now, let's convert the earlier example inputs to refs:

```python
inputs={
    "binary_url": ref("call/build#output/binary_url"),
    "db_conn": ref("github.com/my_org/infra.git/db/@/deploy/staging#output/conn"),
},
```

## Exploring these examples

The first example is a relative ref. It points to the `binary_url` output of the `build` call in the current release of the current package. In this case, it is used to retrieve the URL of the binary build for this
release, ready to deploy it.

The second example is an absolute ref. It points to a database connection string in the `staging` environment.
Both the repository (`github.com/my_org/infra.git`) and package (`db`) are explicitly specified.

The `@` symbol indicates the "latest" release. Specifically, `@/deploy/staging` refers to the latest release
deployed to the staging environment.

## Intent

As well as being able to refer to the current state of deployments, you can also declare and refer to "intent"
for state. This allows GitOps-like interactions with state, allowing you to request changes to inputs, the
release you want deployed in an environment or even to destroy specific resources. These interactions can
even be scripted, so you can automate your operations in ways I might never have considered!

* **Intent** is denoted with a `+` symbol and represents your desired state.
* **State** is denoted with an `@` symbol (optionally followed by a specific version) and represents live state for your releases and environments.

For example:
- `app-backend/+/deploy/prod` indicates the intent to deploy to production
- `app-backend/@/deploy/prod` references the actual state of the production deployment

Users can modify intent as they see fit, whereas state should only be altered via Ocuroot tooling.

## Storage and Accessibility

Refs are designed for maximum accessibility and compatibility:

* They're valid URIs that can be written to a filesystem or included in URLs
* State storage is essentially a JSON document store, making data easy to work with
* The fragment part of the ref allows you to reach into specific JSON fields for precise data access

For example, if a deployment produces contains this state:
```json
{
  ...
  "status": "success",
  "output" {
    "resources": {
      "connectionString": "postgresql://user:pass@db:5432/mydb"
    }
  }
  ...
}
```

You could reference the connection string with:
```
.../@/deploy/prod#output/resources/connectionString
```

## Tags and Links

At present, a version number is a ULID, which aren't the easiest things to type by hand. To make it easier
to interact with state, there is the concept of a "link" from one ref to another. This acts similarly to a
symlink. So you could create a link to the release `github.com/example/my_repo.git/my_package/@01JXJF6QMXF3J4M66WDJTVFX72` from the link `github.com/example/my_repo.git/my_package/@v1`.

This link can be created by "tagging" a release in your output. For example:

```python
ctx.call(
  fn=build, 
  name="build", 
  inputs={
    "previous_count": ref("@/call/build#output/count", default=0),
  },
)

# ...

def build(ctx):
    count = ctx.inputs.previous_count + 1
    return done(
      outputs={"count": count}, 
      tags=["v"+str(count)],
    )
```

In this example, we've declared our "build" call to take in the previous build count from the most recent
run of the build function (most recent being indicated by `@`). This is defaulteed to `0` to handle the first
release for this package.

In the `build` function itself, we increment the count and then return both the count and output a tag for
the release indicating a version number. This will result in a link to our first release of the form: `github.com/example/my_repo.git/my_package/@v1`.

## Using refs in the CLI

Having a consistent URI structure for all state can do wonders for composability and exploration.

Let's say you want to find out which release of a package was currently deployed to staging. Using the
fragment we can go straight to the release ref.

```bash
$ ocuroot state get github.com/example/repo.git/path/to/my-package/@/deploy/staging#release

"github.com/example/repo.git/path/to/my-package/@01JXJF6QMXF3J4M66WDJTVFX72"
```

Or list all the outputs:

```bash
$ ocuroot state get github.com/example/repo.git/path/to/my-package/@/deploy/staging#release

{
  "count": {
    "number": 4
  },
  "env_name": {
    "string": "staging"
  }
}
```

Pipe this to [jq](https://jqlang.org/) and you can do all kinds of manipulation!

You can also match against glob patterns, so you can find which packages in a repo were currently deployed to
staging:

```bash
$ocuroot state match "github.com/example/repo.git/**/@/deploy/staging"

github.com/example/repo.git/path/to/my-package/@/deploy/staging
github.com/example/repo.git/another/package/@/deploy/staging
github.com/example/repo.git/a/third/package/@/deploy/staging
```

## What's next?

I'm in the final stages of preparing the SDK v0.3 release. My aim is to share a version that you can try out very soon. In the meantime, you can follow
Ocuroot on [LinkedIn](https://www.linkedin.com/company/ocuroot), [BlueSky](https://bsky.app/profile/ocuroot.com) or get in touch directly by booking a [demo](/demo).
