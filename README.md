# G3Ops

[![Build Status](https://github.com/jbrunton/g3ops/workflows/ci-build/badge.svg?branch=develop)](https://github.com/jbrunton/g3ops/actions?query=branch%3Adevelop+workflow%3Aci-build)
[![Maintainability](https://api.codeclimate.com/v1/badges/2a4862e062bc388ef8a9/maintainability)](https://codeclimate.com/github/jbrunton/g3ops/maintainability)
[![Test Coverage](https://api.codeclimate.com/v1/badges/2a4862e062bc388ef8a9/test_coverage)](https://codeclimate.com/github/jbrunton/g3ops/test_coverage)

Go-powered GitHub GitOps. This is an opinionated cli for certain common GitOps tasks using GitHub Actions.

I don't claim this will implement industry best practices, but it will make my life a lot easier.

# Examples

## Context

### Get current context

    $ g3ops context get
    sandbox

### Describe current context

    $ g3ops context describe
    name: sandbox
    environments:
      production:
        manifest: /Users/jbrunton/go/src/github.com/jbrunton/g3ops/test/sandbox/manifests/production.yml
      staging:
        manifest: /Users/jbrunton/go/src/github.com/jbrunton/g3ops/test/sandbox/manifests/staging.yml
    services:
      hello:
        manifest: /Users/jbrunton/go/src/github.com/jbrunton/g3ops/test/sandbox/services/hello/manifest.yml
      ping:
        manifest: /Users/jbrunton/go/src/github.com/jbrunton/g3ops/test/sandbox/services/ping/manifest.yml
    ci:
      defaults:
        build:
          env:
            TAG: latest
          command: |
            docker-compose build $BUILD_SERVICE

## Services

### List services

    $ g3ops service ls
    +-------+------------------------------------------+
    | NAME  |                 MANIFEST                 |
    +-------+------------------------------------------+
    | ping  | test/sandbox/services/ping/manifest.yml  |
    | hello | test/sandbox/services/hello/manifest.yml |
    +-------+------------------------------------------+

### Describe service

With existing build:

    $ g3ops service ping describe
    name: ping
    manifest: test/sandbox/services/ping/manifest.xml
    version: 1.2
    build:
      buildId: 123
      buildSha: a1b2c3
      imageTag: ...

If build missing:

    $ g3ops service ping describe
    name: ping
    manifest: test/sandbox/services/ping/manifest.xml
    version: 1.2
    build: <missing> 

## Environments

    g3ops environment describe production
    services:
      ping: 0.2.1
      hello: 1.2

