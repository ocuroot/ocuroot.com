PORT=8080
PROXY_PORT=8081

local_resource(
    'site',
    cmd='templ generate && go run .',
    deps=["."],
    ignore=["**/*_templ.go", "dist/**", ".wrangler/**", "node_modules/**"],
)

local_resource(
  name="cf_worker",
  serve_cmd='wrangler dev --port=' + str(PORT),
  deps=['wrangler.jsonc']
)

local_resource(
    'sync',
    cmd="yarn global add browser-sync",
    serve_cmd="""browser-sync start \
  --files './dist/**' \
  --port {proxy_port} \
  --proxy 'localhost:{port}' \
  --reload-delay 500 \
  --middleware '{middleware}'""".format(proxy_port=PROXY_PORT, port=PORT, middleware='function(req, res, next) { \
    res.setHeader("Cache-Control", "no-cache, no-store, must-revalidate"); \
    return next(); \
  }'),
  resource_deps=["site", "cf_worker"],
)

