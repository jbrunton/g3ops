# G3Ops

Go-powered Github GitOps. This is an opinionated cli for certain common GitOps tasks using Github Actions.

I don't claim this will implement industry best practices, but it will make my life a lot easier.

## Examples

.g3ops.json
{
  "manifests": "./test/sandbox/manifests",
  "services": "test/sandbox/services"
}

g3ops service ls
ping
hello

g3ops environment ls
production ./manifests/production.yml
staging    ./manifests/staging.yml

g3ops environment describe production
services:
  ping: 0.2.1
  hello: 1.2

