# Mildly useful bazel rules

This repository contains some mildly useful bazel rules.
Check each directory to learn more about the rules provided.

So far, you can use:

* **appengine**/defs.bzl - providing support to deploy golang binaries
  on appengine instance. The rule is quite flexible, and just works even
  with complex golang builds (eg, using gRPC, with embedded artifacts,
  with generated templates).

* **qtc**/defs.bzl - providing support for compiling compiling
  [quicktemplate templates](https://github.com/valyala/quicktemplate) into
  golang source code easy to link in go libraries and binaries.

