---
title: Contributing | KubeDB
description: Contributing
menu:
  docs_0.12.0:
    identifier: contributing-cli
    name: Contributing
    parent: welcome
    weight: 10
menu_name: docs_0.12.0
section_menu_id: welcome
url: /docs/0.12.0/welcome/contributing/
aliases:
  - /docs/0.12.0/CONTRIBUTING/
---

# Contribution Guidelines

Want to hack on KubeDB?

AppsCode projects are [Apache 2.0 licensed](https://github.com/kubedb/cli/blob/master/LICENSE) and accept contributions via GitHub pull requests.  This document outlines some of the conventions on development workflow, commit message formatting, contact points and other resources to make it easier to get your contribution accepted.

## Certificate of Origin

By contributing to this project you agree to the Developer Certificate of Origin (DCO). This document was created by the Linux Kernel community and is a
simple statement that you, as a contributor, have the legal right to make the contribution. See the [DCO](https://github.com/kubedb/cli/blob/master/DCO) file for details.

## Developer Guide

We have a [Developer Guide](/docs/setup/developer-guide/overview.md) that outlines everything you need to know from setting up your dev environment to how to build and test KubeDB. If you find something undocumented or incorrect along the way, please feel free to send a Pull Request.

## Getting Help

We use Slack for public discussions. To chit chat with us or the rest of the community, join us in the [Kubernetes Slack team](https://kubernetes.slack.com/messages/C8149MREV/) channel `#kubedb`. To sign up, use our [Slack inviter](http://slack.kubernetes.io/).

To receive product announcements, please join our [mailing list](https://groups.google.com/forum/#!forum/kubedb) or follow us on [Twitter](https://twitter.com/KubeDB). Our mailing list is also used to share design docs shared via Google docs.

## Bugs/Feature request

If you have found a bug with KubeDB or want to request for new features, please [file an issue](https://github.com/kubedb/project/issues/new).

## Submit PR

If you fix a bug or developed a new feature, feel free to submit a PR. In either case, please file a [Github issue](https://github.com/kubedb/project/issues/new) first, so that we can have a discussion on it. This is a rough outline of what a contributor's workflow looks like:

- Create a topic branch from where you want to base your work (usually master).
- Make commits of logical units.
- Push your changes to a topic branch in your fork of the repository.
- Make sure the tests pass, and add any new tests as appropriate.
- Submit a pull request to the original repository.

Thanks for your contributions!

## Spread the word

If you have written blog post or tutorial on KubeDB, please share it with us on [Twitter](https://twitter.com/KubeDB) or the [Kubernetes Slack team](http://slack.kubernetes.io) channel `#kubedb`.
