if [[ -z $VERSION ]]; then 
    echo "Error VERSION env var is not set. Failing to continue.";
    exit 1
fi

docker build -t rdma-ds-init:$VERSION -f src/init/Dockerfile .

docker build -t rdma-ds-server:$VERSION -f src/server/Dockerfile .
