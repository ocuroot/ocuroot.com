ocuroot("0.3.0")

def staging(ctx):
    tag = host.shell("git rev-parse --short HEAD", mute=True).stdout.strip()

    # Build static assets and upload to the staging space
    host.shell("""
    templ generate
    go run .
    wrangler versions upload --preview-alias staging --tag "{}"
    """.format(tag))

    version_list = json.decode(host.shell("wrangler versions list --name=www --json", mute=True).stdout.strip())
    
    version_id = ""
    for version in version_list:
        if not "annotations" in version:
            continue
        if not "workers/tag" in version["annotations"]:
            continue
        
        # Use the last version for this tag, as we may have deployed from this commit
        # multiple times.
        if version["annotations"]["workers/tag"] == tag:
            version_id = version["id"]

    if version_id == "":
        fail("Version not found")

    return done(
        outputs={
            "tag": tag,
            "version_id": version_id,
        }
    )

phase(
    name="staging",
    work=[
        call(
            name="staging",
            fn=staging,
        )
    ]
)

def production(ctx):
    host.shell("wrangler versions deploy --yes {}@100".format(ctx.inputs.version_id))

phase(
    name="production",
    work=[
        call(
            name="production",
            fn=production,
            inputs={
                "approval": input(ref="./custom/approval"),
                "version_id": input(ref="./call/staging#output/version_id"),
            },
        )
    ]
)