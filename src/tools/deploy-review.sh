#!/usr/local/bin/bash

DEPLOY_URL="${CI_ENVIRONMENT_SLUG}.prospector.ie"
echo "Deploy URL: ${DEPLOY_URL}"
levant deploy \
  -var git_sha="${CI_COMMIT_SHORT_SHA}" \
  -var environment_slug="${CI_ENVIRONMENT_SLUG}" \
  -var deploy_url="${DEPLOY_URL}" \
  -address "http://nomad.service.consul:4646" \
  src/tools/templates/prospector-review.hcl
