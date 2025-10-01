---
title: "Complex rollouts with one command - Ocuroot's new Cascade mode"
slug: cascade-mode
excerpt: "With Ocuroot's new cascade mode, you can roll out a complex environment with a single command!"
coverImage:
  src: "/assets/blog/cascade-mode/cover.jpg"
  alt: "A small waterfall in a rocky part of a woodland. The waterfall has multiple levels."
  credit: "Photo by Kennst du schon die Umkreisel App?"
  creditURL: "https://www.pexels.com/photo/waterfall-1080421/"
date: "2025-10-02T10:00:00-04:00"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/cascade-mode/cover.jpg"
---

In the [announcement for v0.3.14](/blog/v0.3.14-release-feedback-requested/), I mentioned that integration with CI
platforms is doable, but not as easy as it could be. This was largely down to having to trigger multiple CI runs
to complete a release or handle intent changes. With [v0.3.16](https://github.com/ocuroot/ocuroot/releases/tag/v0.3.16), Ocuroot has a new mode that allows more work to be done with a single command: cascade!

Cascade mode is a step up from the comprehensive mode that was previously added to the `ocuroot work any` command. While
comprehensive mode would execute all work available on the current commit in the current source repo, cascade mode will
execute any follow on work across *all* known source repos and commits in those repos. Work is batched based on commits,
and clean copies of the source are created in dedicated, temporary directories.

## Simplifying the quickstart

The first place this makes a difference is in the [quickstart](/docs/quickstart). When resolving the network dependency 
for the frontend, you previously had to run three commands:

```bash
ocuroot release new frontend/package.ocu.star
ocuroot release new network/package.ocu.star
ocuroot work any
```

The final `ocuroot work any` command executing any work needed to finish the frontend release. It gets even more complicated
if the frontend and network are released from separate commits, since you need to check out the correct commit to compelete
this work.

```bash
ocuroot release new frontend/package.ocu.star
# ... commit changes here
ocuroot release new network/package.ocu.star
ocuroot work any # Will tell you which commit to checkout
git checkout frontend-tag
ocuroot work any # Will complete the frontend release
```

Suffice to say, with many different resources to release, you could easily end up having to repeat the `work any` command
many times to get all the different commits.

With the new cascade feature, you just need to call two commands:

```bash
ocuroot release new frontend/package.ocu.star
# ...you may commit any changes here
ocuroot release new network/package.ocu.star --cascade
```

## Simplifying CI

With Ocuroot commands able to do more work in one shot, [configuring your CI platform](/docs/usage/ci-integration/) to run Ocuroot becomes much simpler.

Previously, you needed to configure three separate CI workflows:

* `ocuroot release new && ocuroot work trigger` on push to source branches
* `ocuroot work trigger --intent` on push to the intent branch
* `ocuroot work any` run on-demand by Ocuroot itself.

This also required the implementation of a "trigger" function to kick off the on-demand workflow. This required Ocuroot to
be provided with credentials to kick off CI runs automatically. Slightly annoying at the best of times, but configuring things
in just the wrong way could cause runaway execution of runs, burning compute minutes until agressively killed by hand.

With cascade mode, only two workflows are needed:

* `ocuroot release new --cascade` on push to source branches
* `ocuroot work cascade` on push to the intent branch

In both cases, the changes pushed are handled in a single command, with no additional jobs needed. Simpler to configure,
simpler behavior and simple to monitor.

As an added bonus, by consolidating all the work into a single job, changes are applied more quickly, since the CI job
only needs to be queued and set up once.

## Configuring git remotes

Of course, for Ocuroot to clone repos, it needs to know where they are. Since Git remotes can vary quite a lot even for
the same repo (https vs ssh anyone?), there is now the option to provide a set of remotes to try in `repo.ocu.star`.
Even better, because all Ocuroot config is Starlark, you can write logic to build a remote based on where Ocuroot
is running. For example, you can support local execution and GitHub Actions at the same time:

```python
def repo_url():
    env_vars = env()
    # When running outside GitHub Actions, use the origin remote as-is
    if "GH_TOKEN" not in env_vars:
        return host.shell("git remote get-url origin").stdout.strip()
    # Always use https for checkout with the appropriate token on GitHub actions
    return "https://x-access-token:{}@github.com/ocuroot/gh-actions-example.git".format(env_vars["GH_TOKEN"])

remotes([repo_url()])
```

## What's next?

Now that the story for CI integration is a little simpler, I can start looking at automating it! Look out for
initialization commands in future releases to set up your platform of choice for you.

Alongside this automation, I'll also be expanding instructions/automations to cover additional platforms.

CI work aside, I have a few improvements in mind for the general configuration model and telemetry forwarding.

Got something you'd like to see in Ocuroot? Had a chance to try the [quickstart](/docs/quickstart) and have some notes? 
Feel free to raise an [issue](https://github.com/ocuroot/ocuroot/issues) with your ideas!
