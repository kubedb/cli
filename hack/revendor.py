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
DATABASES = ['postgres', 'elasticsearch', 'mysql', 'mongodb', 'memcached', 'redis']
REPO_LIST = DATABASES + ['cli', 'operator', 'apimachinery']
KUTIL_VERSION = 'release-8.0'
KUBEMON_VERSION = 'release-8.0'
FORCED_DEPS = {
    'github.com/cpuguy83/go-md2man': 'v1.0.8',
}


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


def glide_mod(glide_config, changes):
    for dep in glide_config['import']:
        if dep['package'] in changes:
            dep['version'] = changes[dep['package']]
    for pkg, ver in FORCED_DEPS.iteritems():
        found = False
        for dep in glide_config['import']:
            if dep['package'] == pkg:
                dep['version'] = ver
                found = True
                break
        if not found:
            glide_config['import'].append({
                'package': pkg,
                'version': ver,
            })


def glide_write(f, glide_config):
    f.seek(0)
    pkg = glide_config.pop('package')
    out = 'package: ' + pkg + '\n' + yaml.dump(glide_config, default_flow_style=False)
    f.write(out)
    f.truncate()
    glide_config['package'] = pkg


class Kitten(object):
    def __init__(self):
        self.seed = ''.join(random.choice(string.ascii_uppercase + string.digits) for _ in range(6))
        self.master_deps = {}
        for k in REPO_LIST:
            self.master_deps['github.com/kubedb/' + k] = 'master'
        self.master_deps['github.com/appscode/kutil'] = KUTIL_VERSION
        self.master_deps['github.com/appscode/kube-mon'] = KUBEMON_VERSION
        print self.master_deps

    def revendor_db(self, repo_name):
        revendor_branch = 'api-{0}'.format(self.seed)

        repo = libbuild.GOPATH + '/src/github.com/kubedb/' + repo_name
        print(repo)
        print('----------------------------------------------------------------------------------------')
        call('git reset HEAD --hard', cwd=repo)
        call('git clean -xfd', cwd=repo)
        git_checkout('master', cwd=repo)
        call('git pull --rebase origin master', cwd=repo)
        git_checkout(revendor_branch, cwd=repo)
        with open(repo + '/glide.yaml', 'r+') as glide_file:
            glide_config = yaml.load(glide_file)
            glide_mod(glide_config, self.master_deps)
            glide_write(glide_file, glide_config)
            call('glide slow', cwd=repo)
            if git_requires_commit(cwd=repo):
                call('git add --all', cwd=repo)
                call('git commit -s -a -m "Revendor api"', cwd=repo, eoe=False)
                call('git push origin {0}'.format(revendor_branch), cwd=repo)
            else:
                call('git reset HEAD --hard', cwd=repo)

    def revendor_server_binary(self, repo_name):
        revendor_branch = 'api-{0}'.format(self.seed)

        repo = libbuild.GOPATH + '/src/github.com/kubedb/' + repo_name
        print(repo)
        print('----------------------------------------------------------------------------------------')
        call('git reset HEAD --hard', cwd=repo)
        call('git clean -xfd', cwd=repo)
        git_checkout('master', cwd=repo)
        call('git pull --rebase origin master', cwd=repo)
        git_checkout(revendor_branch, cwd=repo)
        with open(repo + '/glide.yaml', 'r+') as glide_file:
            glide_config = yaml.load(glide_file)
            glide_mod(glide_config, self.master_deps)
            glide_write(glide_file, glide_config)
            call('glide slow', cwd=repo)
            if git_requires_commit(cwd=repo):
                call('git add --all', cwd=repo)
                call('git commit -s -a -m "Revendor api"', cwd=repo, eoe=False)
                call('git push origin {0}'.format(revendor_branch), cwd=repo)
            else:
                call('git reset HEAD --hard', cwd=repo)

    def revendor_cli(self):
        repo_name = 'cli'
        revendor_branch = 'api-{0}'.format(self.seed)

        repo = libbuild.GOPATH + '/src/github.com/kubedb/' + repo_name
        print(repo)
        print('----------------------------------------------------------------------------------------')
        call('git reset HEAD --hard', cwd=repo)
        call('git clean -xfd', cwd=repo)
        git_checkout('master', cwd=repo)
        call('git pull --rebase origin master', cwd=repo)
        git_checkout(revendor_branch, cwd=repo)
        with open(repo + '/glide.yaml', 'r+') as glide_file:
            glide_config = yaml.load(glide_file)
            glide_mod(glide_config, self.master_deps)
            glide_write(glide_file, glide_config)
            call('glide slow', cwd=repo)
            if git_requires_commit(cwd=repo):
                call('git add --all', cwd=repo)
                call('git commit -s -a -m "Revendor api"', cwd=repo, eoe=False)
                call('git push origin {0}'.format(revendor_branch), cwd=repo)
            else:
                call('git reset HEAD --hard', cwd=repo)


def revendor(comp=None):
    cat = Kitten()
    if comp is None:
        for name in DATABASES:
            cat.revendor_db(name)
    elif comp == 'all':
        for name in DATABASES:
            cat.revendor_db(name)
        cat.revendor_server_binary('operator')
        cat.revendor_cli()
    elif comp in DATABASES:
        cat.revendor_db(comp)
    elif comp == 'operator':
        cat.revendor_server_binary(comp)
    elif comp == 'cli':
        cat.revendor_cli()


if __name__ == "__main__":
    if len(sys.argv) == 1:
        revendor(None)
    elif len(sys.argv) > 1:
        # http://stackoverflow.com/a/834451
        # http://stackoverflow.com/a/817296
        revendor(*sys.argv[1:])
    else:
        print('Usage ./hack/revendor.py [component]')
