---
title: "How do I test my CI code?"
slug: how-do-i-test-my-ci-code
excerpt: "CI/CD pipelines are code too. In this post we'll explore some testing strategies so you're not always left
waiting for feedback. We'll also look at what Ocuroot is doing to make it even easier!"
coverImage:
  src: "/assets/blog/how-do-i-test-my-ci-code/cover.jpg"
  alt: "A series of test tubes filled with a blue liquid"
  credit: "Photo by Chokniti Khongchum"
  creditURL: "https://www.pexels.com/photo/laboratory-test-tubes-2280549/"
date: "2025-10-16T14:00:00-04:00"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/how-do-i-test-my-ci-code/cover.jpg"
---

A common complaint I've heard about CI systems is around slow feedback loops. Every configuration change needs to
be pushed in Git and run on a server before you know it's correct. For solo projects, this puts your feedback loop
on the order of minutes. For big teams with mandatory code review, you either need someone on-hand for every tweak,
or you're waiting days.

Not only is this slow, but it's also hugely frustrating. I'm sure we've all had to make the follow-up commit: "fix typo". Add on top of that the potential blast radius of making a mistake. Your CI system can modify basically
everything, so how many ways could you bring down production?

These slow feedback loops often come about through necessity. We're forced to run CI jobs only on our CI platforms
for (very valid) security or compliance reasons. This isn't too dissimilar from the other software we build, though.
It's unlikely that your laptop is 100% representative of your production environment, so we break down the problem
and build tests to rapidly verify our code and build confidence before shipping.

What we need is a testing framework for CI configurations.

## What can we test?

When testing a traditional application, we typically want to answer two questions.

1. Does our code build?
1. Does our code do what we expect when we run it?

For the first question, there's a direct parallel for CI configurations. Even in declarative, yaml-based CI
systems, our config has to be correct YAML, and it has to obey the schema for the platform we're using.

There are two parts to the second question. There's the overall behavior of our CI pipeline (does the deploy
step actually deploy?), and the structure of the pipeline (do we test in staging before deploying to production?).

This gives us three aspects of our config to test:

* **Syntax** - Is the config valid?
* **Steps** - Does each step in the pipeline behave correctly?
* **Structure** - Does the pipeline have the right shape?

## How you might do this today

### Syntax

Most CI systems provide some kind of linter for their config. This usually comes in the form of an IDE plugin to
give you inline errors and warnings.

This can be helpful at providing quick feedback while editing and can catch a lot of problems early. However, being
totally separate from your CI platform can result in a surprising disconnect between what the IDE accepts vs the CI platform.

### Steps

A popular option for testing the steps of a CI pipeline is to break each
step out into an isolated script. These scripts can be run separately and depending on language, even have their
own unit tests. This is a common best practice, but requires a lot of context switching and discipline to maintain.
It's all too easy to say "I only need to run a couple of commands, why do I need to have a separate file?", and 
before you know it, you have a patchwork of scripts and config to unpick.

If you're using GitHub Actions, you may have also used [act](https://github.com/nektos/act). This is a third-party
open source tool that emulates the Actions platform and allows you to run your jobs locally. It's a great effort
to improve the Actions experience, and the number of stars reflect this! That said, it's not 100% compatible with
Actions.

For self-deployed CI platforms, there is always the option of running everything locally. But this can be pretty
heavyweight, and if you have a long pipeline, you may not have the ability to skip a long test phase to get to the
step you actually want to run.

### Structure

Visualizing the structure of your pipelines before pushing changes can be a bit of a challenge with many CI
platforms. In the case of GitHub Actions, act can come to the rescue again with it's "graph" mode. This capability
is by no means universal, though.

## The Ocuroot approach

With Ocuroot, we're focusing on testing as a first-class feature. It's baked into the SDK, commands and even
the architecture of the tool.

### Syntax

Because Ocuroot is distributed, the same binary runs locally, on-prem or in the cloud. This way any commands
you run will be evaluating your configuration code in exactly the same way. So any time you interact with
the `ocuroot` tool, it will be providing immediate feedback on the correctness of your syntax and use of the
SDK.

### Steps

Ocuroot comes with a built-in REPL so you can not just run any step in your CI pipeline on demand, but you can
run any helper function as well!

![The Ocuroot REPL](/assets/blog/how-do-i-test-my-ci-code/repl.png)

This allows you to iterate on any steps in your pipeline direct from your terminal with `ocuroot repl <config.ocu.star>`, but also execute any individual
function with a single command, like `ocuroot repl <config.ocu.star> -c "build_step('arg1','arg2')"`.

So you could even create a set of unit tests for your CI config and run them as part of your PR checks!

### Structure

Finally, Ocuroot lets you see your whole release pipeline locally with the `ocuroot preview` command. This starts
a web server that renders the full pipeline in an interactive view, giving you full insight into the structure. It
even takes into account your current set of environments.

![The Ocuroot Preview](/assets/blog/how-do-i-test-my-ci-code/preview.png)

## Give it a try!

Want to try it for yourself? The Ocuroot client is [open source](https://github.com/ocuroot/ocuroot), so there's no cost to experiment! Take a look at
our [quickstart](https://ocuroot.com/ocuroot/quickstart) for a self-guided demo.

If you'd like to keep up to date on what we're building and how you can take advantage of Ocuroot to help your team,
be sure to subscribe to our newsletter!

