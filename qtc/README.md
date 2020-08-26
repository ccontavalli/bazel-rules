The [quicktemplate](https://github.com/valyala/quicktemplate) library provides a templating language
that is compiled into native .go code, that you can just build and link into your application.

To use the bazel rules here:

1. Make sure you import the quicktemplate toolkit in your WORKSPACE. Assuming you are using
   `gazelle` and `go modules`, the easiest way is to run:

       go get -u github.com/valyala/quicktemplate
       bazel run //:gazelle -- update-repos -from_file=go.mod

   If you want to do it manually, you can add:

       go_repository(
           name = "com_github_valyala_quicktemplate",
           importpath = "github.com/valyala/quicktemplate",
           sum = "h1:k0vgK7zlmFzqAoIBIOrhrfmZ6JoTGJlLRPLbkPGr2/M=",
           version = "v1.6.2",
       )

    To your workspace file.

2. Add the rules in this directory to your WORKSPACE:

       http_archive(
           name = "com_github_ccontavalli_bazel_rules",
           strip_prefix = "bazel-rules-4",
           urls = ["https://github.com/ccontavalli/bazel-rules/archive/v4.tar.gz"],
       )


2. Use the rules. In your `BUILD.bazel` file, use something like:

       load("//bazel:qtc.bzl", "qtpl_go_library")
               
       qtpl_go_library(
           name = "templates",
           srcs = [
               "main.qtpl",
               "header.qtpl",
           ],  
           importpath = "github.com/your/normal/import/path",
           visibility = ["//visibility:public"],
           deps = [
               "//my/normal/golang/deps/rpc:grpc-go",
               "@com_github_dustin_go_humanize//:go_default_library",
           ],  
       )       

