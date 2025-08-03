---
title: "Back to basics: why I'm writing end to end tests in bash"
slug: back-to-basics-e2e-tests-in-bash
excerpt: "Sometimes the simplest tools are the most effective. Here's why I chose bash scripts when writing end-to-end tests for Ocuroot."
coverImage:
  src: "/assets/blog/back-to-basics-e2e-tests-in-bash/cover.jpg"
  alt: "A house in the hobbiton movie set"
  credit: "Donovan Kelly"
  creditURL: "https://www.pexels.com/photo/house-in-the-hobbiton-movie-set-17824131/"
date: "2025-07-08T14:48:00-04:00"
author:
  name: Tom Elliott
  picture: "/assets/blog/authors/telliott.jpeg"
ogImage:
  url: "/assets/blog/back-to-basics-e2e-tests-in-bash/cover.jpg"
---

End-to-end (or acceptance) testing is a critical part of the testing process. It provides final
confirmation that all the components you've already tested separately can work together for your users
in a more realistic scenario. There are a ton of tools and frameworks out there to help you build
robust end-to-end tests, but for Ocuroot, I've chosen to write my tests as humble bash scripts.

## What I'm testing

For the upcoming release of Ocuroot, the key user workflows primarily involve interacting with
releases and state, with sequences of commands like:

```bash
# Begin a new release that will pause for approval
ocuroot release new config.ocu.star
# Set an approval intent
ocuroot state set "config.ocu.star/+v1/custom/approval" 1
# Perform any remaining work to complete the release
ocuroot work continue
```

Because everything is stored in state and accessible from the command line, you can verify the
result of a release with commands like:

```bash
# Check if the release has been deployed to production
ocuroot state get "config.ocu.star/@v1/deploy/prod"
```

## Why bash works

You might begin to see why bash makes a natural choice in this situation, since pretty much everything
I need to test can be done at the command line anyway.

Immediate convenience aside, writing tests in bash puts me closer to the experiences of my users.
Being able to chain commands together with pipes, `xargs` and the like can make a simple command
line utility incredibly powerful in the right hands. Being forced to write my tests in this way will
immediately give me feedback on the composability of the Ocuroot CLI.

## What tools I could have used

Ocuroot is written in Go, and the Go testing framework has been long been my go-to for
writing tests. It provides not only a pleasant environment for defining tests and associated
helpers, as well as advanced features like test coverage and benchmarking. For unit testing, it will
continue to be my first choice, but for end-to-end testing, it would present a few challenges. Most notably, I would need to either write a number of verbose exec calls, or write an abstraction layer
to execute Ocuroot commands directly from Go code. This would serve only to complicate the test code
and add extra work to add and test new features. Also, in the case of exec calls, code coverage information would be lost, and in the case of an abstraction layer, there would be code above the
test abstraction that would not be covered.

I also have a decent amount of experience with tools like Playwright, a testing framework geared
towards web UIs. While I'm still in the planning process for a new Ocuroot web UI, I previously did some work on integrating Playwright with CLI tools to make it possible to combine the two forms of UI. Although right now, this seems like overkill.

## What I'm sacrificing

Writing tests with bash isn't a perfect experience. A key sacrifice I'm making is that I won't have
access to test-specific features like summarized results and the aforementioned test coverage.

Many dedicated test frameworks will manage execution of multiple tests at once to either parallelize
them or at the very least not stop after the first failure. This isn't something I intend to spend
a lot of time replicating, so I will be slowed down a little by addressing any failures one at a time.

## A blessing in disguise

I also won't have access to helper libraries that simplify the writing of assertions. This required
a bit of up-front work creating functions, but did allow me to write assertions that are specific
to Ocuroot itself, like this one to check the contents of a ref:

```bash
assert_ref_equals() {  
    local ref_path="$1"
    local expected_value="$2"
    local actual_value=$(ocuroot state get "$ref_path" | jq -r '.')
    local error_message="${3:-"Ref $ref_path does not match expected value, expected $expected_value, got $actual_value"}"
    
    if [ "$actual_value" != "$expected_value" ]; then
        echo "$error_message"
        exit 1
    fi
    return 0
}
```

This results in test segments that look like the following:

```bash
# ...

ocuroot release new ./-/backend/package.ocu.star
assert_equal "0" "$?" "Failed to release backend"

ocuroot work continue
assert_equal "0" "$?" "Failed to continue work on this commit"

assert_deployed "frontend/package.ocu.star" "staging"
assert_deployed "frontend/package.ocu.star" "production"
assert_deployed "frontend/package.ocu.star" "production2"

assert_ref_equals "./-/backend/package.ocu.star/@/deploy/staging#output/credential" "abcd"
assert_ref_equals "./-/frontend/package.ocu.star/@/deploy/staging#output/backend_credential" "abcd"
```

Which has the benefit of being much more readable than shelling out in a language like Go.

## Conclusion

While it may feel "quaint" to write tests in bash, I've found it to be a refreshing change of pace. And when working on a CLI tool, it can actually make for easier to follow and maintain tests.

## What's next?

I'm in the final stages of preparing the SDK v0.3 release. My aim is to share a version that you can try out very soon. In the meantime, you can follow
Ocuroot on [LinkedIn](https://www.linkedin.com/company/ocuroot), [BlueSky](https://bsky.app/profile/ocuroot.com) or get in touch directly by booking a [demo](/demo).
