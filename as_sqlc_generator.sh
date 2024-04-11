#!/bin/sh -x

internaldir="internal"
gendir="pkg/gen/db"
authgendir="pkg/gen/authdb"
echo "./${internaldir}/${gendir}"
ls -al "./${internaldir}/${gendir}"

echo "./$internaldir/$gendir/"
rm -rf "./${internaldir}/${gendir}"
mkdir -p "./${internaldir}/${gendir}"

/opt/homebrew/bin/sqlc generate -f ./sqlc/sqlc.yaml

