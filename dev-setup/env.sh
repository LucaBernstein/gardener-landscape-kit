#!/usr/bin/env bash

# SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
#
# SPDX-License-Identifier: Apache-2.0

export REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

# Git Repositories
export GLK_BASE_REPO_NAME="base"
export GLK_BASE_REPO_PATH="$REPO_ROOT/dev/e2e/$GLK_BASE_REPO_NAME"
export GLK_LANDSCAPE_REPO_NAME="test-landscape"
export GLK_LANDSCAPE_REPO_PATH="$REPO_ROOT/dev/e2e/$GLK_LANDSCAPE_REPO_NAME"
export GLK_CONFIG_PATH="${GLK_BASE_REPO_PATH}/landscapekitconfiguration.yaml"

# Git Server
export GIT_USER_NAME="gitops"
export GIT_USER_PASSWORD="testtest"
export GIT_SERVER_URL="http://git.local.gardener.cloud:6080"

# Registry
export REGISTRY_HOSTNAME="glk-registry.local.gardener.cloud"
export RUNNER_JOB_IMAGE_SOURCE="codeberg.org/oranki/forgejo-runner:3.5.1"
export RUNNER_JOB_IMAGE_LOCAL="${REGISTRY_HOSTNAME}:6001/forgejo-runner-job:3.5.1"

# Kind Cluster
export GLK_KIND_CLUSTER_PREFIX=${GLK_KIND_CLUSTER_PREFIX-glk}
export GLK_CLUSTER_NAME=$GLK_KIND_CLUSTER_PREFIX-single
export GLK_KUBECONFIG=${REPO_ROOT}/dev/kind-$GLK_CLUSTER_NAME-kubeconfig.yaml

