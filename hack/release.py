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

libbuild.REPO_ROOT = expandvars('$GOPATH') + '/src/github.com/kubedb/cli'
DATABASES = ['postgres', 'elasticsearch', 'etcd', 'mysql', 'mongodb', 'memcached', 'redis']
RELEASE_TAGS = {
    'cli': '0.9.0-rc.0',
    'operator': '0.9.0-rc.0',
    'apimachinery': '0.9.0-rc.0',
    'postgres': '0.9.0-rc.0',
    'elasticsearch': '0.9.0-rc.0',
    'etcd': '0.1.0-rc.0',
    'mysql': '0.2.0-rc.0',
    'mongodb': '0.2.0-rc.0',
    'memcached': '0.2.0-rc.0',
    'redis': '0.2.0-rc.0',
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
    def __init__(self):
        self.rel_deps = {}
        self.master_deps = {}
        for k in RELEASE_TAGS:
            self.rel_deps['github.com/kubedb/' + k] = RELEASE_TAGS[k]
            self.master_deps['github.com/kubedb/' + k] = 'master'

        print self.rel_deps
        print self.master_deps

    def release_apimachinery(self):
        repo_name = 'apimachinery'
        tag = RELEASE_TAGS[repo_name]
        version = semver.parse(tag)
        release_branch = 'release-{0}.{1}'.format(version['major'], version['minor'])

        repo = libbuild.GOPATH + '/src/github.com/kubedb/' + repo_name
        print(repo)
        print('----------------------------------------------------------------------------------------')
        call('git clean -xfd', cwd=repo)
        git_checkout('master', cwd=repo)
        call('git pull --rebase origin master', cwd=repo)
        call('glide slow', cwd=repo)
        if git_requires_commit(tag, cwd=repo):
            call('./hack/make.py', cwd=repo)
            call('git add --all', cwd=repo)
            call('git commit -a -m "Prepare release {0}"'.format(tag), cwd=repo, eoe=False)
            call('git push origin master', cwd=repo)
        else:
            call('git reset HEAD --hard', cwd=repo)
        git_checkout(release_branch, cwd=repo)
        call('git merge master', cwd=repo)
        call('git tag -fa {0} -m "Release {0}"'.format(tag), cwd=repo)
        call('git push origin {0} --tags --force'.format(release_branch), cwd=repo)

    def release_db(self, repo_name):
        tag = RELEASE_TAGS[repo_name]
        version = semver.parse(tag)
        release_branch = 'release-{0}.{1}'.format(version['major'], version['minor'])

        repo = libbuild.GOPATH + '/src/github.com/kubedb/' + repo_name
        print(repo)
        print('----------------------------------------------------------------------------------------')
        call('git clean -xfd', cwd=repo)
        git_checkout('master', cwd=repo)
        call('git pull --rebase origin master', cwd=repo)
        with open(repo + '/glide.yaml', 'r+') as glide_file:
            glide_config = yaml.load(glide_file)
            glide_mod(glide_config, self.rel_deps)
            glide_write(glide_file, glide_config)
            call('glide slow', cwd=repo)
            if git_requires_commit(tag, cwd=repo):
                call('./hack/make.py', cwd=repo)
                call('git add --all', cwd=repo)
                call('git commit -a -m "Prepare release {0}"'.format(tag), cwd=repo, eoe=False)
                call('git push origin master', cwd=repo)
            else:
                call('git reset HEAD --hard', cwd=repo)
            git_checkout(release_branch, cwd=repo)
            call('git merge master', cwd=repo)
            call('git tag -fa {0} -m "Release {0}"'.format(tag), cwd=repo)
            call('git push origin {0} --tags --force'.format(release_branch), cwd=repo)
            call('rm -rf dist', cwd=repo)
            call('./hack/release.sh', cwd=repo)
            git_checkout('master', cwd=repo)
            glide_mod(glide_config, self.master_deps)
            glide_write(glide_file, glide_config)
            call('git commit -a -m "Start next dev cycle"', cwd=repo, eoe=False)
            call('git push origin master', cwd=repo)

    def release_server_binary(self, repo_name):
        tag = RELEASE_TAGS[repo_name]
        version = semver.parse(tag)
        release_branch = 'release-{0}.{1}'.format(version['major'], version['minor'])

        repo = libbuild.GOPATH + '/src/github.com/kubedb/' + repo_name
        print(repo)
        print('----------------------------------------------------------------------------------------')
        call('git clean -xfd', cwd=repo)
        git_checkout('master', cwd=repo)
        call('git pull --rebase origin master', cwd=repo)
        with open(repo + '/glide.yaml', 'r+') as glide_file:
            glide_config = yaml.load(glide_file)
            glide_mod(glide_config, self.rel_deps)
            glide_write(glide_file, glide_config)
            call('glide slow', cwd=repo)
            if git_requires_commit(tag, cwd=repo):
                call('./hack/make.py', cwd=repo)
                call('git add --all', cwd=repo)
                call('git commit -a -m "Prepare release {0}"'.format(tag), cwd=repo, eoe=False)
                call('git push origin master', cwd=repo)
            else:
                call('git reset HEAD --hard', cwd=repo)
            git_checkout(release_branch, cwd=repo)
            call('git merge master', cwd=repo)
            call('git tag -fa {0} -m "Release {0}"'.format(tag), cwd=repo)
            call('git push origin {0} --tags --force'.format(release_branch), cwd=repo)
            call('rm -rf dist', cwd=repo)
            call('./hack/release.sh', cwd=repo)
            git_checkout('master', cwd=repo)
            glide_mod(glide_config, self.master_deps)
            glide_write(glide_file, glide_config)
            call('git commit -a -m "Start next dev cycle"', cwd=repo, eoe=False)
            call('git push origin master', cwd=repo)

    def release_cli(self):
        repo_name = 'cli'
        tag = RELEASE_TAGS[repo_name]
        version = semver.parse(tag)
        release_branch = 'release-{0}.{1}'.format(version['major'], version['minor'])

        repo = libbuild.GOPATH + '/src/github.com/kubedb/' + repo_name
        print(repo)
        print('----------------------------------------------------------------------------------------')
        call('git clean -xfd', cwd=repo)
        git_checkout('master', cwd=repo)
        with open(repo + '/glide.yaml', 'r+') as glide_file:
            glide_config = yaml.load(glide_file)
            glide_mod(glide_config, self.rel_deps)
            glide_write(glide_file, glide_config)
            call('glide slow', cwd=repo)
            if git_requires_commit(tag, cwd=repo):
                call('./hack/make.py', cwd=repo)
                call('git add --all', cwd=repo)
                call('git commit -a -m "Prepare release {0}"'.format(tag), cwd=repo, eoe=False)
                call('git push origin master', cwd=repo)
            else:
                call('git reset HEAD --hard', cwd=repo)
            git_checkout(release_branch, cwd=repo)
            call('git merge master', cwd=repo)
            call('git tag -fa {0} -m "Release {0}"'.format(tag), cwd=repo)
            call('git push origin {0} --tags --force'.format(release_branch), cwd=repo)
            call('rm -rf dist', cwd=repo)
            call('env APPSCODE_ENV=prod ./hack/make.py build', cwd=repo)
            git_checkout('master', cwd=repo)
            glide_mod(glide_config, self.master_deps)
            glide_write(glide_file, glide_config)
            call('git commit -a -m "Start next dev cycle"', cwd=repo, eoe=False)
            call('git push origin master', cwd=repo)


def release(comp=None):
    cat = Kitten()
    if comp is None:
        cat.release_apimachinery()
        for name in DATABASES:
            cat.release_db(name)
        cat.release_server_binary('operator')
        cat.release_cli()
    elif comp == 'apimachinery':
        cat.release_apimachinery()
    elif comp in DATABASES:
        cat.release_db(comp)
    elif comp in ['operator']:
        cat.release_server_binary(comp)
    elif comp == 'cli':
        cat.release_cli()


if __name__ == "__main__":
    if len(sys.argv) == 1:
        release(None)
    elif len(sys.argv) > 1:
        # http://stackoverflow.com/a/834451
        # http://stackoverflow.com/a/817296
        release(*sys.argv[1:])
    else:
        print('Usage ./hack/release.py [component]')
