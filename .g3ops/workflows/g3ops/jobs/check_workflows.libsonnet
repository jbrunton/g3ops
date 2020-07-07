local config = import '../config.libsonnet';
local git_config = import '../../common/git.libsonnet';

local check_workflows_step(context) =
  {
    name: 'validate %(name)s workflows' % context,
    env: { 'G3OPS_CONFIG': context.config },
    run: 'g3ops workflows check'
  };

{
  'name': config.job_name,
  'runs-on': 'ubuntu-latest',
  steps: [
    {
      uses: 'actions/checkout@v2'
    },
    {
      uses: 'actions/setup-go@v2',
      with: {
        'go-version': '^1.14.4'
      }
    },
    {
      name: 'install g3ops',
      run: 'go get github.com/jbrunton/g3ops'
    }
  ] + [
    check_workflows_step(context)
    for context in config.g3ops_contexts
  ]
}
