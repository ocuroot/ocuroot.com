---
title: "Scheduling Work"
path: "usage/scheduling-work"
---

Through the course of changes to your source repos and updates to your intent store, Ocuroot will
need to execute various kinds of work to bring the current state into sync.

There are three main ways that work can happen in Ocuroot. Through creating releases, requesting
work, or triggering work.

This page will detail these mechanisms and the commands that can be used to configure them in CI.

## New Releases

New releases will typically be created when code is merged into the main branch of a source repo.

A new release is triggered via the command:

```bash
ocuroot release new path/to/release.ocu.star
```

This process is outlined in more detail on the [Releases](/docs/usage/releases) page.

## Requesting Work

When a release is created, it won't always be able to run to completion. Intent changes must
also be handled outside of a typical release. This is where we need a mechanism to request
any outstanding work.

```bash
ocuroot work any
```

This command will inspect the state and intent stores, and identify any outstanding work that
can be executed on the current commit. If any work is found, it will be executed.

Any outstanding work for other commits may then be "triggered" as detailed in the next section.

## Triggers

A "trigger" is an action taken by Ocuroot to schedule work against a particular commit in a repo.

The command:

`ocuroot work trigger`

Will inspect the state and intent stores, and identify any outstanding work, determining which commits
have work to be done. It will then use a "trigger function" to schedule that work.

A trigger function is defined in the `repo.ocu.star` file, and is called for each commit that has outstanding work.

```python
def do_trigger(commit):
    print("triggering for commit:" + commit)
    # Schedule a job to run the `ocuroot work any` command for this commit.

trigger(do_trigger)
```

In this case, the `do_trigger` function will be called once for each commit with outstanding work.

This would typically then schedule a job to run on this commit and call the `ocuroot work any` command.

The `ocuroot work any` command itself will also issue trigger requests as needed after it runs, to ensure
that all outstanding work can eventually be completed.

If you don't necessarily have access to both the intent and state stores, you can also call:

```bash
ocuroot work trigger --intent
```

Which will trigger work against all commits with current deployments.
