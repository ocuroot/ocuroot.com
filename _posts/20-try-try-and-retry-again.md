---
title: "Try, try and retry again - supporting retries in Ocuroot"
slug: try-try-and-retry-again
excerpt: "Why retries are important and some first steps in building them into Ocuroot"
coverImage:
  src: "/assets/blog/try-try-and-retry-again/cover.jpg"
  alt: "Close-up photography of crumpled paper"
  credit: Photo by Steve Johnson
  creditURL: https://www.pexels.com/photo/close-up-photography-of-crumpled-paper-963048/
date: "2025-08-21T11:00:00-04:00"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/try-try-and-retry-again/cover.jpg"
---

One early feature request from a user of the new Ocuroot SDK was to be able to retry a failed release.

In their case, while testing out some pipeline logic locally, they were getting frustrated having to re-run an
entire release pipeline repeatedly to fix problems in a late phase.

```mermaid
graph LR
    A[Build<br/>30s] --> B[Test<br/>2m]
    B --> C[Staging<br/>10s]
    C --> D[Production<br/>10s]
    
    classDef success fill:#90EE90,stroke:#006400,stroke-width:2px
    classDef failure fill:#FFB6C1,stroke:#DC143C,stroke-width:2px
    classDef pending fill:#D3D3D3,stroke:#696969,stroke-width:2px
    
    class A,B success
    class C failure
    class D pending
```

When you're trying to debug the staging deployment logic, you don't want to have to wait an
extra 90 seconds to verify a one-line fix. What we really wanted was a way to re-run just the
failed step in the pipeline. 

To make this experience a bit less painful, I've introduced a new command:

```bash
ocuroot release retry path/to/release.ocu.star
```

Which will retry the failed work for the most recent release of this file. This way you can just
copy the failed release command and change `new` to `retry`.

You can also specify a release id explicitly if you need to:

```bash
ocuroot release retry path/to/release.ocu.star/@1
```

## But what's it really doing?

Under the hood, retrying a release will:

1. Identify the failed work items
1. Duplicate them, with the new copies set to the `pending` state
1. Mark the failed work items as `failed-retried`
1. Resume the pipeline, picking up the new pending work

This will not only result in the failed work being retried, but also all subsequent phases. 
I also accompanied this update with a UI change, so you can also see all the prior attempts for
a given work item.

![The retry history for a work item](/assets/blog/try-try-and-retry-again/retry-history.jpg)

Clicking into any of these attempts will show you the usual inputs, outputs and logs. Handily, this
update also pulls double-duty, allowing you to see previous *successful* deployments to each environment.

## Limitations

This approach is great for rapid iteration on build, test and deployment logic, but copying the
existing work config does have some drawbacks.

Changes to inputs or the specific functions being used won't be picked up, since these are part
of the existing work config. So as an example from Ocuroot's own release workflow, if we have a
call that failed:

```python
call(
    increment_version, 
    name="increment_version", 
    inputs={
        "prev_prerelease": input(ref="./@/call/increment_version#output/prerelease", default="0.3.4-1"),
        "prev_version": input(ref="./@/call/release#output/version", default="0.3.4"),
    },
)
```

And we change the input defaults:

```python
call(
    increment_version, 
    name="increment_version", 
    inputs={
        "prev_prerelease": input(ref="./@/call/increment_version#output/prerelease", default="0.3.9-1"),
        "prev_version": input(ref="./@/call/release#output/version", default="0.3.9"),
    },
)
```

This won't be picked up by the retry, since the original default was already fixed in the config.

There are also potentially some ergonomic challenges with the current command structure. Changing `new`
to `retry` is simple enough, but you have to mash the left cursor key a bit.

Finally, related to the above ergonmic issue is one of idempotency. The CI platforms that Ocuroot
would run on also have their own retry mechanisms, to deal with third party issues and flaky tests.
What happens if you re-run a job that calls `ocuroot release new`? Should it assume a whole new release,
or retry an existing one?

## What does the future hold?

This feature will no doubt evolve more as time goes on, but the ability to retry failed work is going to
continue to be a key part of Ocuroot's feature set. But the problem of retries goes way beyond fast iteration
during local development. There are questions of how to distinguish between transient issues and consistent failures,
how to distinguish a flaky test from a bad build, and the portability of CI logic with side effects.

Watch this space for updates on these challenges and other CI/CD concerns!