# rdma-ds

This dameon-set will run on every node and will be a RESTful API endpoint for querying RDMA SRIOV data about the node that it is currently running on. There are two processes that exist within this repo. 
1. init - the init processes will startup all sriov devices, this processes should only be run once when starting up the server to enable all SRIOV devices
2. server - the server processes will open up the RESTful API endpoint for querying data about the avialable devices and SRIOV virtual functions (VF's)

## Environmental Variables
The following variables are able to be modefied to change the port that the server is running on:
  - PORT - default 54005 - the port number for the server to run on
Example:
```
export PORT=656565
```

## Building Binaries (for testing or for use)

### Client
To use the client library import `github.com/rit-k8s-rdma/rit-k8s-rdma-ds/src` and utilize the built in supported functions after starting up the init and then the server process.

### Building and Running Init Process
To build the init process:
```
cd src/init
go install
```
A binary should have been built and saved into your `$GOPATH/bin/` directory under the name `init`.

To run it:
```
$GOPATH/init
```
*Note* depending on the amount of VF's you have configured, this could take some time to bind and unbind them (aka sit back and sip your coffee while you wait).

### Building and Running Server Process
The server process should start *after* the init process has successfully ran. To build the server process:
```
cd src/server
go install
```
A binary should have been built and saved into your `$GOPATH/bin/` directory under the name `server`.

To run it:
```
$GOPATH/server
```
This process will startup a server on port 54005 unless you exported another environmental variable for the port. You query it via `localhost:<port-num>` and to get your vf information `localhost:<port-num>/getpfs`.

## Docker
To build the docker image run the following in the root directory of this repo:
```
VERSION=<version> ./build_docker_local.sh
```
This will build two images:
1. `rdma-ds-init:<version>` - this image is the init image that only needs to be run once to set up all the VF's.

To run the `rdma-ds-init` image run the following command:
```
docker run -it --rm --name rdma-ds-init --privileged rdma-ds-init:<version>
```
This will do the following:
  - `-it` - runs in interactive mode
  - `--rm` - removes the container when stopped
  - `--privileged` - this container must be run in privileged mode b/c it is access network resources

2. `rdma-ds-server:<version>` - this image is the server image, which starts up an RESTful API for the VF information

To run the `rdma-ds-server` image run the following command:
```
docker run -it --rm --name test -e PORT=5000 --network host rdma-ds-server:<version>
```
This will do the following:
  - `-it` - runs in interactive mode
  - `--rm` - removes the container when stopped
  - `-e PORT=5000` - specifies the port for the server to run on, if none is specified it will default to port 54005
  - `--network host` - will share the network with your host OS, so you can access the api by going to localhost:5000


### Avoiding Dockerhub
If you want to avoid docker hub completely, you can save the image in a tar and than load it in.
Save command:
```
docker save <image-name>:<image-tag> > <save-name>.tar
```
Load command:
```
docker load < <save-name>.tar
```
Ex:
```
docker save rdma-ds-server:latest > rdma-ds-server.tar
docker load < rdma-ds-server.tar
```

## Kubernetes
In order to run this application in kuberentes you must execute the following command:
```
kubectl apply -f rdma-ds-yaml
```
This will spin up this repo with the dameonset. Note the current implementation you must first put the docker images on the remote nodes before running because this is a private repo.

Deleting dameonset:
```
kubectl delete ds/<dameonset-name> --namespace kube-system
```

Updating dameonset by changing the image:
```
kubectl set image ds/<daemonset-name> <container-name>=<container-new-image> --namespace=kube-system
```

Check updating status
```
kubectl rollout status ds/<daemonset-name> --namespace=kube-system
```
