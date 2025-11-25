---
title: "Navigating the IaC Bootstrap Problem"
slug: iac-bootstrap-problem
excerpt: "When working with IaC platforms, there's usually some infra you need to configure to get them working. With Terraform, you need somewhere to keep your state, for example. How do you configure this infra? Are you stuck doing it manually, or could you automate it?"
coverImage:
  src: "/assets/blog/iac-bootstrap-problem/cover.jpg" 
  alt: "A pair of brown work boots on a mossy tree trunk"
  credit: "Photo by Lum3n"
  creditURL: "https://www.pexels.com/photo/brown-work-boots-167706/"
date: "2025-11-13T11:00:00-04:00"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/iac-bootstrap-problem/cover.jpg"
---

The promise of Infrastructure-As-Code (IAC) is that you always have a complete definition of your stack in 
a versioned source repo. This makes it easier for a team to collaborate on changes, roll back to an earlier
version when something goes wrong or even create a totally new copy of your stack from scratch if needed.

But in order to take advantage of these benefits, you need to set up your IaC tool, and usually that involves
setting up some infrastructure to support it. With Terraform, for example, you need a shared state store. With
open-source Terraform (or OpenTofu), this is most often in an S3-compatible storage bucket.

So where does that bucket come from? The quickest and easiest option would be through "click-ops", using a web UI
to set things up manually. This isn't very reproducible, and good luck tracking changes over time.

It's not just Terraform where we have this problem. Ansible needs a control node setting up, Crossplane needs
Kubernetes. All of this extra infrastructure also needs credentials, so you'll also need to be able to manage
secrets.

## Why does this matter?

In the grand scheme of things, a single storage bucket and a secret or two doesn't seem like a lot to configure,
so you might feel ok with this little bit of manual configuration. And how often will a Terraform state bucket
need to change?

As your organization grows, you might find more and more reasons to make changes. You may want to limit access 
to the bucket to specific members of the team for safety, potentially a lot of work to manage all the roles.

Maybe you add more teams and want to provide them with dedicated state storage for flexibility and to minimize
the blast radius of mistakes. Do you want to set up all of those buckets and credentials manually?

You may be very unlucky and have an outage that requires you to rebuild everything from a backup. Recreating
your bucket manually could slow you down at a critical time.

## What are the alternatives?

As with any problem, there are a lot of ways you could approach bootstrapping your IaC infrastructure. Here are
a few with their own benefits and drawbacks.

### Rigorous documentation

Manual configuration can always benefit from, well, a manual. Having clear documentation of how the IaC
infrastructure is created and playbooks for maintentance and updates will simplify the effort involved.

If your manual steps involve a web UI, there will likely be changes to the UI over time. This will cause
your documentation to "rot" at a quicker pace than API or SDK documentation might. If you're going to rely
on the docs for emergency situations, you better be prepared to exercise them regularly to make sure they're
up-to-date.

Of course, with written documentation there is also the risk of ambiguity, and you need to work hard to make
sure that anyone on the team can apply the documentation, not just the person who wrote it.

### Scripted CLI calls

All decent cloud providers have their own CLI, and this would allow you to script the setup of your IaC
infrastructure. Of course, this is no replacement for a decent IaC tool, so changes may well need to be
applied via separate calls and then integrated into the script after the fact. For anything but the most
simple of cases this can quickly become very difficult to manage.

### Let AI do it

In 2025, I'd be remiss to ignore the options that AI presents. In this case, it's one step up from the
documentation option. As something of a hybrid between documentation and scripting, you can write setup
instructions as prompts to an LLM. You can even use an LLM to help you write those instructions.

Then when you're ready, you can re-run these prompts with the appropriate tools or MCPs configured. This
does leave you needing to configure the right tools, it's own dependency problem. But a bigger issue is
that of non-determinism. There's no guarantee that giving an LLM the same prompt twice will produce the
same results. If you're relying on your IaC infrastructure being configured securely, this could leave
you with some serious problems.

This isn't to mention the risk of your AI making its own decision to delete your IaC infrastructure, leaving
you unable to make changes.

For those of you with management experience, you could see this as a "trust, but verify" problem. If you let
an AI do the work, you can't discount the time needed to double-check what it did.

### Self-managing IaC

In some cases, you can pull the IaC infrastructure into your IaC itself. With Terraform, for example, you could
set up your state bucket in Terraform using local state and then copy the state file into the bucket.

This seems like the best of both worlds at first glance, but can put you into some very tricky situations. If you
delete the bucket from your Terraform module, you would lose all the state and be unable to proceed. Some
operations can result in resources being destroyed and recreated, which again would lose your state.

While it can seem clever, you really don't want to have IaC tools managing the infrastructure critical to their
own operation.

### Tiered infrastructure

With Terraform, configuring your state storage bucket in a Terraform module can still be a viable option if you keep
the original state separate. For bootstrapping, you can use local state to create the store, and share the statefile
in an appropriately secure manner (you probably want to avoid checking it into a git repo in case the state includes
credentials).

This would require a certain amount of care to ensure that nobody is applying changes at the same time, and everyone
has an up to date copy of the statefile. To make this easier, you want to minimise the number of resources being
managed locally.

Once you have your state bucket, you can create the rest of your infrastructure with as many modules as you like. This
splits your infrastructure into two tiers: "bootstrapped" and "managed".

If you end up in a situation where you need multiple state stores, you can insert a new layer above the local "bootstrapped"
state, which can be managed using normal IaC practices.

The key is to minimize the surface area of your lowest tier.

## Conclusion

Bootstrapping IaC can be fiddly and annoying, but depending on your needs, can be solved in quite elegant ways. The real
trick is to not let yourself overthink it and adopt a solution that suits your team for the immediate future.