#!/bin/bash
# Builds binaries for all architectures with secefied version. It then puts it in bin directory

if [[ -z $VERSION ]]; then 
    echo "Error VERSION env var is not set. Failing to continue.";
    exit 1
fi

BINARY_NAME="rdma-ds-v${VERSION}"
SERVER_DIR="src/server"

mkdir -p bin
for GOOS in darwin linux windows; do
    for GOARCH in 386 amd64; do
        TMP_BINARY_NAME="$BINARY_NAME-$GOOS-$GOARCH"
        ( cd $SERVER_DIR && go build -v -o $TMP_BINARY_NAME)
        if [ ! -f "$SERVER_DIR/$TMP_BINARY_NAME" ]; then
            echo "Could not find server binary labelled '$TMP_BINARY_NAME' to move."
            exit 1
        else 
            mv "$SERVER_DIR/$TMP_BINARY_NAME" bin
        fi
        
    done
done

echo "Completed building, binary labelled '$BINARY_NAME' can be found in bin."