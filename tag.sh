#!/bin/bash

set -e

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Usage: tag.sh <version>"
    exit 1
fi

git tag $VERSION
git tag cmd/risor/$VERSION
git tag modules/aws/$VERSION
git tag modules/bcrypt/$VERSION
git tag modules/cli/$VERSION
git tag modules/color/$VERSION
git tag modules/gha/$VERSION
git tag modules/github/$VERSION
git tag modules/goquery/$VERSION
git tag modules/htmltomarkdown/$VERSION
git tag modules/image/$VERSION
git tag modules/isatty/$VERSION
git tag modules/jmespath/$VERSION
git tag modules/kubernetes/$VERSION
git tag modules/pgx/$VERSION
git tag modules/playwright/$VERSION
git tag modules/qrcode/$VERSION
git tag modules/redis/$VERSION
git tag modules/sched/$VERSION
git tag modules/semver/$VERSION
git tag modules/shlex/$VERSION
git tag modules/slack/$VERSION
git tag modules/sql/$VERSION
git tag modules/ssh/$VERSION
git tag modules/tablewriter/$VERSION
git tag modules/template/$VERSION
git tag modules/uuid/$VERSION
git tag modules/vault/$VERSION
git tag modules/yaml/$VERSION
git tag os/s3fs/$VERSION

git push origin $VERSION
git push origin cmd/risor/$VERSION
git push origin modules/aws/$VERSION
git push origin modules/bcrypt/$VERSION
git push origin modules/cli/$VERSION
git push origin modules/color/$VERSION
git push origin modules/gha/$VERSION
git push origin modules/github/$VERSION
git push origin modules/goquery/$VERSION
git push origin modules/htmltomarkdown/$VERSION
git push origin modules/image/$VERSION
git push origin modules/isatty/$VERSION
git push origin modules/jmespath/$VERSION
git push origin modules/kubernetes/$VERSION
git push origin modules/pgx/$VERSION
git push origin modules/playwright/$VERSION
git push origin modules/qrcode/$VERSION
git push origin modules/redis/$VERSION
git push origin modules/sched/$VERSION
git push origin modules/semver/$VERSION
git push origin modules/shlex/$VERSION
git push origin modules/slack/$VERSION
git push origin modules/sql/$VERSION
git push origin modules/ssh/$VERSION
git push origin modules/tablewriter/$VERSION
git push origin modules/template/$VERSION
git push origin modules/uuid/$VERSION
git push origin modules/vault/$VERSION
git push origin modules/yaml/$VERSION
git push origin os/s3fs/$VERSION
