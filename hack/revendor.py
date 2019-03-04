#!/usr/bin/env python


# http://stackoverflow.com/a/14050282
def check_antipackage():
    from sys import version_info
    sys_version = version_info[:2]
    found = True
    if sys_version < (3, 0):
        # 'python 2'
        from pkgutil import find_loader
        found = find_loader('antipackage') is not None
    elif sys_version <= (3, 3):
        # 'python <= 3.3'
        from importlib import find_loader
        found = find_loader('antipackage') is not None
    else:
        # 'python >= 3.4'
        from importlib import util
        found = util.find_spec('antipackage') is not None
    if not found:
        print('Install missing package "antipackage"')
        print('Example: pip install git+https://github.com/ellisonbg/antipackage.git#egg=antipackage')
        from sys import exit
        exit(1)
check_antipackage()

# ref: https://github.com/ellisonbg/antipackage
import antipackage
from github.appscode.libbuild import libbuild

import os
import os.path
import random
import string
import subprocess
import sys
from os.path import expandvars
import yaml
from collections import Counter

libbuild.REPO_ROOT = expandvars('$GOPATH') + '/src/github.com/kubedb/cli'
DATABASES = ['postgres', 'elasticsearch', 'etcd', 'mysql', 'mongodb', 'memcached', 'redis']
REPO_LIST = DATABASES + ['cli', 'operator', 'apimachinery']
REQUIRED_DEPS = [
    {
      "package": "github.com/cpuguy83/go-md2man",
      "version": "v1.0.8"
    },
    {
      "package": "github.com/russross/blackfriday",
      "version": "v1.5.2"
    },
    {
      "package": "github.com/json-iterator/go",
      "version": "1.1.5"
    },
    {
      "package": "github.com/spf13/cobra",
      "version": "v0.0.3"
    },
    {
      "package": "github.com/spf13/pflag",
      "version": "v1.0.3"
    },
    {
      "package": "golang.org/x/text",
      "version": "b19bf474d317b857955b12035d2c5acb57ce8b01"
    },
    {
      "package": "golang.org/x/net",
      "version": "0ed95abb35c445290478a5348a7b38bb154135fd"
    },
    {
      "package": "golang.org/x/sys",
      "version": "95c6576299259db960f6c5b9b69ea52422860fce"
    },
    {
      "package": "golang.org/x/crypto",
      "version": "de0752318171da717af4ce24d0a2e8626afaeb11"
    },
    {
      "package": "github.com/golang/protobuf",
      "version": "v1.1.0"
    },
    {
      "package": "github.com/davecgh/go-spew",
      "version": "v1.1.1"
    },
    {
      "package": "k8s.io/kube-openapi",
      "version": "c59034cc13d587f5ef4e85ca0ade0c1866ae8e1d"
    },
    {
      "package": "gopkg.in/yaml.v2",
      "version": "v2.2.1"
    },
    {
      "package": "github.com/gorilla/websocket",
      "version": "v1.4.0"
    },
    {
      "package": "gopkg.in/square/go-jose.v2",
      "version": "v2.2.1"
    },
    {
      "package": "github.com/imdario/mergo",
      "version": "v0.3.5"
    },
    {
      "package": "github.com/mitchellh/mapstructure",
      "version": "v1.1.2"
    },
    {
      "package": "github.com/go-ini/ini",
      "version": "v1.40.0"
    },
    {
      "package": "gopkg.in/ini.v1",
      "version": "v1.40.0"
    },
    {
      "package": "sigs.k8s.io/yaml",
      "version": "v1.1.0"
    },
    {
      "package": "github.com/prometheus/client_golang",
      "version": "v0.9.2"
    },
    {
      "package": "k8s.io/utils",
      "version": "66066c83e385e385ccc3c964b44fd7dcd413d0ed"
    }
]
DEP_LIST = [
    {
      "package": "github.com/cpuguy83/go-md2man",
      "version": "v1.0.8"
    },
    {
      "package": "github.com/json-iterator/go",
      "version": "1.1.5"
    },
    {
      "package": "github.com/coreos/prometheus-operator",
      "version": "v0.29.0"
    },
    {
      "package": "k8s.io/api",
      "version": "kubernetes-1.13.0"
    },
    {
      "package": "k8s.io/apiextensions-apiserver",
      "version": "kubernetes-1.13.0"
    },
    {
      "package": "k8s.io/apimachinery",
      "repo": "https://github.com/kmodules/apimachinery.git",
      "vcs": "git",
      "version": "ac-1.13.0"
    },
    {
      "package": "k8s.io/apiserver",
      "repo": "https://github.com/kmodules/apiserver.git",
      "vcs": "git",
      "version": "ac-1.13.0"
    },
    {
      "package": "k8s.io/client-go",
      "version": "v10.0.0"
    },
    {
      "package": "k8s.io/cli-runtime",
      "version": "kubernetes-1.13.0"
    },
    {
      "package": "k8s.io/kubernetes",
      "version": "v1.13.0"
    },
    {
      "package": "k8s.io/kube-aggregator",
      "version": "kubernetes-1.13.0"
    },
    {
      "package": "k8s.io/metrics",
      "version": "kubernetes-1.13.0"
    },
    {
      "package": "kmodules.xyz/client-go",
      "version": "release-10.0"
    },
    {
      "package": "kmodules.xyz/webhook-runtime",
      "version": "release-10.0"
    },
    {
      "package": "kmodules.xyz/custom-resources",
      "version": "release-10.0"
    },
    {
      "package": "kmodules.xyz/monitoring-agent-api",
      "version": "release-10.0"
    },
    {
      "package": "kmodules.xyz/objectstore-api",
      "version": "release-10.0"
    },
    {
      "package": "kmodules.xyz/offshoot-api",
      "version": "release-10.0"
    },
    {
      "package": "kmodules.xyz/openshift",
      "version": "release-10.0"
    },
    {
      "package": "github.com/graymeta/stow",
      "repo": "https://github.com/appscode/stow.git",
      "vcs": "git",
      "version": "master"
    },
    {
      "package": "github.com/Azure/azure-sdk-for-go",
      "version": "v21.3.0"
    },
    {
      "package": "github.com/Azure/go-autorest",
      "version": "v11.1.0"
    },
    {
      "package": "github.com/aws/aws-sdk-go",
      "version": "v1.14.12"
    },
    {
      "package": "google.golang.org/api/storage/v1",
      "version": "3639d6d93f377f39a1de765fa4ef37b3c7ca8bd9"
    },
    {
      "package": "cloud.google.com/go",
      "version": "v0.23.0"
    },
    {
      "package": "github.com/spf13/afero",
      "version": "v1.1.2"
    },
    {
      "package": "github.com/appscode/osm",
      "version": "0.10.0"
    },
    {
      "package": "github.com/kubepack/onessl",
      "version": "0.11.0"
    }
]
DELETE_LIST=[
    "github.com/openshift/api",
    "github.com/openshift/client-go",
    "github.com/openshift/origin",
    "github.com/appscode/ocutil"
]


def die(status):
    if status:
        sys.exit(status)


def call(cmd, stdin=None, cwd=libbuild.REPO_ROOT, eoe=True):
    print(cwd + ' $ ' + cmd)
    status = subprocess.call([expandvars(cmd)], shell=True, stdin=stdin, cwd=cwd)
    if eoe:
        die(status)
    else:
        return status


def check_output(cmd, stdin=None, cwd=libbuild.REPO_ROOT):
    print(cwd + ' $ ' + cmd)
    return subprocess.check_output([expandvars(cmd)], shell=True, stdin=stdin, cwd=cwd)


def git_branch_exists(branch, cwd=libbuild.REPO_ROOT):
    return call('git show-ref --quiet refs/heads/{0}'.format(branch), eoe=False, cwd=cwd) == 0


def git_checkout(branch, cwd=libbuild.REPO_ROOT):
    call('git fetch --all --prune', cwd=cwd)
    call('git fetch --tags', cwd=cwd)
    if git_branch_exists(branch, cwd):
        call('git checkout {0}'.format(branch), cwd=cwd)
    else:
        call('git checkout -b {0}'.format(branch), cwd=cwd)


def git_requires_commit(cwd=libbuild.REPO_ROOT):
    changed_files = check_output('git diff --name-only', cwd=cwd).strip().split('\n')
    return Counter(changed_files) != Counter(['glide.lock'])


def sortDep(val):
    return val['package']


def glide_mod(glide_config, changes):
    for dep in glide_config['import']:
        if dep['package'] in changes:
            dep['version'] = changes[dep['package']]
    for x in REQUIRED_DEPS:
        for idx, dep in enumerate(glide_config['import']):
            found = False
            if dep['package'] == x['package']:
                glide_config['import'][idx] = x
                found = True
                break
        if not found:
            glide_config['import'].append(x)
    for x in DEP_LIST:
        for idx, dep in enumerate(glide_config['import']):
            if dep['package'] == x['package']:
                glide_config['import'][idx] = x
                break
    for package in DELETE_LIST:
        for idx, dep in enumerate(glide_config['import']):
            if dep['package'] == package:
                del glide_config['import'][idx]
                break
    glide_config['import'].sort(key=sortDep)


def glide_write(f, glide_config):
    f.seek(0)
    pkg = glide_config.pop('package')
    out = 'package: ' + pkg + '\n' + yaml.dump(glide_config, default_flow_style=False)
    f.write(out)
    f.truncate()
    glide_config['package'] = pkg


class DepFixer(object):
    def __init__(self):
        self.seed = ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(6))
        self.master_deps = {}
        for k in REPO_LIST:
            self.master_deps['github.com/kubedb/' + k] = 'master'
        print self.master_deps

    def revendor_repo(self, repo_name):
        revendor_branch = 'api-{0}'.format(self.seed)

        repo = libbuild.GOPATH + '/src/github.com/kubedb/' + repo_name
        print(repo)
        print('----------------------------------------------------------------------------------------')
        call('git reset HEAD --hard', cwd=repo)
        call('git clean -xfd', cwd=repo)
        git_checkout('master', cwd=repo)
        call('git pull --rebase origin master', cwd=repo)
        git_checkout(revendor_branch, cwd=repo)
        # https://stackoverflow.com/a/6759339/244009
        call("find " + repo + "/apis -type f -exec sed -i -e 's/k8s.io\\/apimachinery\\/pkg\\/api\\/testing\\/roundtrip/k8s.io\\/apimachinery\\/pkg\\/api\\/apitesting\\/roundtrip/g' {} \;", eoe=False)
        with open(repo + '/glide.yaml', 'r+') as glide_file:
            glide_config = yaml.load(glide_file)
            glide_mod(glide_config, self.master_deps)
            glide_write(glide_file, glide_config)
            call('glide slow', cwd=repo)
            if git_requires_commit(cwd=repo):
                call('git add --all', cwd=repo)
                call('git commit -s -a -m "Revendor dependencies"', cwd=repo, eoe=False)
                call('git push origin {0}'.format(revendor_branch), cwd=repo)
            else:
                call('git reset HEAD --hard', cwd=repo)


def revendor(comp=None):
    cat = DepFixer()
    if comp is None:
        for name in DATABASES:
            cat.revendor_repo(name)
    elif comp == 'all':
        cat.revendor_repo('apimachinery')
        for name in DATABASES:
            cat.revendor_repo(name)
        cat.revendor_repo('operator')
        cat.revendor_repo('cli')
    elif comp in DATABASES:
        cat.revendor_repo(comp)
    elif comp == 'operator':
        cat.revendor_repo(comp)
    elif comp == 'cli':
        cat.revendor_repo(comp)
    elif comp == 'apimachinery':
        cat.revendor_repo(comp)


if __name__ == "__main__":
    if len(sys.argv) == 1:
        revendor(None)
    elif len(sys.argv) > 1:
        # http://stackoverflow.com/a/834451
        # http://stackoverflow.com/a/817296
        revendor(*sys.argv[1:])
    else:
        print('Usage ./hack/revendor.py [component]')
