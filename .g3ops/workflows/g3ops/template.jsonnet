local config = import 'config.libsonnet';
local git_config = import '../common/git.libsonnet';
local check_workflows_job = import 'jobs/check_workflows.libsonnet';
local check_manifest_job = import 'jobs/check_manifest.libsonnet';
local release_job = import 'jobs/release.libsonnet';

local workflow = {
  name: config.workflow_name,
  on: {
    pull_request: {
      branches: [git_config.main_branch]
    },
    push: {
      branches: [git_config.main_branch]
    }
  },
  jobs: {
    check_workflows: check_workflows_job,
    check_manifest: check_manifest_job,
    release: release_job
  },
};

std.manifestYamlDoc(workflow)
