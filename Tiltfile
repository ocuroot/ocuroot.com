PORT=8080
PROXY_PORT=8081

local_resource(
    'site',
    cmd='templ generate',
    serve_cmd='go run . -dev=true -dev-port=' + str(PORT),
    deps=["."],
    ignore=["**/*_templ.go"],
)

local_resource(
    'sync',
    cmd="yarn global add browser-sync",
    serve_cmd="""browser-sync start \
  --files './**/*_templ.go' \
  --port {proxy_port} \
  --proxy 'localhost:{port}' \
  --reload-delay 500 \
  --middleware '{middleware}'""".format(proxy_port=PROXY_PORT, port=PORT, middleware='function(req, res, next) { \
    res.setHeader("Cache-Control", "no-cache, no-store, must-revalidate"); \
    return next(); \
  }'),
)

