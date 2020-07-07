local git_config = import '../../../../../.g3ops/workflows/config/git.libsonnet';

local hello_world_job = {
  name: 'hello world',
  'runs-on': 'ubuntu-latest',
  steps: [
    {
      run: 'echo Hello, World!'
    },
  ],
};

local workflow = {
  on: {
    pull_request: {
      branches: [git_config.main_branch]
    },
    push: {
      branches: [git_config.main_branch]
    }
  },
  jobs: {
    hello: hello_world_job
  }
};

std.manifestYamlDoc(workflow)
