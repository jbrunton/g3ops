{
  "runs-on": "ubuntu-latest",
  needs: ["check_workflows"],
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
        GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}"
      },
      run: "g3ops ci check release-manifest" }
  ]
}
