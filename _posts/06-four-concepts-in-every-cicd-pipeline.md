---
title: "Four concepts you see in every CI/CD pipeline"
slug: 06-four-concepts-in-every-cicd-pipeline
excerpt: "CI/CD pipelines come in all shapes and sizes, but you can describe almost every step using four key concepts."
coverImage:
  src: "/assets/blog/four-concepts-in-every-cicd-pipeline/cover.jpg"
  alt: A set of metal pipelines in a forest
  credit: "Photo by Wolfgang Weiser"
  creditURL: https://www.pexels.com/photo/view-of-pipelines-in-a-forest-18784617/
date: "2024-12-18T10:00:00.000Z"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/four-concepts-in-every-cicd-pipeline/cover.jpg"
---

Traditional CI tools often focus on concepts like steps and pipelines. These provide the basic building blocks
that allow teams to craft complex workflows for various use cases.

However, these building blocks essentially mean most CI platforms are variations on a "thing doer", a tool that
just executes scripts in response to events.

Having spoken to a number of engineering teams from different companies, I've seen a clear pattern in the kinds of
CI/CD pipelines that teams usually create. When teams construct their CI/CD pipelines, they will typically end up 
with very similar patterns based on a few key concepts. These concepts are **builds**, **deployments**, **environments**
and **tests**.

## Definitions

### Builds

Builds are operations to transform source code into some kind of artifact that can be deployed. This could be a binary, a container image
a tarball or even just a tag in a version control system.

A build may be combined with a deployment into a single step, but if you're dealing with a lot of environments this could result in a
performance hit.

### Deployments

Deployments are operations that make an artifact available in an environment. This could involve applying yaml to Kubernetes, creating a VM in
the cloud or uploading a few files via SFTP.

Deployments may create one resource or many instances, they could also operate on multiple environments in a single operation.

When dealing with many environments, it's important to consider what happens when a deployment fails. Do you redeploy a previous build, do you 
leave things in a potentially bad state and alert a human?

### Environments

Environments are the target for a deployment. They represent a single "location" where a build may be deployed. This could be a staging environment,
a production environment, or even a production environment specific to a region or a single customer.

Most organizations start with a small number of environments, like staging+production or even just a single production environment. But as they grow,
it's likely that additional environments will be needed for different purposes. This can be most pronounced for B2B organizations that have varied
compliance requirements and large customers with specific demands.

### Tests

Tests are operations to verify a collection of source code, an artifact or an existing deployment behave as expected. There are many kinds of tests, which
we won't go into in detail here, but for the context of this article, it's worth noting that some kinds of tests require an environment to run in, and some
can run standalone.

A lot of tests will be intended to verify the result of a build or deployment, so can run as a part of those operations. Alternatively, I've seen some teams
run tests continuously in a single environment to verify its integrity. You could also consider drift detection (checking for manual changes made outside
of the pipeline) as another type of continuously run test.

## How Ocuroot uses these concepts

Ocuroot is intended to work directly with these four concepts to provide an opinionated framework for CI/CD, reducing the number of decisions your team needs
to make when building out your workflow.

This approach also allows Ocuroot to provide some key features (like rollbacks and populating new environments) "for free" while still giving you flexibility
to customize your workflows as you see fit.

## Want to try Ocuroot?

We're currently in a closed alpha. If you're interested in trying Ocuroot for yourself, [book a demo](https://ro.am/ocuroot/) so we can get you started.