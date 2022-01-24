#!/bin/bash

# Sets up bash profile for go root, path, and bin
# Note that it is assumed your GOPATH will be placed in your home directory
if [ "$SHELL" == "/bin/bash" ]; then
    echo 'export GOROOT=~/go' >> ~/.bashrc
    echo 'export GOPATH=$HOME' >> ~/.bashrc
    echo 'export GOBIN=${GOPATH}/bin' >> ~/.bashrc
    echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.bashrc
    source ~/.bashrc
elif [ "$SHELL" == "*zsh*" ]; then
# Sets up zsh profile for go root, path, and bin if your default shell is zsh
    echo 'export GOROOT=~/go' >> ~/.zshrc
    echo 'export GOPATH=$HOME' >> ~/.zshrc
    echo 'export GOBIN=${GOPATH}/bin' >> ~/.zshrc
    echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.zshrc
    source ~/.zshrc
else
    echo 'export GOROOT=~/go' >> ~/.env
    echo 'export GOPATH=$HOME' >> ~/.env
    echo 'export GOBIN=${GOPATH}/bin' >> ~/.env
    echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.env
    source ~/.env
fi

