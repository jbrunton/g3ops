local git_config = import '../../../../../.g3ops/workflows/common/git.libsonnet';

// local build_job = {
//   steps: [
//     {
//       name: 'commit',
//       run: |||
//         git config --global user.name "%(user)s"
//         git config --global user.email "%(email)s"
//         g3ops commit build ${{ matrix.service }}
//         git push origin:%(main_branch)s
//       ||| % git_config
//     }
//   ]
// };

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
