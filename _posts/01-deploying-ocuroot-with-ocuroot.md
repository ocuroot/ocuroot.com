---
title: "Deploying Ocuroot with Ocuroot"
slug: 01-deploying-ocuroot-with-ocuroot
excerpt: "One of the best things about building a developer tool is that I can use it while building it! "
coverImage:
  src: "/assets/blog/deploying-ocuroot-with-ocuroot/cover.png"
  alt: Screenshot of Ocuroot documentation site, showing the introduction page.
date: "2024-10-28T10:00:00.000Z"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/deploying-ocuroot-with-ocuroot/cover.png"
---

One of the best things about building a developer tool is that I can use it while building it! Over the last few weeks, I’ve been setting up a cloud instance of Ocuroot as a preview for everyone, and I’ve been using Ocuroot to deploy itself.

This creates a bit of a challenge in getting things up and running. How do I deploy the initial build of the Ocuroot server without having the server available to store the state. This is where bootstrap mode saves the day!

The Ocuroot client can use a sqlite file as a virtual server to manage the initial set of builds and deployments. Then, when everything’s up and running, the export command uploads the bootstrap database to the remote server and it takes over.

This is making my life easier right now, but if you wanted to deploy a self-hosted instance of Ocuroot in your cloud of choice, this could get you up and running in a few minutes by bootstrapping a template repository!