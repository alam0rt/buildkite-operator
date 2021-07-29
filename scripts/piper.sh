#!/usr/bin/env bash

set -eou pipefail


warp() {
  pushd . &>/dev/null
  cd "$1"
  trap 'popd &>/dev/null' INT EXIT
}

base_git_dir="$(git rev-parse --show-toplevel)"
repo="$(git remote get-url origin)"
name="$(basename $(pwd))"
org="${1?please provide an organization name}"
token="${2?please provide an accesstoken name}"


warp "$base_git_dir"

mkdir -p .buildkite &>/dev/null

cat <<'EOF' > .buildkite/pipeline.yaml
steps:
  - label: ":hammer: Example Script"
    command: "script.sh"
    artifact_paths: "artifacts/*"
    agents:
      queue: "${BUILDKITE_AGENT_META_DATA_QUEUE:-default}"
EOF

cat <<EOF
apiVersion: pipeline.buildkite.alam0rt.io/v1alpha1
kind: Pipeline
metadata:
  name: ${name}-test
spec:
  organization: ${org}
  accessTokenRef: accesstoken-sample
  repository: ${repo}
  configuration: |
    steps:
      - command: buildkite-agent pipeline upload .buildkite/pipeline.yaml
EOF
