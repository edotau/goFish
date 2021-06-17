#!/bin/bash
#set -eou pipefail
release_list=https://golang.org/dl/
SOURCE=https://storage.googleapis.com/golang
# TODO Refine logic for setting up GOPATH in other locations
export DIR=$HOME
cd $DIR
# Figures out what OS you are using
os=$(uname -s | tr "[:upper:]" "[:lower:]")

# Determines the arch distribution of your current machine using switch cases
arch_probe="uname -m"
arch=$($arch_probe)
case "$arch" in
    i*)
        arch=386
        ;;
    x*)
        arch=amd64
        ;;
    aarch64)
        #arch=armv6l
        arch=arm64
        ;;
    armv7l)
        # Go project does not provide a binary release for armv71
        echo 'armv7l is not supported, using armv6l'
        arch=armv6l
        ;;
esac

# Performs a quick query for golang releases, excluding beta distributions
queryReleaseList() {
    local fetch="$*"
    if hash "jq" 2>/dev/null; then
        local rl="$release_list?mode=json"
        $fetch "$rl" | jq -r '.[].files[].version' | sort | uniq | grep -v -E 'go[0-9\.]+(beta|rc)'| sed -e 's/go//' | sort -V
    else
        $fetch "$release_list" | grep -v -E 'go[0-9\.]+(beta|rc)' | grep -E -o 'go[0-9\.]+' | grep -E -o '[0-9]\.[0-9]+(\.[0-9]+)?' | sort -V | uniq
    fi
}
# From the query to fetch all golang releases, we pull the latest stable release
fetchUpdate() {
    local last=
    local fetch=
    if hash "wget" 2>/dev/null; then
	    fetch="wget -qO-"
    elif hash "curl" 2>/dev/null; then
	    fetch="curl --silent"
    else
        echo "Error: local wget and/or curl is not functional..."
        exit 3
    fi
    last=$(queryReleaseList "$fetch" | tail -1)
    if echo "$last" | grep -q -E '[0-9]\.[0-9]+(\.[0-9]+)?'; then
	    release=$last
    fi
}

# Performs the function call to pull the latest stable release
fetchUpdate

# Set our variables
label=go$release.$os-$arch
FILENAME=$label.tar.gz
URL=$SOURCE/$FILENAME

# Download and decompress golang stable release
wget $URL
tar xf $FILENAME -C $HOME; rm $FILENAME

if ! [ -e "/usr/local/go" ] ; then
    mv go/ /usr/local/
fi

# Sets up bash profile for go root, path, and bin
# Note that it is assumed your GOPATH will be placed in your home directory
if [ "$SHELL" == "/bin/bash" ]; then
    echo 'export GOROOT=/usr/local/go' >> ~/.bashrc
    echo 'export GOPATH=$HOME' >> ~/.bashrc
    echo 'export GOBIN=${GOPATH}/bin' >> ~/.bashrc
    echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.bashrc
    source ~/.bashrc
elif [ "$SHELL" == "*zsh*" ]; then
# Sets up zsh profile for go root, path, and bin if your default shell is zsh
    echo 'export GOROOT=/usr/local/go' >> ~/.zshrc
    echo 'export GOPATH=$HOME' >> ~/.zshrc
    echo 'export GOBIN=${GOPATH}/bin' >> ~/.zshrc
    echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.zshrc
    source ~/.zshrc
else
    echo 'export GOROOT=/usr/local/go' >> ~/.env
    echo 'export GOPATH=$HOME' >> ~/.env
    echo 'export GOBIN=${GOPATH}/bin' >> ~/.env
    echo 'export PATH=$GOPATH/bin:$GOROOT/bin:$PATH' >> ~/.env
    source ~/.env
fi

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

cd $goFishDir
# Contingent on passing all go tests, we finally install/build all goFish binary programs
go mod init; go mod tidy
go test ./...
go install ./...
ln -s $GOBIN
