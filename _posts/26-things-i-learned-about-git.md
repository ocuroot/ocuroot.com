---
title: "Three things I've learned about Git while building a CI/CD tool"
slug: things-i-learned-about-git
excerpt: "Over the years, I got into a narrow routine when using Git for day-to-day development. Now I'm building tools around Git, I'm discovering all kinds of features I'd never noticed before. Here are my favorite three."
coverImage:
  src: "/assets/blog/things-i-learned-about-git/cover.svg" 
  alt: "The Git Logo"
date: "2025-11-13T11:00:00-04:00"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/things-i-learned-about-git/cover.svg"
---

I've been using Git pretty much daily for well over a decade, but over time, I've settled into rhythms
and routines that only covered a small subset of what it could do. Clone a repo here and there, add some files, commit and push. A little branch manipulation when working in a team or contributing to open source. Why would I ever need more? And why should I worry about the internals?

Of course, now I'm building tools that integrate with Git more closely, I've had to dive a little deeper. Digging into the documentation and the internals, I've discovered a few features and fun little
quirks of Git that I really wish I'd known sooner. Here's a quick overview of my top three.

## Worktrees

[Git worktrees](https://git-scm.com/docs/git-worktree) allow you to work on multiple branches at once.
Imagine you're working on a new feature on a dedicated branch (short-lived, of course!). You get paged
with an urgent issue and rolling back won't cut it, you need to quickly ship a patch! At minimum, this context switch would involve checking out a different branch. Of course, you probably have some uncommitted changes so might need to run `git stash`. Worse, if your change involves messing with the `.gitignore` you could end up with a bunch of dangling changes when you switch branches.

I've known some engineers who had multiple local clones of a repo for parallel work, but if you have
a particularly large monorepo, that can get pretty unwieldy.

With worktrees, you can quickly create a new working directory on a specific branch with `git worktree add <path> <branch>`. So for our example above, you might run:

```bash
git worktree add ../hotfix main
```

You can then open the `hotfix` directory and get to work, without touching the directory for your feature work.

When you're done, you can remove the worktree with `git worktree remove`. Or just delete
the directory and Git will eventually clean up for you (or be forced to with `git worktree prune`).

## The porcelain flag

Git commands can be divided into two categories:

* Porcelain: designed for humans, with easy(ish) to read output and friendly(ish) error messages.
* Plumbing: lower level commands that do low level work and provide output in a more machine-friendly format.

When you're writing scripts and building tooling around git, the plumbing commands are probably what you should be using. A full list of commands can be found in the [Git docs](https://git-scm.com/docs) with the plumbing commands in their own category at the end.

The plumbing commands can be a little cumbersome at times, though, and you probably don't always want
to learn a whole new set of commands to script things you do with the porcelain commands every day.
This is where the `--porcelain` flag comes in.

Some of the porcelain commands have a `--porcelain` flag that converts their output into a machine-readable format. So you tell porcelain commands to use the porcelain setting so they behave
more like plumbing. Isn't that just a delightful kind of confusing?

The subset of commands with this flag is pretty limited, and includes `git status`, `git blame`, `git push` and `git worktree list`. Looking at `git status`, you might see the following as typical output:

```bash
$ git status
On branch main
Your branch is up to date with 'origin/main'.

Changes to be committed:
  (use "git restore --staged <file>..." to unstage)
        new file:   _posts/26-things-i-learned-about-git.md
        new file:   static/assets/blog/things-i-learned-about-git/cover.svg
```

Not super user friendly, but if we add the porcelain flag it becomes much easier to parse:

```bash
$ git status --porcelain
A  _posts/26-things-i-learned-about-git.md
A  static/assets/blog/things-i-learned-about-git/cover.svg
```

As a bonus, there's even a *v2* mode with more info:

```bash
$ git status --porcelain=v2
1 A. N... 000000 100644 100644 0000000000000000000000000000000000000000 5b3c9adde4dccc7053cc778abe42c1290595faa0 _posts/26-things-i-learned-about-git.md
1 A. N... 000000 100644 100644 0000000000000000000000000000000000000000 a6b38eee2e639377e44819ed5084f52bf736cdf5 static/assets/blog/things-i-learned-about-git/cover.svg
```

## File protocol

When you're writing integration tests that involve Git, you may want an upstream repo you can
clone. You could do this with a GitHub repo, but then you need to muck around with credentials, repo
creation and cleanup (especially given recent [repo count limitations](https://github.blog/changelog/2025-03-27-repository-ownership-limits/)).

Thankfully, there's a way to clone a repo from your local disk! Along with the ssh and http protocols,
Git also supports a *file* protocol.

We can create a "remote" on the local disk in a temp directory:

```bash
export GIT_REMOTE_DIR="$(mktemp -d)"
cd "$GIT_REMOTE_DIR"
git init --bare
```

The `--bare` flag tells Git that we don't want a working tree. Our directory will now contain only the administrative files for the repo, everything that would ordinarily go in the `.git` directory.

We can now clone this bare repo into another working directory 

```bash
git clone "file://$GIT_REMOTE_DIR" git-remote-clone
cd git-remote-clone
git checkout -b my-branch
echo "Generic greeting here" > README.md
git add README.md
git commit -m "Add README"
git push origin my-branch
```

This supports pretty much everything you can do with a "real" remote repo and is plenty for most basic
Git interactions. However, I have seen some contention issues if multiple Git commands are run
concurrently or in rapid succession, so there are still times where a more complete Git server can be
useful for testing.