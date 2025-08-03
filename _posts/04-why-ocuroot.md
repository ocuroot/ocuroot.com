---
title: "Why Ocuroot?"
slug: 04-why-ocuroot
excerpt: "In June of this year, I quit my job to start Ocuroot. It's a huge change, and people have asked my why I decided to take
this leap. I've asked myself that far more often. But there's a simple answer."
coverImage:
  src: "/assets/blog/why-ocuroot/cover.jpg"
  alt: Person in yellow gloves checking a green plant.
  credit: "Photo by Tima Miroshnichenko"
  creditURL: https://www.pexels.com/photo/person-in-yellow-gloves-checking-a-green-plant-6511168/
date: "2024-12-02T10:00:00.000Z"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/why-ocuroot/cover.jpg"
---

In June of this year, I quit my job to start Ocuroot. It's a huge change, and more than a few people have asked me why I decided to take
this leap. I've asked myself that far more often. But there's a simple answer. 

CI/CD is hard at scale.

From 2019-2024 I was running the Engineering Productivity group at Yext, a mid-sized B2B technology company, with the developer
workflow being one of our key responsibilities.

As the company grew, I noticed that we were adding more and more production environments. Points-of-presence for serving content to consumers.
Backup sites for disaster recovery. And most notably, a full replica of our stack in the EU for GDPR and data sovereignty requirements.

With the extra complexity of all these environments multiplied by thousands of microservices, it became harder and harder to keep track of what versions of our code were running where. This slowed down not only software delivery, but also our ability to respond to audits and handle
critical incidents.

I want to build the CI/CD tool that I've always wished I had. With Ocuroot, I'm building an opinionated tool based on managing state rather than
pipelines-as-DAGs, following the common relationships between commits, builds and deployments that many of us take for granted.

If you've experienced the pain of chasing down builds and deployments across many production environments, or are daunted an upcoming project
to scale into a new region, I'd love you to join me on this journey!

I'm building in the open and putting out regular updates via [LinkedIn](https://www.linkedin.com/company/ocuroot) and [Bluesky](https://bsky.app/profile/ocuroot.com). Follow along on your platform of choice!