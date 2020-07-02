local git_opts = {
  main_branch: 'master',
  user: 'jbrunton-ci-minion',
  email: 'jbrunton-ci-minion@outlook.com'
};

local build_job = {
  steps: [
    {
      name: 'commit',
      run: |||
        git config --global user.name "%(user)%s"
        git config --global user.email "%(email)%s"
        g3ops commit build ${{ matrix.service }}
        git push origin:%(main_branch)s
      ||| % git_opts
    }
  ]
};

local workflow = {
  on: {
    pull_request: {
      branches: [git_opts.main_branch]
    },
    push: {
      branches: [git_opts.main_branch]
    }
  },
  jobs: [
    build_job
  ],
};

std.manifestYamlDoc(workflow)
