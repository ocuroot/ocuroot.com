---
title: "SDK v0.3: Managing dependencies across deployments"
slug: ocuroot-sdk-v0-3-dependencies
excerpt: "An update on how Ocuroot SDK v0.3 will handle dependencies, enabling seamless asset sharing and cross-deployment references."
coverImage:
  src: "/assets/blog/ocuroot-sdk-v0-3-dependencies/cover.jpg"
  alt: "Connected building blocks representing dependencies"
  credit: "Photo by Polina Tankilevitch"
  creditURL: "https://www.pexels.com/photo/a-person-in-black-and-gray-jacket-handing-over-boxes-to-a-person-4440774/"
date: "2025-05-29T10:00:00-04:00"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/ocuroot-sdk-v0-3-dependencies/cover.jpg"
---

Over the past few months, I've been sharing updates about v0.3.0 of the Ocuroot SDK, focusing on [function chaining](/blog/ocuroot-sdk-v0-3-preview) and [simplified phase definitions](/blog/ocuroot-sdk-v0-3-simplifying-phase-definitions). These help you create a clear
release pipeline for a single package, but what if you need to coordinate
dependencies between multiple packages?

Perhaps your application deployment needs to know which VM it should SSH into
to download and start your binary. Your application may need to be provided with
a connection string for your database, or the string is stored elsewhere and you just want to make sure the database is created before your application starts. You may even need to know some details of the currently deployed version of your 
application to handle the migration correctly.

This post will detail some of the patterns in the upcoming SDK version that will
enable these use cases, plus plenty of others!

## Sharing assets during a release

Before we discuss cross-package dependencies, let's look at how we can share
data between phases of a single release.

As mentioned in the previous post, there are now no separate "build" functions for
a release, instead we support "deploy" work for deployments to a specific environment, and generic "calls" for arbitrary work that can happen before, after or alongside deployments:

```python
# Build our binary
phase(
    name="build",
    work=[call(fn=build, name="build")],
),

# Deploy to staging
phase(
    name="staging",
    work=[
        deploy(
            up=do_up, 
            environment=staging_env,
        )
    ],
),

# Obtain user approval to deploy to production
phase(
    name="prod approval",
    work=[call(name="prod approval", gate=approval())],
),

# Deploy to production
phase(
    name="prod",
    work=[
        deploy(
            up=do_up, 
            environment=prod_env,
        )
    ],
),
```

What's missing from this example is a connection between the build call and the
deployments. At a minimum, we'd need to know where the binary or container for this
release is stored.

This is where input dependencies come in:

```python
# ...
phase(
    name="staging",
    work=[
        deploy(
            up=do_up, 
            environment=staging_env,
            inputs={
                "binary_url": input.self().call("build").output("binary_url"),
            },
        )
    ],
),
# ...
```

The `inputs` parameter for the deployment defines a set of input values to be passed
to the `up` function. In this case, we're retrieving the binary URL from the build
call. Breaking this apart:

* The `input.self()` expression refers to the current package
* The `call("build")` expression refers to the `build` call in the current package
* The `output("binary_url")` expression refers to the `binary_url` output of the `build` call

The `binary_url` output should be returned from the `build` call:

```python
# Build our binary
def build(ctx):
    binary_url = build_binary()
    return done(outputs={"binary_url": binary_url})
```

This input can then be referenced in the `up` function we specified:

```python
# Deploy the binary to the appropriate environment
def do_up(ctx):
    deploy_binary(ctx.environment, ctx.inputs.binary_url)
    return done()
```

## Cross-deployment references and dependencies

Now we have a model for sharing data within a release, we can extend it to reference other releases!

Where `self()` refers to the current package, we can use `repo()` and `package()`
to specify other packages, such as:

 * `input.repo("github.com/my_org/infra").package("db")` to refer to the db package in the infra repo.
 * `input.package("backend")` to refer to the backend package in the current repo.

We can then use the `deploy` function to refer to a deployment rather than a call,
and pull details about the database into our deployment.

```python
deploy(
    up=do_up, 
    environment=staging_env,
    inputs={
        "binary_url": input.self().call("build").output("binary_url"),
        "db_conn": input.repo("github.com/my_org/infra").package("db").deploy(staging_env).output("conn"),
    },
)
```

We can even make the connection string a secret when deploying the database to 
ensure it is encrypted in storage and not exposed in logs:

```python
def deploy_db(ctx):
    conn = create_db_here(ctx.environment)
    return done(outputs={"conn": secret(conn)})
```

This example assumes that you want to use Ocuroot for service discovery, and the
odds are you will already have some solution for this. If you're in the HashiStack
you may be using Consul or Vault. If you use Kubernetes, you can use configmaps,
secrets and internal DNS. But these resources will also need setting up. Input
dependencies in Ocuroot can bridge that gap, so you can build entire environments
from scratch!

## What's next?

I'm in the final stages of preparing the SDK v0.3 release. My aim is to share an
example client that you can try out very soon. In the meantime, you can follow
Ocuroot on [LinkedIn](https://www.linkedin.com/company/ocuroot), [BlueSky](https://bsky.app/profile/ocuroot.com) or get in touch directly by booking a [demo](/demo).

