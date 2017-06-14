# kubedb

## Installing

Lets install `kubedb` CLI using `go get` from source code.

Following command will install the latest version of the library from master.

```bash
go get github.com/k8sdb/cli/...
```

## Using KubeDB
Want to learn how to use KubeDB? Please start [here](docs/user-guide/tutorial.md).

`kubedb` CLI is used to manipulate kubedb ThirdPartyResource objects.

We will go through each of the commands and will see how these commands interact with TPR objects for kubedb databases.

* [kubedb init](docs/user-guide/task/init.md) to deploy unified operator.
* [kubedb create](docs/user-guide/task/create.md) to create a database object.
* [kubedb describe](docs/user-guide/task/describe.md) to describe a supported object.
* [kubedb get](docs/user-guide/task/get.md) to get/list supported object(s).
* [kubedb edit](docs/user-guide/task/edit.md) to edit supported object(s).
* [kubedb delete](docs/user-guide/task/delete.md) to delete supported object(s).

## Versioning Policy
There are 2 parts to versioning policy:
 - Operator & cli version: KubeDB follows semver versioning policy. Until 1.0 release is done, there might be breaking changes between point releases of the operator. Please always check the release notes for upgrade instructions.
 - TPR version: kubedb.com/v1alpha1 is considered in alpha. This means breaking changes to the YAML format might happen among different releases of the operator.

---

**The kubedb operator & cli collects anonymous usage statistics to help us learning
how the software is being used and how we can improve it. To disable stats collection,
run the operator with the flag** `--analytics=false`.

---

## Contribution guidelines
Want to help improve KubeDB? Please start [here](https://github.com/k8sdb/cli/tree/master/docs/contribution).

## Support
If you have any questions, you can reach out to us.
* [Slack](https://slack.appscode.com)
* [Forum](https://discuss.appscode.com)
* [Twitter](https://twitter.com/AppsCodeHQ)
* [Website](https://appscode.com)
