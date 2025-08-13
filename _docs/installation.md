---
title: "Installation"
path: "installation"
---

The Ocuroot client is provided as an open source tool that provides you with everything you
need to run Ocuroot-enabled releases from your CI platform of choice.

## Binary downloads

URLs for client binaries are available under each release on [GitHub](https://github.com/ocuroot/ocuroot/releases).

Note that these binaries are currently not signed, so on macOS you will either need to download via `curl` or `wget`
or allow the unsigned binary in Settings.

## From Source

If you have [Go](https://go.dev/) installed, you can build Ocuroot directly from the source repo:

```bash
go install github.com/ocuroot/ocuroot@latest
```

