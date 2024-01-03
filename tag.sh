#!/bin/bash

set -e

VERSION=$1

if [ -z "$VERSION" ]; then
    echo "Usage: tag.sh <version>"
    exit 1
fi

git tag $VERSION
git tag os/s3fs/$VERSION
git tag modules/uuid/$VERSION
git tag modules/pgx/$VERSION
git tag modules/image/$VERSION
git tag modules/aws/$VERSION
git tag modules/kubernetes/$VERSION
git tag cmd/risor/$VERSION

git push origin $VERSION
git push origin os/s3fs/$VERSION
git push origin modules/uuid/$VERSION
git push origin modules/pgx/$VERSION
git push origin modules/sql/$VERSION
git push origin modules/image/$VERSION
git push origin modules/aws/$VERSION
git push origin modules/kubernetes/$VERSION
git push origin cmd/risor/$VERSION
