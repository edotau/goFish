#!/bin/bash

DIR=$1
# Sets up bash profile for go root, path, and bin
# Note that it is assumed your GOPATH will be placed in your home directory
echo 'export GOROOT=${DIR}/go' >> ~/.bashrc
echo 'export GOPATH=$HOME' >> ~/.bashrc
echo 'export GOBIN=${GOPATH}/bin' >> ~/.bashrc
echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.bashrc
source ~/.bashrc

