# ocuroot.com

Source for the Ocuroot website, hosted at https://www.ocuroot.com.

## Implementation

The site is built using [Templ](https://templ.guide/), a component-based templating system for Go. Templates are rendered to files using a custom generator.

UI components are provided by the Ocuroot UI module: https://github.com/ocuroot/ui.

## Testing

The site can be run locally with [Tilt](https://tilt.dev/).

```bash
tilt up
```

This will build the site, start a local instance of the Cloudflare Worker and serve via a proxy.

Changes to the source will result in a rebuild and the proxy will
automatically reload the page.

## Deployment

Ocuroot is used to deploy the site to a Cloudflare Worker.
Content is first uploaded to a staging area, with a manual approval
to promote to production.

The configuration for this release process is at [release.ocu.star](release.ocu.star).