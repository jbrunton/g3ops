local workflow = {
  on: {
    pull_request: {
      branches: ['master']
    },
    push: {
      branches: ['master']
    }
  },
  jobs: [
    {
      name: "hello-world",
      steps: [
        {
          run: "echo hello, world!"
        },
      ],
    },
  ],
};

std.manifestYamlDoc(workflow)
