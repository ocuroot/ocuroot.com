---
title: "Breaking out of YAML for CI/CD"
slug: breaking-out-of-yaml-for-cicd
excerpt: "How imperative code in CI/CD configurations can unlock powerful workflows that declarative YAML simply can't handle."
coverImage:
  src: "/assets/blog/breaking-out-of-yaml/cover.jpg"
  alt: "A pair of hands holding bars at the window of a stone jail"
date: "2025-04-04T12:00:00.000Z"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/breaking-out-of-yaml/cover.jpg"
---

YAML has become something of a de-facto language for CI/CD tooling, with tools like GitHub Actions, CircleCI and Travis all primarily configured using YAML. Beyond CI/CD YAML has significant penetration in the DevOps space as a whole, with YAML being
the primary interface for configuring Kubernetes.

Despite its popularity, YAML is not a perfect solution. At heart, YAML is a data serialization language and was never intended to represent complex logic. This has resulted in different providers effectively building their own DSL on top of YAML to allow for branching logic and sharing variables. The end result being a sometimes steep learning curve for each provider, and occasionally quirky pipeline structures as engineers try to work around limitations. At the very least, you're going to be reaching for bash or Python scripts to handle anything but the simplest of logic.

Born of something of a personal frustration with this situation, Ocuroot is breaking out of the YAML jail, and even other declarative languages like HCL, with configuration that you can write entirely in imperative code - specifically [Starlark](https://github.com/bazelbuild/starlark).

## Why Starlark?

Starlark is a dialect of Python originally created for the [Bazel](https://bazel.build) build system. It is a pure subset of Python designed to allow for simplicity and repeatability, with a high degree of control over the capabilities of third-party code written in the language.

I've seen Starlark used in a number of interesting projects before, not least Bazel itself, but also [Tilt](https://tilt.dev).
Tilt is a great example of how you can use Starlark to expand the capabilities of a tool. In a previous role, my team migrated
to Tilt from an in-house solution, and we were able to quickly script up support for our old configuration file format directly in our Tiltfiles.

What sealed the deal was the fined grained control over capabilities of the language. Out of the box, you can't access the filesystem or even get the current time. This allows Ocuroot to provide extremely limited capabilities when calculating pipeline structure, making the results predictable and safe to run on a shared host. At the same time, a worker executing builds and deployments can be provided full access to the host if desired.

There's also an excellent [Go implementation of Starlark](https://github.com/google/starlark-go), that was surprisingly easy to
get to grips with!

## Pros and cons of imperative configuration

Like any solution, writing your configuration in imperative code has its own advantages and disadvantages.

Among its advantages are a high degree of flexibility, the ability to break down complex logic into small functions (even spread over multiple files) and a level of familiarity that allows engineers to develop their pipelines in a form that feels more "natural" to them.

On the other hand, imperative code is less predictable by nature (subject to the dreaded halting problem). This can require robust testing and verification of shared logic. The price of the additional freedom with imperative code is an associated risk of writing
spaghetti!

To blunt the edges of some of these disadvantages and embrace the nature of living in an imperative world, Ocuroot includes
features like [pipeline rendering](/blog/07-enabling-pipeline-visualization) and a [REPL](https://docs.ocuroot.com/cli/commands/repl). Future items on the roadmap include a more complete testing suite, so you can verify your shared logic.

