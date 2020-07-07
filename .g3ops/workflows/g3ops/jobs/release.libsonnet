{
  needs: "check_manifest",
  "if": "${{ needs.check_manifest.outputs.releaseRequired == true }}",
  "runs-on": "ubuntu-latest",
  steps: [
    { uses: "actions/checkout@v2" },
    { run: "git config --global user.name 'jbrunton-ci-minion'" },
    { run: "git config --global user.email 'jbrunton-ci-minion@outlook.com'" },
    { run: "git fetch --unshallow" },
    { run: "go build" },
    { run: "npm install" },
    { run: "npm run release -- $RELEASE_NAME",
      env: {
        GITHUB_TOKEN: "${{ secrets.GITHUB_TOKEN }}",
        RELEASE_NAME: "${{ steps.check_manifest.outputs.releaseName }}"
      }
    }
  ]
}
