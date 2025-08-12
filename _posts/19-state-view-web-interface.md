---
title: "Ocuroot's new state view - a work in progress"
slug: state-view-web-interface
excerpt: "Augmenting command-line state manipulation with a web UI for more intuitive navigation."
coverImage:
  src: "/assets/blog/state-view-web-interface/cover.jpg"
  alt: "A modern web interface showing organized state data"
  credit: Photo by Andrea Piacquadio
  creditURL: https://www.pexels.com/photo/portrait-photo-of-woman-holding-up-a-magnifying-glass-over-her-eye-3771107/
date: "2025-08-06T11:00:00-04:00"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/state-view-web-interface/cover.jpg"
---

The [new state model](https://www.ocuroot.com/blog/ocuroot-sdk-v0-3-refs) for Ocuroot gives a lot of power
to developers, providing flexible access to all the details of each release and environment. But this has the
potential to be overwhelming, requiring developers to learn the structure ahead of time.

To help ease users into the world of Ocuroot state, I've added a command to render it all in a web interface:

```bash
ocuroot state view
```

### The terminal state experience

At present, navigation of state in the terminal is enabled through the `state match` and `state get` commands.

To find all current deployments in the production environment, you could run:

```bash
$ ocuroot state match **/@*/deploy/production
github.com/example/repo/-/path/to/my-package/@/deploy/production
github.com/example/repo/-/another/package/@/deploy/production
github.com/example/repo/-/a/third/package/@/deploy/production
```

Then to inspect state of one of those refs, you could run:

```bash
$ ocuroot state get github.com/example/repo/-/path/to/my-package/@/deploy/production
```

You could even stream all the deployment refs and extract the URL outputs for each service:

```bash
$ ocuroot state match **/@/deploy/production | \
xargs -n 1 ocuroot state get | \
jq -r '.output[] | .url'
```

This is great if you're familiar with the structure and have the command-line chops to build these pipelines,
but that's not everyone. Even if you are a command-line wizard, sometimes it can help to get the lay of the
land visually first.

The existing state commands aren't going anywhere, but are now augmented with the `state view` command.

### The new state view

The new state view is started by running `ocuroot state view` in any Ocuroot-configured repo. This will
start the state server on a local port.

From this server, you can view the progress of a specific release:

![release view](/assets/blog/state-view-web-interface/release.png)

And then click through to see the details of each function call or deployment:

![work view](/assets/blog/state-view-web-interface/work.png)

Or start from environments and see their configuration and current deployments:

![environment view](/assets/blog/state-view-web-interface/environment.png)

### Roadmap and learnings

As you can probably tell, it still needs some polish. Plus there's still a ways to go on supporting 
all the state types and flows to make this truly useful, including providing ways to edit intent and
monitor outstanding work.

But even at this stage it's helped me discover some things. With a heightened ability to explore state, 
I've noticed a few rough edges that I'm planning to address ASAP.

As an example, release refs including commit hashes can get very long, and can't be ordered. 
Short git hashes won't be much help either. To combat this, I'm planning to number releases
monotonically, and have a separate ref to indicate their commit.

So refs like:

```
github.com/example/repo/-/path/to/my-package/@fa56a23554a75a7ab334f841c5f61f952e52930c.2/deploy/production
```

Would be split into:

```
github.com/example/repo/-/path/to/my-package/@11/deploy/production
github.com/example/repo/-/path/to/my-package/@11/commit/fa56a23554a75a7ab334f841c5f61f952e52930c
```

Note the larger number, since this would encompass all releases of the package, not just those for
this commit.

### Added bonus

While working through the UI updates needed to support the state view, I also took the
time to add a *real dark mode*. Because why not?

![dark mode](/assets/blog/state-view-web-interface/dark-mode.png)

## What's next?

I'm in the final stages of preparing the SDK v0.3 release. My aim is to share a version that you can try out very soon. In the meantime, you can follow
Ocuroot on [LinkedIn](https://www.linkedin.com/company/ocuroot), [BlueSky](https://bsky.app/profile/ocuroot.com) or get in touch directly by booking a [demo](/demo).