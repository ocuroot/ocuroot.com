---
title: "Local testing: the missing piece for rapid development cycles in CI/CD"
slug: 08-local-testing-missing-piece
excerpt: "How a good local testing workflow can accelerate your development cycle, approaches you can take today and how Ocuroot can help in the future."
coverImage:
  src: "/assets/blog/local-testing-keystone/cover.jpg"
  alt: "A man in a hoodie holding a laptop smiling with streams of code in the background"
date: "2025-03-21T15:00:00.000Z"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/local-testing-keystone/cover.png"
---

Coding tasks are at their core, a feedback loop. You write some code, run it to see if it works, and fix it if it doesn't.
The second part of this loop, running the code is typically highly automated, and the faster you can make it, the better the
developer experience will be.

If you're left waiting for even a few seconds for your code to run, your mind wanders, and it's easy to get distracted by Slack or social media. Worse, if you're left waiting for a few minutes, your frustration can easily boil over when you finally get the feedback that you made a one-character typo and have go wait all over again.

I remember in 2014, seeing Apple's keynote announcing [Swift Playgrounds](https://developer.apple.com/swift-playground/) and being
wowed by the potential of seeing the results of your code appear with every character you typed. Since this, the value of a fast
feedback loop has not been lost on my, and I've sought it out in my own work wherever possible.

The CI/CD feedback loop, unfortunately, often feels like it's stuck in the past. You can frequently find yourself having to commit
and push your changes for any kind of feedback from your CI platform. For a complex pipeline, this can take a long time, and if
deployments are involved, you can find yourself in a bad state from which it's hard to recover. Add in a strict code review process,
and it can take *days* to validate a change.

## Why is the CI/CD feedback loop so slow?

At its heart, CI/CD feedback loops are often slow because you can *only* run your pipelines on your CI/CD platform. If you've worked with GitHub Actions, you'll likely be familiar with this pattern, having to modify YAML configs and
push to the origin to execute the pipeline and see results.

This is often because of the potential differences between your CI runners and your local machine. Builds, tests and
deployments often require a large number of different tools, some of which may be OS-specific or vary significantly from
version to version. This can make it hard to replicate the runner environment effectively.

So it's easier to require that your pipeline runs only on pre-configured runners. After all, you can't say "it works on
my machine" if you can't *run* it on your machine.

In larger organizations, there may also be security or compliance reasons for restricting pipeline execution to approved
machines. Deploying to production, in particular requires access to sensitive resources, which you can't share with every
developer.

## How can CI/CD fail?

Slow feedback loops are most frustrating when things go wrong, so improving the situation will likely be based on detecting
errors more quickly.

There are infinite ways a CI/CD pipeline can fail, but they can mostly be broken down into two categories: *syntax errors* and
*logic errors*. Syntax errors are mistakes in coding that prevent the code from being run at all - like bad indents in a YAML file. Logic errors are mistakes that result in incorrect behavior - like testing one version but deploying another.

CI/CD code can also be split into two categories, *configuration* and *scripts*. Configuration defines the structure of a pipeline,
as well as how and where it runs. Scripts, on the other hand, are executed when each step in the pipeline actually runs.

This gives us four categories of errors:

- Configuration syntax errors
- Configuration logic errors
- Script syntax errors
- Script logic errors

## Tightening the loop

Taking these sources of errors into account, there are a number of ways you can speed
up the feedback loop for your CI/CD pipeline - regardless of the platform you're using.

### IDE plugins and linters

Most CI/CD platforms provide linters and even IDE plugins. This will help catch syntax
errors in your configuration (and in many cases, your scripts) as the code is being written.

Even if your platform doesn't provide a dedicated plugin for your IDE of choice, there
may be a generic option. For common languages like YAML, there are many linters available.

### Separate scripts from pipeline config

If you can make your build, deployment and test scripts standalone, you will no longer be
limited to running them on your CI/CD platform. This means they can be tested in isolation
on a local machine, speeding up the feedback loop.

As a bonus, your pipeline config then becomes simpler, so there is a smaller surface area
for issues you can't catch without pushing code.

Of course, running these scripts on your local machine will potentially hit issues with
missing (or incompatible) dependencies and tools. This is where containerization and build
systems like Bazel can help.

### Optimize your builds and deploys

Of course, there will still be situations where you aren't able to catch problems before
pushing your code to main. So you will want the final pipelines to run quickly and catch problems as early as possible. For example, unit tests can be run during the code review
phase even before merging.

Designing your pipelines to front-load with your fastest tests will give you early warning
for many issues, and caching dependencies can help speed up builds.

#### Manage your queues

A common bottleneck in executing your scripts will be the queue of pending jobs. Fully
hosted problems usually manage this for you, but if you're self-hosting or at least hosting
your own runners, you'll need to make sure you have enough capacity to service the load.

CI/CD load is often cyclical, with the bulk of the work happening during the day while your
developers are actively committing changes. So you may be able to scale cloud runners down overnight or take advantage of autoscaling/serverless runners.

If you're able to prioritize some workloads over others, you may want to prioritize the earlier stages of your pipelines, to catch more of those early issues in your scripts.

#### Run your CI/CD platform locally

Of course, a brute-force approach to testing CI/CD workflows locally is to
run a local copy of your whole CI/CD platform. Self-hosted platforms can often
be run on a laptop (although some might be pretty resource intensive). For GitHub Actions, whose server is usually fully managed, there is [Act](https://github.com/nektos/act), an open source tool to emulate the GitHub environment. A valiant effort (and a great tool!), but one that often has to deal with a moving target.

## The Ocuroot approach

Ocuroot is being built from the ground up with a focus on providing rapid feedback loops, primarily through local testing.

Logic for your pipeline structure, builds, deploys and tests are written as functions in [Starlark](https://github.com/bazelbuild/starlark), 
a Python-like language that is easy to lint and validate. This means that many syntax
errors can be caught locally.

Since your scripts can be written in a common, cross-platform language, the Ocuroot CLI tool also allows you to execute them locally. So you can test your build and deployment
scripts using inputs appropriate to local execution.

If you've already separated your build, deployment and test code into separate scripts, no
problem! Ocuroot can act as a thin wrapper for whatever tools and scripts you're already
using.

When it comes to managing dependencies, Ocuroot allows you to continue to use containerization and the build systems you love, but also provides a module system for
reusing code and providing clean abstractions for your tools. For example, a Starlark module for Terraform can automatically download and cache the correct version of the Terraform tool
as needed, and provide function calls to run common commands like `plan` and `apply` using
configurations appropriate to your organization.

Finally, Ocuroot's package structure allows fine-grained control of when builds
are triggered, so you can avoid unnecessary work clogging up your build queues.

It's still early days for Ocuroot, and there are many goals on the roadmap to further enhance your feedback loops, including:

- On-demand execution of your scripts on a remote runner for testing
- Testing framework for shared modules
- Custom IDE plugins