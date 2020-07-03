local config = import 'config.libsonnet';
local git_config = import '../common/git.libsonnet';

local workflow_check_job = {
  'name': config.job_name,
  'runs-on': 'ubuntu-latest',
  steps: [
    {
      name: 'install g3ops',
      run: 'go get github.com/jbrunton/g3ops'
    },
    {
      name: 'validate workflows',
      run: 'g3ops workflows check'
    }
  ]
};

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
  jobs: [
    workflow_check_job
  ],
};

std.manifestYamlDoc(workflow)
