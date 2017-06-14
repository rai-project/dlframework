#!/bin/sh
set -ex

jq -s '.[0] * .[1]' dlframework.swagger.json swagger-info.json > dlframework.versioned.swagger.json
