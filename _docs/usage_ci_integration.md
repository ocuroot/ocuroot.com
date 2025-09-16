---
title: "CI Integration"
path: "usage/ci-integration"
---

As detailed in the [Scheduling Work](/docs/usage/scheduling-work) page, Ocuroot integrates with CI
platforms via a set of commands that can be executed in response to commits or changes to intent.

This page details concrete configuration that can be applied to various CI platforms to integrate with
Ocuroot.

## General approach

We're always aiming for Ocuroot to be tooling-agnostic, and your choice of CI platform is no exception!
Integrating Ocuroot with your platform of choice is as simple as configuring execution of three key commands.

### Releases

Releases are where everything starts, and you will typically want to trigger a release on every commit to the
default branch of your source repo(s).

When a new push is received, you'll want to configure your platform to run:

```bash
# Create a new release
ocuroot release new release.ocu.star

# Trigger any follow on work
ocuroot work trigger
```

Where *release.ocu.star* is the path to an Ocuroot config for the appropriate release. In a monorepo environment,
you may want to run `ocuroot release new` multiple times for different files. To help with controlling which commands
are run based on which files change, we've provided a free tool: [ifdiff](https://github.com/ocuroot/ifdiff).

### Triggering follow-on work

When `ocuroot release new` finishes, there may be follow-on work to do in other commits or even other repos.
At this point, Ocuroot will attempt to "trigger" the work. So that it knows what to do, you need to configure
a [trigger](/docs/reference/sdk/repo/) function in your `repo.ocu.star`.

Your trigger function takes in a commit hash, and should result in the CI platform executing the `ocuroot work any` command for that commit. Your implementation function must be registered by calling `trigger`.

For example, a `repo.ocu.star` file may contain:

```python
def do_trigger(commit):
    # Example remote
    remote = "git@github.com/example/example.git"

    # Clone the repo
    wd = shell("mktemp -d /tmp/ocuroot-work-any-$$").stdout.strip()
    shell("git clone {} {}".format(remote, wd))
    # Checkout the commit
    shell("git checkout {}".format(commit), dir=wd)
    # Run the command
    shell("ocuroot work any --comprehensive", dir=wd)

trigger(do_trigger)
```

This example runs locally for illustrative purposes, clones the repo into a temp directory, checks out the
target commit and then runs `ocuroot work any` on it to complete any necessary work.

### Handling Intent updates

Finally, you need to be able to handle changes to the intent repo, like manual modification of environments or custom
state.

Assuming you have a trigger function configured in all your source repos, you just need to make sure the following
command is run on push to the appropriate branch of your intent repo:

```bash
ocuroot work trigger --intent
```

This will opportunistically trigger work for the current commit of each source repo. For this reason it is important
to ensure that your trigger functions do not depend on running in the context of the source repo itself.

## GitHub Actions

### Releases

Below is an example workflow config to execute release on pushes to the "main" branch of your repo.

You would put this file under `.github/workflows/ocuroot-release.yml`.

```yaml
name: Ocuroot Release

on:
  push:
    branches: [main]

jobs:
  release:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3

      - name: Download ocuroot binary
        run: |
          ocuroot_tar=$(curl -L -s ${{ vars.OCUROOT_TAR_URL }})
          echo "$ocuroot_tar" | tar -xzf - ocuroot

      # Set up git to report state changes as coming from Ocuroot
      - name: Configure Git
        run: |
          git config --global user.email "ocuroot@example.com"
          git config --global user.name "Ocuroot"

      - name: Release all packages
        env:
          GH_TOKEN: ${{ secrets.PAT_TOKEN }}
        run: |
          ./ocuroot release new release.ocu.star

      - name: Trigger following work
        env:
          GH_TOKEN: ${{ secrets.PAT_TOKEN }} # Use the PAT token so we can trigger workflows
        run: ./ocuroot work trigger
```

There are two additional pieces of configuration to be aware of.

`vars.OCUROOT_TAR_URL` should contain a download URL for the desired release of the Ocuroot binary. See the [releases](https://github.com/ocuroot/ocuroot/releases) page for details of the latest assets.

`secrets.PAT_TOKEN` should contain a [GitHub PAT](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) token with access to the Git repo(s) containins state and intent, as well as permission to trigger workflows. GitHub Actions does provide a default token, but this will
not have the required permissions.

### Triggering follow-on work

In GitHub actions, we can use workflow dispatch to schedule jobs. The following trigger function can
be added to your `repo.ocu.star` to make the appropriate API call.

```python
def do_trigger(commit):
    print("Triggering work for repo at commit " + commit)
    
    # Get environment variables
    env_vars = host.env()
    
    if "GH_TOKEN" in env_vars:
        gh_token = env_vars["GH_TOKEN"]
        
        # Repository owner and name from the repo URL
        # EDIT THESE FOR YOUR SOURCE REPO
        owner = "example"
        repo = "example"
        
        # GitHub API endpoint for workflow dispatch
        workflow_id = "ocuroot-work-any.yml"
        url = "https://api.github.com/repos/{}/{}/actions/workflows/{}/dispatches".format(owner, repo, workflow_id)
        
        # Payload with the commit to check out
        payload = json.encode({"ref": "main", "inputs": {"commit_sha": commit}})
        
        # Headers for authentication
        headers = {
            "Accept": "application/vnd.github+json",
            "Authorization": "token " + gh_token,
            "X-GitHub-Api-Version": "2022-11-28"
        }
        
        print("Triggering workflow via GitHub API")
        response = http.post(url=url, body=payload, headers=headers)
        
        if response["status_code"] == 204:
            print("Successfully triggered workflow")
        else:
            print("Failed to trigger workflow. Status code: " + str(response["status_code"]))
            print("Response: " + response["body"])
    else:
        print("GH_TOKEN not available. Cannot trigger GitHub workflow.")

trigger(do_trigger)
```

This needs a workflow to execute, so add the following config to `.github/workflows/ocuroot-work-any.yml`.

```yaml
name: Ocuroot Work Any

on:
  workflow_dispatch:
    inputs:
      commit_sha:
        description: "Commit SHA to check out"
        required: true
        type: string

run-name: Continue on ${{ github.event.inputs.commit_sha }}

jobs:
  continue:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3
        with:
          ref: ${{ github.event.inputs.commit_sha }}

      - name: Download ocuroot binary
        run: |
          ocuroot_tar=$(curl -L -s ${{ vars.OCUROOT_TAR_URL }})
          echo "$ocuroot_tar" | tar -xzf - ocuroot

      # Set up git to report state changes as coming from Ocuroot
      - name: Configure Git
        run: |
          git config --global user.email "ocuroot@example.com"
          git config --global user.name "Ocuroot"

      # Continue any work for this commit
      - name: Run ocuroot work any --comprehensive
        env:
          GH_TOKEN: ${{ secrets.PAT_TOKEN }}
        run: ./ocuroot work any --comprehensive
```

### Handing Intent updates

Lastly, you'll need to add a workflow to your **intent repo** to handle intent changes:

```yaml
# This holds an example for a workflow that
# should be stored in the state branch
name: Ocuroot Apply Intent

on:
  push:
    branches: [intent] # EDIT FOR THE APPROPRIATE BRANCH
    paths:
      - "**/\\+*"
      - "**/\\+*/**"

jobs:
  intent:
    runs-on: ubuntu-latest
    steps:
      - name: Checkout repository
        uses: actions/checkout@v3

      - name: Download ocuroot binary
        run: |
          ocuroot_tar=$(curl -L -s ${{ vars.OCUROOT_TAR_URL }})
          echo "$ocuroot_tar" | tar -xzf - ocuroot

      # Set up git to report state changes as coming from Ocuroot
      - name: Configure Git
        run: |
          git config --global user.email "ocuroot@example.com"
          git config --global user.name "Ocuroot"

      - name: Trigger work from intent
        env:
          GH_TOKEN: ${{ secrets.PAT_TOKEN }}
        run: ./ocuroot work trigger --intent
```

## Coming soon

We're aiming to provide instructions for other platforms over time, including, but not limited to:

* Jenkins
* CircleCI

If you'd like to see instructions for your CI platform of choice here, feel free to raise a [GitHub issue](https://github.com/ocuroot/ocuroot/issues).

