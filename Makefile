# Copyright AppsCode Inc. and Contributors
#
# Licensed under the AppsCode Community License 1.0.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://github.com/appscode/licenses/raw/1.0.0/AppsCode-Community-1.0.0.md
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

SHELL=/bin/bash -o pipefail

GO_PKG   := kubedb.dev
REPO     := $(notdir $(shell pwd))
BIN      := kubectl-dba
COMPRESS ?= no

# This version-strategy uses git tags to set the version string
git_branch       := $(shell git rev-parse --abbrev-ref HEAD)
git_tag          := $(shell git describe --exact-match --abbrev=0 2>/dev/null || echo "")
commit_hash      := $(shell git rev-parse --verify HEAD)
commit_timestamp := $(shell date --date="@$$(git show -s --format=%ct)" --utc +%FT%T)

VERSION          := $(shell git describe --tags --always --dirty)
version_strategy := commit_hash
ifdef git_tag
	VERSION := $(git_tag)
	version_strategy := tag
else
	ifeq (,$(findstring $(git_branch),master HEAD))
		ifneq (,$(patsubst release-%,,$(git_branch)))
			VERSION := $(git_branch)
			version_strategy := branch
		endif
	endif
endif

###
### These variables should not need tweaking.
###

SRC_PKGS := cmd pkg
SRC_DIRS := $(SRC_PKGS) hack/gendocs # directories which hold app source (not vendored)

DOCKER_PLATFORMS := linux/amd64 linux/arm linux/arm64
BIN_PLATFORMS    := $(DOCKER_PLATFORMS) windows/amd64 darwin/amd64 darwin/arm64

# Used internally.  Users should pass GOOS and/or GOARCH.
OS   := $(if $(GOOS),$(GOOS),$(shell go env GOOS))
ARCH := $(if $(GOARCH),$(GOARCH),$(shell go env GOARCH))

GO_VERSION       ?= 1.25
BUILD_IMAGE      ?= ghcr.io/appscode/golang-dev:$(GO_VERSION)

OUTBIN = bin/$(BIN)-$(OS)-$(ARCH)
ifeq ($(OS),windows)
  OUTBIN := bin/$(BIN)-$(OS)-$(ARCH).exe
  BIN := $(BIN).exe
endif

# Directories that we need created to build/test.
BUILD_DIRS  := bin/$(OS)_$(ARCH)     \
               .go/bin/$(OS)_$(ARCH) \
               .go/cache             \
               hack/config           \
               $(HOME)/.credentials  \
               $(HOME)/.kube         \
               $(HOME)/.minikube

DOCKER_REPO_ROOT := /go/src/$(GO_PKG)/$(REPO)

# If you want to build all binaries, see the 'all-build' rule.
# If you want to build all containers, see the 'all-container' rule.
# If you want to build AND push all containers, see the 'all-push' rule.
all: fmt build

# For the following OS/ARCH expansions, we transform OS/ARCH into OS_ARCH
# because make pattern rules don't match with embedded '/' characters.

build-%:
	@$(MAKE) build                        \
	    --no-print-directory              \
	    GOOS=$(firstword $(subst _, ,$*)) \
	    GOARCH=$(lastword $(subst _, ,$*))

all-build: $(addprefix build-, $(subst /,_, $(BIN_PLATFORMS)))
ifeq ($(COMPRESS),yes)
	@cd bin; \
	sha256sum $(patsubst $(BIN)-windows-%.tar.gz,$(BIN)-windows-%.zip, $(addsuffix .tar.gz, $(addprefix $(BIN)-, $(subst /,-, $(BIN_PLATFORMS))))) > $(BIN)-checksums.txt
endif

version: version-PROD version-DBG
	@echo IMAGE=$(IMAGE)
	@echo BIN=$(BIN)
	@echo version=$(VERSION)
	@echo version_strategy=$(version_strategy)
	@echo git_tag=$(git_tag)
	@echo git_branch=$(git_branch)
	@echo commit_hash=$(commit_hash)
	@echo commit_timestamp=$(commit_timestamp)
version-%:
	@echo TAG_$*=$(TAG_$*)

gen:
	@true

fmt: $(BUILD_DIRS)
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    $(BUILD_IMAGE)                                          \
	    /bin/bash -c "                                          \
	        REPO_PKG=$(GO_PKG)                                  \
	        ./hack/fmt.sh $(SRC_DIRS)                           \
	    "

build: $(OUTBIN)

# The following structure defeats Go's (intentional) behavior to always touch
# result files, even if they have not changed.  This will still run `go` but
# will not trigger further work if nothing has actually changed.

$(OUTBIN): .go/$(OUTBIN).stamp
	@true

# This will build the binary under ./.go and update the real binary iff needed.
.PHONY: .go/$(OUTBIN).stamp
.go/$(OUTBIN).stamp: $(BUILD_DIRS)
	@echo "making $(OUTBIN)"
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    $(BUILD_IMAGE)                                          \
	    /bin/bash -c "                                          \
	        ARCH=$(ARCH)                                        \
	        OS=$(OS)                                            \
	        VERSION=$(VERSION)                                  \
	        version_strategy=$(version_strategy)                \
	        git_branch=$(git_branch)                            \
	        git_tag=$(git_tag)                                  \
	        commit_hash=$(commit_hash)                          \
	        commit_timestamp=$(commit_timestamp)                \
	        ./hack/build.sh                                     \
	    "
	@if ! cmp -s .go/bin/$(OS)_$(ARCH)/$(BIN) $(OUTBIN); then   \
	    mv .go/bin/$(OS)_$(ARCH)/$(BIN) $(OUTBIN);              \
	    date >$@;                                               \
	fi
ifeq ($(COMPRESS),yes)
ifeq ($(OS),windows)
	@echo "compressing $(OUTBIN)";                               \
	cd bin;                                                      \
	zip -j $(subst .exe,,$(BIN))-$(OS)-$(ARCH).zip $(subst .exe,,$(BIN))-$(OS)-$(ARCH).exe ../LICENSE.md
else
	@echo "compressing $(OUTBIN)";                               \
	cd bin;                                                      \
	tar -czvf $(BIN)-$(OS)-$(ARCH).tar.gz $(BIN)-$(OS)-$(ARCH) ../LICENSE.md
endif
endif
	@echo

.PHONY: test
test: unit-tests

unit-tests: $(BUILD_DIRS)
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    $(BUILD_IMAGE)                                          \
	    /bin/bash -c "                                          \
	        ARCH=$(ARCH)                                        \
	        OS=$(OS)                                            \
	        VERSION=$(VERSION)                                  \
	        ./hack/test.sh $(SRC_PKGS)                          \
	    "

ADDTL_LINTERS   := gofmt,goimports,unparam

.PHONY: lint
lint: $(BUILD_DIRS)
	@echo "running linter"
	@docker run                                                 \
	    -i                                                      \
	    --rm                                                    \
	    -u $$(id -u):$$(id -g)                                  \
	    -v $$(pwd):/src                                         \
	    -w /src                                                 \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin                \
	    -v $$(pwd)/.go/bin/$(OS)_$(ARCH):/go/bin/$(OS)_$(ARCH)  \
	    -v $$(pwd)/.go/cache:/.cache                            \
	    --env HTTP_PROXY=$(HTTP_PROXY)                          \
	    --env HTTPS_PROXY=$(HTTPS_PROXY)                        \
	    --env GO111MODULE=on                                    \
	    --env GOFLAGS="-mod=vendor"                             \
	    $(BUILD_IMAGE)                                          \
	    golangci-lint run --enable $(ADDTL_LINTERS) --timeout=10m --exclude-files="generated.*\.go$\" --exclude-dirs-use-default

$(BUILD_DIRS):
	@mkdir -p $@

.PHONY: dev
dev: gen fmt build

.PHONY: verify
verify: verify-gen verify-modules

.PHONY: verify-modules
verify-modules:
	GO111MODULE=on go mod tidy
	GO111MODULE=on go mod vendor
	@if !(git diff --exit-code HEAD); then \
		echo "go module files are out of date"; exit 1; \
	fi

.PHONY: verify-gen
verify-gen: gen fmt
	@if !(git diff --exit-code HEAD); then \
		echo "files are out of date, run make gen fmt"; exit 1; \
	fi

.PHONY: add-license
add-license:
	@echo "Adding license header"
	@docker run --rm 	                                 \
		-u $$(id -u):$$(id -g)                           \
		-v /tmp:/.cache                                  \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                   \
		-w $(DOCKER_REPO_ROOT)                           \
		--env HTTP_PROXY=$(HTTP_PROXY)                   \
		--env HTTPS_PROXY=$(HTTPS_PROXY)                 \
		$(BUILD_IMAGE)                                   \
		ltag -t "./hack/license" --excludes "vendor contrib" -v

.PHONY: check-license
check-license:
	@echo "Checking files for license header"
	@docker run --rm 	                                 \
		-u $$(id -u):$$(id -g)                           \
		-v /tmp:/.cache                                  \
		-v $$(pwd):$(DOCKER_REPO_ROOT)                   \
		-w $(DOCKER_REPO_ROOT)                           \
		--env HTTP_PROXY=$(HTTP_PROXY)                   \
		--env HTTPS_PROXY=$(HTTPS_PROXY)                 \
		$(BUILD_IMAGE)                                   \
		ltag -t "./hack/license" --excludes "vendor contrib" --check -v

.PHONY: ci
ci: verify check-license lint build unit-tests #cover

.PHONY: qa
qa:
	@if [ "$$APPSCODE_ENV" = "prod" ]; then                                              \
		echo "Nothing to do in prod env. Are you trying to 'release' binaries to prod?"; \
		exit 1;                                                                          \
	fi
	@if [ "$(version_strategy)" = "tag" ]; then               \
		echo "Are you trying to 'release' binaries to prod?"; \
		exit 1;                                               \
	fi
	@$(MAKE) clean all-build --no-print-directory

.PHONY: release
release:
	@if [ "$$APPSCODE_ENV" != "prod" ]; then      \
		echo "'release' only works in PROD env."; \
		exit 1;                                   \
	fi
	@if [ "$(version_strategy)" != "tag" ]; then                    \
		echo "apply tag to release binaries and/or docker images."; \
		exit 1;                                                     \
	fi
	@$(MAKE) clean all-build --no-print-directory

.PHONY: clean
clean:
	rm -rf .go bin

.PHONY: gen-krew-manifest
gen-krew-manifest:
	@if [ "$(version_strategy)" != "tag" ]; then \
		echo "apply tag to generate krew manifest."; \
		exit 1; \
	fi
	@sed \
      -e 's|{VERSION}|$(VERSION)|g' \
      -e 's|{SHA256SUM_DARWIN_AMD64}|$(shell sha256sum bin/$(BIN)-darwin-amd64.tar.gz | cut -d' ' -f1)|g' \
      -e 's|{SHA256SUM_DARWIN_ARM64}|$(shell sha256sum bin/$(BIN)-darwin-arm64.tar.gz | cut -d' ' -f1)|g' \
      -e 's|{SHA256SUM_LINUX_AMD64}|$(shell sha256sum bin/$(BIN)-linux-amd64.tar.gz | cut -d' ' -f1)|g' \
      -e 's|{SHA256SUM_LINUX_ARM}|$(shell sha256sum bin/$(BIN)-linux-arm.tar.gz | cut -d' ' -f1)|g' \
      -e 's|{SHA256SUM_LINUX_ARM64}|$(shell sha256sum bin/$(BIN)-linux-arm64.tar.gz | cut -d' ' -f1)|g' \
      -e 's|{SHA256SUM_WINDOWS_AMD64}|$(shell sha256sum bin/$(BIN)-windows-amd64.zip | cut -d' ' -f1)|g' \
      hack/krew/plugin.yaml
