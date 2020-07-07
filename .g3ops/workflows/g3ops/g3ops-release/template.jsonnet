local config = import 'config.libsonnet';
local git_config = import '../../config/git.libsonnet';

local check_manifest_job = {
  "runs-on": "ubuntu-latest",
  outputs: {
    releaseRequired: "${{ steps.check.outputs.releaseRequired }}",
    releaseName: "${{ steps.check.outputs.releaseName }}"
  },
  steps: [
    { uses: "actions/checkout@v2" },
    { uses: "actions/setup-go@v2",
      with: { "go-version": "^1.14.4" }},
    { name: "install g3ops",
      run: "go get github.com/jbrunton/g3ops" },
    { name: "check manifest",
      id: "check",
      env: {
        GITHUB_TOKEN: "${{ secrets.%(token_name)s }}" % config.jobs.check_manifest
      },
      run: "g3ops ci check release-manifest" }
  ]
};

local release_job = {
  needs: "check_manifest",
  "if": "${{ needs.check_manifest.outputs.releaseRequired == true }}",
  "runs-on": "ubuntu-latest",
  steps: [
    { uses: "actions/checkout@v2" },
    { run: "git config --global user.name '%(user)s'" % git_config },
    { run: "git config --global user.email '%(email)s'" % git_config },
    { run: "git fetch --unshallow" },
    { run: "go build" },
    { run: "npm install" },
    { run: "npm run release -- $RELEASE_NAME",
      env: {
        GITHUB_TOKEN: "${{ secrets.%(token_name)s }}" % config.jobs.release,
        RELEASE_NAME: "${{ steps.check_manifest.outputs.releaseName }}"
      }
    }
  ]
};

local workflow = {
  name: "g3ops-release",
  on: {
    push: {
      branches: [git_config.main_branch]
    }
  },
  jobs: {
    check_manifest: check_manifest_job,
    release: release_job
  },
};

std.manifestYamlDoc(workflow)
