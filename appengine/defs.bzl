load("@io_bazel_rules_go//go:def.bzl", "GoLibrary", "GoPath", "GoSource", "go_path")

def _go_appengine_deploy_path_impl(ctx):
    config = ctx.file.config
    gp = ctx.attr.path[GoPath]
    args = []
    if ctx.attr.gomod:
        args.append("-gomod=" + ctx.file.gomod.path)
    if ctx.attr.gosum:
        args.append("-gosum=" + ctx.file.gosum.path)

    output = ctx.actions.declare_file(gp.gopath_file.path + "/src/" + ctx.attr.entry + "/app.yaml")
    ctx.actions.run(
        mnemonic = "deploy",
        executable = ctx.executable._deployer,
        inputs = [gp.gopath_file, config, ctx.file.gomod, ctx.file.gosum],
        arguments = [
            "-config=" + config.path,
            "-entry=" + ctx.attr.entry,
            "-path=" + gp.gopath_file.path,
            "-output=" + output.path,
        ] + args,
        outputs = [output],
    )
    return [DefaultInfo(files = depset([output]))]

go_appengine_deploy_path = rule(
    implementation = _go_appengine_deploy_path_impl,
    attrs = {
        "path": attr.label(
            mandatory = True,
            providers = [GoPath],
        ),
        "gomod": attr.label(
            allow_single_file = True,
        ),
        "gosum": attr.label(
            allow_single_file = True,
        ),
        "entry": attr.string(
            mandatory = True,
            doc = "The entry point where your application is located (eg, github.com/ccontavalli/myapp)",
        ),
        "config": attr.label(
            allow_single_file = True,
            mandatory = True,
            doc = "The app.yaml configuration file to use for gcloud app deploy",
        ),
        "_deployer": attr.label(
            default = Label("//appengine/deployer"),
            allow_single_file = True,
            executable = True,
            cfg = "host",
        ),
    },
)

def go_appengine_deploy(name, entry, deps, config, **kwargs):
    # Packs all the dependencies in some place.
    go_path(name = name + "-dir", deps = deps)
    go_appengine_deploy_path(name = name, path = ":" + name + "-dir", entry = entry, config = config, **kwargs)
