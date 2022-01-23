#!/bin/bash
#set -eou pipefail
release_list=https://golang.org/dl/
SOURCE=https://storage.googleapis.com/golang

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


