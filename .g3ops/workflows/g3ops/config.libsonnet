{
  workflow_name: 'g3ops',
  job_name: 'check-workflows',
  g3ops_contexts: [
    { name: 'g3ops', config: '.g3ops/config.yml' },
    { name: 'sandbox', config: 'test/sandbox/.g3ops/config.yml' }
  ],
}
