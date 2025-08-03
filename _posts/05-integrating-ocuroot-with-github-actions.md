---
title: "Integrating Ocuroot with GitHub Actions"
slug: 05-integrating-ocuroot-with-github-actions
excerpt: "How you can use Ocuroot to enhance your existing CI platform"
coverImage:
  src: "/assets/blog/integrating-ocuroot-with-github-actions/cover.jpg"
  alt: A construction worker laying bricks
  credit: "Photo by Yura Forrat"
  creditURL: https://www.pexels.com/photo/brick-workers-at-construction-site-11429199/
date: "2024-12-09T17:00:00.000Z"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/integrating-ocuroot-with-github-actions/cover.jpg"
---

When setting up automated builds, tests, and deployments, you need somewhere to run them unattended. This is where Continuous Integration (CI) platforms come in. 
They provide a way to run scripts and commands on a schedule, or in response to events like source code changes.

Ocuroot can be run on top of your existing CI platform, so you can take advantage of existing configuration and integration with your source repositories. This allows you to focus on the specifics of your build and deployment process, without having to worry about the underlying infrastructure.

## Key Workflows

Integrating with a CI system requires implementation of three key workflows:

1. **Review**: Automated checks for pull requests
2. **Deliver**: Builds and registration of merged changes into Ocuroot state
3. **Sync**: Deploying builds to environments when ready

## Documentation

Check out our [GitHub Actions integration guide](https://docs.ocuroot.com/ci/github-actions) for details of how it all works. The guide includes example workflow files and detailed configuration instructions that will help you get up and running quickly.

For a broader overview of CI integration options, visit our [CI documentation](https://docs.ocuroot.com/ci/introduction).

## Want to try Ocuroot?

We're currently in a closed alpha. If you're interested in trying Ocuroot for yourself, [book a demo](https://ro.am/ocuroot/) so we can get you started.
