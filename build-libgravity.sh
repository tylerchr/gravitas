#!/bin/sh

# builds libgravity.dylib
# this solutions feels pretty hacky but I'm not sure of a better way right now

BASEDIR=$(dirname "$0")

pushd ${BASEDIR}/gravity
make lib
popd
