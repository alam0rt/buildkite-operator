#!/usr/bin/env bash

set -eou pipefail


warp() {
  pushd .
  git rev-parse --show-toplevel
  trap 'popd' INT EXIT
}

repo="$(git remote get-url origin)"
name="$(basename $(pwd))"
org="${1?please provide an organization name}"
token="${2?please provide an accesstoken name}"

warp

mkdir -p .buildkite
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
  step:
    - command: "buildkite-agent pipeline upload .buildkite/pipeline.yaml"
EOF
