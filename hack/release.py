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
from github.appscode.pysemver import semver

import os
import os.path
import subprocess
import sys
from os.path import expandvars
import yaml
from collections import Counter

libbuild.REPO_ROOT = expandvars('$GOPATH') + '/src/github.com/k8sdb/cli'

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
    call('git clean -xfd', cwd=cwd)
    call('git fetch --all --prune', cwd=cwd)
    call('git fetch --tags', cwd=cwd)
    if git_branch_exists(branch, cwd):
        call('git checkout {0}'.format(branch), cwd=cwd)
    else:
        call('git checkout -b {0}'.format(branch), cwd=cwd)


def git_requires_commit(tag, cwd=libbuild.REPO_ROOT):
    status = call('git rev-parse {0} >/dev/null 2>&1'.format(tag), eoe=False, cwd=cwd)
    if status == 0:
        return False
    changed_files = check_output('git diff --name-only', cwd=cwd).strip().split('\n')
    return Counter(changed_files) != Counter(['glide.lock'])


def glide_mod(glide_config, changes):
    for dep in glide_config['import']:
        if dep['package'] in changes:
            dep['version'] = changes[dep['package']]


def glide_write(f, glide_config):
    f.seek(0)
    pkg = glide_config.pop('package')
    out = 'package: ' + pkg + '\n' + yaml.dump(glide_config, default_flow_style=False)
    f.write(out)
    f.truncate()
    glide_config['package'] = pkg


class Kitten(object):
    def __init__(self, tag):
        self.tag = tag
        self.version = semver.parse(self.tag)
        self.next_version = semver.bump_minor(self.tag)
        self.release_branch = 'release-{0}.{1}'.format(self.version['major'], self.version['minor'])
        self.rel_deps = {
            'github.com/k8sdb/apimachinery': self.release_branch,
            'github.com/k8sdb/postgres': self.release_branch,
            'github.com/k8sdb/elasticsearch': self.release_branch,
        }
        self.master_deps = {
            'github.com/k8sdb/apimachinery': 'master',
            'github.com/k8sdb/postgres': 'master',
            'github.com/k8sdb/elasticsearch': 'master',
        }

    def release_apimachinery(self):
        repo = libbuild.GOPATH + '/src/github.com/k8sdb/apimachinery'
        print(repo)
        print('----------------------------------------------------------------------------------------')
        git_checkout('master', cwd=repo)
        call('glide slow', cwd=repo)
        if git_requires_commit(self.tag, cwd=repo):
            call('./hack/make.py', cwd=repo)
            call('git commit -a -m "Prepare release {0}"'.format(self.tag), cwd=repo, eoe=False)
            call('git push origin master', cwd=repo)
        else:
            call('git reset HEAD --hard', cwd=repo)
        git_checkout(self.release_branch, cwd=repo)
        call('git merge master', cwd=repo)
        call('git push origin {0}'.format(self.release_branch), cwd=repo)

    def release_db(self, repo_name, short_code):
        repo = libbuild.GOPATH + '/src/github.com/k8sdb/' + repo_name
        print(repo)
        print('----------------------------------------------------------------------------------------')
        git_checkout('master', cwd=repo)
        with open(repo + '/glide.yaml', 'r+') as glide_file:
            glide_config = yaml.load(glide_file)
            glide_mod(glide_config, self.rel_deps)
            glide_write(glide_file, glide_config)
            call('glide slow', cwd=repo)
            if git_requires_commit(self.tag, cwd=repo):
                call('./hack/make.py', cwd=repo)
                call('git commit -a -m "Prepare release {0}"'.format(self.tag), cwd=repo, eoe=False)
                call('git push origin master', cwd=repo)
            else:
                call('git reset HEAD --hard', cwd=repo)
            git_checkout(self.release_branch, cwd=repo)
            call('git merge master', cwd=repo)
            call('git tag -fa {0} -m "Release {0}"'.format(self.tag), cwd=repo)
            call('git push origin {0} --tags --force'.format(self.release_branch), cwd=repo)
            call('rm -rf dist', cwd=repo)
            call('./hack/docker/{0}-operator/setup.sh'.format(short_code), cwd=repo)
            call('env APPSCODE_ENV=prod ./hack/docker/{0}-operator/setup.sh release'.format(short_code), cwd=repo)
            git_checkout('master', cwd=repo)
            glide_mod(glide_config, self.master_deps)
            glide_write(glide_file, glide_config)
            call('git commit -a -m "Start {0} dev cycle"'.format(self.next_version), cwd=repo, eoe=False)
            call('git push origin master', cwd=repo)

    def release_cli(self):
        repo = libbuild.GOPATH + '/src/github.com/k8sdb/cli'
        print(repo)
        print('----------------------------------------------------------------------------------------')
        git_checkout('master', cwd=repo)
        with open(repo + '/glide.yaml', 'r+') as glide_file:
            glide_config = yaml.load(glide_file)
            glide_mod(glide_config, self.rel_deps)
            glide_write(glide_file, glide_config)
            call('glide slow', cwd=repo)
            if git_requires_commit(self.tag, cwd=repo):
                call('./hack/make.py', cwd=repo)
                call('git commit -a -m "Prepare release {0}"'.format(self.tag), cwd=repo, eoe=False)
                call('git push origin master', cwd=repo)
            else:
                call('git reset HEAD --hard', cwd=repo)
            git_checkout(self.release_branch, cwd=repo)
            call('git merge master', cwd=repo)
            call('git tag -fa {0} -m "Release {0}"'.format(self.tag), cwd=repo)
            call('git push origin {0} --tags --force'.format(self.release_branch), cwd=repo)
            call('rm -rf dist', cwd=repo)
            call('env APPSCODE_ENV=prod ./hack/make.py build', cwd=repo)
            git_checkout('master', cwd=repo)
            glide_mod(glide_config, self.master_deps)
            glide_write(glide_file, glide_config)
            call('git commit -a -m "Start {0} dev cycle"'.format(self.next_version), cwd=repo, eoe=False)
            call('git push origin master', cwd=repo)


def release(tag=None):
    cat = Kitten(tag)
    cat.release_apimachinery()
    cat.release_db('postgres', 'pg')
    cat.release_db('elasticsearch', 'es')
    cat.release_cli()


if __name__ == "__main__":
    if len(sys.argv) > 1:
        # http://stackoverflow.com/a/834451
        # http://stackoverflow.com/a/817296
        release(*sys.argv[1:])
    else:
        print('Usage ./hack/release.py 0.3.0')
