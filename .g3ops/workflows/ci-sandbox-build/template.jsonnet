local git_config = import '../config/git.libsonnet';

local manifest_check_job = {
  name: 'sandbox_manifest_check',
  'runs-on': 'ubuntu-latest',
  'if': "github.event_name == 'push'",
  outputs: {
    buildMatrix: '${{ steps.check.outputs.buildMatrix }}',
    buildRequired: '${{ steps.check.outputs.buildRequired }}'
  },
  steps: [
    { uses: 'actions/checkout@v2' },
    { uses: 'actions/setup-go@v2',
      with: { 'go-version': '^1.14.4' } },
    { name: 'install g3ops',
      run: 'go get github.com/jbrunton/g3ops' },
    { name: 'check manifest',
      id: 'check',
      run: 'g3ops set-outputs build-matrix' },
  ]
};

local build_job = {
  name: 'sandbox_build',
  'runs-on': 'ubuntu-latest',
  needs: 'manifest_check',
  'if': '${{ needs.manifest_check.outputs.buildRequired == true }}',
  strategy: {
    matrix: '${{ fromJson(needs.manifest_check.outputs.buildMatrix) }}',
  },
  env: {
    G3OPS_DOCKER_ACCESS_TOKEN: '${{ secrets.G3OPS_DOCKER_ACCESS_TOKEN }}',
    G3OPS_DOCKER_USERNAME: '${{ secrets.G3OPS_DOCKER_USERNAME }}'
  },
  steps: [
    { uses: 'actions/checkout@v2',
      with: { token: '${{ secrets.G3OPS_ADMIN_ACCESS_TOKEN }}' } },
    { uses: 'actions/setup-go@v2',
      with: { 'go-version': '^1.14.4' } },
    { name: 'install g3ops',
      run: 'go get github.com/jbrunton/g3ops' },
    { name: 'build',
      run: 'g3ops service build ${{ matrix.service }}' },
    { name: 'commit',
      run: |||
        git config --global user.name 'jbrunton-ci-minion'
        git config --global user.email 'jbrunton-ci-minion@outlook.com'
        g3ops commit build ${{ matrix.service }}
        git push origin:master
      ||| }
  ]
};

local workflow = {
  name: 'ci-sandbox-build',
  env: {
    G3OPS_CONTEXT: 'test/sandbox/.g3ops/config.yml'
  },
  on: {
    pull_request: {
      branches: [git_config.main_branch]
    },
    push: {
      branches: [git_config.main_branch]
    }
  },
  jobs: {
    manifest_check: manifest_check_job,
    build: build_job
  },
};

std.manifestYamlDoc(workflow)
