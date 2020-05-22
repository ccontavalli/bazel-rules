# golang AppEngine rules

This directory provides the `go_appengine_deploy` rule, to deploy your `golang` code on Google App Engine (GAE).

# How to use the rule

1. Install the `gcloud` command on your system (or container) [following the instructions here](https://cloud.google.com/sdk/install),
   and login to your account as necessary. The bazel rules in this repository
   are not hermetic and cannot be, as they rely on you supplying your credentials, and gcloud being set up correctly
   on your system.

2. Make sure your repository is already configured to use go, as per instructions in the
   [bazel golang repository](https://github.com/bazelbuild/rules_go/releases).

3. Configure your `WORKSPACE` to use this repository:

       http_archive(
           name = "com_github_ccontavalli_bazel_rules_appengine",
           sha256 = "3b8da2664d6508b0483adf7359ddf8fa02a4da7c99627ca5e7294ca4e62e39b0",
           strip_prefix = "bazel-rules-2",
           urls = ["https://github.com/ccontavalli/bazel-rules/archive/v2.tar.gz"],
       )

4. Use the rule to define a deploy target. In this target you specify one or more `go` dependencies
   and a few other dependencies. The rule will take care of copying all their dependencies, recursively,
   into a dedicated directory, and then call `gcloud app deploy` from there. For example:

       load("@com_github_ccontavalli_bazel_rules_appengine//appengine:defs.bzl", "go_appengine_deploy")

       go_appengine_deploy(
           # Name of the rule. Go use SCons if you don't know what this is.
           name = "deploy"

           # Once all the dependencies have recursively been copied, this is
           # the directory where `gcloud app deploy` should be called from.
           # It is the directory where the app.yaml file will be created, and
           # where the go.mod, go.sum files will be created.
           entry = "github.com/ccontavalli/whatever",

           # The app.yaml file to use, see below.
           # IMPORTANT: the rule makes no attempt to reconcile the go version in
           # your WORKSPACE with that in the app.yaml file.
           # If you use go1.13 SDK, make sure you have go113 in the app.yaml file.
           config = "deploy/app.yaml",
        
           # (optional) if you use go modules with bazel (you should!) specify
           # the path of go.mod and go.sum to have them deployed to cloud.
           gomod = "go.mod",
           gosum = "go.sum",

           # List of dependencies to deploy. Must be golang targets.
           # To include data dependencies or similar, have your go targets
           # depend on them. (hint: it's not the deployment that needs those
           # files, it's your binary).
           deps = [":your_go_web_app"],
       )

5. Grab some pop corn, and run `bazel run :deploy` to deploy your go binary to GAE.
