#!/bin/bash

GITHUB=~/src/github.com

if ! [ -e "$GITHUB" ] ; then
    echo "mkdir -p $GITHUB"
    mkdir -p $GITHUB
fi

# Create GOBIN DIR for executables if it doesn't already exist
if ! [ -e "$GOBIN" ] ; then
    echo "mkdir -p $GOBIN"
    mkdir -p $GOBIN
fi

# Set directories to install goFish
goFishDir=${GITHUB}/edotau
mkdir -p ${goFishDir}
mv ${GITHUB}/goFish/ ${goFishDir}/

