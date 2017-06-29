# Contributing guidelines

## Developer Guide
We have a [Developer's Guide](../developer-guide/README.md) that outlines everything you need to know from setting up your
dev environment to how to get faster Pull Request reviews. If you find something undocumented or incorrect along the way,
please feel free to send a Pull Request.

## Filing issues
If you have a question about KubeDB or have a problem using it, please start with contacting us in Github issues or Slack.
If that doesn't answer your questions, or if you think you found a bug, please [file an issue](https://github.com/k8sdb/operator/issues/new).

## Submit PR
If you fix a bug or developed a new feature feel free to submit a PR.

1. Fork the projects
1. Add Your changes
1. Add Test Cases to justify the changes
1. Run the tests
1. Build the project
1. Run e2e tests
1. Submit PR. And you are go. Contact @tamal @aerokite for review.

### Adding Dependency
If your patch Depends on new packagees, add that to vendor with glide.

### Building Operator
Read [build instruction](../developer-guide/build.md) to build operators.
