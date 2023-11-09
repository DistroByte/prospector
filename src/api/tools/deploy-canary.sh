#!/usr/local/bin/bash

levant deploy \
  -var git_sha="${CI_COMMIT_SHORT_SHA}" \
  -address "http://nomad.service.consul:4646" \
  tools/templates/api-prod.hcl
