---
title: "How Big Terraform Modules Slow You Down"
slug: big-terraform-modules
excerpt: "It might be tempting to put all your infrastructure into a single Terraform module for simplicity. But as you scale, this approach will start slowing you down. Here's what can do wrong and how you can address it."
coverImage:
  src: "/assets/blog/big-terraform-modules/cover.jpg" 
  alt: "A little red car overloaded by luggage"
  credit: "Photo by David Henry"
  creditURL: "https://www.pexels.com/photo/luggage-in-a-little-red-car-10246432/"
date: "2025-10-30T11:00:00-04:00"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/big-terraform-modules/cover.jpg"
---

Terraform (and by extension OpenTofu) is a fantastic tool. With just a few HCL files, you can create infrastructure resources
across multiple platforms at once, and update them at will. Plus it's declarative, so what you get is exactly what's in your
module code, even if someone goes in and changes something manually without you knowing.

But Terraform comes with its own challenges, particularly as you scale. In this post, we'll be looking at what can happen
when your Terraform module contains too many resources.

## What could go wrong?

It can get **slow**. In order to make decisions about what actions to take, Terraform has to refresh the remote state for every resource in the module for every `plan` or `apply`. When you have just a few resources, you might not even notice, but as
complexity increases, you'll be waiting longer and longer. Imagine the frustration when you want to add a tag to a single resource, but you have to wait for it to check 1000 others to even see the plan!

It can seem like all or nothing when working with a large module. If you want to remove and recreate a resource, you need
to destroy everything and start again, right? Not quite, but we'll get onto this later.

There are also **security** concerns. With a single module and state file, you need to give your whole team access to everything.
Not only do you need to be able to access the resources themselves to plan and make changes, but Terraform state files often
include sensitive values. I've seen this lead to tight controls over Terraform state, preventing engineering teams from running
`plan` operations until they submit a PR. At best, this leads to long feedback loops, at worst, you need to pull in multiple teams to get anything done.

Finally, **testing** becomes more difficult. When you have a complex, tightly coupled Terraform module, the only way to see if
it works is to apply the whole thing and see if it works together. Iterating on a single resource becomes a case of cobbling
together something that looks right then throwing it over the wall and waiting for alarms to go off.

## Mitigations

Terraform does have some features that mitigate the all-or-nothing nature of large modules. When running `plan`, `apply` or
`destroy`, you can specify a `-target` to operate only on specific resources. In the case of a destroy, this will also remove
the target's dependencies. This can be useful when something goes wrong, allowing you to remove and recreate a resource in
isolation.

When applying changes to a large module, it can also be useful if you know you only made changes to a single resource. Specify
a target and you avoid having to check every resource for changes.

However, this approach has its drawbacks. You need to be very deliberate when specifying a target. If you make a mistake with an `apply`, you could end up with un-applied changes to resources that could trip someone else up later. If you make a mistake with
a `destroy` you could get rid of something critical, or worse, *all your infrastructure*.

This also doesn't address the security or testing concerns.

## Better Scoping

A far better approach is to have multiple modules that you can work on independantly. Breaking down your Terraform modules into well-considered scopes gives you much more flexibility. You can enforce clearer boundaries for teams to work together, give different teams different permissions and iterate more quickly.

Hashicorp have a great [example](https://developer.hashicorp.com/terraform/tutorials/modules/pattern-module-creation#explore-a-scoping-example) showing how a relatively simply application can be scoped into 
multiple modules. Your breakdown may vary, and will require consideration of service boundaries, deployment patterns and
ownership.

Of course, there will be dependencies between these modules. VMs or Kubernetes clusters need to be attached to networks,
applications need to know where to run and how to talk to their database. You could connect these manually (tedious), or
push values to a secret or key-value store (more things to set up), or you can reach directly into Terraform state!

Terraform includes `terraform_remote_state` which glues together modules by making an output from one into a data source for
another.

```terraform
data "terraform_remote_state" "source_module" {
  backend = "s3"
  config = {
    bucket = "my-terraform-state"
    key    = "prefix/terraform.tfstate"
    region = "us-east-1"
  }
}

resource "aws_instance" "example" {
  ami           = data.terraform_remote_state.source_module.outputs.ami_id
  instance_type = "t2.micro"
}
```

This expands the permissions needed a little, but a team would only need access to their own modules and immediate dependencies.
Another reason to be very deliberate in setting your service boundaries!

Want more power? Tools like Terragrunt and Spacelift can help orchestrate these dependencies and provide a set of mechanisms
for managing many Terraform modules at a high-level.

## The Ocuroot Approach

Ocuroot is being built to make this even easier, allowing you to connect outputs from one module to the inputs of another and trigger automated updates when those inputs changes. So no matter how complicated your architecture is, you can break it down into manageable chunks without losing key relationships.

The best part? It isn't just for Terraform! Using Pulimu? Or a combination of Ansible, Terraform and Helm? You can connect
them all together with a little config code.

Want to try it for yourself? The Ocuroot client is [open source](https://github.com/ocuroot/ocuroot), so there's no cost to experiment! Take a look at
our [quickstart](https://ocuroot.com/ocuroot/quickstart) for a self-guided demo.

If you'd like to keep up to date on what we're building and how you can take advantage of Ocuroot to help your team,
be sure to subscribe to our newsletter!

